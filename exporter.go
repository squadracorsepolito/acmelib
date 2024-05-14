package acmelib

import (
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/squadracorsepolito/acmelib/dbc"
)

// ExportNetwork exports the given [Network] to DBC.
// It will create a directory with the given base path and network name.
// Into the directory, it will create a DBC file for each [Bus] of the network.
func ExportNetwork(network *Network, basePath string) error {
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
	ExportBus(w, bus)
}

// ExportBus exports the given [Bus] to DBC.
// It writes the content of the result DBC file into the [io.Writer].
func ExportBus(w io.Writer, bus *Bus) {
	exp := newExporter()
	dbcFile := exp.exportBus(bus)
	dbc.Write(w, dbcFile, false)
}

type exporter struct {
	dbcFile *dbc.File

	currDBCMsg *dbc.Message

	attNames     map[string]bool
	nodeAttNames map[string]bool
	msgAttNames  map[string]bool
	sigAttNames  map[string]bool
}

func newExporter() *exporter {
	return &exporter{
		dbcFile: new(dbc.File),

		currDBCMsg: nil,

		attNames:     make(map[string]bool),
		nodeAttNames: make(map[string]bool),
		msgAttNames:  make(map[string]bool),
		sigAttNames:  make(map[string]bool),
	}
}

func (e *exporter) addDBCComment(comment *dbc.Comment) {
	e.dbcFile.Comments = append(e.dbcFile.Comments, comment)
}

func (e *exporter) exportAttributeValue(attVal *AttributeValue, dbcAttKind dbc.AttributeKind, dbcAttVal *dbc.AttributeValue) {
	att := attVal.attribute
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

	switch att.Kind() {
	case AttributeKindString:
		dbcAttVal.Type = dbc.AttributeValueString
		dbcAttVal.ValueString = attVal.value.(string)

	case AttributeKindInteger:
		intAtt, err := att.ToInteger()
		if err != nil {
			panic(err)
		}
		if intAtt.isHexFormat {
			dbcAttVal.Type = dbc.AttributeValueHex
			dbcAttVal.ValueHex = uint32(attVal.value.(int))
		} else {

			dbcAttVal.Type = dbc.AttributeValueInt
			dbcAttVal.ValueInt = attVal.value.(int)
		}

	case AttributeKindFloat:
		dbcAttVal.Type = dbc.AttributeValueFloat
		dbcAttVal.ValueFloat = attVal.value.(float64)

	case AttributeKindEnum:
		enumAtt, err := att.ToEnum()
		if err != nil {
			panic(err)
		}
		dbcAttVal.Type = dbc.AttributeValueInt

		valIdx := 0
		strVal := attVal.value.(string)
		for idx, val := range enumAtt.Values() {
			if strVal == val {
				valIdx = idx
				break
			}
		}
		dbcAttVal.ValueInt = valIdx
	}
}

func (e *exporter) exportAttribute(att Attribute, dbcAtt *dbc.Attribute) {
	attName := clearSpaces(att.Name())
	dbcAtt.Name = attName

	dbcAttDef := new(dbc.AttributeDefault)
	dbcAttDef.AttributeName = attName

	switch att.Kind() {
	case AttributeKindString:
		strAtt, err := att.ToString()
		if err != nil {
			panic(err)
		}
		dbcAtt.Type = dbc.AttributeString

		dbcAttDef.Type = dbc.AttributeDefaultString
		dbcAttDef.ValueString = strAtt.defValue

	case AttributeKindInteger:
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

	case AttributeKindFloat:
		floatAtt, err := att.ToFloat()
		if err != nil {
			panic(err)
		}
		dbcAtt.Type = dbc.AttributeFloat
		dbcAtt.MinFloat = floatAtt.min
		dbcAtt.MaxFloat = floatAtt.max

		dbcAttDef.Type = dbc.AttributeDefaultString
		dbcAttDef.ValueFloat = floatAtt.defValue

	case AttributeKindEnum:
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

func (e *exporter) exportBus(bus *Bus) *dbc.File {
	if bus.desc != "" {
		e.addDBCComment(&dbc.Comment{
			Kind: dbc.CommentGeneral,
			Text: bus.desc,
		})
	}

	for _, attVal := range bus.AttributeValues() {
		dbcAttVal := new(dbc.AttributeValue)
		e.exportAttributeValue(attVal, dbc.AttributeGeneral, dbcAttVal)
		e.dbcFile.AttributeValues = append(e.dbcFile.AttributeValues, dbcAttVal)
	}

	e.dbcFile.BitTiming = &dbc.BitTiming{
		Baudrate: uint32(bus.baudrate),
	}

	e.exportNodes(bus.Nodes())

	return e.dbcFile
}

func (e *exporter) exportNodes(nodes []*Node) {
	dbcNodes := new(dbc.Nodes)

	for _, node := range nodes {
		nodeName := clearSpaces(node.name)

		if node.desc != "" {
			e.addDBCComment(&dbc.Comment{
				Kind:     dbc.CommentNode,
				Text:     node.desc,
				NodeName: nodeName,
			})
		}

		for _, attVal := range node.AttributeValues() {
			dbcAttVal := new(dbc.AttributeValue)
			dbcAttVal.NodeName = nodeName
			e.exportAttributeValue(attVal, dbc.AttributeNode, dbcAttVal)
			e.dbcFile.AttributeValues = append(e.dbcFile.AttributeValues, dbcAttVal)
		}

		dbcNodes.Names = append(dbcNodes.Names, nodeName)

		for _, msg := range node.Messages() {
			e.exportMessage(msg)
		}
	}

	e.dbcFile.Nodes = dbcNodes
}

func (e *exporter) exportMessage(msg *Message) {
	dbcMsg := new(dbc.Message)

	if msg.desc != "" {
		e.addDBCComment(&dbc.Comment{
			Kind:      dbc.CommentMessage,
			Text:      msg.desc,
			MessageID: uint32(msg.CANID()),
		})
	}

	dbcMsg.ID = uint32(msg.CANID())

	attValues := msg.AttributeValues()
	if msg.cycleTime != 0 {
		attValues = append(attValues, newAttributeValue(msgCycleTimeAtt, msg.cycleTime))
	}
	if msg.delayTime != 0 {
		attValues = append(attValues, newAttributeValue(msgDelayTimeAtt, msg.delayTime))
	}
	if msg.startDelayTime != 0 {
		attValues = append(attValues, newAttributeValue(msgStartDelayTimeAtt, msg.startDelayTime))
	}
	if msg.sendType != MessageSendTypeUnset {
		attValues = append(attValues, newAttributeValue(msgSendTypeAtt, string(msg.sendType)))
	}
	for _, attVal := range attValues {
		dbcAttVal := new(dbc.AttributeValue)
		dbcAttVal.MessageID = dbcMsg.ID
		e.exportAttributeValue(attVal, dbc.AttributeMessage, dbcAttVal)
		e.dbcFile.AttributeValues = append(e.dbcFile.AttributeValues, dbcAttVal)
	}

	dbcMsg.Name = clearSpaces(msg.name)
	dbcMsg.Size = uint32(msg.sizeByte)
	dbcMsg.Transmitter = msg.senderNode.name

	e.currDBCMsg = dbcMsg

	for _, sig := range msg.Signals() {
		e.exportSignal(sig)
	}

	e.dbcFile.Messages = append(e.dbcFile.Messages, dbcMsg)
}

func (e *exporter) exportSignal(sig Signal) {
	parMsg := sig.ParentMessage()
	msgID := parMsg.CANID()
	sigName := clearSpaces(sig.Name())

	if sig.Desc() != "" {
		e.addDBCComment(&dbc.Comment{
			Kind:       dbc.CommentSignal,
			Text:       sig.Desc(),
			MessageID:  uint32(msgID),
			SignalName: sigName,
		})
	}

	attValues := sig.AttributeValues()
	if sig.SendType() != SignalSendTypeUnset {
		attValues = append(attValues, newAttributeValue(sigSendTypeAtt, string(sig.SendType())))
	}
	for _, attVal := range attValues {
		dbcAttVal := new(dbc.AttributeValue)
		dbcAttVal.MessageID = uint32(msgID)
		dbcAttVal.SignalName = sigName
		e.exportAttributeValue(attVal, dbc.AttributeSignal, dbcAttVal)
		e.dbcFile.AttributeValues = append(e.dbcFile.AttributeValues, dbcAttVal)
	}

	dbcSig := new(dbc.Signal)
	dbcSig.Name = sigName

	if len(parMsg.Receivers()) == 0 {
		dbcSig.Receivers = []string{dbc.DummyNode}
	} else {
		receiverNames := []string{}
		for _, rec := range parMsg.Receivers() {
			receiverNames = append(receiverNames, clearSpaces(rec.name))
		}
		dbcSig.Receivers = receiverNames
	}

	if sig.ParentMultiplexerSignal() != nil {
		dbcSig.IsMultiplexed = true
	}

	switch sig.Kind() {
	case SignalKindStandard:
		stdSig, err := sig.ToStandard()
		if err != nil {
			panic(err)
		}
		e.exportStandardSignal(stdSig, dbcSig)
		e.currDBCMsg.Signals = append(e.currDBCMsg.Signals, dbcSig)

	case SignalKindEnum:
		enumSig, err := sig.ToEnum()
		if err != nil {
			panic(err)
		}
		e.exportEnumSignal(enumSig, dbcSig)
		e.currDBCMsg.Signals = append(e.currDBCMsg.Signals, dbcSig)

	case SignalKindMultiplexer:
		muxSig, err := sig.ToMultiplexer()
		if err != nil {
			panic(err)
		}
		e.exportMultiplexerSignal(muxSig, dbcSig)
	}
}

func (e *exporter) getStartBit(startBit int, byteOrder MessageByteOrder) (uint32, dbc.SignalByteOrder) {
	if byteOrder == MessageByteOrderLittleEndian {
		return uint32(startBit), dbc.SignalLittleEndian
	}

	tmpStartBit := startBit + 7 - 2*(startBit%8)
	return uint32(tmpStartBit), dbc.SignalBigEndian
}

func (e *exporter) exportStandardSignal(stdSig *StandardSignal, dbcSig *dbc.Signal) {
	dbcSig.Size = uint32(stdSig.GetSize())

	startBit, byteOrder := e.getStartBit(stdSig.GetStartBit(), stdSig.parentMsg.byteOrder)
	dbcSig.StartBit = startBit
	dbcSig.ByteOrder = byteOrder

	if stdSig.typ.signed {
		dbcSig.ValueType = dbc.SignalSigned
	} else {
		dbcSig.ValueType = dbc.SignalUnsigned
	}

	dbcSig.Min = stdSig.min
	dbcSig.Max = stdSig.max
	dbcSig.Offset = stdSig.offset
	dbcSig.Factor = stdSig.scale

	unit := stdSig.unit
	if unit != nil {
		dbcSig.Unit = unit.symbol
	}
}

func (e *exporter) exportEnumSignal(enumSig *EnumSignal, dbcSig *dbc.Signal) {
	dbcSig.Size = uint32(enumSig.GetSize())

	startBit, byteOrder := e.getStartBit(enumSig.GetStartBit(), enumSig.parentMsg.byteOrder)
	dbcSig.StartBit = startBit
	dbcSig.ByteOrder = byteOrder

	dbcSig.ValueType = dbc.SignalUnsigned

	dbcSig.Min = 0
	dbcSig.Max = float64(enumSig.enum.maxIndex)
	dbcSig.Offset = 0
	dbcSig.Factor = 1

	dbcValEnc := new(dbc.ValueEncoding)
	dbcValEnc.Kind = dbc.ValueEncodingSignal
	dbcValEnc.MessageID = uint32(enumSig.parentMsg.id)
	dbcValEnc.SignalName = clearSpaces(enumSig.Name())

	for _, val := range enumSig.enum.Values() {
		dbcValEnc.Values = append(dbcValEnc.Values, &dbc.ValueDescription{
			ID:   uint32(val.index),
			Name: val.name,
		})
	}

	e.dbcFile.ValueEncodings = append(e.dbcFile.ValueEncodings, dbcValEnc)
}

func (e *exporter) exportMultiplexerSignal(muxSig *MultiplexerSignal, dbcSig *dbc.Signal) {
	dbcSig.Size = uint32(muxSig.GetGroupCountSize())

	startBit, byteOrder := e.getStartBit(muxSig.GetStartBit(), muxSig.parentMsg.byteOrder)
	dbcSig.StartBit = startBit
	dbcSig.ByteOrder = byteOrder

	dbcSig.IsMultiplexor = true

	dbcSig.ValueType = dbc.SignalUnsigned

	dbcSig.Min = 0
	dbcSig.Max = float64(muxSig.groupCount - 1)
	dbcSig.Offset = 0
	dbcSig.Factor = 1

	e.currDBCMsg.Signals = append(e.currDBCMsg.Signals, dbcSig)

	isExtended := false
	nestedMux := dbcSig.IsMultiplexed

	sigNames := []string{}
	sigGroupIDs := make(map[string][]int)
	for id, group := range muxSig.GetSignalGroups() {
		for _, tmpSig := range group {
			tmpSigName := clearSpaces(tmpSig.Name())
			groupIDs, ok := sigGroupIDs[tmpSigName]
			if !ok {
				sigGroupIDs[tmpSigName] = []int{id}
			}

			if tmpSig.Kind() == SignalKindMultiplexer {
				nestedMux = true
			}

			if len(groupIDs) == 0 {
				sigNames = append(sigNames, tmpSigName)
				e.exportSignal(tmpSig)
				e.currDBCMsg.Signals[len(e.currDBCMsg.Signals)-1].MuxSwitchValue = uint32(id)
				continue
			}

			sigGroupIDs[tmpSigName] = append(sigGroupIDs[tmpSigName], id)
			isExtended = true
		}
	}

	if !isExtended && !nestedMux {
		return
	}

	for _, tmpSigName := range sigNames {
		groupIDs := sigGroupIDs[tmpSigName]

		if !nestedMux && len(groupIDs) == 1 {
			continue
		}

		dbcExtMux := new(dbc.ExtendedMux)
		dbcExtMux.MessageID = uint32(muxSig.parentMsg.id)
		dbcExtMux.MultiplexorName = clearSpaces(muxSig.name)
		dbcExtMux.MultiplexedName = tmpSigName

		from := groupIDs[0]
		next := from
		for i := 0; i < len(groupIDs)-1; i++ {
			curr := groupIDs[i]
			next = groupIDs[i+1]

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
			To:   uint32(next),
		})

		e.dbcFile.ExtendedMuxes = append(e.dbcFile.ExtendedMuxes, dbcExtMux)
	}
}
