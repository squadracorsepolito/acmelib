package acmelib

import (
	"fmt"
	"io"
	"strings"

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

func (e *mdExporter) getLink(text string) string {
	return md.Link(text, "#"+strings.ReplaceAll(text, " ", ""))
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
		e.w.LF()
	}

	e.w.PlainTextf("Node ID: %s", md.Bold(fmt.Sprintf("%d", node.id)))
	e.w.LF()

	for _, msg := range node.Messages() {
		e.exportMessage(msg)
	}

	e.w.HorizontalRule()
}

func (e *mdExporter) exportMessage(msg *Message) {
	e.w.H4(msg.name)

	if len(msg.desc) > 0 {
		e.w.PlainText(msg.desc)
		e.w.LF()
	}

	e.w.PlainTextf("CAN-ID: %s", md.Bold(fmt.Sprintf("%d", msg.id)))
	e.w.LF()

	e.w.PlainTextf("Size: %s bytes", md.Bold(fmt.Sprintf("%d", msg.sizeByte)))
	e.w.LF()

	if msg.cycleTime > 0 {
		e.w.PlainTextf("Cycle Time: %s ms", md.Bold(fmt.Sprintf("%d", msg.cycleTime)))
	} else {
		e.w.PlainText("Cycle Time: -")
	}
	e.w.LF()

	recStr := "Receivers: "
	for idx, rec := range msg.Receivers() {
		recLink := e.getLink(rec.name)

		if idx == 0 {
			recStr += recLink
			continue
		}

		recStr = fmt.Sprintf("%s, %s", recLink, recLink)
	}
	e.w.PlainText(recStr)
	e.w.LF()

	sigTable := md.TableSet{
		Header: []string{"Name", "Start Bit", "Size", "Min", "Max", "Offset", "Scale", "Unit", "Description"},
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

	e.sigTableRow = append(e.sigTableRow, fmt.Sprintf("%d", sig.GetSize()))

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
	e.sigTableRow = append(e.sigTableRow, fmt.Sprintf("%g", stdSig.min))

	e.sigTableRow = append(e.sigTableRow, fmt.Sprintf("%g", stdSig.max))

	e.sigTableRow = append(e.sigTableRow, fmt.Sprintf("%g", stdSig.offset))

	e.sigTableRow = append(e.sigTableRow, fmt.Sprintf("%g", stdSig.scale))

	unitSymbol := "-"
	if stdSig.unit != nil {
		unitSymbol = stdSig.unit.symbol
	}
	e.sigTableRow = append(e.sigTableRow, unitSymbol)
}

func (e *mdExporter) exportEnumSignal(enumSig *EnumSignal) {
	e.sigTableRow = append(e.sigTableRow, "0")

	e.sigTableRow = append(e.sigTableRow, fmt.Sprintf("%d", enumSig.GetSize()))

	e.sigTableRow = append(e.sigTableRow, "0")

	e.sigTableRow = append(e.sigTableRow, "1")

	e.sigTableRow = append(e.sigTableRow, "-")
}

func (e *mdExporter) exportMultiplexerSignal(_ *MultiplexerSignal) {
	e.sigTableRow = append(e.sigTableRow, "-")

	e.sigTableRow = append(e.sigTableRow, "-")

	e.sigTableRow = append(e.sigTableRow, "-")

	e.sigTableRow = append(e.sigTableRow, "-")

	e.sigTableRow = append(e.sigTableRow, "-")
}
