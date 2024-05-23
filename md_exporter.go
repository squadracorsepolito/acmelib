package acmelib

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	md "github.com/nao1215/markdown"
)

// ExportToMarkdown exports the given [Network] to a markdown document.
// It writes the markdown document to the given [io.Writer].
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
	e.w.Importantf("This markdown document is generated by %s", md.Link("acmelib", "https://github.com/squadracorsepolito/acmelib")).LF()

	e.w.H1(net.name)

	if len(net.desc) > 0 {
		e.w.PlainText(net.desc).LF()
	}

	for _, bus := range net.Buses() {
		e.exportBus(bus)
	}
}

func (e *mdExporter) exportBus(bus *Bus) {
	e.w.H2(bus.name)

	if len(bus.desc) > 0 {
		e.w.PlainText(bus.desc).LF()
	}

	baudrateStr := "-"
	if bus.baudrate != 0 {
		baudrateStr = md.Bold(strconv.Itoa(bus.baudrate))
	}
	e.w.PlainTextf("Baudrate: %s", baudrateStr).LF()

	for _, node := range bus.Nodes() {
		e.exportNode(node)
	}
}

func (e *mdExporter) exportNode(node *NodeInterface) {
	e.w.H3(node.node.name)

	if len(node.node.desc) > 0 {
		e.w.PlainText(node.node.desc).LF()
	}

	e.w.PlainTextf("Node ID: %s", md.Bold(fmt.Sprintf("%d", node.node.id))).LF()

	for _, msg := range node.Messages() {
		e.exportMessage(msg)
	}

	e.w.HorizontalRule()
}

func (e *mdExporter) exportMessage(msg *Message) {
	e.w.H4(msg.name)

	if len(msg.desc) > 0 {
		e.w.PlainText(msg.desc).LF()
	}

	e.w.PlainTextf("CAN-ID: %s", md.Bold(fmt.Sprintf("%d", msg.id))).LF()
	e.w.PlainTextf("Size: %s bytes", md.Bold(fmt.Sprintf("%d", msg.sizeByte))).LF()
	e.w.PlainTextf("Byte Order: %s", md.Bold(msg.byteOrder.String())).LF()

	cycleTimeStr := "-"
	if msg.cycleTime > 0 {
		cycleTimeStr = fmt.Sprintf("%s ms", md.Bold(strconv.Itoa(msg.cycleTime)))
	}
	e.w.PlainTextf("Cycle Time: %s", cycleTimeStr).LF()

	recStr := "Receivers: "
	for idx, recInt := range msg.Receivers() {
		recLink := e.getLink(recInt.node.name)

		if idx == 0 {
			recStr += recLink
			continue
		}

		recStr = fmt.Sprintf("%s, %s", recLink, recLink)
	}
	e.w.PlainText(recStr).LF()

	sigTable := md.TableSet{
		Header: []string{"Name", "Start Bit", "Size", "Min", "Max", "Offset", "Scale", "Unit", "Description"},
		Rows:   [][]string{},
	}
	for _, sig := range msg.Signals() {
		e.exportSignal(sig)
		sigTable.Rows = append(sigTable.Rows, e.sigTableRow)
		e.sigTableRow = []string{}
	}
	e.w.CustomTable(sigTable, md.TableOptions{AutoWrapText: false, AutoFormatHeaders: false})
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
	e.sigTableRow = append(e.sigTableRow, fmt.Sprintf("%g", stdSig.typ.min))

	e.sigTableRow = append(e.sigTableRow, fmt.Sprintf("%g", stdSig.typ.max))

	e.sigTableRow = append(e.sigTableRow, fmt.Sprintf("%g", stdSig.typ.offset))

	e.sigTableRow = append(e.sigTableRow, fmt.Sprintf("%g", stdSig.typ.scale))

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
