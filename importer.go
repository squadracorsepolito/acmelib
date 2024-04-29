package acmelib

import (
	"fmt"

	"github.com/FerroO2000/acmelib/dbc"
)

func ImportDBCFile(dbcFile *dbc.File) (*Bus, error) {
	importer := newImporter()
	return importer.importFile(dbcFile)
}

type dbcFileLocator interface {
	Location() *dbc.Location
}

type importer struct {
	bus *Bus

	nodeDesc map[string]string
	msgDesc  map[MessageCANID]string
	sigDesc  map[string]string

	signalEnums map[string]*SignalEnum

	extMuxMsg   map[uint32]bool
	dbcExtMuxes map[string]*dbc.ExtendedMux
}

func newImporter() *importer {
	return &importer{
		bus: nil,

		nodeDesc: make(map[string]string),
		msgDesc:  make(map[MessageCANID]string),
		sigDesc:  make(map[string]string),

		signalEnums: make(map[string]*SignalEnum),

		extMuxMsg:   make(map[uint32]bool),
		dbcExtMuxes: make(map[string]*dbc.ExtendedMux),
	}
}

func (i *importer) errorf(dbcLoc dbcFileLocator, err error) error {
	return fmt.Errorf("%s : %w", dbcLoc.Location(), err)
}

func (i *importer) getSignalKey(dbcMsgID uint32, sigName string) string {
	return fmt.Sprintf("%d_%s", dbcMsgID, sigName)
}

func (i *importer) getMuxSignalKey(dbcMsgID uint32, muxorName, muxedName string) string {
	return fmt.Sprintf("%d_%s_%s", dbcMsgID, muxorName, muxedName)
}

func (i *importer) importFile(dbcFile *dbc.File) (*Bus, error) {
	bus := NewBus(dbcFile.Location().Filename)
	i.bus = bus

	i.importComments(dbcFile.Comments)

	for _, dbcValEnc := range dbcFile.ValueEncodings {
		if err := i.importValueEncoding(dbcValEnc); err != nil {
			return nil, err
		}
	}

	i.importExtMuxes(dbcFile.ExtendedMuxes)

	if err := i.importNodes(dbcFile.Nodes); err != nil {
		return nil, err
	}

	for _, dbcMsg := range dbcFile.Messages {
		if err := i.importMessage(dbcMsg); err != nil {
			return nil, err
		}
	}

	return bus, nil
}

func (i *importer) importComments(dbcComments []*dbc.Comment) {
	for _, dbcComm := range dbcComments {
		switch dbcComm.Kind {
		case dbc.CommentGeneral:
			i.bus.SetDesc(dbcComm.Text)

		case dbc.CommentNode:
			i.nodeDesc[dbcComm.NodeName] = dbcComm.Text

		case dbc.CommentMessage:
			i.msgDesc[MessageCANID(dbcComm.MessageID)] = dbcComm.Text

		case dbc.CommentSignal:
			key := i.getSignalKey(dbcComm.MessageID, dbcComm.SignalName)
			i.sigDesc[key] = dbcComm.Text
		}
	}
}

func (i *importer) importValueEncoding(dbcValEnc *dbc.ValueEncoding) error {
	if dbcValEnc.Kind != dbc.ValueEncodingSignal {
		return nil
	}

	sigName := dbcValEnc.SignalName
	sigEnum := NewSignalEnum(fmt.Sprintf("%s_Enum", sigName))
	for _, dbcVal := range dbcValEnc.Values {
		if err := sigEnum.AddValue(NewSignalEnumValue(dbcVal.Name, int(dbcVal.ID))); err != nil {
			return i.errorf(dbcVal, err)
		}
	}

	i.signalEnums[i.getSignalKey(dbcValEnc.MessageID, sigName)] = sigEnum

	return nil
}

func (i *importer) importExtMuxes(dbcExtMuxes []*dbc.ExtendedMux) {
	for _, tmpExtMux := range dbcExtMuxes {
		key := i.getSignalKey(tmpExtMux.MessageID, tmpExtMux.MultiplexedName)
		i.dbcExtMuxes[key] = tmpExtMux
		i.extMuxMsg[tmpExtMux.MessageID] = true
	}
}

func (i *importer) importNodes(dbcNodes *dbc.Nodes) error {
	for idx, nodeName := range dbcNodes.Names {
		tmpNode := NewNode(nodeName, NodeID(idx))

		if desc, ok := i.nodeDesc[nodeName]; ok {
			tmpNode.SetDesc(desc)
		}

		if err := i.bus.AddNode(tmpNode); err != nil {
			return i.errorf(dbcNodes, err)
		}
	}

	if err := i.bus.AddNode(NewNode(dbc.DummyNode, 1024)); err != nil {
		return i.errorf(dbcNodes, err)
	}

	return nil
}

func (i *importer) importMessage(dbcMsg *dbc.Message) error {
	msg := NewMessage(dbcMsg.Name, int(dbcMsg.Size))

	msgID := MessageCANID(dbcMsg.ID)
	msg.SetCANID(msgID)

	if desc, ok := i.msgDesc[msgID]; ok {
		msg.SetDesc(desc)
	}

	receivers := make(map[string]bool)
	muxSignals := []*dbc.Signal{}
	for _, dbcSig := range dbcMsg.Signals {
		if dbcSig.IsMultiplexor {
			muxSignals = append(muxSignals, dbcSig)
		}

		for _, rec := range dbcSig.Receivers {
			receivers[rec] = true
		}
	}

	for recName := range receivers {
		if recName == dbc.DummyNode {
			continue
		}

		recNode, err := i.bus.GetNodeByName(recName)
		if err != nil {
			return i.errorf(dbcMsg, err)
		}

		msg.AddReceiver(recNode)
	}

	sendNode, err := i.bus.GetNodeByName(dbcMsg.Transmitter)
	if err != nil {
		return i.errorf(dbcMsg, err)
	}

	if err := sendNode.AddMessage(msg); err != nil {
		return i.errorf(dbcMsg, err)
	}

	muxSigCount := len(muxSignals)
	if muxSigCount == 0 {
		for _, dbcSig := range dbcMsg.Signals {
			tmpSig, err := i.importSignal(dbcSig, dbcMsg.ID)
			if err != nil {
				return err
			}

			startBit := int(dbcSig.StartBit)
			if tmpSig.ByteOrder() == SignalByteOrderBigEndian {
				startBit -= (tmpSig.GetSize() - 1)
			}

			if err := msg.InsertSignal(tmpSig, startBit); err != nil {
				return i.errorf(dbcSig, err)
			}
		}

		return nil
	} else if muxSigCount == 1 {
		lastStartBit := 0
		lastSize := 0
		for _, dbcSig := range dbcMsg.Signals {
			if dbcSig.IsMultiplexed {
				lastStartBit = int(dbcSig.StartBit)
				lastSize = int(dbcSig.Size)
			}
		}

		dbcMuxSig := muxSignals[0]
		muxSigStartBit := int(dbcMuxSig.StartBit)
		groupSize := lastStartBit + lastSize - muxSigStartBit - int(dbcMuxSig.Size)

		muxSig, err := NewMultiplexerSignal(dbcMuxSig.Name, calcValueFromSize(int(dbcMuxSig.Size)), groupSize)
		if err != nil {
			i.errorf(dbcMuxSig, err)
		}

		for _, dbcSig := range dbcMsg.Signals {
			if dbcSig.Name == muxSig.name {
				continue
			}

			tmpSig, err := i.importSignal(dbcSig, dbcMsg.ID)
			if err != nil {
				return err
			}

			startBit := int(dbcSig.StartBit)
			if tmpSig.ByteOrder() == SignalByteOrderBigEndian {
				startBit -= (tmpSig.GetSize() - 1)
			}

			relStartBit := startBit - int(dbcMuxSig.StartBit) - int(dbcMuxSig.Size)

			if dbcSig.IsMultiplexed {
				groupIDs := []int{}

				dbcExtMux, ok := i.dbcExtMuxes[i.getSignalKey(dbcMsg.ID, dbcSig.Name)]
				if ok {
					for _, valRange := range dbcExtMux.Ranges {
						for j := valRange.From; j <= valRange.To; j++ {
							groupIDs = append(groupIDs, int(j))
						}
					}

					if len(groupIDs) == muxSig.groupCount {
						groupIDs = []int{}
					}

				} else {
					groupIDs = append(groupIDs, int(dbcSig.MuxSwitchValue))
				}

				if err := muxSig.InsertSignal(tmpSig, relStartBit, groupIDs...); err != nil {
					return i.errorf(dbcSig, err)
				}

				continue
			}

			if startBit > muxSigStartBit && startBit <= lastSize {
				if err := muxSig.InsertSignal(tmpSig, relStartBit); err != nil {
					return i.errorf(dbcSig, err)
				}
				continue
			}

			if err := msg.InsertSignal(tmpSig, startBit); err != nil {
				return i.errorf(dbcSig, err)
			}
		}

		if err := msg.InsertSignal(muxSig, muxSigStartBit); err != nil {
			return i.errorf(dbcMuxSig, err)
		}
	}

	return nil
}

func (i *importer) importSignal(dbcSig *dbc.Signal, dbcMsgID uint32) (Signal, error) {
	var sig Signal

	sigName := dbcSig.Name
	sigKey := i.getSignalKey(dbcMsgID, sigName)

	if sigEnum, ok := i.signalEnums[sigKey]; ok {
		enumSig, err := NewEnumSignal(sigName, sigEnum)
		if err != nil {
			return nil, i.errorf(dbcSig, err)
		}
		sig = enumSig

	} else {
		signed := false
		if dbcSig.ValueType == dbc.SignalSigned {
			signed = true
		}

		sigType, err := NewIntegerSignalType(fmt.Sprintf("%s_Type", sigName), int(dbcSig.Size), signed)
		if err != nil {
			return nil, i.errorf(dbcSig, err)
		}

		stdSig, err := NewStandardSignal(sigName, sigType)
		if err != nil {
			return nil, i.errorf(dbcSig, err)
		}

		if err := stdSig.SetPhysicalValues(dbcSig.Min, dbcSig.Max, dbcSig.Offset, dbcSig.Factor); err != nil {
			return nil, i.errorf(dbcSig, err)
		}

		if dbcSig.Unit != "" {
			stdSig.SetUnit(NewSignalUnit(fmt.Sprintf("%s_Unit", sigName), SignalUnitKindCustom, dbcSig.Unit))
		}

		sig = stdSig
	}

	if desc, ok := i.sigDesc[sigKey]; ok {
		sig.SetDesc(desc)
	}

	if dbcSig.ByteOrder == dbc.SignalBigEndian {
		sig.SetByteOrder(SignalByteOrderBigEndian)
	}

	return sig, nil
}
