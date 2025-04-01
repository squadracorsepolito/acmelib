package acmelib

import (
	"fmt"
	"io"
	"slices"

	"github.com/squadracorsepolito/acmelib/dbc"
)

// ImportDBCFile imports a DBC file passed as [io.Reader] and converts it
// to a [Bus]. The given filename will be used as the name of the bus.
func ImportDBCFile(filename string, r io.Reader) (*Bus, error) {
	dbcFile, err := dbc.Parse(filename, r, false)
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
	msgDesc  map[MessageID]string
	sigDesc  map[string]string

	nodeInts map[string]*NodeInterface
	messages map[MessageID]*Message
	signals  map[string]Signal

	flagSigType *SignalType
	signalTypes map[string]*SignalType

	signalUnits map[string]*SignalUnit

	signalEnumRegistry []*SignalEnum
	signalEnums        map[string]*SignalEnum

	dbcExtMuxes map[string]*dbc.ExtendedMux
}

func newImporter() *importer {
	return &importer{
		bus: nil,

		nodeDesc: make(map[string]string),
		msgDesc:  make(map[MessageID]string),
		sigDesc:  make(map[string]string),

		nodeInts: make(map[string]*NodeInterface),
		messages: make(map[MessageID]*Message),
		signals:  make(map[string]Signal),

		flagSigType: NewFlagSignalType("flag_t"),
		signalTypes: make(map[string]*SignalType),

		signalUnits: make(map[string]*SignalUnit),

		signalEnumRegistry: []*SignalEnum{},
		signalEnums:        make(map[string]*SignalEnum),

		dbcExtMuxes: make(map[string]*dbc.ExtendedMux),
	}
}

func (i *importer) errorf(dbcLoc dbcFileLocator, err error) error {
	return fmt.Errorf("%s : %w", dbcLoc.Location(), err)
}

func (i *importer) getSignalKey(dbcMsgID uint32, sigName string) string {
	return fmt.Sprintf("%d_%s", dbcMsgID, sigName)
}

func (i *importer) getSignalTypeKey(dbcSig *dbc.Signal) string {
	signStr := "u"
	if dbcSig.ValueType == dbc.SignalSigned {
		signStr = "s"
	}
	return fmt.Sprintf("%s%d_%g-%g_%g_%g", signStr, dbcSig.Size, dbcSig.Min, dbcSig.Max, dbcSig.Factor, dbcSig.Offset)
}

func (i *importer) importFile(dbcFile *dbc.File) (*Bus, error) {
	bus := NewBus(dbcFile.Location().Filename)
	i.bus = bus

	i.importComments(dbcFile.Comments)

	for _, dbcValTable := range dbcFile.ValueTables {
		if err := i.importValueTable(dbcValTable); err != nil {
			return nil, err
		}
	}

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

	dummyNode, err := bus.GetNodeInterfaceByNodeName(dbc.DummyNode)
	if err != nil {
		panic(err)
	}

	if len(dummyNode.SentMessages()) == 0 {
		if err := bus.RemoveNodeInterface(dummyNode.node.entityID); err != nil {
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
			i.msgDesc[MessageID(dbcComm.MessageID)] = dbcComm.Text

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
			return i.errorf(dbcAtt, &ErrIsRequired{Item: "attribute default"})
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
			if att.Type() == AttributeTypeEnum {
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

			if att.Type() == AttributeTypeFloat {
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
			if err := i.bus.AssignAttribute(att, value); err != nil {
				return i.errorf(dbcAttVal, err)
			}

		case dbc.AttributeNode:
			if node, ok := i.nodeInts[dbcAttVal.NodeName]; ok {
				if err := node.node.AssignAttribute(att, value); err != nil {
					return i.errorf(dbcAttVal, err)
				}
			}

		case dbc.AttributeMessage:
			if msg, ok := i.messages[MessageID(dbcAttVal.MessageID)]; ok {
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
						msg.SetSendType(messageSendTypeFromDBC(value.(string)))
					}

					break
				}

				if err := msg.AssignAttribute(att, value); err != nil {
					return i.errorf(dbcAttVal, err)
				}
			}

		case dbc.AttributeSignal:
			if sig, ok := i.signals[i.getSignalKey(dbcAttVal.MessageID, dbcAttVal.SignalName)]; ok {
				attType, ok := specialAttributeTypes[attName]
				if ok {
					switch attType {
					case specialAttributeSigStartValue:
						switch value.(type) {
						case float64:
							sig.SetStartValue(value.(float64))

						case int:
							sig.SetStartValue(float64(value.(int)))
						}

					case specialAttributeSigSendType:
						sig.SetSendType(signalSendTypeFromDBC(value.(string)))
					}

					break
				}

				if err := sig.AssignAttribute(att, value); err != nil {
					return i.errorf(dbcAttVal, err)
				}
			}

		}
	}

	return nil
}

func (i *importer) importValueTable(dbcValTable *dbc.ValueTable) error {
	sigEnum := NewSignalEnum(dbcValTable.Name)

	for _, dbcVal := range dbcValTable.Values {
		if err := sigEnum.AddValue(NewSignalEnumValue(dbcVal.Name, int(dbcVal.ID))); err != nil {
			return i.errorf(dbcVal, err)
		}
	}

	i.signalEnumRegistry = append(i.signalEnumRegistry, sigEnum)

	return nil
}

func (i *importer) importValueEncoding(dbcValEnc *dbc.ValueEncoding) error {
	if dbcValEnc.Kind != dbc.ValueEncodingSignal {
		return nil
	}

	values := dbcValEnc.Values
	slices.SortFunc(values, func(a, b *dbc.ValueDescription) int { return int(a.ID) - int(b.ID) })

	sigEnum := new(SignalEnum)
	inRegistry := false
	for _, tmpSigEnum := range i.signalEnumRegistry {
		if len(values) != tmpSigEnum.values.size() {
			continue
		}

		for idx, tmpVal := range tmpSigEnum.Values() {
			if tmpVal.name != values[idx].Name || tmpVal.index != int(values[idx].ID) {
				break
			}

			if idx == (len(values) - 1) {
				inRegistry = true
			}
		}

		if inRegistry {
			sigEnum = tmpSigEnum
			break
		}
	}

	sigName := dbcValEnc.SignalName
	if !inRegistry {
		sigEnum = NewSignalEnum(fmt.Sprintf("%s_Enum", sigName))
		for _, dbcVal := range dbcValEnc.Values {
			if err := sigEnum.AddValue(NewSignalEnumValue(dbcVal.Name, int(dbcVal.ID))); err != nil {
				return i.errorf(dbcVal, err)
			}
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
		if nodeName == dbc.DummyNode {
			continue
		}

		tmpNode := NewNode(nodeName, NodeID(idx), 1)

		if desc, ok := i.nodeDesc[nodeName]; ok {
			tmpNode.SetDesc(desc)
		}

		tmpNodeInt := tmpNode.Interfaces()[0]

		if err := i.bus.AddNodeInterface(tmpNodeInt); err != nil {
			return i.errorf(dbcNodes, err)
		}

		i.nodeInts[tmpNode.name] = tmpNodeInt
	}

	if err := i.bus.AddNodeInterface(NewNode(dbc.DummyNode, 1024, 1).Interfaces()[0]); err != nil {
		return i.errorf(dbcNodes, err)
	}

	return nil
}

func (i *importer) importMessage(dbcMsg *dbc.Message) error {
	msgID := MessageID(dbcMsg.ID)
	msg := NewMessage(dbcMsg.Name, msgID, int(dbcMsg.Size))

	if err := msg.SetStaticCANID(CANID(msgID)); err != nil {
		return i.errorf(dbcMsg, err)
	}

	i.messages[msgID] = msg

	if desc, ok := i.msgDesc[msgID]; ok {
		msg.SetDesc(desc)
	}

	receivers := make(map[string]bool)
	muxSignals := []*dbc.Signal{}
	muxSigNames := make(map[string]int)

	slices.SortFunc(dbcMsg.Signals, func(a, b *dbc.Signal) int { return int(a.StartBit) - int(b.StartBit) })

	var currByteOrder dbc.SignalByteOrder
	for idx, dbcSig := range dbcMsg.Signals {
		if dbcSig.IsMultiplexor {
			muxSignals = append(muxSignals, dbcSig)
			muxSigNames[dbcSig.Name] = len(muxSignals) - 1
		}

		for _, rec := range dbcSig.Receivers {
			receivers[rec] = true
		}

		if idx == 0 {
			currByteOrder = dbcSig.ByteOrder
			continue
		}

		if dbcSig.ByteOrder != currByteOrder {
			return i.errorf(dbcSig, fmt.Errorf("byte_order: should be the same for all the signals within the message"))
		}
	}

	if currByteOrder == dbc.SignalBigEndian {
		msg.SetByteOrder(MessageByteOrderBigEndian)
	}

	for recName := range receivers {
		if recName == dbc.DummyNode {
			continue
		}

		recNode, err := i.bus.GetNodeInterfaceByNodeName(recName)
		if err != nil {
			return i.errorf(dbcMsg, err)
		}

		msg.AddReceiver(recNode)
	}

	sendNode, err := i.bus.GetNodeInterfaceByNodeName(dbcMsg.Transmitter)
	if err != nil {
		return i.errorf(dbcMsg, err)
	}

	if err := sendNode.AddSentMessage(msg); err != nil {
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

		lastMuxedStartPos := -1
		stdSignals := []*importerSignal{}

		muxedSignals := []*importerSignal{}
		for _, dbcSig := range dbcMsg.Signals {
			if dbcSig.Name == dbcMuxSig.Name {
				continue
			}

			tmpSig, err := i.importSignal(dbcSig, dbcMsg.ID)
			if err != nil {
				return err
			}

			tmpStartPos := i.getSignalStartBit(dbcSig)

			if dbcSig.IsMultiplexed {
				muxedSignals = append(muxedSignals, &importerSignal{sig: tmpSig, dbcSig: dbcSig, startPos: tmpStartPos})

				if tmpStartPos > lastMuxedStartPos {
					lastMuxedStartPos = tmpStartPos
				}

				continue
			}

			stdSignals = append(stdSignals, &importerSignal{sig: tmpSig, dbcSig: dbcSig, startPos: tmpStartPos})
		}

		muxorStartPos := i.getSignalStartBit(dbcMuxSig)

		for _, stdSig := range stdSignals {
			tmpStartPos := stdSig.startPos
			if tmpStartPos > muxorStartPos && tmpStartPos < lastMuxedStartPos {
				muxedSignals = append(muxedSignals, stdSig)
				continue
			}

			if err := msg.InsertSignal(stdSig.sig, tmpStartPos); err != nil {
				return i.errorf(stdSig.dbcSig, err)
			}
		}

		muxSig, err := i.importMuxSignal(dbcMuxSig, dbcMsg.ID, muxedSignals)
		if err != nil {
			return err
		}

		if err := msg.InsertSignal(muxSig, muxorStartPos); err != nil {
			return i.errorf(dbcMuxSig, err)
		}

		return nil
	}

	muxedSigGroups := make([][]*importerSignal, muxSigCount)
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
				return i.errorf(dbcSig, &ErrIsRequired{Item: "extended multiplexing"})
			}

			muxIdx, ok := muxSigNames[dbcExtMux.MultiplexorName]
			if !ok {
				return i.errorf(dbcExtMux, &NameError{Name: dbcExtMux.MultiplexorName, Err: ErrNotFound})
			}

			muxedSigGroups[muxIdx] = append(muxedSigGroups[muxIdx], &importerSignal{
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

		muxedSigGroups[muxIdx] = append(muxedSigGroups[muxIdx], &importerSignal{
			sig:    muxSig,
			dbcSig: dbcMuxSig,
		})
	}

	return nil
}

func (i *importer) getSignalStartBit(dbcSig *dbc.Signal) int {
	startBit := int(dbcSig.StartBit)
	if dbcSig.ByteOrder == dbc.SignalLittleEndian {
		return startBit
	}
	return startBit + 7 - 2*(startBit%8)
}

type importerSignal struct {
	sig      Signal
	dbcSig   *dbc.Signal
	startPos int
}

func (i *importer) importMuxSignal(dbcMuxSig *dbc.Signal, dbcMsgID uint32, muxedSignals []*importerSignal) (*MultiplexerSignal, error) {
	groupSize := 0
	muxedEndBit := 0
	for _, tmpMuxedSig := range muxedSignals {
		tmpEndBit := tmpMuxedSig.sig.GetSize() + i.getSignalStartBit(tmpMuxedSig.dbcSig)
		if tmpEndBit > muxedEndBit {
			muxedEndBit = tmpEndBit
		}
	}

	muxSigStartBit := i.getSignalStartBit(dbcMuxSig)
	muxSigSize := int(dbcMuxSig.Size)

	if muxedEndBit > 0 {
		groupSize = muxedEndBit - muxSigStartBit - muxSigSize
	}

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

		} else if tmpDBCSig.IsMultiplexed {
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

		dbcSigSize := int(dbcSig.Size)
		if enumSig.GetSize() < dbcSigSize {
			sigEnum.SetMinSize(dbcSigSize)
		}

		sig = enumSig

	} else {
		sigType, err := i.importSignalType(dbcSig)
		if err != nil {
			return nil, err
		}

		stdSig, err := NewStandardSignal(sigName, sigType)
		if err != nil {
			return nil, i.errorf(dbcSig, err)
		}

		symbol := dbcSig.Unit
		if symbol != "" {
			if sigUnit, ok := i.signalUnits[symbol]; ok {
				stdSig.SetUnit(sigUnit)
			} else {
				sigUnit := NewSignalUnit(symbol, SignalUnitKindCustom, symbol)
				stdSig.SetUnit(sigUnit)
				i.signalUnits[symbol] = sigUnit
			}
		}

		sig = stdSig
	}

	if desc, ok := i.sigDesc[sigKey]; ok {
		sig.SetDesc(desc)
	}

	i.signals[sigKey] = sig

	return sig, nil
}

func (i *importer) importSignalType(dbcSig *dbc.Signal) (*SignalType, error) {
	signed := false
	if dbcSig.ValueType == dbc.SignalSigned {
		signed = true
	}

	sigSize := int(dbcSig.Size)
	if sigSize == 1 && !signed {
		return i.flagSigType, nil
	}

	sigTypeKey := i.getSignalTypeKey(dbcSig)
	if sigType, ok := i.signalTypes[sigTypeKey]; ok {
		return sigType, nil
	}

	sigType := new(SignalType)
	if isDecimal(dbcSig.Factor) || isDecimal(dbcSig.Max) || isDecimal(dbcSig.Min) || isDecimal(dbcSig.Offset) {
		decSigType, err := NewDecimalSignalType(sigTypeKey, sigSize, signed)
		if err != nil {
			return nil, i.errorf(dbcSig, err)
		}
		sigType = decSigType
	} else {
		intSigType, err := NewIntegerSignalType(sigTypeKey, sigSize, signed)
		if err != nil {
			return nil, i.errorf(dbcSig, err)
		}
		sigType = intSigType
	}

	sigType.SetMin(dbcSig.Min)
	sigType.SetMax(dbcSig.Max)
	sigType.SetScale(dbcSig.Factor)
	sigType.SetOffset(dbcSig.Offset)

	i.signalTypes[sigTypeKey] = sigType

	return sigType, nil
}
