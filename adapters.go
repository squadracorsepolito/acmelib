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
