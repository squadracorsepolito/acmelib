package dbc

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Parser struct {
	s *scanner

	usePrev bool
	curr    *token

	foundVer    bool
	foundNewSym bool
	foundBitTim bool
	foundNode   bool
}

func (p *Parser) parseUint(val string) (uint32, error) {
	res, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(res), nil
}

func (p *Parser) parseHexInt(val string) (int, error) {
	if !strings.HasPrefix(val, "0x") && !strings.HasPrefix(val, "0X") {
		return 0, errors.New("invalid hex number")
	}
	res, err := strconv.ParseUint(val[2:], 16, 32)
	if err != nil {
		return 0, err
	}
	return int(res), nil
}

func (p *Parser) parseInt(val string) (int, error) {
	res, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}
	return int(res), nil
}

func (p *Parser) parseDouble(val string) (float64, error) {
	return strconv.ParseFloat(val, 64)
}

func NewParser(file []byte) *Parser {
	return &Parser{
		s: newScanner(bytes.NewReader(file), string(file)),

		usePrev: false,

		foundVer:    false,
		foundNewSym: false,
		foundBitTim: false,
		foundNode:   false,
	}
}

func (p *Parser) scan() *token {
	if p.usePrev {
		p.usePrev = false
		return p.curr
	}

	token := p.s.scan()
	if token.isSpace() {
		token = p.s.scan()
	}
	p.curr = token

	return token
}

func (p *Parser) unscan() {
	p.usePrev = true
}

func (p *Parser) errorf(msg string) error {
	return fmt.Errorf(`%s; got "%s" at line %d, column %d`, msg, p.curr.value, p.curr.line, p.curr.col)
}

func (p *Parser) expectSyntax(kind syntaxKind) error {
	if !p.scan().isSyntax(kind) {
		return p.errorf(fmt.Sprintf("expected %q", getSyntaxRune(kind)))
	}
	return nil
}

func (p *Parser) Parse() (*DBC, error) {
	ast := new(DBC)
	t := p.scan()

	for !t.isEOF() {
		switch t.kind {
		case tokenError:
			return nil, p.errorf("syntax error")

		case tokenKeyword:
			keywordKind := getKeywordKind(t.value)
			switch keywordKind {
			case keywordVersion:
				ver, err := p.parseVersion()
				if err != nil {
					return nil, err
				}
				ast.Version = ver

			case keywordNewSymbols:
				ns, err := p.parseNewSymbols()
				if err != nil {
					return nil, err
				}
				ast.NewSymbols = ns

			case keywordBitTiming:
				bt, err := p.parseBitTiming()
				if err != nil {
					return nil, err
				}
				ast.BitTiming = bt

			case keywordNode:
				node, err := p.parseNodes()
				if err != nil {
					return nil, err
				}
				ast.Nodes = node

			case keywordValueTable:
				vt, err := p.parseValueTable()
				if err != nil {
					return nil, err
				}
				ast.ValueTables = append(ast.ValueTables, vt)

			case keywordMessage:
				message, err := p.parseMessage()
				if err != nil {
					return nil, err
				}
				ast.Messages = append(ast.Messages, message)

			case keywordMessageTransmitter:
				mt, err := p.parseMessageTransmitter()
				if err != nil {
					return nil, err
				}
				ast.MessageTransmitters = append(ast.MessageTransmitters, mt)

			case keywordEnvVar:
				envVar, err := p.parseEnvVar()
				if err != nil {
					return nil, err
				}
				ast.EnvVars = append(ast.EnvVars, envVar)

			case keywordEnvVarData:
				evData, err := p.parseEnvVarData()
				if err != nil {
					return nil, err
				}
				ast.EnvVarDatas = append(ast.EnvVarDatas, evData)

			case keywordSignalType:
				sigType, sigTypeRef, err := p.parseSignalType()
				if err != nil {
					return nil, err
				}
				if sigTypeRef != nil {
					ast.SignalTypeRefs = append(ast.SignalTypeRefs, sigTypeRef)
				} else {
					ast.SignalTypes = append(ast.SignalTypes, sigType)
				}

			case keywordComment:
				com, err := p.parseComment()
				if err != nil {
					return nil, err
				}
				ast.Comments = append(ast.Comments, com)

			case keywordAttribute:
				att, err := p.parseAttribute()
				if err != nil {
					return nil, err
				}
				ast.Attributes = append(ast.Attributes, att)

			case keywordAttributeDefault:
				attDef, err := p.parseAttributeDefault()
				if err != nil {
					return nil, err
				}
				ast.AttributeDefaults = append(ast.AttributeDefaults, attDef)

			case keywordAttributeValue:
				attVal, err := p.parseAttributeValue()
				if err != nil {
					return nil, err
				}
				ast.AttributeValues = append(ast.AttributeValues, attVal)

			case keywordValueEncoding:
				valEnc, err := p.parseValueEncoding()
				if err != nil {
					return nil, err
				}
				ast.ValueEncodings = append(ast.ValueEncodings, valEnc)

			case keywordSignalGroup:
				sigGroup, err := p.parseSignalGroup()
				if err != nil {
					return nil, err
				}
				ast.SignalGroups = append(ast.SignalGroups, sigGroup)

			case keywordSignalValueType:
				sigExtValType, err := p.parseSignalExtValueType()
				if err != nil {
					return nil, err
				}
				ast.SignalExtValueTypes = append(ast.SignalExtValueTypes, sigExtValType)

			case keywordExtendedMux:
				extMux, err := p.parseExtendedMux()
				if err != nil {
					return nil, err
				}
				ast.ExtendedMuxes = append(ast.ExtendedMuxes, extMux)
			}

		default:
			return nil, p.errorf("unexpected token")
		}

		t = p.scan()
	}

	return ast, nil
}

func (p *Parser) parseVersion() (string, error) {
	if p.foundVer {
		return "", p.errorf("duplicated version")
	}
	p.foundVer = true

	t := p.scan()
	if !t.isString() {
		return "", p.errorf("expected version")
	}
	return t.value, nil
}

func (p *Parser) parseNewSymbols() (*NewSymbols, error) {
	if p.foundNewSym {
		return nil, p.errorf("duplicated new symbols")
	}
	p.foundNewSym = true

	ns := new(NewSymbols)

	if err := p.expectSyntax(syntaxColon); err != nil {
		return nil, err
	}

	posValues := make(map[string]struct{})
	for _, val := range newSymbolsValues {
		posValues[val] = struct{}{}
	}

	for {
		t := p.scan()
		if t.isEOF() {
			break
		}

		if t.isKeyword(keywordBitTiming) {
			p.unscan()
			break
		}

		if t.kind == tokenKeyword || t.isIdent() {
			if _, ok := posValues[t.value]; !ok {
				return nil, p.errorf("invalid new symbol")
			}

			ns.Symbols = append(ns.Symbols, t.value)
		}
	}

	return ns, nil
}

func (p *Parser) parseBitTiming() (*BitTiming, error) {
	if p.foundBitTim {
		return nil, p.errorf("duplicated bit timing")
	}

	bt := new(BitTiming)

	if err := p.expectSyntax(syntaxColon); err != nil {
		return nil, err
	}

	t := p.scan()
	if t.isKeyword(keywordNode) {
		p.unscan()
		return bt, nil
	} else if !t.isNumber() {
		return nil, p.errorf("expected bit timing baudrate")
	}

	baudrate, err := p.parseUint(t.value)
	if err != nil {
		return nil, err
	}
	bt.Baudrate = baudrate

	if err := p.expectSyntax(syntaxColon); err != nil {
		return nil, err
	}

	t = p.scan()
	if !t.isNumber() {
		return nil, p.errorf("expected bit timing for register 1")
	}
	btr1, err := p.parseUint(t.value)
	if err != nil {
		return nil, err
	}
	bt.BitTimingReg1 = btr1

	if err := p.expectSyntax(syntaxComma); err != nil {
		return nil, err
	}

	t = p.scan()
	if !t.isNumber() {
		return nil, p.errorf("expected bit timing for register 2")
	}
	btr2, err := p.parseUint(t.value)
	if err != nil {
		return nil, err
	}
	bt.BitTimingReg2 = btr2

	return bt, nil
}

func (p *Parser) parseNodeName() (string, error) {
	t := p.scan()
	if !t.isIdent() {
		return "", p.errorf("expected node name")
	}
	return t.value, nil
}

func (p *Parser) parseNodes() (*Nodes, error) {
	if p.foundNode {
		return nil, p.errorf("duplicated node definition")
	}
	p.foundNode = true

	node := new(Nodes)

	if err := p.expectSyntax(syntaxColon); err != nil {
		return nil, err
	}

	for {
		t := p.scan()
		if !t.isIdent() {
			p.unscan()
			break
		}
		node.Names = append(node.Names, t.value)
	}

	return node, nil
}

func (p *Parser) parseValueDescription() (*ValueDescription, error) {
	valDesc := new(ValueDescription)

	t := p.scan()
	if !t.isNumber() {
		p.unscan()
		return valDesc, nil
	}

	valID, err := p.parseUint(t.value)
	if err != nil {
		return nil, p.errorf("cannot parse value description id as uint")
	}
	valDesc.ID = valID

	t = p.scan()
	if !t.isString() {
		return nil, p.errorf("expected value description name after id")
	}
	valDesc.Name = t.value

	return valDesc, nil
}

func (p *Parser) parseValueTable() (*ValueTable, error) {
	vt := new(ValueTable)

	t := p.scan()
	if !t.isIdent() {
		return nil, p.errorf("expected value table name")
	}
	vt.Name = t.value

	for {
		t := p.scan()
		p.unscan()
		if !t.isNumber() {
			break
		}

		vd, err := p.parseValueDescription()
		if err != nil {
			return nil, err
		}

		vt.Values = append(vt.Values, vd)
	}

	if err := p.expectSyntax(syntaxSemicolon); err != nil {
		return nil, err
	}

	return vt, nil
}

func (p *Parser) parseMessageID() (uint32, error) {
	t := p.scan()
	if !t.isNumber() {
		return 0, p.errorf("expected message id")
	}
	id, err := p.parseUint(t.value)
	if err != nil {
		return 0, p.errorf("cannot parse message id as uint")
	}
	return id, nil
}

func (p *Parser) parseMessage() (*Message, error) {
	msg := new(Message)

	id, err := p.parseMessageID()
	if err != nil {
		return nil, err
	}
	msg.ID = id

	t := p.scan()
	if !t.isIdent() {
		return nil, p.errorf("expected message name")
	}
	msg.Name = t.value

	if err := p.expectSyntax(syntaxColon); err != nil {
		return nil, err
	}

	t = p.scan()
	if !t.isNumber() {
		return nil, p.errorf("expected message size")
	}
	size, err := p.parseUint(t.value)
	if err != nil {
		return nil, p.errorf("cannot parse message size as uint")
	}
	msg.Size = size

	t = p.scan()
	if !t.isIdent() {
		return nil, p.errorf("expected message transmitter")
	}
	msg.Transmitter = t.value

	for {
		if !p.scan().isKeyword(keywordSignal) {
			p.unscan()
			break
		}
		sig, err := p.parseSignal()
		if err != nil {
			return nil, err
		}
		msg.Signals = append(msg.Signals, sig)
	}

	return msg, nil
}

func (p *Parser) parseSignalName() (string, error) {
	t := p.scan()
	if !t.isIdent() {
		return "", p.errorf("expected signal name")
	}
	return t.value, nil
}

func (p *Parser) parseSignal() (*Signal, error) {
	sig := new(Signal)

	name, err := p.parseSignalName()
	if err != nil {
		return nil, err
	}
	sig.Name = name

	t := p.scan()
	if t.isMuxIndicator() {
		if t.value[len(t.value)-1] == 'M' {
			sig.IsMultiplexor = true
		}
		if t.value[0] == 'm' {
			strNum := ""
			if sig.IsMultiplexor {
				strNum = t.value[1 : len(t.value)-1]
			} else {
				strNum = t.value[1:]
			}
			switchNum, err := p.parseUint(strNum)
			if err != nil {
				return nil, p.errorf("cannot parse signal multiplexer switch number as uint")
			}
			sig.IsMultiplexed = true
			sig.MuxSwitchValue = switchNum
		}
	} else {
		p.unscan()
	}

	if err := p.expectSyntax(syntaxColon); err != nil {
		return nil, err
	}

	t = p.scan()
	if !t.isNumber() {
		return nil, p.errorf("expected signal start bit")
	}
	startBit, err := p.parseUint(t.value)
	if err != nil {
		return nil, p.errorf("cannot parse signal start bit as uint")
	}
	sig.StartBit = startBit

	if err := p.expectSyntax(syntaxPipe); err != nil {
		return nil, err
	}

	t = p.scan()
	if !t.isNumber() {
		return nil, p.errorf("expected signal size")
	}
	size, err := p.parseUint(t.value)
	if err != nil {
		return nil, p.errorf("cannot parse signal size as uint")
	}
	sig.Size = size

	if err := p.expectSyntax(syntaxAt); err != nil {
		return nil, err
	}

	t = p.scan()
	if !t.isNumber() {
		return nil, p.errorf("expected signal byte order")
	}
	byteOrder, err := p.parseUint(t.value)
	if err != nil {
		return nil, p.errorf("cannot parse signal byte order as uint")
	}
	if byteOrder == 0 {
		sig.ByteOrder = SignalBigEndian
	} else if byteOrder == 1 {
		sig.ByteOrder = SignalLittleEndian
	} else {
		return nil, p.errorf("signal byte order must be 0 or 1")
	}

	t = p.scan()
	syntKind := getSyntaxKind(t.value)
	if t.kind != tokenSyntax || (syntKind != syntaxPlus && syntKind != syntaxMinus) {
		return nil, p.errorf(`expected "+" or "-"`)
	}
	if t.value == "+" {
		sig.ValueType = SignalUnsigned
	} else {
		sig.ValueType = SignalSigned
	}

	if err := p.expectSyntax(syntaxLeftParen); err != nil {
		return nil, err
	}

	t = p.scan()
	if !t.isNumber() {
		return nil, p.errorf("expected signal factor")
	}
	factor, err := p.parseDouble(t.value)
	if err != nil {
		return nil, p.errorf("cannot parse signal factor as double")
	}
	sig.Factor = factor

	if err := p.expectSyntax(syntaxComma); err != nil {
		return nil, err
	}

	t = p.scan()
	if !t.isNumber() {
		return nil, p.errorf("expected signal offset")
	}
	offset, err := p.parseDouble(t.value)
	if err != nil {
		return nil, p.errorf("cannot parse signal offset as double")
	}
	sig.Offset = offset

	if err := p.expectSyntax(syntaxRightParen); err != nil {
		return nil, err
	}

	if err := p.expectSyntax(syntaxLeftSquareBrace); err != nil {
		return nil, err
	}

	t = p.scan()
	if !t.isNumber() {
		return nil, p.errorf("expected signal minimum")
	}
	min, err := p.parseDouble(t.value)
	if err != nil {
		return nil, p.errorf("cannot parse signal minimum as double")
	}
	sig.Min = min

	if err := p.expectSyntax(syntaxPipe); err != nil {
		return nil, err
	}

	t = p.scan()
	if !t.isNumber() {
		return nil, p.errorf("expected signal maximum")
	}
	max, err := p.parseDouble(t.value)
	if err != nil {
		return nil, p.errorf("cannot parse signal maximum as double")
	}
	sig.Max = max

	if err := p.expectSyntax(syntaxRightSquareBrace); err != nil {
		return nil, err
	}

	t = p.scan()
	if !t.isString() {
		return nil, p.errorf("expected signal unit")
	}
	sig.Unit = t.value

	t = p.scan()
	if !t.isIdent() {
		return nil, p.errorf("expected signal receiver")
	}
	sig.Receivers = append(sig.Receivers, t.value)
	for {
		if !p.scan().isSyntax(syntaxComma) {
			p.unscan()
			break
		}
		t = p.scan()
		if !t.isIdent() {
			return nil, p.errorf("expected signal receiver")
		}
		sig.Receivers = append(sig.Receivers, t.value)
	}

	return sig, nil
}

func (p *Parser) parseMessageTransmitter() (*MessageTransmitter, error) {
	mt := new(MessageTransmitter)

	msgID, err := p.parseMessageID()
	if err != nil {
		return nil, err
	}
	mt.MessageID = msgID

	for {
		t := p.scan()
		p.unscan()
		if !t.isIdent() {
			break
		}

		node, err := p.parseNodeName()
		if err != nil {
			return nil, err
		}
		mt.Transmitters = append(mt.Transmitters, node)
	}

	if err := p.expectSyntax(syntaxSemicolon); err != nil {
		return nil, err
	}

	return mt, nil
}

func (p *Parser) parseEnvVarName() (string, error) {
	t := p.scan()
	if !t.isIdent() {
		return "", p.errorf("expected envvar name")
	}
	return t.value, nil
}

func (p *Parser) parseEnvVar() (*EnvVar, error) {
	envVar := new(EnvVar)

	t := p.scan()
	if !t.isIdent() {
		return nil, p.errorf("expected envvar name")
	}
	envVar.Name = t.value

	if err := p.expectSyntax(syntaxColon); err != nil {
		return nil, err
	}

	t = p.scan()
	if !t.isNumber() {
		return nil, p.errorf("expected envvar type")
	}
	typ, err := p.parseUint(t.value)
	if err != nil {
		return nil, p.errorf("cannot parse envvar type as uint")
	}
	if typ == 0 {
		envVar.Type = EnvVarInt
	} else if typ == 1 {
		envVar.Type = EnvVarFloat
	} else if typ == 2 {
		envVar.Type = EnvVarString
	} else {
		return nil, p.errorf("envvar type must be 0, 1 or 2")
	}

	if err := p.expectSyntax(syntaxLeftSquareBrace); err != nil {
		return nil, err
	}

	t = p.scan()
	if !t.isNumber() {
		return nil, p.errorf("expected envvar minimum value")
	}
	min, err := p.parseDouble(t.value)
	if err != nil {
		return nil, p.errorf("cannot parse envvar minimum value as double")
	}
	envVar.Min = min

	if err := p.expectSyntax(syntaxPipe); err != nil {
		return nil, err
	}

	t = p.scan()
	if !t.isNumber() {
		return nil, p.errorf("expected envvar maximum value")
	}
	max, err := p.parseDouble(t.value)
	if err != nil {
		return nil, p.errorf("cannot parse envvar maximum value as double")
	}
	envVar.Max = max

	if err := p.expectSyntax(syntaxLeftSquareBrace); err != nil {
		return nil, err
	}

	t = p.scan()
	if !t.isString() {
		return nil, p.errorf("expected envvar unit")
	}
	envVar.Unit = t.value

	t = p.scan()
	if !t.isNumber() {
		return nil, p.errorf("expected envvar initial value")
	}
	initialVal, err := p.parseDouble(t.value)
	if err != nil {
		return nil, p.errorf("cannot parse envvar initial value as double")
	}
	envVar.InitialValue = initialVal

	t = p.scan()
	if !t.isNumber() {
		return nil, p.errorf("expected envvar id")
	}
	id, err := p.parseUint(t.value)
	if err != nil {
		return nil, p.errorf("cannot parse envvar id as uint")
	}
	envVar.ID = id

	t = p.scan()
	if !t.isIdent() {
		return nil, p.errorf("expected envvar access type")
	}
	accTyp, foundAccTyp := envVarAccessTypes[t.value]
	if !foundAccTyp {
		return nil, p.errorf("unknown envvar access type")
	}
	envVar.AccessType = accTyp

	nodeName, err := p.parseNodeName()
	if err != nil {
		return nil, err
	}
	envVar.AccessNodes = append(envVar.AccessNodes, nodeName)

	for {
		t = p.scan()
		if !t.isSyntax(syntaxComma) {
			p.unscan()
			break
		}

		nodeName, err := p.parseNodeName()
		if err != nil {
			return nil, err
		}
		envVar.AccessNodes = append(envVar.AccessNodes, nodeName)
	}

	if err := p.expectSyntax(syntaxSemicolon); err != nil {
		return nil, err
	}

	return envVar, nil
}

func (p *Parser) parseEnvVarData() (*EnvVarData, error) {
	evData := new(EnvVarData)

	evName, err := p.parseEnvVarName()
	if err != nil {
		return nil, err
	}
	evData.EnvVarName = evName

	if err := p.expectSyntax(syntaxColon); err != nil {
		return nil, err
	}

	t := p.scan()
	if !t.isNumber() {
		return nil, p.errorf("expected envvar data size")
	}
	dataSize, err := p.parseUint(t.value)
	if err != nil {
		return nil, p.errorf("cannot parse envvar data size as uint")
	}
	evData.DataSize = dataSize

	if err := p.expectSyntax(syntaxSemicolon); err != nil {
		return nil, err
	}

	return evData, nil
}

func (p *Parser) parseSignalType() (*SignalType, *SignalTypeRef, error) {
	sigType := new(SignalType)
	sigTypeRef := new(SignalTypeRef)

	t := p.scan()
	p.unscan()
	switch t.kind {
	case tokenIdent:
		t = p.scan()
		if !t.isIdent() {
			return nil, nil, p.errorf("expected signal type name")
		}
		sigType.TypeName = t.value

		if err := p.expectSyntax(syntaxColon); err != nil {
			return nil, nil, err
		}

		t = p.scan()
		if !t.isNumber() {
			return nil, nil, p.errorf("expected signal start bit")
		}

		if err := p.expectSyntax(syntaxPipe); err != nil {
			return nil, nil, err
		}

		t = p.scan()
		if !t.isNumber() {
			return nil, nil, p.errorf("expected signal size")
		}
		size, err := p.parseUint(t.value)
		if err != nil {
			return nil, nil, p.errorf("cannot parse signal size as uint")
		}
		sigType.Size = size

		if err := p.expectSyntax(syntaxAt); err != nil {
			return nil, nil, err
		}

		t = p.scan()
		if !t.isNumber() {
			return nil, nil, p.errorf("expected signal byte order")
		}
		byteOrder, err := p.parseUint(t.value)
		if err != nil {
			return nil, nil, p.errorf("cannot parse signal byte order as uint")
		}
		if byteOrder == 0 {
			sigType.ByteOrder = SignalBigEndian
		} else if byteOrder == 1 {
			sigType.ByteOrder = SignalLittleEndian
		} else {
			return nil, nil, p.errorf("signal byte order must be 0 or 1")
		}

		t = p.scan()
		syntKind := getSyntaxKind(t.value)
		if t.kind != tokenSyntax || (syntKind != syntaxPlus && syntKind != syntaxMinus) {
			return nil, nil, p.errorf(`expected "+" or "-"`)
		}
		if t.value == "+" {
			sigType.ValueType = SignalUnsigned
		} else {
			sigType.ValueType = SignalSigned
		}

		if err := p.expectSyntax(syntaxLeftParen); err != nil {
			return nil, nil, err
		}

		t = p.scan()
		if !t.isNumber() {
			return nil, nil, p.errorf("expected signal factor")
		}
		factor, err := p.parseDouble(t.value)
		if err != nil {
			return nil, nil, p.errorf("cannot parse signal factor as double")
		}
		sigType.Factor = factor

		if err := p.expectSyntax(syntaxComma); err != nil {
			return nil, nil, err
		}

		t = p.scan()
		if !t.isNumber() {
			return nil, nil, p.errorf("expected signal offset")
		}
		offset, err := p.parseDouble(t.value)
		if err != nil {
			return nil, nil, p.errorf("cannot parse signal offset as double")
		}
		sigType.Offset = offset

		if err := p.expectSyntax(syntaxRightParen); err != nil {
			return nil, nil, err
		}

		if err := p.expectSyntax(syntaxLeftSquareBrace); err != nil {
			return nil, nil, err
		}

		t = p.scan()
		if !t.isNumber() {
			return nil, nil, p.errorf("expected signal minimum")
		}
		min, err := p.parseDouble(t.value)
		if err != nil {
			return nil, nil, p.errorf("cannot parse signal minimum as double")
		}
		sigType.Min = min

		if err := p.expectSyntax(syntaxPipe); err != nil {
			return nil, nil, err
		}

		t = p.scan()
		if !t.isNumber() {
			return nil, nil, p.errorf("expected signal maximum")
		}
		max, err := p.parseDouble(t.value)
		if err != nil {
			return nil, nil, p.errorf("cannot parse signal maximum as double")
		}
		sigType.Max = max

		if err := p.expectSyntax(syntaxRightSquareBrace); err != nil {
			return nil, nil, err
		}

		t = p.scan()
		if !t.isString() {
			return nil, nil, p.errorf("expected signal unit")
		}
		sigType.Unit = t.value

		t = p.scan()
		if !t.isNumber() {
			return nil, nil, p.errorf("expected signal default value")
		}
		defVal, err := p.parseDouble(t.value)
		if err != nil {
			return nil, nil, p.errorf("cannot parse signal default value as double")
		}
		sigType.DefaultValue = defVal

		if err := p.expectSyntax(syntaxComma); err != nil {
			return nil, nil, err
		}

		t = p.scan()
		if !t.isIdent() {
			return nil, nil, p.errorf("expected signal value table name")
		}
		sigType.ValueTableName = t.value

		if err := p.expectSyntax(syntaxSemicolon); err != nil {
			return nil, nil, err
		}

	case tokenNumber:
		msgID, err := p.parseMessageID()
		if err != nil {
			return nil, nil, err
		}
		sigTypeRef.MessageID = msgID

		sigName, err := p.parseSignalName()
		if err != nil {
			return nil, nil, err
		}
		sigTypeRef.SignalName = sigName

		t = p.scan()
		if !t.isIdent() {
			return nil, nil, p.errorf("expected signal type name")
		}
		sigType.TypeName = t.value

		if err := p.expectSyntax(syntaxSemicolon); err != nil {
			return nil, nil, err
		}

		return nil, sigTypeRef, nil

	default:
		return nil, nil, p.errorf("expected signal type name or message id")
	}

	return sigType, nil, nil
}

func (p *Parser) parseComment() (*Comment, error) {
	com := new(Comment)

	t := p.scan()
	switch t.kind {
	case tokenString:
		com.Kind = CommentGeneral
		p.unscan()

	case tokenKeyword:
		keywordKind := getKeywordKind(t.value)
		switch keywordKind {
		case keywordNode:
			com.Kind = CommentNode
			nodeName, err := p.parseNodeName()
			if err != nil {
				return nil, err
			}
			com.NodeName = nodeName

		case keywordMessage:
			com.Kind = CommentMessage
			msgID, err := p.parseMessageID()
			if err != nil {
				return nil, err
			}
			com.MessageID = msgID

		case keywordSignal:
			com.Kind = CommentSignal
			msgID, err := p.parseMessageID()
			if err != nil {
				return nil, err
			}
			com.MessageID = msgID
			sigName, err := p.parseSignalName()
			if err != nil {
				return nil, err
			}
			com.SignalName = sigName

		case keywordEnvVar:
			com.Kind = CommentEnvVar
			envvarName, err := p.parseEnvVarName()
			if err != nil {
				return nil, err
			}
			com.EnvVarName = envvarName

		default:
			return nil, p.errorf("expected node, message, signal or envvar keyword")
		}

	default:
		return nil, p.errorf("expected string or keyword")
	}

	t = p.scan()
	if !t.isString() {
		return nil, p.errorf("expected comment text string")
	}
	com.Text = t.value

	if err := p.expectSyntax(syntaxSemicolon); err != nil {
		return nil, err
	}

	return com, nil
}

func (p *Parser) parseAttributeName() (string, error) {
	t := p.scan()
	if !t.isString() {
		return "", p.errorf("expected attribute name")
	}
	if strings.ContainsRune(t.value, ' ') ||
		strings.ContainsRune(t.value, '\t') ||
		strings.ContainsRune(t.value, '\n') {
		return "", p.errorf("attribute name cannot contain whitespaces")
	}

	return t.value, nil
}

func (p *Parser) parseAttribute() (*Attribute, error) {
	att := new(Attribute)

	t := p.scan()
	switch t.kind {
	case tokenString:
		p.unscan()
		att.Kind = AttributeGeneral

	case tokenKeyword:
		keywordKind := getKeywordKind(t.value)
		switch keywordKind {
		case keywordNode:
			att.Kind = AttributeNode

		case keywordMessage:
			att.Kind = AttributeMessage

		case keywordSignal:
			att.Kind = AttributeSignal

		case keywordEnvVar:
			att.Kind = AttributeEnvVar

		default:
			return nil, p.errorf("expected node, message, signal or envvar keyword")
		}

	default:
		return nil, p.errorf("expected string or keyword")
	}

	attName, err := p.parseAttributeName()
	if err != nil {
		return nil, err
	}
	att.Name = attName

	t = p.scan()
	if t.kind != tokenKeyword {
		return nil, p.errorf("expected attribute type keyword")
	}
	keywordKind := getKeywordKind(t.value)
	switch keywordKind {
	case keywordAttributeInt:
		att.Type = AttributeInt
		t = p.scan()
		if !t.isNumber() {
			return nil, p.errorf("expected int attribute min value")
		}
		minInt, err := p.parseInt(t.value)
		if err != nil {
			return nil, p.errorf("cannot parse int attribute min value as int")
		}
		att.MinInt = minInt
		t = p.scan()
		if !t.isNumber() {
			return nil, p.errorf("expected int attribute max value")
		}
		maxInt, err := p.parseInt(t.value)
		if err != nil {
			return nil, p.errorf("cannot parse int attribute max value as int")
		}
		att.MaxInt = maxInt

	case keywordAttributeHex:
		att.Type = AttributeHex
		t = p.scan()
		if !t.isNumber() {
			return nil, p.errorf("expected hex attribute min value")
		}
		minHex, err := p.parseHexInt(t.value)
		if err != nil {
			return nil, p.errorf("cannot parse hex attribute min value as int")
		}
		att.MinHex = minHex
		t = p.scan()
		if !t.isNumber() {
			return nil, p.errorf("expected hex attribute max value")
		}
		maxHex, err := p.parseHexInt(t.value)
		if err != nil {
			return nil, p.errorf("cannot parse hex attribute max value as int")
		}
		att.MaxHex = maxHex

	case keywordAttributeFloat:
		att.Type = AttributeFloat
		t = p.scan()
		if !t.isNumber() {
			return nil, p.errorf("expected float attribute min value")
		}
		minFloat, err := p.parseDouble(t.value)
		if err != nil {
			return nil, p.errorf("cannot parse float attribute min value as double")
		}
		att.MinFloat = minFloat
		t = p.scan()
		if !t.isNumber() {
			return nil, p.errorf("expected float attribute max value")
		}
		maxFloat, err := p.parseDouble(t.value)
		if err != nil {
			return nil, p.errorf("cannot parse float attribute max value as double")
		}
		att.MaxFloat = maxFloat

	case keywordAttributeString:
		att.Type = AttributeString

	case keywordAttributeEnum:
		att.Type = AttributeEnum
		t = p.scan()
		if !t.isString() {
			return nil, p.errorf("expected enum attribute values")
		}
		att.EnumValues = append(att.EnumValues, t.value)
		for {
			if !p.scan().isSyntax(syntaxComma) {
				p.unscan()
				break
			}
			t = p.scan()
			if !t.isString() {
				return nil, p.errorf("expected enum attribute values")
			}
			att.EnumValues = append(att.EnumValues, t.value)
		}

	default:
		return nil, p.errorf("expected attribute type keyword to be INT, HEX, FLOAT, STRING or ENUM")
	}

	if err := p.expectSyntax(syntaxSemicolon); err != nil {
		return nil, err
	}

	return att, nil
}

func (p *Parser) parseAttributeDefault() (*AttributeDefault, error) {
	attDef := new(AttributeDefault)

	attName, err := p.parseAttributeName()
	if err != nil {
		return nil, err
	}
	attDef.AttributeName = attName

	t := p.scan()
	if t.isString() {
		attDef.ValueString = t.value
		attDef.Type = AttributeDefaultString

	} else if t.isNumber() {
		if strings.HasPrefix(t.value, "0x") || strings.HasPrefix(t.value, "0X") {
			hexVal, err := p.parseHexInt(t.value)
			if err != nil {
				return nil, p.errorf("cannot parse hex attribute default value as int")
			}
			attDef.ValueHex = hexVal
			attDef.Type = AttributeDefaultHex

		} else if strings.Contains(t.value, ".") {
			floatVal, err := p.parseDouble(t.value)
			if err != nil {
				return nil, p.errorf("cannot parse float attribute default value as double")
			}
			attDef.ValueFloat = floatVal
			attDef.Type = AttributeDefaultFloat

		} else {
			invVal, err := p.parseInt(t.value)
			if err != nil {
				return nil, p.errorf("cannot parse int attribute default value as int")
			}
			attDef.ValueInt = invVal
			attDef.Type = AttributeDefaultInt
		}

	} else {
		return nil, p.errorf("expected attribute default value")
	}

	if err := p.expectSyntax(syntaxSemicolon); err != nil {
		return nil, err
	}

	return attDef, nil
}

func (p *Parser) parseAttributeValue() (*AttributeValue, error) {
	attVal := new(AttributeValue)

	t := p.scan()
	if !t.isString() {
		return nil, p.errorf("expected attribute value")
	}
	attVal.AttributeName = t.value

	switch t := p.scan(); {
	case t.isString() || t.isNumber():
		attVal.AttributeKind = AttributeGeneral
		p.unscan()

	case t.kind == tokenKeyword:
		keywordKind := getKeywordKind(t.value)
		switch keywordKind {
		case keywordNode:
			attVal.AttributeKind = AttributeNode
			nodeName, err := p.parseNodeName()
			if err != nil {
				return nil, err
			}
			attVal.NodeName = nodeName

		case keywordMessage:
			attVal.AttributeKind = AttributeMessage
			msgID, err := p.parseMessageID()
			if err != nil {
				return nil, err
			}
			attVal.MessageID = msgID

		case keywordSignal:
			attVal.AttributeKind = AttributeSignal
			msgID, err := p.parseMessageID()
			if err != nil {
				return nil, err
			}
			attVal.MessageID = msgID
			sigName, err := p.parseSignalName()
			if err != nil {
				return nil, err
			}
			attVal.SignalName = sigName

		case keywordEnvVar:
			attVal.AttributeKind = AttributeEnvVar
			envvarName, err := p.parseEnvVarName()
			if err != nil {
				return nil, err
			}
			attVal.EnvVarName = envvarName

		default:
			return nil, p.errorf("expected node, message, signal or envvar keyword")
		}

	default:
		return nil, p.errorf("expected string, number or keyword")
	}

	t = p.scan()
	if t.isString() {
		attVal.ValueString = t.value
		attVal.Type = AttributeValueString

	} else if t.isNumber() {
		if strings.HasPrefix(t.value, "0x") || strings.HasPrefix(t.value, "0X") {
			hexVal, err := p.parseHexInt(t.value)
			if err != nil {
				return nil, p.errorf("cannot parse hex attribute value as int")
			}
			attVal.ValueHex = hexVal
			attVal.Type = AttributeValueHex

		} else if strings.Contains(t.value, ".") {
			floatVal, err := p.parseDouble(t.value)
			if err != nil {
				return nil, p.errorf("cannot parse float attribute value as double")
			}
			attVal.ValueFloat = floatVal
			attVal.Type = AttributeValueFloat

		} else {
			invVal, err := p.parseInt(t.value)
			if err != nil {
				return nil, p.errorf("cannot parse int attribute value as int")
			}
			attVal.ValueInt = invVal
			attVal.Type = AttributeValueInt
		}

	} else {
		return nil, p.errorf("expected attribute value")
	}

	if err := p.expectSyntax(syntaxSemicolon); err != nil {
		return nil, err
	}

	return attVal, nil
}

func (p *Parser) parseValueEncoding() (*ValueEncoding, error) {
	valEnc := new(ValueEncoding)

	t := p.scan()
	p.unscan()
	switch t.kind {
	case tokenIdent:
		valEnc.Kind = ValueEncodingEnvVar
		envVarName, err := p.parseEnvVarName()
		if err != nil {
			return nil, err
		}
		valEnc.EnvVarName = envVarName

	case tokenNumber:
		valEnc.Kind = ValueEncodingSignal
		msgID, err := p.parseMessageID()
		if err != nil {
			return nil, err
		}
		valEnc.MessageID = msgID
		sigName, err := p.parseSignalName()
		if err != nil {
			return nil, err
		}
		valEnc.SignalName = sigName

	default:
		return nil, p.errorf("expected value encoding message id or envvar name")
	}

	for {
		t := p.scan()
		p.unscan()
		if !t.isNumber() {
			break
		}

		vd, err := p.parseValueDescription()
		if err != nil {
			return nil, err
		}

		valEnc.Values = append(valEnc.Values, vd)
	}

	if err := p.expectSyntax(syntaxSemicolon); err != nil {
		return nil, err
	}

	return valEnc, nil
}

func (p *Parser) parseSignalGroup() (*SignalGroup, error) {
	sigGroup := new(SignalGroup)

	msgID, err := p.parseMessageID()
	if err != nil {
		return nil, err
	}
	sigGroup.MessageID = msgID

	t := p.scan()
	if !t.isIdent() {
		return nil, p.errorf("expected signal group name")
	}
	sigGroup.GroupName = t.value

	t = p.scan()
	if !t.isNumber() {
		return nil, p.errorf("expected signal group repetitions")
	}
	r, err := p.parseUint(t.value)
	if err != nil {
		return nil, p.errorf("cannot parse signal group repetitions as uint")
	}
	sigGroup.Repetitions = r

	if err := p.expectSyntax(syntaxColon); err != nil {
		return nil, err
	}

	for {
		t = p.scan()
		if !t.isIdent() {
			break
		}
		sigGroup.SignalNames = append(sigGroup.SignalNames, t.value)
	}

	if err := p.expectSyntax(syntaxSemicolon); err != nil {
		return nil, err
	}

	return sigGroup, nil
}

func (p *Parser) parseSignalExtValueType() (*SignalExtValueType, error) {
	valType := new(SignalExtValueType)

	msgID, err := p.parseMessageID()
	if err != nil {
		return nil, err
	}
	valType.MessageID = msgID

	sigName, err := p.parseSignalName()
	if err != nil {
		return nil, err
	}
	valType.SignalName = sigName

	t := p.scan()
	if !t.isNumber() {
		return nil, p.errorf("expected signal extended value type")
	}
	vt, err := p.parseUint(t.value)
	if err != nil {
		return nil, p.errorf("cannot parse signal extended value type as uint")
	}

	if vt == 0 {
		valType.ExtValueType = SignalExtValueTypeInteger
	} else if vt == 1 {
		valType.ExtValueType = SignalExtValueTypeFloat
	} else if vt == 2 {
		valType.ExtValueType = SignalExtValueTypeDouble
	} else {
		return nil, p.errorf("signal extended value type must be 0, 1 or 2")
	}

	if err := p.expectSyntax(syntaxSemicolon); err != nil {
		return nil, err
	}

	return valType, nil
}

func (p *Parser) parseExtendedMuxRange() (*ExtendedMuxRange, error) {
	extMuxR := new(ExtendedMuxRange)

	t := p.scan()
	if !t.isNumberRange() {
		return nil, p.errorf("expected extended mux range")
	}
	tmpRange := strings.Split(t.value, "-")
	from, err := p.parseUint(tmpRange[0])
	if err != nil {
		return nil, p.errorf("cannot parse extended mux range as uint")
	}
	extMuxR.From = from
	to, err := p.parseUint(tmpRange[1])
	if err != nil {
		return nil, p.errorf("cannot parse extended mux range as uint")
	}
	extMuxR.To = to

	return extMuxR, nil
}

func (p *Parser) parseExtendedMux() (*ExtendedMux, error) {
	extMux := new(ExtendedMux)

	msgID, err := p.parseMessageID()
	if err != nil {
		return nil, err
	}
	extMux.MessageID = msgID

	t := p.scan()
	if !t.isIdent() {
		return nil, p.errorf("expected extended mux multiplexed signal name")
	}
	extMux.MultiplexedName = t.value

	t = p.scan()
	if !t.isIdent() {
		return nil, p.errorf("expected extended mux multiplexor signal name")
	}
	extMux.MultiplexorName = t.value

	r, err := p.parseExtendedMuxRange()
	if err != nil {
		return nil, err
	}
	extMux.Ranges = append(extMux.Ranges, r)

	for {
		t = p.scan()
		if !t.isSyntax(syntaxComma) {
			p.unscan()
			break
		}

		r, err := p.parseExtendedMuxRange()
		if err != nil {
			return nil, err
		}
		extMux.Ranges = append(extMux.Ranges, r)
	}

	if err := p.expectSyntax(syntaxSemicolon); err != nil {
		return nil, err
	}

	return extMux, nil
}
