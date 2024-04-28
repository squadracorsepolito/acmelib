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
}

func newImporter() *importer {
	return &importer{
		bus: nil,

		nodeDesc: make(map[string]string),
		msgDesc:  make(map[MessageCANID]string),
		sigDesc:  make(map[string]string),

		signalEnums: make(map[string]*SignalEnum),
	}
}

func (i *importer) errorf(dbcLoc dbcFileLocator, err error) error {
	return fmt.Errorf("%s : %w", dbcLoc.Location(), err)
}

func (i *importer) getSignalKey(dbcMsgID uint32, sigName string) string {
	return fmt.Sprintf("%d_%s", dbcMsgID, sigName)
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

	} else if dbcSig.IsMultiplexor {
		// mux signal
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
