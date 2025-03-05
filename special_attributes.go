package acmelib

import "github.com/squadracorsepolito/acmelib/dbc"

type specialAttributeType int

const (
	specialAttributeMsgCycleTime specialAttributeType = iota
	specialAttributeMsgDelayTime
	specialAttributeMsgStartDelayTime
	specialAttributeMsgSendType

	specialAttributeSigStartValue
	specialAttributeSigSendType
)

var specialAttributeNames = map[specialAttributeType]string{
	specialAttributeMsgCycleTime:      dbc.MsgCycleTimeName,
	specialAttributeMsgDelayTime:      dbc.MsgDelayTimeName,
	specialAttributeMsgStartDelayTime: dbc.MsgStartDelayTimeName,
	specialAttributeMsgSendType:       dbc.MsgSendTypeName,

	specialAttributeSigStartValue: dbc.SigStartValueName,
	specialAttributeSigSendType:   dbc.SigSendTypeName,
}

var specialAttributeTypes = map[string]specialAttributeType{
	dbc.MsgCycleTimeName:      specialAttributeMsgCycleTime,
	dbc.MsgDelayTimeName:      specialAttributeMsgDelayTime,
	dbc.MsgStartDelayTimeName: specialAttributeMsgStartDelayTime,
	dbc.MsgSendTypeName:       specialAttributeMsgSendType,

	dbc.SigStartValueName: specialAttributeSigStartValue,
	dbc.SigSendTypeName:   specialAttributeSigSendType,
}

var (
	msgCycleTimeAtt, _      = NewIntegerAttribute(dbc.MsgCycleTimeName, dbc.MsgCycleTimeMin, dbc.MsgCycleTimeMin, dbc.MsgCycleTimeMax)
	msgDelayTimeAtt, _      = NewIntegerAttribute(dbc.MsgDelayTimeName, dbc.MsgDelayTimeMin, dbc.MsgDelayTimeMin, dbc.MsgDelayTimeMax)
	msgStartDelayTimeAtt, _ = NewIntegerAttribute(dbc.MsgStartDelayTimeName, dbc.MsgStartDelayTimeMin, dbc.MsgStartDelayTimeMin, dbc.MsgStartDelayTimeMax)
	msgSendTypeAtt, _       = NewEnumAttribute(dbc.MsgSendTypeName, dbc.MsgSendTypeValues...)

	sigStartValueAtt, _ = NewFloatAttribute(dbc.SigStartValueName, dbc.SigStartValueMin, dbc.SigStartValueMin, dbc.SigStartValueMax)
	sigSendTypeAtt, _   = NewEnumAttribute(dbc.SigSendTypeName, dbc.SigSendTypeValues...)
)

func messageSendTypeFromDBC(str string) MessageSendType {
	switch str {
	case dbc.MsgSendTypeValues[1]:
		return MessageSendTypeCyclic
	case dbc.MsgSendTypeValues[2]:
		return MessageSendTypeCyclicIfActive
	case dbc.MsgSendTypeValues[3]:
		return MessageSendTypeCyclicAndTriggered
	case dbc.MsgSendTypeValues[4]:
		return MessageSendTypeCyclicIfActiveAndTriggered
	default:
		return MessageSendTypeUnset
	}
}

func messageSendTypeToDBC(sendType MessageSendType) string {
	switch sendType {
	case MessageSendTypeCyclic:
		return dbc.MsgSendTypeValues[1]
	case MessageSendTypeCyclicIfActive:
		return dbc.MsgSendTypeValues[2]
	case MessageSendTypeCyclicAndTriggered:
		return dbc.MsgSendTypeValues[3]
	case MessageSendTypeCyclicIfActiveAndTriggered:
		return dbc.MsgSendTypeValues[4]
	default:
		return dbc.MsgSendTypeValues[0]
	}
}

func signalSendTypeFromDBC(str string) SignalSendType {
	switch str {
	case dbc.SigSendTypeValues[1]:
		return SignalSendTypeCyclic
	case dbc.SigSendTypeValues[2]:
		return SignalSendTypeOnWrite
	case dbc.SigSendTypeValues[3]:
		return SignalSendTypeOnWriteWithRepetition
	case dbc.SigSendTypeValues[4]:
		return SignalSendTypeOnChange
	case dbc.SigSendTypeValues[5]:
		return SignalSendTypeOnChangeWithRepetition
	case dbc.SigSendTypeValues[6]:
		return SignalSendTypeIfActive
	case dbc.SigSendTypeValues[7]:
		return SignalSendTypeIfActiveWithRepetition
	default:
		return SignalSendTypeUnset
	}
}

func signalSendTypeToDBC(sendType SignalSendType) string {
	switch sendType {
	case SignalSendTypeCyclic:
		return dbc.SigSendTypeValues[1]
	case SignalSendTypeOnWrite:
		return dbc.SigSendTypeValues[2]
	case SignalSendTypeOnWriteWithRepetition:
		return dbc.SigSendTypeValues[3]
	case SignalSendTypeOnChange:
		return dbc.SigSendTypeValues[4]
	case SignalSendTypeOnChangeWithRepetition:
		return dbc.SigSendTypeValues[5]
	case SignalSendTypeIfActive:
		return dbc.SigSendTypeValues[6]
	case SignalSendTypeIfActiveWithRepetition:
		return dbc.SigSendTypeValues[7]
	default:
		return dbc.SigSendTypeValues[0]
	}
}
