package acmelib

type specialAttributeType int

const (
	specialAttributeMsgCycleTime specialAttributeType = iota
	specialAttributeMsgDelayTime
	specialAttributeMsgStartDelayTime
	specialAttributeMsgSendType

	specialAttributeSigSendType
)

var specialAttributeNames = map[specialAttributeType]string{
	specialAttributeMsgCycleTime:      "GenMsgCycleTime",
	specialAttributeMsgDelayTime:      "GenMsgDelayTime",
	specialAttributeMsgStartDelayTime: "GenMsgStartDelayTime",
	specialAttributeMsgSendType:       "GenMsgSendType",

	specialAttributeSigSendType: "GenSigSendType",
}

var specialAttributeTypes = map[string]specialAttributeType{
	"GenMsgCycleTime":      specialAttributeMsgCycleTime,
	"GenMsgDelayTime":      specialAttributeMsgDelayTime,
	"GenMsgStartDelayTime": specialAttributeMsgStartDelayTime,
	"GenMsgSendType":       specialAttributeMsgSendType,

	"GenSigSendType": specialAttributeSigSendType,
}

var (
	msgCycleTimeAtt, _      = NewIntegerAttribute(specialAttributeNames[specialAttributeMsgCycleTime], 0, 0, 1000)
	msgDelayTimeAtt, _      = NewIntegerAttribute(specialAttributeNames[specialAttributeMsgDelayTime], 0, 0, 1000)
	msgStartDelayTimeAtt, _ = NewIntegerAttribute(specialAttributeNames[specialAttributeMsgStartDelayTime], 0, 0, 1000)
)

var msgSendTypeValues = []string{
	"NoMsgSendType",
	"Cyclic",
	"CyclicIfActive",
	"CyclicAndTriggered",
	"CyclicIfActiveAndTriggered",
}

func messageSendTypeFromString(str string) MessageSendType {
	switch str {
	case msgSendTypeValues[1]:
		return MessageSendTypeCyclic
	case msgSendTypeValues[2]:
		return MessageSendTypeCyclicIfActive
	case msgSendTypeValues[3]:
		return MessageSendTypeCyclicAndTriggered
	case msgSendTypeValues[4]:
		return MessageSendTypeCyclicIfActiveAndTriggered
	}
	return MessageSendTypeUnset
}

var msgSendTypeAtt, _ = NewEnumAttribute(specialAttributeNames[specialAttributeMsgSendType], msgSendTypeValues...)

var sigSendTypeValues = []string{
	"NoSigSendType",
	"Cyclic",
	"OnWrite",
	"OnWriteWithRepetition",
	"OnChange",
	"OnChangeWithRepetition",
	"IfActive",
	"IfActiveWithRepetition",
}

func signalSendTypeFromString(str string) SignalSendType {
	switch str {
	case sigSendTypeValues[1]:
		return SignalSendTypeCyclic
	case sigSendTypeValues[2]:
		return SignalSendTypeOnWrite
	case sigSendTypeValues[3]:
		return SignalSendTypeOnWriteWithRepetition
	case sigSendTypeValues[4]:
		return SignalSendTypeOnChange
	case sigSendTypeValues[5]:
		return SignalSendTypeOnChangeWithRepetition
	case sigSendTypeValues[6]:
		return SignalSendTypeIfActive
	case sigSendTypeValues[6]:
		return SignalSendTypeIfActiveWithRepetition
	}
	return SignalSendTypeUnset
}

var sigSendTypeAtt, _ = NewEnumAttribute(specialAttributeNames[specialAttributeSigSendType], sigSendTypeValues...)
