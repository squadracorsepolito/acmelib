package acmelib

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/FerroO2000/acmelib/dbc"
)

var (
	msgCycleTimeAtt, _      = NewIntegerAttribute("GenMsgCycleTime", 0, 0, 1000)
	msgDelayTimeAtt, _      = NewIntegerAttribute("GenMsgDelayTime", 0, 0, 1000)
	msgStartDelayTimeAtt, _ = NewIntegerAttribute("GenMsgStartDelayTime", 0, 0, 1000)
)

var msgSendTypeAtt, _ = NewEnumAttribute("GenMsgSendType",
	string(MessageSendTypeUnset), string(MessageSendTypeCyclic), string(MessageSendTypeCyclicIfActive),
	string(MessageSendTypeCyclicIfTriggered), string(MessageSendTypeCyclicIfActiveAndTriggered),
)

var sigSendTypeAtt, _ = NewEnumAttribute("GenSigSendType", string(SignalSendTypeUnset), string(SignalSendTypeCyclic),
	string(SignalSendTypeOnWrite), string(SignalSendTypeOnWriteWithRepetition), string(SignalSendTypeOnChange),
	string(SignalSendTypeOnChangeWithRepetition), string(SignalSendTypeIfActive), string(SignalSendTypeIfActiveWithRepetition),
)

func LoadDBC(networkName string, dbcFilenames ...string) (*Network, error) {
	net := NewNetwork(networkName)

	for _, filename := range dbcFilenames {
		extension := filepath.Ext(filename)
		if extension != dbc.FileExtension {
			return nil, fmt.Errorf(`file "%s" must have "%s" extension`, filename, dbc.FileExtension)
		}

		file, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		parser := dbc.NewParser(filepath.Base(filename), file)
		dbcFile, err := parser.Parse()
		if err != nil {
			return nil, err
		}

		adapter := newDBCAdapter()
		bus, err := adapter.adaptFile(dbcFile)
		if err != nil {
			return nil, err
		}

		if err := net.AddBus(bus); err != nil {
			return nil, err
		}
	}

	return net, nil
}

type dbcAdapter struct {
	busComment   string
	nodeComments map[string]string
	msgComments  map[uint32]string
	sigComments  map[string]string
}

func newDBCAdapter() *dbcAdapter {
	return &dbcAdapter{
		busComment:   "",
		nodeComments: make(map[string]string),
		msgComments:  make(map[uint32]string),
		sigComments:  make(map[string]string),
	}
}

func (a *dbcAdapter) adaptFile(dbcFile *dbc.File) (*Bus, error) {
	bus := NewBus(dbcFile.Location.Filename)

	for _, comment := range dbcFile.Comments {
		a.mapComment(comment)
	}
	if a.busComment != "" {
		bus.SetDesc(a.busComment)
	}

	for _, node := range a.adaptNodes(dbcFile.Nodes) {
		if err := bus.AddNode(node); err != nil {
			return nil, err
		}
	}

	dummyNode := NewNode(dbc.DummyNode, 1024)
	dummyNodeUsed := false
	for _, dbcMsg := range dbcFile.Messages {
		msg, err := a.adaptMessage(dbcMsg)
		if err != nil {
			return nil, err
		}

		if dbcMsg.Transmitter == dbc.DummyNode {
			if err := dummyNode.AddMessage(msg); err != nil {
				return nil, err
			}
			dummyNodeUsed = true
			continue
		}

		senderNode, err := bus.GetNodeByName(dbcMsg.Transmitter)
		if err != nil {
			return nil, err
		}
		if err := senderNode.AddMessage(msg); err != nil {
			return nil, err
		}

		receiverMap := make(map[string]bool)
		for _, dbcSig := range dbcMsg.Signals {
			for _, rec := range dbcSig.Receivers {
				receiverMap[rec] = true
			}
		}
		for rec := range receiverMap {
			if rec == dbc.DummyNode {
				msg.AddReceiver(dummyNode)
				continue
			}
			recNode, err := bus.GetNodeByName(rec)
			if err != nil {
				return nil, err
			}
			msg.AddReceiver(recNode)
		}
	}

	if dummyNodeUsed {
		if err := bus.AddNode(dummyNode); err != nil {
			return nil, err
		}
	}

	return bus, nil
}

func (a *dbcAdapter) mapComment(dbcComment *dbc.Comment) {
	switch dbcComment.Kind {
	case dbc.CommentGeneral:
		a.busComment = dbcComment.Text

	case dbc.CommentNode:
		a.nodeComments[dbcComment.NodeName] = dbcComment.Text

	case dbc.CommentMessage:
		a.msgComments[dbcComment.MessageID] = dbcComment.Text

	case dbc.CommentSignal:
		key := fmt.Sprintf("%d_%s", dbcComment.MessageID, dbcComment.SignalName)
		a.sigComments[key] = dbcComment.Text
	}
}

func (a *dbcAdapter) adaptNodes(dbcNodes *dbc.Nodes) []*Node {
	nodes := make([]*Node, len(dbcNodes.Names))
	for i := 0; i < len(dbcNodes.Names); i++ {
		nodeName := dbcNodes.Names[i]
		tmpNode := NewNode(nodeName, NodeID(i))
		if comment, ok := a.nodeComments[nodeName]; ok {
			tmpNode.SetDesc(comment)
		}
		nodes[i] = tmpNode
	}
	return nodes
}

func (a *dbcAdapter) adaptMessage(dbcMsg *dbc.Message) (*Message, error) {
	msg := NewMessage(dbcMsg.Name, int(dbcMsg.Size))
	msg.SetCANID(MessageCANID(dbcMsg.ID))
	if comment, ok := a.msgComments[dbcMsg.ID]; ok {
		msg.SetDesc(comment)
	}

	for _, dbcSig := range dbcMsg.Signals {
		sig, err := a.adaptSignal(dbcMsg.ID, dbcSig)
		if err != nil {
			return nil, err
		}

		if err := msg.InsertSignal(sig, int(dbcSig.StartBit)); err != nil {
			return nil, err
		}
	}

	return msg, nil
}

func (a *dbcAdapter) adaptSignal(msgID uint32, dbcSig *dbc.Signal) (Signal, error) {
	sigTyp, err := a.adaptSignalType(dbcSig)
	if err != nil {
		return nil, err
	}

	sig, err := NewStandardSignal(dbcSig.Name, sigTyp)
	if err != nil {
		return nil, err
	}

	if comment, ok := a.sigComments[fmt.Sprintf("%d_%s", msgID, dbcSig.Name)]; ok {
		sig.SetDesc(comment)
	}

	if err := sig.SetPhysicalValues(dbcSig.Min, dbcSig.Max, dbcSig.Offset, dbcSig.Factor); err != nil {
		return nil, err
	}

	if dbcSig.Unit != "" {
		sigUnit := NewSignalUnit(fmt.Sprintf("unit_%s", dbcSig.Name), SignalUnitKindCustom, dbcSig.Unit)
		sig.SetUnit(sigUnit)
	}

	return sig, nil
}

func (a *dbcAdapter) adaptSignalType(dbcSig *dbc.Signal) (*SignalType, error) {
	if dbcSig.Size == 1 {
		return NewFlagSignalType(fmt.Sprintf("flag_type_%s", dbcSig.Name)), nil
	}

	var signed bool
	switch dbcSig.ValueType {
	case dbc.SignalSigned:
		signed = true
	case dbc.SignalUnsigned:
		signed = false
	}

	switch dbcSig.ByteOrder {
	case dbc.SignalBigEndian:
		return NewCustomSignalType(fmt.Sprintf("custom_type_%s", dbcSig.Name),
			int(dbcSig.Size), signed, dbcSig.Min, dbcSig.Max)
	}

	return NewIntegerSignalType(fmt.Sprintf("int_type_%s", dbcSig.Name), int(dbcSig.Size), signed)
}

func WriteDBC(network *Network, basePath string) error {
	for _, bus := range network.Buses() {
		adapter := newAdapter()
		writer := dbc.NewWriter()

		filename := ""
		if len(basePath) > 0 {
			filename = fmt.Sprintf("%s/%s/%s%s", basePath, network.name, bus.name, dbc.FileExtension)
		} else {
			filename = fmt.Sprintf("%s/%s%s", network.name, bus.name, dbc.FileExtension)
		}

		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = f.WriteString(writer.Write(adapter.adaptBus(bus)))
		if err != nil {
			return err
		}
	}

	return nil
}

type netToDBCAdapter struct {
	dbcFile *dbc.File

	currDBCMsg *dbc.Message

	attNames     map[string]bool
	nodeAttNames map[string]bool
	msgAttNames  map[string]bool
	sigAttNames  map[string]bool
}

func newAdapter() *netToDBCAdapter {
	return &netToDBCAdapter{
		dbcFile: new(dbc.File),

		currDBCMsg: nil,

		attNames:     make(map[string]bool),
		nodeAttNames: make(map[string]bool),
		msgAttNames:  make(map[string]bool),
		sigAttNames:  make(map[string]bool),
	}
}

func (a *netToDBCAdapter) addDBCComment(comment *dbc.Comment) {
	a.dbcFile.Comments = append(a.dbcFile.Comments, comment)
}

func (a *netToDBCAdapter) adaptAttributeValue(attVal *AttributeValue, dbcAttKind dbc.AttributeKind, dbcAttVal *dbc.AttributeValue) {
	att := attVal.attribute
	attName := att.Name()

	dbcAttVal.AttributeName = attName

	hasAtt := false
	switch dbcAttKind {
	case dbc.AttributeGeneral:
		if _, ok := a.attNames[attName]; ok {
			hasAtt = true
		} else {
			a.attNames[attName] = true
		}
		dbcAttVal.AttributeKind = dbc.AttributeGeneral

	case dbc.AttributeNode:
		if _, ok := a.nodeAttNames[attName]; ok {
			hasAtt = true
		} else {
			a.nodeAttNames[attName] = true
		}
		dbcAttVal.AttributeKind = dbc.AttributeNode

	case dbc.AttributeMessage:
		if _, ok := a.msgAttNames[attName]; ok {
			hasAtt = true
		} else {
			a.msgAttNames[attName] = true
		}
		dbcAttVal.AttributeKind = dbc.AttributeMessage

	case dbc.AttributeSignal:
		if _, ok := a.sigAttNames[attName]; ok {
			hasAtt = true
		} else {
			a.sigAttNames[attName] = true
		}
		dbcAttVal.AttributeKind = dbc.AttributeSignal
	}

	if !hasAtt {
		dbcAtt := new(dbc.Attribute)
		dbcAtt.Kind = dbcAttKind
		a.adaptAttribute(att, dbcAtt)
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
			dbcAttVal.ValueHex = attVal.value.(int)
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

func (a *netToDBCAdapter) adaptAttribute(att Attribute, dbcAtt *dbc.Attribute) {
	attName := att.Name()
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
			dbcAtt.MinHex = intAtt.min
			dbcAtt.MaxHex = intAtt.max

			dbcAttDef.Type = dbc.AttributeDefaultHex
			dbcAttDef.ValueHex = intAtt.defValue
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

	a.dbcFile.Attributes = append(a.dbcFile.Attributes, dbcAtt)
	a.dbcFile.AttributeDefaults = append(a.dbcFile.AttributeDefaults, dbcAttDef)
}

func (a *netToDBCAdapter) adaptBus(bus *Bus) *dbc.File {
	if bus.desc != "" {
		a.addDBCComment(&dbc.Comment{
			Kind: dbc.CommentGeneral,
			Text: bus.desc,
		})
	}

	for _, attVal := range bus.AttributeValues() {
		dbcAttVal := new(dbc.AttributeValue)
		a.adaptAttributeValue(attVal, dbc.AttributeGeneral, dbcAttVal)
		a.dbcFile.AttributeValues = append(a.dbcFile.AttributeValues, dbcAttVal)
	}

	a.dbcFile.BitTiming = &dbc.BitTiming{
		Baudrate: uint32(bus.baudrate),
	}

	a.adaptNodes(bus.Nodes())

	return a.dbcFile
}

func (a *netToDBCAdapter) adaptNodes(nodes []*Node) {
	dbcNodes := new(dbc.Nodes)

	for _, node := range nodes {
		if node.desc != "" {
			a.addDBCComment(&dbc.Comment{
				Kind:     dbc.CommentNode,
				Text:     node.desc,
				NodeName: node.name,
			})
		}

		for _, attVal := range node.AttributeValues() {
			dbcAttVal := new(dbc.AttributeValue)
			dbcAttVal.NodeName = node.name
			a.adaptAttributeValue(attVal, dbc.AttributeNode, dbcAttVal)
			a.dbcFile.AttributeValues = append(a.dbcFile.AttributeValues, dbcAttVal)
		}

		dbcNodes.Names = append(dbcNodes.Names, node.name)

		for _, msg := range node.Messages() {
			a.adaptMessage(msg)
		}
	}

	a.dbcFile.Nodes = dbcNodes
}

func (a *netToDBCAdapter) adaptMessage(msg *Message) {
	dbcMsg := new(dbc.Message)

	if msg.desc != "" {
		a.addDBCComment(&dbc.Comment{
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
		a.adaptAttributeValue(attVal, dbc.AttributeMessage, dbcAttVal)
		a.dbcFile.AttributeValues = append(a.dbcFile.AttributeValues, dbcAttVal)
	}

	dbcMsg.Name = msg.name
	dbcMsg.Size = uint32(msg.sizeByte)
	dbcMsg.Transmitter = msg.senderNode.name

	a.currDBCMsg = dbcMsg

	for _, sig := range msg.Signals() {
		a.adaptSignal(sig)
	}

	a.dbcFile.Messages = append(a.dbcFile.Messages, dbcMsg)
}

func (a *netToDBCAdapter) adaptSignal(sig Signal) {
	parMsg := sig.ParentMessage()
	msgID := parMsg.CANID()

	if sig.Desc() != "" {
		a.addDBCComment(&dbc.Comment{
			Kind:       dbc.CommentSignal,
			Text:       sig.Desc(),
			MessageID:  uint32(msgID),
			SignalName: sig.Name(),
		})
	}

	attValues := sig.AttributeValues()
	if sig.SendType() != SignalSendTypeUnset {
		attValues = append(attValues, newAttributeValue(sigSendTypeAtt, string(sig.SendType())))
	}
	for _, attVal := range attValues {
		dbcAttVal := new(dbc.AttributeValue)
		dbcAttVal.MessageID = uint32(msgID)
		dbcAttVal.SignalName = sig.Name()
		a.adaptAttributeValue(attVal, dbc.AttributeSignal, dbcAttVal)
		a.dbcFile.AttributeValues = append(a.dbcFile.AttributeValues, dbcAttVal)
	}

	dbcSig := new(dbc.Signal)
	dbcSig.Name = sig.Name()

	if len(parMsg.Receivers()) == 0 {
		dbcSig.Receivers = []string{dbc.DummyNode}
	} else {
		receiverNames := []string{}
		for _, rec := range parMsg.Receivers() {
			receiverNames = append(receiverNames, rec.name)
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
		a.adaptStandardSignal(stdSig, dbcSig)
		a.currDBCMsg.Signals = append(a.currDBCMsg.Signals, dbcSig)

	case SignalKindEnum:
		enumSig, err := sig.ToEnum()
		if err != nil {
			panic(err)
		}
		a.adaptEnumSignal(enumSig, dbcSig)
		a.currDBCMsg.Signals = append(a.currDBCMsg.Signals, dbcSig)

	case SignalKindMultiplexer:
		muxSig, err := sig.ToMultiplexer()
		if err != nil {
			panic(err)
		}
		a.adaptMultiplexerSignal(muxSig, dbcSig)
	}
}

func (a *netToDBCAdapter) adaptStandardSignal(stdSig *StandardSignal, dbcSig *dbc.Signal) {
	dbcSig.Size = uint32(stdSig.GetSize())

	switch stdSig.byteOrder {
	case SignalByteOrderLittleEndian:
		dbcSig.ByteOrder = dbc.SignalLittleEndian
		dbcSig.StartBit = uint32(stdSig.GetStartBit())

	case SignalByteOrderBigEndian:
		dbcSig.ByteOrder = dbc.SignalBigEndian
		dbcSig.StartBit = uint32(stdSig.GetStartBit()) + dbcSig.Size - 1
	}

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

func (a *netToDBCAdapter) adaptEnumSignal(enumSig *EnumSignal, dbcSig *dbc.Signal) {
	dbcSig.Size = uint32(enumSig.GetSize())

	switch enumSig.byteOrder {
	case SignalByteOrderLittleEndian:
		dbcSig.ByteOrder = dbc.SignalLittleEndian
		dbcSig.StartBit = uint32(enumSig.GetStartBit())

	case SignalByteOrderBigEndian:
		dbcSig.ByteOrder = dbc.SignalBigEndian
		dbcSig.StartBit = uint32(enumSig.GetStartBit()) + dbcSig.Size - 1
	}

	dbcSig.ValueType = dbc.SignalUnsigned

	dbcSig.Min = 0
	dbcSig.Max = float64(enumSig.enum.maxIndex)
	dbcSig.Offset = 0
	dbcSig.Factor = 1

	dbcValEnc := new(dbc.ValueEncoding)
	dbcValEnc.Kind = dbc.ValueEncodingSignal
	dbcValEnc.MessageID = uint32(enumSig.parentMsg.id)
	dbcValEnc.SignalName = enumSig.Name()

	for _, val := range enumSig.enum.Values() {
		dbcValEnc.Values = append(dbcValEnc.Values, &dbc.ValueDescription{
			ID:   uint32(val.index),
			Name: val.name,
		})
	}

	a.dbcFile.ValueEncodings = append(a.dbcFile.ValueEncodings, dbcValEnc)
}

func (a *netToDBCAdapter) adaptMultiplexerSignal(muxSig *MultiplexerSignal, dbcSig *dbc.Signal) {
	dbcSig.Size = uint32(muxSig.GetGroupCountSize())

	switch muxSig.byteOrder {
	case SignalByteOrderLittleEndian:
		dbcSig.ByteOrder = dbc.SignalLittleEndian
		dbcSig.StartBit = uint32(muxSig.GetStartBit())

	case SignalByteOrderBigEndian:
		dbcSig.ByteOrder = dbc.SignalBigEndian
		dbcSig.StartBit = uint32(muxSig.GetStartBit()) + dbcSig.Size - 1
	}

	dbcSig.IsMultiplexor = true

	dbcSig.ValueType = dbc.SignalUnsigned

	dbcSig.Min = 0
	dbcSig.Max = float64(muxSig.groupCount - 1)
	dbcSig.Offset = 0
	dbcSig.Factor = 1

	a.currDBCMsg.Signals = append(a.currDBCMsg.Signals, dbcSig)

	isExtended := false
	nestedMux := dbcSig.IsMultiplexed

	sigNames := []string{}
	sigGroupIDs := make(map[string][]int)
	for id, group := range muxSig.GetSignalGroups() {
		for _, tmpSig := range group {
			tmpSigName := tmpSig.Name()
			groupIDs, ok := sigGroupIDs[tmpSigName]
			if !ok {
				sigGroupIDs[tmpSigName] = []int{id}
			}

			if tmpSig.Kind() == SignalKindMultiplexer {
				nestedMux = true
			}

			if len(groupIDs) == 0 {
				sigNames = append(sigNames, tmpSigName)
				a.adaptSignal(tmpSig)
				a.currDBCMsg.Signals[len(a.currDBCMsg.Signals)-1].MuxSwitchValue = uint32(id)
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
		dbcExtMux.MultiplexorName = muxSig.name
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

		a.dbcFile.ExtendedMuxes = append(a.dbcFile.ExtendedMuxes, dbcExtMux)
	}
}
