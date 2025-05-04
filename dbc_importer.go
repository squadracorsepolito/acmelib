package acmelib

import (
	"cmp"
	"errors"
	"fmt"
	"io"
	"slices"

	"github.com/squadracorsepolito/acmelib/dbc"
	"github.com/squadracorsepolito/acmelib/internal/collection"
)

// ImportDBCFile imports a DBC file passed as [io.Reader] and converts it
// to a [Bus]. The given filename will be used as the name of the bus.
func ImportDBCFile(filename string, r io.Reader) (*Bus, error) {
	dbcFile, err := dbc.Parse(filename, r, false)
	if err != nil {
		return nil, err
	}
	importer := newDBCImporter()
	return importer.importFile(dbcFile)
}

type dbcFileLocator interface {
	Location() *dbc.Location
}

type dbcImporter struct {
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

func newDBCImporter() *dbcImporter {
	return &dbcImporter{
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

func (i *dbcImporter) errorf(dbcLoc dbcFileLocator, err error) error {
	return fmt.Errorf("%s : %w", dbcLoc.Location(), err)
}

func (i *dbcImporter) getSignalKey(dbcMsgID uint32, sigName string) string {
	return fmt.Sprintf("%d_%s", dbcMsgID, sigName)
}

func (i *dbcImporter) getSignalTypeKey(dbcSig *dbc.Signal) string {
	signStr := "u"
	if dbcSig.ValueType == dbc.SignalSigned {
		signStr = "s"
	}
	return fmt.Sprintf("%s%d_%g-%g_%g_%g", signStr, dbcSig.Size, dbcSig.Min, dbcSig.Max, dbcSig.Factor, dbcSig.Offset)
}

func (i *dbcImporter) importFile(dbcFile *dbc.File) (*Bus, error) {
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

func (i *dbcImporter) importComments(dbcComments []*dbc.Comment) {
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

func (i *dbcImporter) importAttributes(dbcAtts []*dbc.Attribute, dbcAttDefs []*dbc.AttributeDefault, dbcAttVals []*dbc.AttributeValue) error {
	dbcAttDefMap := make(map[string]*dbc.AttributeDefault)
	for _, dbcAttDef := range dbcAttDefs {
		dbcAttDefMap[dbcAttDef.AttributeName] = dbcAttDef
	}

	attributes := make(map[string]Attribute)
	for _, dbcAtt := range dbcAtts {
		dbcAttDef, ok := dbcAttDefMap[dbcAtt.Name]
		if !ok {
			return i.errorf(dbcAtt, &IsRequiredError{Item: "attribute default"})
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

func (i *dbcImporter) importValueTable(dbcValTable *dbc.ValueTable) error {
	sigEnum := NewSignalEnum(dbcValTable.Name)

	for _, dbcVal := range dbcValTable.Values {
		if _, err := sigEnum.AddValue(int(dbcVal.ID), dbcVal.Name); err != nil {
			return i.errorf(dbcVal, err)
		}
	}

	i.signalEnumRegistry = append(i.signalEnumRegistry, sigEnum)

	return nil
}

func (i *dbcImporter) importValueEncoding(dbcValEnc *dbc.ValueEncoding) error {
	if dbcValEnc.Kind != dbc.ValueEncodingSignal {
		return nil
	}

	values := dbcValEnc.Values
	slices.SortFunc(values, func(a, b *dbc.ValueDescription) int {
		return cmp.Compare(a.ID, b.ID)
	})

	var sigEnum *SignalEnum
	inRegistry := false
	for _, tmpSigEnum := range i.signalEnumRegistry {
		if len(values) != len(tmpSigEnum.values) {
			continue
		}

		for idx, tmpVal := range tmpSigEnum.values {
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
		sigEnum = NewSignalEnum(fmt.Sprintf("%s_enum", sigName))
		for _, dbcVal := range dbcValEnc.Values {
			if _, err := sigEnum.AddValue(int(dbcVal.ID), dbcVal.Name); err != nil {
				return i.errorf(dbcVal, err)
			}
		}

	}

	i.signalEnums[i.getSignalKey(dbcValEnc.MessageID, sigName)] = sigEnum

	return nil
}

func (i *dbcImporter) importExtMuxes(dbcExtMuxes []*dbc.ExtendedMux) {
	for _, tmpExtMux := range dbcExtMuxes {
		key := i.getSignalKey(tmpExtMux.MessageID, tmpExtMux.MultiplexedName)
		i.dbcExtMuxes[key] = tmpExtMux
	}
}

func (i *dbcImporter) importNodes(dbcNodes *dbc.Nodes) error {
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

func (i *dbcImporter) importMessage(dbcMsg *dbc.Message) error {
	msgID := MessageID(dbcMsg.ID)
	msg := NewMessage(dbcMsg.Name, msgID, int(dbcMsg.Size))

	// Set the message ID as static
	if err := msg.SetStaticCANID(CANID(msgID)); err != nil {
		return i.errorf(dbcMsg, err)
	}

	i.messages[msgID] = msg

	if desc, ok := i.msgDesc[msgID]; ok {
		msg.SetDesc(desc)
	}

	// Sort the signals by start bit
	slices.SortFunc(dbcMsg.Signals, func(a, b *dbc.Signal) int {
		return cmp.Compare(a.StartBit, b.StartBit)
	})

	muxorOnly := []*dbc.Signal{}
	muxedOnly := []*dbc.Signal{}
	muxorMuxed := collection.NewQueue[*dbc.Signal]()

	receivers := make(map[string]struct{})

	for _, dbcSig := range dbcMsg.Signals {
		// Add the receivers
		for _, rec := range dbcSig.Receivers {
			receivers[rec] = struct{}{}
		}

		// Check if the signal is standard or enum
		if !dbcSig.IsMultiplexed && !dbcSig.IsMultiplexor {
			// Import and insert the signal into the message
			sig, err := i.importSignal(dbcSig, dbcMsg.ID)
			if err != nil {
				return err
			}

			if err := msg.InsertSignal(sig, i.getSignalStartPos(dbcSig)); err != nil {
				return i.errorf(dbcSig, err)
			}

			continue
		}

		// Check if the signal is a muxor and it is not multiplexed
		if dbcSig.IsMultiplexor && !dbcSig.IsMultiplexed {
			muxorOnly = append(muxorOnly, dbcSig)
			continue
		}

		// Check if the signal is multiplexed and it is not a muxor
		if dbcSig.IsMultiplexed && !dbcSig.IsMultiplexor {
			muxedOnly = append(muxedOnly, dbcSig)
			continue
		}

		// The signal is both muxor and multiplexed
		muxorMuxed.Push(dbcSig)
	}

	muxLayers := make(map[string]*MultiplexedLayer)

	// Add a new multiplexed layer for each muxor signal to the message layout
	for _, dbcMuxorSig := range muxorOnly {
		muxLayer, err := i.importMuxorSignal(msg.layout, dbcMuxorSig, dbcMsg.ID)
		if err != nil {
			return err
		}

		muxLayers[dbcMuxorSig.Name] = muxLayer
	}

	// Add muxor signals that are also multiplexed (nested multiplexing)
	prevSigName := ""
	for muxorMuxed.Size() > 0 {
		dbcSig := muxorMuxed.Pop()

		// Check if the current signal is missing multiplexing information
		muxorName := dbcSig.Name
		if muxorName == prevSigName {
			return i.errorf(dbcSig, newIsRequiredError("extended multiplexing"))
		}
		prevSigName = muxorName

		// Get the multiplexing information
		dbcExtMux, ok := i.dbcExtMuxes[i.getSignalKey(dbcMsg.ID, muxorName)]
		if !ok {
			return i.errorf(dbcSig, newIsRequiredError("extended multiplexing"))
		}

		// Check if the multiplexing information is valid
		if len(dbcExtMux.Ranges) != 1 || dbcExtMux.Ranges[0].From != dbcExtMux.Ranges[0].To {
			return i.errorf(dbcExtMux, errors.New("muxor signal must have only one range"))
		}

		// Check if the current muxor can be inserted into a multiplexed layer.
		// If it cannot, reinsert it into the queue
		parentMuxorName := dbcExtMux.MultiplexorName
		parentMuxLayer, ok := muxLayers[parentMuxorName]
		if !ok {
			muxorMuxed.Push(dbcSig)
			continue
		}

		// Get the corresponding layout from the multiplexed layer
		layoutID := int(dbcExtMux.Ranges[0].From)
		muxedLayout := parentMuxLayer.GetLayout(layoutID)
		if muxedLayout == nil {
			return i.errorf(dbcExtMux, ErrOutOfBounds)
		}

		muxLayer, err := i.importMuxorSignal(muxedLayout, dbcSig, dbcMsg.ID)
		if err != nil {
			return err
		}

		muxLayers[muxorName] = muxLayer
	}

	// Add the multiplexed signals
	for _, dbcMuxedSig := range muxedOnly {
		// Check if the signal is missing a muxor
		if len(muxorOnly) == 0 {
			return i.errorf(dbcMuxedSig, errors.New("missing a muxor signal"))
		}

		sigKey := i.getSignalKey(dbcMsg.ID, dbcMuxedSig.Name)
		dbcExtMux, hasExtMux := i.dbcExtMuxes[sigKey]

		// By default, set the current multiplexed layer to the first one
		muxLayer := muxLayers[muxorOnly[0].Name]

		if len(muxorOnly) > 1 {
			// In this case the extended multiplexing information is required
			if !hasExtMux {
				return i.errorf(dbcMuxedSig, newIsRequiredError("extended multiplexing"))
			}

			// Set the current multiplexed layer to the one corresponding to
			// the extended multiplexing information
			tmpMuxLayer, ok := muxLayers[dbcExtMux.MultiplexorName]
			if !ok {
				return i.errorf(dbcMuxedSig, ErrNotFound)
			}
			muxLayer = tmpMuxLayer
		}

		// Get the layout IDs
		layoutIDs := []int{}
		if hasExtMux {
			for _, r := range dbcExtMux.Ranges {
				for i := r.From; i <= r.To; i++ {
					layoutIDs = append(layoutIDs, int(i))
				}
			}
		} else {
			// This case can only happen when there is only one muxor
			layoutIDs = append(layoutIDs, int(dbcMuxedSig.MuxSwitchValue))
		}

		sig, err := i.importSignal(dbcMuxedSig, dbcMsg.ID)
		if err != nil {
			return err
		}

		if err := muxLayer.InsertSignal(sig, i.getSignalStartPos(dbcMuxedSig), layoutIDs...); err != nil {
			return i.errorf(dbcMuxedSig, err)
		}
	}

	// Add the receivers
	for recName := range receivers {
		if recName == dbc.DummyNode {
			continue
		}

		recNode, err := i.bus.GetNodeInterfaceByNodeName(recName)
		if err != nil {
			return i.errorf(dbcMsg, err)
		}

		if err := msg.AddReceiver(recNode); err != nil {
			return i.errorf(dbcMsg, err)
		}
	}

	// Add the message to the sender node
	sendNode, err := i.bus.GetNodeInterfaceByNodeName(dbcMsg.Transmitter)
	if err != nil {
		return i.errorf(dbcMsg, err)
	}

	if err := sendNode.AddSentMessage(msg); err != nil {
		return i.errorf(dbcMsg, err)
	}

	return nil
}

func (i *dbcImporter) getSignalStartPos(dbcSig *dbc.Signal) int {
	startBit := int(dbcSig.StartBit)
	if dbcSig.ByteOrder == dbc.SignalLittleEndian {
		return startBit
	}
	return StartPosFromBigEndian(startBit)
}

// finishSignal updates the description and endianness of the given signal if needed.
func (i *dbcImporter) finishSignal(sig Signal, dbcSig *dbc.Signal, sigKey string) {
	if desc, ok := i.sigDesc[sigKey]; ok {
		sig.SetDesc(desc)
	}

	if dbcSig.ByteOrder == dbc.SignalBigEndian {
		sig.SetEndianness(EndiannessBigEndian)
	}

	i.signals[sigKey] = sig
}

// importMuxorSignal imports a muxor signal and returns the corresponding multiplexed layer.
func (i *dbcImporter) importMuxorSignal(layout *SignalLayout, dbcSig *dbc.Signal, dbcMsgID uint32) (*MultiplexedLayer, error) {
	muxorName := dbcSig.Name

	layoutCount := getValueFromSize(int(dbcSig.Size))
	muxLayer, err := layout.AddMultiplexedLayer(muxorName, i.getSignalStartPos(dbcSig), layoutCount)
	if err != nil {
		return nil, i.errorf(dbcSig, err)
	}

	i.finishSignal(muxLayer.muxor, dbcSig, i.getSignalKey(dbcMsgID, muxorName))

	return muxLayer, nil
}

func (i *dbcImporter) importSignal(dbcSig *dbc.Signal, dbcMsgID uint32) (Signal, error) {
	var sig Signal

	sigName := dbcSig.Name
	sigKey := i.getSignalKey(dbcMsgID, sigName)

	// Check if the signal has a value table (enum signal)
	if sigEnum, ok := i.signalEnums[sigKey]; ok {
		enumSig, err := NewEnumSignal(sigName, sigEnum)
		if err != nil {
			return nil, i.errorf(dbcSig, err)
		}

		// If the size of the signal is bigger than the size of the enum,
		// set the size of the enum as fixed with the given size
		dbcSigSize := int(dbcSig.Size)
		if enumSig.Size() < dbcSigSize {
			sigEnum.SetFixedSize(true)
			if err := sigEnum.UpdateSize(dbcSigSize); err != nil {
				return nil, i.errorf(dbcSig, err)
			}
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

	i.finishSignal(sig, dbcSig, sigKey)

	return sig, nil
}

func (i *dbcImporter) importSignalType(dbcSig *dbc.Signal) (*SignalType, error) {
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
