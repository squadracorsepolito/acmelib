package acmelib

import (
	"fmt"
	"os"

	"github.com/FerroO2000/acmelib/dbc"
)

func LoadDBC(networkName string, dbcFilenames ...string) (*Network, error) {
	net := NewNetwork(networkName)

	for _, filename := range dbcFilenames {
		file, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		parser := dbc.NewParser(filename, file)
		dbcFile, err := parser.Parse()
		if err != nil {
			return nil, err
		}

		if err := handleDBCFile(net, dbcFile); err != nil {
			return nil, err
		}
	}

	return net, nil
}

func handleDBCFile(net *Network, file *dbc.File) error {
	bus := NewBus("dbc bus")

	if err := net.AddBus(bus); err != nil {
		return err
	}

	for idx, nodeName := range file.Nodes.Names {
		tmpNode := NewNode(nodeName, NodeID(idx))
		if err := bus.AddNode(tmpNode); err != nil {
			return err
		}
	}

	dummyNode := NewNode("dummyNode", 1024)
	dummyNodeUsed := false

	for _, dbcMsg := range file.Messages {
		tmpMsg := NewMessage(dbcMsg.Name, int(dbcMsg.Size))
		tmpMsg.SetID(MessageID(dbcMsg.ID))

		if dbcMsg.Transmitter == dbc.DummyNode {
			dummyNodeUsed = true
			if err := dummyNode.AddMessage(tmpMsg); err != nil {
				return err
			}

		} else {
			node, err := bus.GetNodeByName(dbcMsg.Transmitter)
			if err != nil {
				return err
			}

			if err := node.AddMessage(tmpMsg); err != nil {
				return err
			}
		}

		receiverList := make(map[string]bool)
		for _, dbcSig := range dbcMsg.Signals {
			signed := false
			switch dbcSig.ValueType {
			case dbc.SignalSigned:
				signed = true
			case dbc.SignalUnsigned:
				signed = false
			}
			order := SignalTypeOrderLittleEndian
			switch dbcSig.ByteOrder {
			case dbc.SignalLittleEndian:
				order = SignalTypeOrderLittleEndian
			case dbc.SignalBigEndian:
				order = SignalTypeOrderBigEndian
			}

			tmpSigType, err := NewCustomSignalType(fmt.Sprintf("%s_type", dbcSig.Name), int(dbcSig.Size), signed, order, dbcSig.Min, dbcSig.Max)
			if err != nil {
				return err
			}

			tmpSig, err := NewStandardSignal(dbcSig.Name, tmpSigType)
			if err != nil {
				return err
			}

			if err := tmpSig.SetPhysicalValues(dbcSig.Min, dbcSig.Max, dbcSig.Offset, dbcSig.Factor); err != nil {
				return err
			}

			if dbcSig.Unit != "" {
				tmpSig.SetUnit(NewSignalUnit(fmt.Sprintf("%s_unit", dbcSig.Name), SignalUnitKindCustom, dbcSig.Unit))
			}

			if err := tmpMsg.InsertSignal(tmpSig, int(dbcSig.StartBit)); err != nil {
				return err
			}

			for _, rec := range dbcSig.Receivers {
				receiverList[rec] = true
			}
		}

		for recName := range receiverList {
			node, err := bus.GetNodeByName(recName)
			if err != nil {
				return err
			}
			tmpMsg.AddReceiver(node)
		}
	}

	if dummyNodeUsed {
		if err := bus.AddNode(dummyNode); err != nil {
			return err
		}
	}

	return nil
}
