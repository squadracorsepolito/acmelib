package acmelib

import (
	"fmt"
	"io"

	md "github.com/nao1215/markdown"
)

func ExportToMarkdown(network *Network, w io.Writer) error {
	mdWriter := md.NewMarkdown(w)
	exporter := newMDExporter(mdWriter)
	exporter.exportNetwork(network)
	return mdWriter.Build()
}

type mdExporter struct {
	w *md.Markdown

	sigTableRow []string
}

func newMDExporter(mdWriter *md.Markdown) *mdExporter {
	return &mdExporter{
		w: mdWriter,

		sigTableRow: []string{},
	}
}

func (e *mdExporter) exportNetwork(net *Network) {
	e.w.H1(net.name)

	if len(net.desc) > 0 {
		e.w.PlainText(net.desc)
	}

	for _, bus := range net.Buses() {
		e.exportBus(bus)
	}
}

func (e *mdExporter) exportBus(bus *Bus) {
	e.w.H2(bus.name)

	if len(bus.desc) > 0 {
		e.w.PlainText(bus.desc)
	}

	if bus.baudrate != 0 {
		e.w.BulletList(fmt.Sprintf("Baudrate: %d", bus.baudrate))
	}

	for _, node := range bus.Nodes() {
		e.exportNode(node)
	}
}

func (e *mdExporter) exportNode(node *Node) {
	e.w.H3(node.name)

	if len(node.desc) > 0 {
		e.w.PlainText(node.desc)
	}

	for _, msg := range node.Messages() {
		e.exportMessage(msg)
	}
}

func (e *mdExporter) exportMessage(msg *Message) {
	e.w.H4(msg.name)

	if len(msg.desc) > 0 {
		e.w.PlainText(msg.desc)
	}

	sigTable := md.TableSet{
		Header: []string{"Name", "Start Bit", "Type", "Min", "Max", "Offset", "Scale", "Description"},
		Rows:   [][]string{},
	}
	for _, sig := range msg.Signals() {
		e.exportSignal(sig)
		sigTable.Rows = append(sigTable.Rows, e.sigTableRow)
		e.sigTableRow = []string{}
	}
	e.w.Table(sigTable)
}

func (e *mdExporter) exportSignal(sig Signal) {
	e.sigTableRow = append(e.sigTableRow, sig.Name())

	e.sigTableRow = append(e.sigTableRow, fmt.Sprintf("%d", sig.GetStartBit()))

	switch sig.Kind() {
	case SignalKindStandard:
		stdSig, err := sig.ToStandard()
		if err != nil {
			panic(err)
		}
		e.exportStandardSignal(stdSig)

	case SignalKindEnum:
		enumSig, err := sig.ToEnum()
		if err != nil {
			panic(err)
		}
		e.exportEnumSignal(enumSig)

	case SignalKindMultiplexer:
		muxSig, err := sig.ToMultiplexer()
		if err != nil {
			panic(err)
		}
		e.exportMultiplexerSignal(muxSig)
	}

	desc := sig.Desc()
	if len(desc) == 0 {
		desc = "-"
	}
	e.sigTableRow = append(e.sigTableRow, desc)
}

func (e *mdExporter) exportStandardSignal(stdSig *StandardSignal) {
	e.sigTableRow = append(e.sigTableRow, stdSig.typ.name)

	e.sigTableRow = append(e.sigTableRow, fmt.Sprintf("%g", stdSig.min))

	e.sigTableRow = append(e.sigTableRow, fmt.Sprintf("%g", stdSig.max))

	e.sigTableRow = append(e.sigTableRow, fmt.Sprintf("%g", stdSig.offset))

	e.sigTableRow = append(e.sigTableRow, fmt.Sprintf("%g", stdSig.scale))
}

func (e *mdExporter) exportEnumSignal(enumSig *EnumSignal) {
	e.sigTableRow = append(e.sigTableRow, enumSig.enum.name)

	e.sigTableRow = append(e.sigTableRow, "0")

	e.sigTableRow = append(e.sigTableRow, fmt.Sprintf("%d", enumSig.GetSize()))

	e.sigTableRow = append(e.sigTableRow, "0")

	e.sigTableRow = append(e.sigTableRow, "1")
}

func (e *mdExporter) exportMultiplexerSignal(mucSig *MultiplexerSignal) {
	e.sigTableRow = append(e.sigTableRow, "-")

	e.sigTableRow = append(e.sigTableRow, "-")

	e.sigTableRow = append(e.sigTableRow, "-")

	e.sigTableRow = append(e.sigTableRow, "-")

	e.sigTableRow = append(e.sigTableRow, "-")
}
