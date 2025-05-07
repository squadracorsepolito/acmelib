package acmelib

import (
	"io"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/squadracorsepolito/acmelib/dbc"
	"github.com/squadracorsepolito/acmelib/internal/collection"
)

// ExportDBCNetwork exports the given [Network] to DBC.
// It will create a directory with the given base path and network name.
// Into the directory, it will create a DBC file for each [Bus] of the network.
func ExportDBCNetwork(network *Network, basePath string) error {
	netName := clearSpaces(network.name)
	dirPath := netName
	if len(basePath) > 0 {
		dirPath = filepath.Join(basePath, netName)
	}

	err := os.MkdirAll(dirPath, 0666)
	if err != nil {
		return err
	}

	buses := network.Buses()
	wg := &sync.WaitGroup{}
	wg.Add(len(buses))

	for _, bus := range buses {
		f, err := os.Create(filepath.Join(dirPath, clearSpaces(bus.name)+dbc.FileExtension))
		if err != nil {
			return err
		}
		defer f.Close()

		go exportBusAsync(f, bus, wg)
	}

	wg.Wait()

	return nil
}

func exportBusAsync(w io.Writer, bus *Bus, wg *sync.WaitGroup) {
	defer wg.Done()
	ExportDBCBus(w, bus)
}

// ExportDBCBus exports the given [Bus] to DBC.
// It writes the content of the result DBC file into the [io.Writer].
func ExportDBCBus(w io.Writer, bus *Bus) {
	exp := newDBCExporter()
	dbcFile := exp.exportBus(bus)
	dbc.Write(w, dbcFile, false)
}

type dbcExporter struct {
	dbcFile *dbc.File

	attNames     map[string]bool
	nodeAttNames map[string]bool
	msgAttNames  map[string]bool
	sigAttNames  map[string]bool

	sigEnums map[EntityID]*SignalEnum
}

func newDBCExporter() *dbcExporter {
	return &dbcExporter{
		dbcFile: new(dbc.File),

		attNames:     make(map[string]bool),
		nodeAttNames: make(map[string]bool),
		msgAttNames:  make(map[string]bool),
		sigAttNames:  make(map[string]bool),

		sigEnums: make(map[EntityID]*SignalEnum),
	}
}

func (e *dbcExporter) addDBCComment(comment *dbc.Comment) {
	e.dbcFile.Comments = append(e.dbcFile.Comments, comment)
}

func (e *dbcExporter) exportAttributeAssignment(attAss *AttributeAssignment, dbcAttKind dbc.AttributeKind, dbcAttVal *dbc.AttributeValue) {
	att := attAss.attribute
	attName := clearSpaces(att.Name())

	dbcAttVal.AttributeName = attName

	hasAtt := false
	switch dbcAttKind {
	case dbc.AttributeGeneral:
		if _, ok := e.attNames[attName]; ok {
			hasAtt = true
		} else {
			e.attNames[attName] = true
		}
		dbcAttVal.AttributeKind = dbc.AttributeGeneral

	case dbc.AttributeNode:
		if _, ok := e.nodeAttNames[attName]; ok {
			hasAtt = true
		} else {
			e.nodeAttNames[attName] = true
		}
		dbcAttVal.AttributeKind = dbc.AttributeNode

	case dbc.AttributeMessage:
		if _, ok := e.msgAttNames[attName]; ok {
			hasAtt = true
		} else {
			e.msgAttNames[attName] = true
		}
		dbcAttVal.AttributeKind = dbc.AttributeMessage

	case dbc.AttributeSignal:
		if _, ok := e.sigAttNames[attName]; ok {
			hasAtt = true
		} else {
			e.sigAttNames[attName] = true
		}
		dbcAttVal.AttributeKind = dbc.AttributeSignal
	}

	if !hasAtt {
		dbcAtt := new(dbc.Attribute)
		dbcAtt.Kind = dbcAttKind
		e.exportAttribute(att, dbcAtt)
	}

	switch att.Type() {
	case AttributeTypeString:
		dbcAttVal.Type = dbc.AttributeValueString
		dbcAttVal.ValueString = attAss.value.(string)

	case AttributeTypeInteger:
		intAtt, err := att.ToInteger()
		if err != nil {
			panic(err)
		}
		if intAtt.isHexFormat {
			dbcAttVal.Type = dbc.AttributeValueHex
			dbcAttVal.ValueHex = uint32(attAss.value.(int))
		} else {
			dbcAttVal.Type = dbc.AttributeValueInt
			dbcAttVal.ValueInt = attAss.value.(int)
		}

	case AttributeTypeFloat:
		dbcAttVal.Type = dbc.AttributeValueFloat
		dbcAttVal.ValueFloat = attAss.value.(float64)

	case AttributeTypeEnum:
		enumAtt, err := att.ToEnum()
		if err != nil {
			panic(err)
		}
		dbcAttVal.Type = dbc.AttributeValueInt

		valIdx := 0
		strVal := attAss.value.(string)
		for idx, val := range enumAtt.Values() {
			if strVal == val {
				valIdx = idx
				break
			}
		}
		dbcAttVal.ValueInt = valIdx
	}
}

func (e *dbcExporter) exportAttribute(att Attribute, dbcAtt *dbc.Attribute) {
	attName := clearSpaces(att.Name())
	dbcAtt.Name = attName

	dbcAttDef := new(dbc.AttributeDefault)
	dbcAttDef.AttributeName = attName

	switch att.Type() {
	case AttributeTypeString:
		strAtt, err := att.ToString()
		if err != nil {
			panic(err)
		}
		dbcAtt.Type = dbc.AttributeString

		dbcAttDef.Type = dbc.AttributeDefaultString
		dbcAttDef.ValueString = strAtt.defValue

	case AttributeTypeInteger:
		intAtt, err := att.ToInteger()
		if err != nil {
			panic(err)
		}

		if intAtt.isHexFormat {
			dbcAtt.Type = dbc.AttributeHex
			dbcAtt.MinHex = uint32(intAtt.min)
			dbcAtt.MaxHex = uint32(intAtt.max)

			dbcAttDef.Type = dbc.AttributeDefaultHex
			dbcAttDef.ValueHex = uint32(intAtt.defValue)
		} else {
			dbcAtt.Type = dbc.AttributeInt
			dbcAtt.MinInt = intAtt.min
			dbcAtt.MaxInt = intAtt.max

			dbcAttDef.Type = dbc.AttributeDefaultInt
			dbcAttDef.ValueInt = intAtt.defValue
		}

	case AttributeTypeFloat:
		floatAtt, err := att.ToFloat()
		if err != nil {
			panic(err)
		}
		dbcAtt.Type = dbc.AttributeFloat
		dbcAtt.MinFloat = floatAtt.min
		dbcAtt.MaxFloat = floatAtt.max

		dbcAttDef.Type = dbc.AttributeDefaultString
		dbcAttDef.ValueFloat = floatAtt.defValue

	case AttributeTypeEnum:
		enumAtt, err := att.ToEnum()
		if err != nil {
			panic(err)
		}
		dbcAtt.Type = dbc.AttributeEnum
		dbcAtt.EnumValues = enumAtt.Values()

		dbcAttDef.Type = dbc.AttributeDefaultString
		dbcAttDef.ValueString = enumAtt.defValue
	}

	e.dbcFile.Attributes = append(e.dbcFile.Attributes, dbcAtt)
	e.dbcFile.AttributeDefaults = append(e.dbcFile.AttributeDefaults, dbcAttDef)
}

func (e *dbcExporter) exportBus(bus *Bus) *dbc.File {
	if bus.desc != "" {
		e.addDBCComment(&dbc.Comment{
			Kind: dbc.CommentGeneral,
			Text: bus.desc,
		})
	}

	for _, attVal := range bus.AttributeAssignments() {
		dbcAttVal := new(dbc.AttributeValue)
		e.exportAttributeAssignment(attVal, dbc.AttributeGeneral, dbcAttVal)
		e.dbcFile.AttributeValues = append(e.dbcFile.AttributeValues, dbcAttVal)
	}

	e.exportNodeInterfaces(bus.NodeInterfaces())

	for _, sigEnum := range e.sigEnums {
		e.exportSignalEnum(sigEnum)
	}

	return e.dbcFile
}

func (e *dbcExporter) exportNodeInterfaces(nodeInts []*NodeInterface) {
	dbcNodes := new(dbc.Nodes)

	for _, nodeInt := range nodeInts {
		nodeName := clearSpaces(nodeInt.node.name)

		if nodeInt.node.desc != "" {
			e.addDBCComment(&dbc.Comment{
				Kind:     dbc.CommentNode,
				Text:     nodeInt.node.desc,
				NodeName: nodeName,
			})
		}

		for _, attVal := range nodeInt.node.AttributeAssignments() {
			dbcAttVal := new(dbc.AttributeValue)
			dbcAttVal.NodeName = nodeName
			e.exportAttributeAssignment(attVal, dbc.AttributeNode, dbcAttVal)
			e.dbcFile.AttributeValues = append(e.dbcFile.AttributeValues, dbcAttVal)
		}

		dbcNodes.Names = append(dbcNodes.Names, nodeName)

		for _, msg := range nodeInt.SentMessages() {
			e.exportMessage(msg)
		}
	}

	e.dbcFile.Nodes = dbcNodes
}

func (e *dbcExporter) exportMessage(msg *Message) {
	dbcMsg := new(dbc.Message)

	msgID := uint32(msg.GetCANID())
	dbcMsg.ID = msgID

	// Handle message description
	if msg.desc != "" {
		e.addDBCComment(&dbc.Comment{
			Kind:      dbc.CommentMessage,
			Text:      msg.desc,
			MessageID: msgID,
		})
	}

	// Handle message attributes
	attAssignments := msg.AttributeAssignments()
	if msg.cycleTime != 0 {
		attAssignments = append(attAssignments, newAttributeAssignment(msgCycleTimeAtt, msg, msg.cycleTime))
	}
	if msg.delayTime != 0 {
		attAssignments = append(attAssignments, newAttributeAssignment(msgDelayTimeAtt, msg, msg.delayTime))
	}
	if msg.startDelayTime != 0 {
		attAssignments = append(attAssignments, newAttributeAssignment(msgStartDelayTimeAtt, msg, msg.startDelayTime))
	}
	if msg.sendType != MessageSendTypeUnset {
		attAssignments = append(attAssignments, newAttributeAssignment(msgSendTypeAtt, msg, messageSendTypeToDBC(msg.sendType)))
	}
	for _, attAss := range attAssignments {
		dbcAttVal := new(dbc.AttributeValue)
		dbcAttVal.MessageID = dbcMsg.ID
		e.exportAttributeAssignment(attAss, dbc.AttributeMessage, dbcAttVal)
		e.dbcFile.AttributeValues = append(e.dbcFile.AttributeValues, dbcAttVal)
	}

	dbcMsg.Name = clearSpaces(msg.name)
	dbcMsg.Size = uint32(msg.sizeByte)
	dbcMsg.Transmitter = msg.senderNodeInt.node.name

	// Get the receivers
	receivers := msg.Receivers()
	dbcReceivers := make([]string, 0, len(receivers))
	if len(receivers) == 0 {
		dbcReceivers = append(dbcReceivers, dbc.DummyNode)
	}
	for _, rec := range receivers {
		dbcReceivers = append(dbcReceivers, clearSpaces(rec.node.name))
	}

	// Handle the message signals
	e.exportSignalLayout(msg.layout, dbcMsg, dbcReceivers, e.getExtMuxNeeded(msg.layout))

	e.dbcFile.Messages = append(e.dbcFile.Messages, dbcMsg)
}

func (e *dbcExporter) getSignalStartBit(sig Signal) (uint32, dbc.SignalByteOrder) {
	startPos := sig.StartPos()

	if sig.Endianness() == EndiannessLittleEndian {
		return uint32(startPos), dbc.SignalLittleEndian
	}

	return uint32(StartPosFromBigEndian(startPos)), dbc.SignalBigEndian
}

// exportSignal fills in the given dbc.Signal with common information
// and adds it to the given dbc.Message.
func (e *dbcExporter) exportSignal(sig Signal, dbcSig *dbc.Signal, dbcMsg *dbc.Message) {
	// Handle signal description
	if sig.Desc() != "" {
		e.addDBCComment(&dbc.Comment{
			Kind:       dbc.CommentSignal,
			Text:       sig.Desc(),
			MessageID:  dbcMsg.ID,
			SignalName: dbcSig.Name,
		})
	}

	// Handle signal attributes
	attAssignments := sig.AttributeAssignments()
	if sig.StartValue() != 0 {
		attAssignments = append(attAssignments, newAttributeAssignment(sigStartValueAtt, sig, sig.StartValue()))
	}
	if sig.SendType() != SignalSendTypeUnset {
		attAssignments = append(attAssignments, newAttributeAssignment(sigSendTypeAtt, sig, signalSendTypeToDBC(sig.SendType())))
	}
	for _, attAss := range attAssignments {
		dbcAttVal := new(dbc.AttributeValue)
		dbcAttVal.MessageID = dbcMsg.ID
		dbcAttVal.SignalName = dbcSig.Name
		e.exportAttributeAssignment(attAss, dbc.AttributeSignal, dbcAttVal)
		e.dbcFile.AttributeValues = append(e.dbcFile.AttributeValues, dbcAttVal)
	}

	startBit, byteOrder := e.getSignalStartBit(sig)
	dbcSig.StartBit = startBit
	dbcSig.ByteOrder = byteOrder

	dbcSig.Size = uint32(sig.Size())

	switch sig.Kind() {
	case SignalKindStandard:
		stdSig, err := sig.ToStandard()
		if err != nil {
			panic(err)
		}
		e.exportStandardSignal(stdSig, dbcSig)

	case SignalKindEnum:
		enumSig, err := sig.ToEnum()
		if err != nil {
			panic(err)
		}
		e.exportEnumSignal(enumSig, dbcMsg.ID, dbcSig)

	case SignalKindMuxor:
		muxorSig, err := sig.ToMuxor()
		if err != nil {
			panic(err)
		}
		e.exportMuxorSignal(muxorSig, dbcSig)
	}

	dbcMsg.Signals = append(dbcMsg.Signals, dbcSig)
}

func (e *dbcExporter) exportStandardSignal(stdSig *StandardSignal, dbcSig *dbc.Signal) {
	if stdSig.typ.signed {
		dbcSig.ValueType = dbc.SignalSigned
	} else {
		dbcSig.ValueType = dbc.SignalUnsigned
	}

	dbcSig.Min = stdSig.typ.min
	dbcSig.Max = stdSig.typ.max
	dbcSig.Offset = stdSig.typ.offset
	dbcSig.Factor = stdSig.typ.scale

	unit := stdSig.unit
	if unit != nil {
		dbcSig.Unit = unit.symbol
	}
}

func (e *dbcExporter) exportEnumSignal(enumSig *EnumSignal, dbcMsgID uint32, dbcSig *dbc.Signal) {
	dbcSig.ValueType = dbc.SignalUnsigned

	enum := enumSig.enum

	dbcSig.Min = 0
	dbcSig.Max = float64(enum.maxIndex)
	dbcSig.Offset = 0
	dbcSig.Factor = 1

	dbcValEnc := new(dbc.ValueEncoding)
	dbcValEnc.Kind = dbc.ValueEncodingSignal
	dbcValEnc.MessageID = dbcMsgID
	dbcValEnc.SignalName = clearSpaces(enumSig.Name())

	values := enum.values

	// Set min value if first value is not 0
	if len(values) > 0 && values[0].index != 0 {
		dbcSig.Min = float64(values[0].index)
	}

	dbcValEnc.Values = e.getDBCValueDescription(values)

	e.dbcFile.ValueEncodings = append(e.dbcFile.ValueEncodings, dbcValEnc)
	e.sigEnums[enum.entityID] = enumSig.enum
}

func (e *dbcExporter) getDBCValueDescription(enumValues []*SignalEnumValue) []*dbc.ValueDescription {
	dbcValDesc := make([]*dbc.ValueDescription, len(enumValues))
	for idx, val := range enumValues {
		dbcValDesc[idx] = &dbc.ValueDescription{
			ID:   uint32(val.index),
			Name: val.name,
		}
	}
	return dbcValDesc
}

func (e *dbcExporter) exportSignalEnum(enum *SignalEnum) {
	e.dbcFile.ValueTables = append(e.dbcFile.ValueTables, &dbc.ValueTable{
		Name:   enum.name,
		Values: e.getDBCValueDescription(enum.Values()),
	})
}

func (e *dbcExporter) exportMuxorSignal(muxorSig *MuxorSignal, dbcSig *dbc.Signal) {
	dbcSig.IsMultiplexor = true
	dbcSig.Min = 0
	dbcSig.Max = float64(muxorSig.layoutCount - 1)
	dbcSig.Factor = 1
	dbcSig.Offset = 0
}

// exportSignalLayout exports a signal layout.
// It creates a dbc.Signal for each signal in the layout,
// but it does not add it to the given dbc.Message since that is done by the exportSignal method.
func (e *dbcExporter) exportSignalLayout(layout *SignalLayout, dbcMsg *dbc.Message, dbcReceivers []string, extMuxNeeded bool) {
	for sig := range layout.ibst.InOrder() {
		sigName := clearSpaces(sig.Name())

		// Check if the signal has already been processed
		if slices.ContainsFunc(dbcMsg.Signals, func(s *dbc.Signal) bool { return s.Name == sigName }) {
			continue
		}

		dbcSig := new(dbc.Signal)
		dbcSig.Name = sigName
		dbcSig.Receivers = dbcReceivers
		e.exportSignal(sig, dbcSig, dbcMsg)

		if !layout.fromMultiplexedLayer() {
			continue
		}

		// The signal is multiplexed
		dbcSig.IsMultiplexed = true
		parentMuxorName := layout.parentMuxLayer.muxor.name

		// Check if the signal is a muxor. If so, extended multiplexing is needed
		if sig.Kind() == SignalKindMuxor {
			layoutID := layout.id
			dbcSig.MuxSwitchValue = uint32(layoutID)
			e.exportExtendedMux(sigName, parentMuxorName, []int{layoutID}, dbcMsg.ID)
			continue
		}

		// The signal is not a muxor
		layoutIDs, ok := layout.parentMuxLayer.singalLayoutIDs.Get(sig.EntityID())
		if !ok || len(layoutIDs) == 0 {
			continue
		}

		dbcSig.MuxSwitchValue = uint32(layoutIDs[0])

		// Add extended multiplexing if needed
		if len(layoutIDs) > 1 || extMuxNeeded {
			e.exportExtendedMux(sigName, parentMuxorName, layoutIDs, dbcMsg.ID)
		}
	}

	for muxLayer := range layout.muxLayers.Values() {
		for _, muxLayout := range muxLayer.iterLayouts() {
			e.exportSignalLayout(muxLayout, dbcMsg, dbcReceivers, extMuxNeeded)
		}
	}
}

// getExtMuxNeeded states whether extended multiplexing is needed.
// It is needed if a message has more than one muxor signal.
func (e *dbcExporter) getExtMuxNeeded(layout *SignalLayout) bool {
	count := 0

	s := collection.NewStack[*SignalLayout]()
	s.Push(layout)

	for !s.IsEmpty() {
		currLayout := s.Pop()
		count += currLayout.muxLayers.Size()

		if count > 1 {
			return true
		}

		for muxLayer := range currLayout.muxLayers.Values() {
			for _, muxLayout := range muxLayer.iterLayouts() {
				s.Push(muxLayout)
			}
		}
	}

	return count > 1
}

// exportExtendedMux creates a dbc.ExtendedMux.
func (e *dbcExporter) exportExtendedMux(sigName, muxorSigName string, layoutIDs []int, dbcMsgID uint32) {
	dbcExtMux := new(dbc.ExtendedMux)
	dbcExtMux.MessageID = dbcMsgID
	dbcExtMux.MultiplexedName = clearSpaces(sigName)
	dbcExtMux.MultiplexorName = clearSpaces(muxorSigName)

	// Aggregate ranges
	from := layoutIDs[0]
	next := from
	for i := range len(layoutIDs) - 1 {
		curr := layoutIDs[i]
		next = layoutIDs[i+1]

		if next == curr+1 {
			continue
		}

		dbcExtMux.Ranges = append(dbcExtMux.Ranges, &dbc.ExtendedMuxRange{
			From: uint32(from),
			To:   uint32(curr),
		})

		from = next
	}
	dbcExtMux.Ranges = append(dbcExtMux.Ranges, &dbc.ExtendedMuxRange{
		From: uint32(from),
		To:   uint32(layoutIDs[len(layoutIDs)-1]),
	})

	e.dbcFile.ExtendedMuxes = append(e.dbcFile.ExtendedMuxes, dbcExtMux)
}
