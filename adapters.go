package acmelib

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/FerroO2000/acmelib/dbc"
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
	msg.SetID(MessageID(dbcMsg.ID))
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
			int(dbcSig.Size), signed, SignalTypeOrderBigEndian, dbcSig.Min, dbcSig.Max)
	}

	return NewIntegerSignalType(fmt.Sprintf("int_type_%s", dbcSig.Name), int(dbcSig.Size), signed)
}

func WriteDBC(network *Network, basePath string) error {
	for _, bus := range network.Buses() {
		adapter := newAdapter()
		writer := dbc.NewWriter()

		f, err := os.Create(fmt.Sprintf("%s/%s/%s.%s", basePath, network.name, bus.name, dbc.FileExtension))
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

type adapter struct {
	dbcFile *dbc.File

	currDBCMsg *dbc.Message
}

func newAdapter() *adapter {
	return &adapter{
		dbcFile: new(dbc.File),
	}
}

func (a *adapter) addDBCComment(comment *dbc.Comment) {
	a.dbcFile.Comments = append(a.dbcFile.Comments, comment)
}

func (a *adapter) adaptBus(bus *Bus) *dbc.File {
	if bus.desc != "" {
		a.addDBCComment(&dbc.Comment{
			Kind: dbc.CommentGeneral,
			Text: bus.desc,
		})
	}

	a.dbcFile.NewSymbols = &dbc.NewSymbols{
		Symbols: dbc.GetNewSymbols(),
	}

	a.dbcFile.BitTiming = &dbc.BitTiming{
		Baudrate: uint32(bus.baudrate),
	}

	a.adaptNodes(bus.Nodes())

	return a.dbcFile
}

func (a *adapter) adaptNodes(nodes []*Node) {
	dbcNodes := new(dbc.Nodes)

	for _, node := range nodes {
		if node.desc != "" {
			a.addDBCComment(&dbc.Comment{
				Kind:     dbc.CommentNode,
				Text:     node.desc,
				NodeName: node.name,
			})
		}

		dbcNodes.Names = append(dbcNodes.Names, node.name)

		for _, msg := range node.Messages() {
			a.adaptMessage(msg)
		}
	}

	a.dbcFile.Nodes = dbcNodes
}

func (a *adapter) adaptMessage(msg *Message) {
	dbcMsg := new(dbc.Message)

	if msg.desc != "" {
		a.addDBCComment(&dbc.Comment{
			Kind:      dbc.CommentMessage,
			Text:      msg.desc,
			MessageID: uint32(msg.ID()),
		})
	}

	dbcMsg.ID = uint32(msg.ID())
	dbcMsg.Name = msg.name
	dbcMsg.Size = uint32(msg.sizeByte)
	dbcMsg.Transmitter = msg.senderNode.name

	a.currDBCMsg = dbcMsg

	receiverNames := []string{}
	for _, rec := range msg.Receivers() {
		receiverNames = append(receiverNames, rec.name)
	}
	for _, sig := range msg.Signals() {
		a.adaptSignal(sig, receiverNames...)
	}

	a.dbcFile.Messages = append(a.dbcFile.Messages, dbcMsg)
}

func (a *adapter) adaptSignal(sig Signal, receiverNames ...string) {
	parMsg, err := sig.Parent().ToParentMessage()
	if err != nil {
		panic(err)
	}
	msgID := parMsg.ID()

	if sig.Desc() != "" {
		a.addDBCComment(&dbc.Comment{
			Kind:       dbc.CommentSignal,
			Text:       sig.Desc(),
			MessageID:  uint32(msgID),
			SignalName: sig.Name(),
		})
	}

	dbcSig := new(dbc.Signal)

	dbcSig.Name = sig.Name()
	dbcSig.Size = uint32(sig.GetSize())
	dbcSig.StartBit = uint32(sig.GetStartBit())

	if len(receiverNames) == 0 {
		dbcSig.Receivers = []string{dbc.DummyNode}
	} else {
		dbcSig.Receivers = receiverNames
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

		dbcValEnc := new(dbc.ValueEncoding)
		dbcValEnc.Kind = dbc.ValueEncodingSignal
		dbcValEnc.MessageID = uint32(msgID)
		dbcValEnc.SignalName = sig.Name()

		for _, val := range enumSig.enum.Values() {
			dbcValEnc.Values = append(dbcValEnc.Values, &dbc.ValueDescription{
				ID:   uint32(val.index),
				Name: val.name,
			})
		}

		a.dbcFile.ValueEncodings = append(a.dbcFile.ValueEncodings, dbcValEnc)
		a.currDBCMsg.Signals = append(a.currDBCMsg.Signals, dbcSig)

	case SignalKindMultiplexer:
		// muxSig, err := sig.ToMultiplexer()
		// if err != nil {
		// 	panic(err)
		// }
		// a.adaptMultiplexerSignal(muxSig)
	}

}

func (a *adapter) adaptStandardSignal(stdSig *StandardSignal, dbcSig *dbc.Signal) {
	switch stdSig.typ.order {
	case SignalTypeOrderLittleEndian:
		dbcSig.ByteOrder = dbc.SignalLittleEndian
	case SignalTypeOrderBigEndian:
		dbcSig.ByteOrder = dbc.SignalBigEndian
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

func (a *adapter) adaptEnumSignal(enumSig *EnumSignal, dbcSig *dbc.Signal) {
	dbcSig.ByteOrder = dbc.SignalLittleEndian
	dbcSig.ValueType = dbc.SignalUnsigned

	dbcSig.Min = 0
	dbcSig.Max = float64(enumSig.enum.maxIndex)
	dbcSig.Offset = 0
	dbcSig.Factor = 1
}

// func (a *adapter) adaptMultiplexerSignal(muxSig *MultiplexerSignal1) {
// 	dbcMuxorSig := new(dbc.Signal)

// 	dbcMuxorSig.Name = muxSig.Name()

// 	dbcMuxorSig.IsMultiplexor = true

// 	dbcMuxorSig.Size = uint32(muxSig.GetSize())
// 	dbcMuxorSig.StartBit = uint32(muxSig.GetStartBit())

// 	dbcMuxorSig.ByteOrder = dbc.SignalLittleEndian
// 	dbcMuxorSig.ValueType = dbc.SignalUnsigned

// 	selectValues := 1 << muxSig.SelectSize()

// 	dbcMuxorSig.Factor = 1
// 	dbcMuxorSig.Offset = 0
// 	dbcMuxorSig.Min = 0
// 	dbcMuxorSig.Max = float64(selectValues)

// 	a.currDBCMsg.Signals = append(a.currDBCMsg.Signals, dbcMuxorSig)

// 	for i := 0; i < selectValues; i++ {
// 		for _, muxedSig := range muxSig.GetSelectedMuxSignals(i) {
// 			a.adaptSignal(muxedSig)
// 		}
// 	}
// }
