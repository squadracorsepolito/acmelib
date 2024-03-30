package dbc

import (
	"fmt"
	"strconv"
	"strings"
)

func writeSlice[T any](slice []T, writeFn func(T), newLineFn func()) {
	for idx, val := range slice {
		writeFn(val)
		if idx == len(slice)-1 {
			newLineFn()
		}
	}
}

type Writer struct {
	f *strings.Builder
}

func NewWriter() *Writer {
	return &Writer{
		f: &strings.Builder{},
	}

}

func (w *Writer) print(format string, a ...any) {
	fmt.Fprintf(w.f, format, a...)
}

func (w *Writer) println(format string, a ...any) {
	fmt.Fprintf(w.f, format+"\n", a...)
}

func (w *Writer) newLine() {
	w.f.WriteRune('\n')
}

func (w *Writer) formatDouble(val float64) string {
	return strconv.FormatFloat(val, 'f', -1, 64)
}

func (w *Writer) formatString(val string) string {
	return "\"" + val + "\""
}

func (w *Writer) formatInt(val int) string {
	return strconv.FormatInt(int64(val), 10)
}

func (w *Writer) formatHexInt(val int) string {
	return "0x" + strconv.FormatInt(int64(val), 16)
}

func (w *Writer) formatUint(val uint32) string {
	return strconv.FormatUint(uint64(val), 10)
}

func (w *Writer) Write(ast *DBC) string {
	if ast.Version != "" {
		w.writeVersion(ast.Version)
	}

	if ast.NewSymbols != nil {
		w.writeNewSymbols(ast.NewSymbols)
	}

	if ast.BitTiming != nil {
		w.writeBitTiming(ast.BitTiming)
	}

	if ast.Nodes != nil {
		w.writeNodes(ast.Nodes)
	}

	writeSlice(ast.ValueTables, w.writeValueTable, w.newLine)
	writeSlice(ast.Messages, w.writeMessage, w.newLine)
	writeSlice(ast.MessageTransmitters, w.writeMessageTransmitter, w.newLine)
	writeSlice(ast.EnvVars, w.writeEnvVar, w.newLine)
	writeSlice(ast.EnvVarDatas, w.writeEnvVarData, w.newLine)
	writeSlice(ast.SignalTypes, w.writeSignalType, w.newLine)
	writeSlice(ast.Comments, w.writeComment, w.newLine)
	writeSlice(ast.Attributes, w.writeAttribute, w.newLine)
	writeSlice(ast.AttributeDefaults, w.writeAttributeDefault, w.newLine)
	writeSlice(ast.AttributeValues, w.writeAttributeValue, w.newLine)
	writeSlice(ast.ValueEncodings, w.writeValueEncoding, w.newLine)
	writeSlice(ast.SignalTypeRefs, w.writeSignalTypeRef, w.newLine)
	writeSlice(ast.SignalGroups, w.writeSignalGroup, w.newLine)
	writeSlice(ast.SignalExtValueTypes, w.writeSignalExtValueType, w.newLine)
	writeSlice(ast.ExtendedMuxes, w.writeExtendedMux, w.newLine)

	return w.f.String()
}

func (w *Writer) writeVersion(ver string) {
	w.println("%s %s", getKeyword(keywordVersion), w.formatString(ver))
	w.newLine()
}

func (w *Writer) writeNewSymbols(ns *NewSymbols) {
	w.println("%s:", getKeyword(keywordNewSymbols))
	for _, symbol := range ns.Symbols {
		w.println("\t%s", symbol)
	}
	w.newLine()
}

func (w *Writer) writeBitTiming(bitTime *BitTiming) {
	w.print("%s:", getKeyword(keywordBitTiming))
	defer w.newLine()

	if bitTime.Baudrate == 0 {
		w.newLine()
		return
	}

	w.println("%s : %s, %s",
		w.formatUint(bitTime.Baudrate),
		w.formatUint(bitTime.BitTimingReg1),
		w.formatUint(bitTime.BitTimingReg2),
	)
}

func (w *Writer) writeNodes(nodes *Nodes) {
	w.print("%s:", getKeyword(keywordNode))
	for _, node := range nodes.Names {
		w.print(" %s", node)
	}
	w.newLine()
	w.newLine()
}

func (w *Writer) writeValueDescription(valDesc *ValueDescription) {
	w.print(" %s %s", w.formatUint(valDesc.ID), w.formatString(valDesc.Name))
}

func (w *Writer) writeValueTable(valTable *ValueTable) {
	w.print("%s %s", getKeyword(keywordValueTable), valTable.Name)
	for _, valDesc := range valTable.Values {
		w.writeValueDescription(valDesc)
	}
	w.println(";")
}

func (w *Writer) writeMessage(msg *Message) {
	w.println("%s %s %s : %s %s",
		getKeyword(keywordMessage),
		w.formatUint(msg.ID),
		msg.Name,
		w.formatUint(msg.Size),
		msg.Transmitter,
	)

	for _, sig := range msg.Signals {
		w.writeSignal(sig)
	}

	w.newLine()
}

func (w *Writer) writeSignal(sig *Signal) {
	w.print("\t%s %s", getKeyword(keywordSignal), sig.Name)

	if sig.IsMultiplexed && sig.IsMultiplexor {
		w.print(" m%sM", w.formatUint(sig.MuxSwitchValue))
	} else if sig.IsMultiplexed {
		w.print(" m%s", w.formatUint(sig.MuxSwitchValue))
	} else if sig.IsMultiplexor {
		w.print(" M")
	}

	w.print(" : %s|%s@", w.formatUint(sig.StartBit), w.formatUint(sig.Size))

	switch sig.ByteOrder {
	case SignalBigEndian:
		w.print("0")
	case SignalLittleEndian:
		w.print("1")
	}

	switch sig.ValueType {
	case SignalUnsigned:
		w.print("+")
	case SignalSigned:
		w.print("-")
	}

	w.print(" (%s,%s)", w.formatDouble(sig.Factor), w.formatDouble(sig.Offset))
	w.print(" [%s|%s]", w.formatDouble(sig.Min), w.formatDouble(sig.Max))
	w.print(" %s", w.formatString(sig.Unit))

	for idx, receiver := range sig.Receivers {
		if idx > 0 {
			w.print(",")
		}
		w.print(" %s", receiver)
	}
	w.newLine()
}

func (w *Writer) writeMessageTransmitter(msgTx *MessageTransmitter) {
	w.print("%s %s :", getKeyword(keywordMessageTransmitter), w.formatUint(msgTx.MessageID))
	for _, tx := range msgTx.Transmitters {
		w.print(" %s", tx)
	}
	w.println(";")
}

func (w *Writer) writeEnvVar(envVar *EnvVar) {
	w.print("%s %s : ", getKeyword(keywordEnvVar), envVar.Name)

	switch envVar.Type {
	case EnvVarInt:
		w.print("0")
	case EnvVarFloat:
		w.print("1")
	case EnvVarString:
		w.print("2")
	}

	w.print(" [%s|%s] ", w.formatDouble(envVar.Min), w.formatDouble(envVar.Max))
	w.print("%s %s %s ",
		w.formatString(envVar.Unit),
		w.formatDouble(envVar.InitialValue),
		w.formatUint(envVar.ID),
	)

	for accType, acc := range envVarAccessTypes {
		if acc == envVar.AccessType {
			w.print("%s", accType)
		}
	}

	for _, node := range envVar.AccessNodes {
		w.print(" , %s", node)
	}

	w.println(";")
}

func (w *Writer) writeEnvVarData(envVarData *EnvVarData) {
	w.println("%s %s : %s ;",
		getKeyword(keywordEnvVarData),
		envVarData.EnvVarName,
		w.formatUint(envVarData.DataSize),
	)
}

func (w *Writer) writeSignalType(sigTyp *SignalType) {
	w.print("%s %s : %s@", getKeyword(keywordSignalType), sigTyp.TypeName, w.formatUint(sigTyp.Size))

	switch sigTyp.ByteOrder {
	case SignalLittleEndian:
		w.print("1 ")
	case SignalBigEndian:
		w.print("0 ")
	}

	switch sigTyp.ValueType {
	case SignalUnsigned:
		w.print("+")
	case SignalSigned:
		w.print("-")
	}

	w.print(" (%s,%s)", w.formatDouble(sigTyp.Factor), w.formatDouble(sigTyp.Offset))
	w.print(" [%s|%s]", w.formatDouble(sigTyp.Min), w.formatDouble(sigTyp.Max))
	w.print(" %s %s , %s;", w.formatString(sigTyp.Unit), w.formatDouble(sigTyp.DefaultValue), sigTyp.ValueTableName)
	w.newLine()
}

func (w *Writer) writeComment(com *Comment) {
	w.print("%s ", getKeyword(keywordComment))

	switch com.Kind {
	case CommentNode:
		w.print("%s %s ", getKeyword(keywordNode), com.NodeName)
	case CommentMessage:
		w.print("%s %s ", getKeyword(keywordMessage), w.formatUint(com.MessageID))
	case CommentSignal:
		w.print("%s %s %s ", getKeyword(keywordSignal), w.formatUint(com.MessageID), com.SignalName)
	case CommentEnvVar:
		w.print("%s %s ", getKeyword(keywordEnvVar), com.EnvVarName)
	}

	w.println("%s;", w.formatString(com.Text))
}

func (w *Writer) writeAttribute(att *Attribute) {
	w.print("%s ", getKeyword(keywordAttribute))

	switch att.Kind {
	case AttributeNode:
		w.print("%s ", getKeyword(keywordNode))
	case AttributeMessage:
		w.print("%s ", getKeyword(keywordMessage))
	case AttributeSignal:
		w.print("%s ", getKeyword(keywordSignal))
	case AttributeEnvVar:
		w.print("%s ", getKeyword(keywordEnvVar))
	}

	w.print(`"%s" `, att.Name)

	switch att.Type {
	case AttributeInt:
		w.print("%s %s %s", getKeyword(keywordAttributeInt), w.formatInt(att.MinInt), w.formatInt(att.MaxInt))
	case AttributeHex:
		w.print("%s %s %s", getKeyword(keywordAttributeHex), w.formatHexInt(att.MinHex), w.formatHexInt(att.MaxHex))
	case AttributeString:
		w.print("%s", getKeyword(keywordAttributeString))
	case AttributeFloat:
		w.print("%s %s %s", getKeyword(keywordAttributeFloat), w.formatDouble(att.MinFloat), w.formatDouble(att.MaxFloat))
	case AttributeEnum:
		w.print("%s", getKeyword(keywordAttributeEnum))
		for idx, enumVal := range att.EnumValues {
			if idx != 0 {
				w.print(",")
			}
			w.print(" %s", w.formatString(enumVal))
		}
	}

	w.println(";")
}

func (w *Writer) writeAttributeDefault(attDef *AttributeDefault) {
	w.print(`%s "%s" `, getKeyword(keywordAttributeDefault), attDef.AttributeName)

	switch attDef.Type {
	case AttributeDefaultInt:
		w.print(w.formatInt(attDef.ValueInt))
	case AttributeDefaultHex:
		w.print(w.formatHexInt(attDef.ValueHex))
	case AttributeDefaultFloat:
		w.print(w.formatDouble(attDef.ValueFloat))
	case AttributeDefaultString:
		w.print(w.formatString(attDef.ValueString))
	}

	w.println(";")
}

func (w *Writer) writeAttributeValue(attVal *AttributeValue) {
	w.print(`%s "%s "`, getKeyword(keywordAttributeValue), attVal.AttributeName)

	switch attVal.AttributeKind {
	case AttributeNode:
		w.print("%s %s ", getKeyword(keywordNode), attVal.NodeName)
	case AttributeMessage:
		w.print("%s %s ", getKeyword(keywordMessage), w.formatUint(attVal.MessageID))
	case AttributeSignal:
		w.print("%s %s %s ", getKeyword(keywordSignal), w.formatUint(attVal.MessageID), attVal.SignalName)
	case AttributeEnvVar:
		w.print("%s %s ", getKeyword(keywordEnvVar), attVal.EnvVarName)
	}

	switch attVal.Type {
	case AttributeValueInt:
		w.print(w.formatInt(attVal.ValueInt))
	case AttributeValueHex:
		w.print(w.formatHexInt(attVal.ValueHex))
	case AttributeValueFloat:
		w.print(w.formatDouble(attVal.ValueFloat))
	case AttributeValueString:
		w.print(w.formatString(attVal.ValueString))
	}

	w.println(";")
}

func (w *Writer) writeValueEncoding(valEnc *ValueEncoding) {
	w.print("%s ", getKeyword(keywordValueEncoding))

	switch valEnc.Kind {
	case ValueEncodingSignal:
		w.print("%s %s", w.formatUint(valEnc.MessageID), valEnc.SignalName)
	case ValueEncodingEnvVar:
		w.print("%s", valEnc.EnvVarName)
	}

	for _, valDesc := range valEnc.Values {
		w.writeValueDescription(valDesc)
	}
	w.println(";")
}

func (w *Writer) writeSignalTypeRef(sigTypRef *SignalTypeRef) {
	w.println("%s %s %s : %s;",
		getKeyword(keywordSignalType),
		w.formatUint(sigTypRef.MessageID),
		sigTypRef.SignalName,
		sigTypRef.TypeName,
	)
}

func (w *Writer) writeSignalGroup(sigGroup *SignalGroup) {
	w.print("%s %s %s %s :",
		getKeyword(keywordSignalGroup),
		w.formatUint(sigGroup.MessageID),
		sigGroup.GroupName,
		w.formatUint(sigGroup.Repetitions),
	)

	for _, sigName := range sigGroup.SignalNames {
		w.print(" %s", sigName)
	}
	w.println(";")
}

func (w *Writer) writeSignalExtValueType(sigExtValTyp *SignalExtValueType) {
	w.print("%s %s %s ",
		getKeyword(keywordSignalValueType),
		w.formatUint(sigExtValTyp.MessageID),
		sigExtValTyp.SignalName,
	)

	switch sigExtValTyp.ExtValueType {
	case SignalExtValueTypeInteger:
		w.print("0")
	case SignalExtValueTypeFloat:
		w.print("1")
	case SignalExtValueTypeDouble:
		w.print("2")
	}
	w.println(";")
}

func (w *Writer) writeExtendedMux(extMux *ExtendedMux) {
	w.print("%s %s", getKeyword(keywordExtendedMux), w.formatUint(extMux.MessageID))
	w.print(" %s %s", extMux.MultiplexedName, extMux.MultiplexorName)

	for idx, r := range extMux.Ranges {
		if idx != 0 {
			w.print(",")
		}
		w.print(" %s-%s", w.formatUint(r.From), w.formatUint(r.To))
	}
	w.println(";")
}
