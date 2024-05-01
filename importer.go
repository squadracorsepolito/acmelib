package acmelib

import (
	"fmt"
	"io"
	"slices"

	"github.com/FerroO2000/acmelib/dbc"
)

func ImportDBCFile(filename string, r io.Reader) (*Bus, error) {
	dbcFile, err := dbc.Parse(filename, r)
	if err != nil {
		return nil, err
	}
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

	nodes    map[string]*Node
	messages map[MessageCANID]*Message
	signals  map[string]Signal

	signalEnums map[string]*SignalEnum

	dbcExtMuxes map[string]*dbc.ExtendedMux
}

func newImporter() *importer {
	return &importer{
		bus: nil,

		nodeDesc: make(map[string]string),
		msgDesc:  make(map[MessageCANID]string),
		sigDesc:  make(map[string]string),

		nodes:    make(map[string]*Node),
		messages: make(map[MessageCANID]*Message),
		signals:  make(map[string]Signal),

		signalEnums: make(map[string]*SignalEnum),

		dbcExtMuxes: make(map[string]*dbc.ExtendedMux),
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

	i.importExtMuxes(dbcFile.ExtendedMuxes)

	if err := i.importNodes(dbcFile.Nodes); err != nil {
		return nil, err
	}

	for _, dbcMsg := range dbcFile.Messages {
		if err := i.importMessage(dbcMsg); err != nil {
			return nil, err
		}
	}

	if err := i.importAttributes(dbcFile.Attributes, dbcFile.AttributeDefaults, dbcFile.AttributeValues); err != nil {
		return nil, err
	}

	dummyNode, err := bus.GetNodeByName(dbc.DummyNode)
	if err != nil {
		panic(err)
	}

	if len(dummyNode.Messages()) == 0 {
		if err := bus.RemoveNode(dummyNode.entityID); err != nil {
			panic(err)
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

func (i *importer) importAttributes(dbcAtts []*dbc.Attribute, dbcAttDefs []*dbc.AttributeDefault, dbcAttVals []*dbc.AttributeValue) error {
	dbcAttDefMap := make(map[string]*dbc.AttributeDefault)
	for _, dbcAttDef := range dbcAttDefs {
		dbcAttDefMap[dbcAttDef.AttributeName] = dbcAttDef
	}

	attributes := make(map[string]Attribute)
	for _, dbcAtt := range dbcAtts {
		dbcAttDef, ok := dbcAttDefMap[dbcAtt.Name]
		if !ok {
			return i.errorf(dbcAtt, &ErrIsRequired{Thing: "attribute default"})
		}

		var att Attribute
		switch dbcAtt.Type {
		case dbc.AttributeString:
			att = NewStringAttribute(dbcAtt.Name, dbcAttDef.ValueString)

		case dbc.AttributeInt:
			intAtt, err := NewIntegerAttribute(dbcAtt.Name, dbcAttDef.ValueInt, dbcAtt.MinInt, dbcAtt.MaxInt)
			if err != nil {
				return i.errorf(dbcAtt, err)
			}
			att = intAtt

		case dbc.AttributeHex:
			hexAtt, err := NewIntegerAttribute(dbcAtt.Name, int(dbcAttDef.ValueHex), int(dbcAtt.MinHex), int(dbcAtt.MaxHex))
			if err != nil {
				return i.errorf(dbcAtt, err)
			}
			hexAtt.SetFormatHex()
			att = hexAtt

		case dbc.AttributeFloat:
			floatAtt, err := NewFloatAttribute(dbcAtt.Name, dbcAttDef.ValueFloat, dbcAtt.MinFloat, dbcAtt.MaxFloat)
			if err != nil {
				return i.errorf(dbcAtt, err)
			}
			att = floatAtt

		case dbc.AttributeEnum:
			enumAtt, err := NewEnumAttribute(dbcAtt.Name, dbcAtt.EnumValues...)
			if err != nil {
				return i.errorf(dbcAtt, err)
			}
			att = enumAtt
		}

		attributes[att.Name()] = att
	}

	for _, dbcAttVal := range dbcAttVals {
		attName := dbcAttVal.AttributeName

		att, ok := attributes[attName]
		if !ok {
			continue
		}

		var value any
		switch dbcAttVal.Type {
		case dbc.AttributeValueString:
			value = dbcAttVal.ValueString

		case dbc.AttributeValueInt:
			if att.Kind() == AttributeKindEnum {
				enumAtt, err := att.ToEnum()
				if err != nil {
					panic(err)
				}

				strVal, err := enumAtt.GetValueAtIndex(dbcAttVal.ValueInt)
				if err != nil {
					i.errorf(dbcAttVal, err)
				}
				value = strVal

				break
			}

			if att.Kind() == AttributeKindFloat {
				value = float64(dbcAttVal.ValueInt)
				break
			}

			value = dbcAttVal.ValueInt

		case dbc.AttributeValueHex:
			value = int(dbcAttVal.ValueHex)

		case dbc.AttributeValueFloat:
			value = dbcAttVal.ValueFloat
		}

		switch dbcAttVal.AttributeKind {
		case dbc.AttributeGeneral:
			if err := i.bus.AddAttributeValue(att, value); err != nil {
				return i.errorf(dbcAttVal, err)
			}

		case dbc.AttributeNode:
			if node, ok := i.nodes[dbcAttVal.NodeName]; ok {
				if err := node.AddAttributeValue(att, value); err != nil {
					return i.errorf(dbcAttVal, err)
				}
			}

		case dbc.AttributeMessage:
			if msg, ok := i.messages[MessageCANID(dbcAttVal.MessageID)]; ok {
				attType, ok := specialAttributeTypes[attName]
				if ok {
					switch attType {
					case specialAttributeMsgCycleTime:
						msg.SetCycleTime(value.(int))

					case specialAttributeMsgDelayTime:
						msg.SetDelayTime(value.(int))

					case specialAttributeMsgStartDelayTime:
						msg.SetStartDelayTime(value.(int))

					case specialAttributeMsgSendType:
						msg.SetSendType(messageSendTypeFromString(value.(string)))
					}

					break
				}

				if err := msg.AddAttributeValue(att, value); err != nil {
					return i.errorf(dbcAttVal, err)
				}
			}

		case dbc.AttributeSignal:
			if sig, ok := i.signals[i.getSignalKey(dbcAttVal.MessageID, dbcAttVal.SignalName)]; ok {
				attType, ok := specialAttributeTypes[attName]
				if ok {
					switch attType {
					case specialAttributeSigSendType:
						sig.SetSendType(signalSendTypeFromString(value.(string)))
					}

					break
				}

				if err := sig.AddAttributeValue(att, value); err != nil {
					return i.errorf(dbcAttVal, err)
				}
			}

		}
	}

	return nil
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

		i.nodes[tmpNode.name] = tmpNode
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

	i.messages[msgID] = msg

	if desc, ok := i.msgDesc[msgID]; ok {
		msg.SetDesc(desc)
	}

	receivers := make(map[string]bool)
	muxSignals := []*dbc.Signal{}
	muxSigNames := make(map[string]int)

	slices.SortFunc(dbcMsg.Signals, func(a, b *dbc.Signal) int { return int(a.StartBit) - int(b.StartBit) })

	for _, dbcSig := range dbcMsg.Signals {
		if dbcSig.IsMultiplexor {
			muxSignals = append(muxSignals, dbcSig)
			muxSigNames[dbcSig.Name] = len(muxSignals) - 1
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

			if err := msg.InsertSignal(tmpSig, i.getSignalStartBit(dbcSig)); err != nil {
				return i.errorf(dbcSig, err)
			}
		}
		return nil
	}

	if muxSigCount == 1 {
		dbcMuxSig := muxSignals[0]

		muxedSignals := []*importerMuxedSignal{}
		for _, dbcSig := range dbcMsg.Signals {
			if dbcSig.Name == dbcMuxSig.Name {
				continue
			}

			tmpSig, err := i.importSignal(dbcSig, dbcMsg.ID)
			if err != nil {
				return err
			}

			if dbcSig.IsMultiplexed {
				muxedSignals = append(muxedSignals, &importerMuxedSignal{sig: tmpSig, dbcSig: dbcSig})
				continue
			}

			if err := msg.InsertSignal(tmpSig, i.getSignalStartBit(dbcSig)); err != nil {
				return i.errorf(dbcSig, err)
			}
		}

		muxSig, err := i.importMuxSignal(dbcMuxSig, dbcMsg.ID, muxedSignals)
		if err != nil {
			return err
		}

		if err := msg.InsertSignal(muxSig, i.getSignalStartBit(dbcMuxSig)); err != nil {
			return i.errorf(dbcMuxSig, err)
		}

		return nil
	}

	muxedSigGroups := make([][]*importerMuxedSignal, muxSigCount)
	for _, dbcSig := range dbcMsg.Signals {
		if _, ok := muxSigNames[dbcSig.Name]; ok {
			continue
		}

		tmpSig, err := i.importSignal(dbcSig, dbcMsg.ID)
		if err != nil {
			return err
		}

		if dbcSig.IsMultiplexed {
			dbcExtMux, ok := i.dbcExtMuxes[i.getSignalKey(dbcMsg.ID, dbcSig.Name)]
			if !ok {
				return i.errorf(dbcSig, &ErrIsRequired{Thing: "extended multiplexing"})
			}

			muxIdx, ok := muxSigNames[dbcExtMux.MultiplexorName]
			if !ok {
				return i.errorf(dbcExtMux, &NameError{Name: dbcExtMux.MultiplexorName, Err: ErrNotFound})
			}

			muxedSigGroups[muxIdx] = append(muxedSigGroups[muxIdx], &importerMuxedSignal{
				sig:    tmpSig,
				dbcSig: dbcSig,
			})

			continue
		}

		if err := msg.InsertSignal(tmpSig, i.getSignalStartBit(dbcSig)); err != nil {
			return i.errorf(dbcSig, err)
		}
	}

	for j := muxSigCount - 1; j >= 0; j-- {
		dbcMuxSig := muxSignals[j]

		muxSig, err := i.importMuxSignal(dbcMuxSig, dbcMsg.ID, muxedSigGroups[j])
		if err != nil {
			return err
		}

		dbcExtMux, ok := i.dbcExtMuxes[i.getSignalKey(dbcMsg.ID, dbcMuxSig.Name)]
		if !ok {
			if err := msg.InsertSignal(muxSig, i.getSignalStartBit(dbcMuxSig)); err != nil {
				return i.errorf(dbcMuxSig, err)
			}
			continue
		}

		muxIdx, ok := muxSigNames[dbcExtMux.MultiplexorName]
		if !ok {
			return i.errorf(dbcExtMux, &NameError{Name: dbcExtMux.MultiplexorName, Err: ErrNotFound})
		}

		muxedSigGroups[muxIdx] = append(muxedSigGroups[muxIdx], &importerMuxedSignal{
			sig:    muxSig,
			dbcSig: dbcMuxSig,
		})
	}

	return nil
}

func (i *importer) getSignalStartBit(dbcSig *dbc.Signal) int {
	startBit := int(dbcSig.StartBit)
	size := int(dbcSig.Size)
	if dbcSig.ByteOrder == dbc.SignalBigEndian {
		startBit -= size - 1
	}
	return startBit
}

type importerMuxedSignal struct {
	sig    Signal
	dbcSig *dbc.Signal
}

func (i *importer) importMuxSignal(dbcMuxSig *dbc.Signal, dbcMsgID uint32, muxedSignals []*importerMuxedSignal) (*MultiplexerSignal, error) {
	lastMuxedSig := muxedSignals[len(muxedSignals)-1]

	lastSize := lastMuxedSig.sig.GetSize()

	lastStartBit := int(lastMuxedSig.dbcSig.StartBit)
	if lastMuxedSig.dbcSig.ByteOrder == dbc.SignalBigEndian {
		lastStartBit -= lastSize - 1
	}

	muxSigStartBit := i.getSignalStartBit(dbcMuxSig)

	muxSigSize := int(dbcMuxSig.Size)
	groupSize := lastStartBit + lastSize - muxSigStartBit - muxSigSize

	muxSig, err := NewMultiplexerSignal(dbcMuxSig.Name, calcValueFromSize(muxSigSize), groupSize)
	if err != nil {
		return nil, i.errorf(dbcMuxSig, err)
	}

	for _, muxedSig := range muxedSignals {
		tmpSig := muxedSig.sig
		tmpDBCSig := muxedSig.dbcSig

		tmpStartBit := i.getSignalStartBit(tmpDBCSig)
		relStartBit := tmpStartBit - muxSigStartBit - muxSigSize

		groupIDs := []int{}
		dbcExtMux, ok := i.dbcExtMuxes[i.getSignalKey(dbcMsgID, tmpSig.Name())]
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
			groupIDs = append(groupIDs, int(tmpDBCSig.MuxSwitchValue))
		}

		if err := muxSig.InsertSignal(tmpSig, relStartBit, groupIDs...); err != nil {
			return nil, i.errorf(tmpDBCSig, err)
		}
	}

	sigKey := i.getSignalKey(dbcMsgID, muxSig.name)

	if desc, ok := i.sigDesc[sigKey]; ok {
		muxSig.SetDesc(desc)
	}

	if dbcMuxSig.ByteOrder == dbc.SignalBigEndian {
		muxSig.SetByteOrder(SignalByteOrderBigEndian)
	}

	i.signals[sigKey] = muxSig

	return muxSig, nil
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

	i.signals[sigKey] = sig

	return sig, nil
}
