package dbc

const (
	// MsgCycleTimeName is the name of the well known attribute for message cycle time.
	MsgCycleTimeName = "GenMsgCycleTime"
	// MsgCycleTimeMin is the min value of the well known attribute for message cycle time.
	MsgCycleTimeMin = 0
	// MsgCycleTimeMax is the max value of the well known attribute for message cycle time.
	MsgCycleTimeMax = 3600000
)

const (
	// MsgDelayTimeName is the name of the well known attribute for message delay time.
	MsgDelayTimeName = "GenMsgDelayTime"
	// MsgDelayTimeMin is the min value of the well known attribute for message delay time.
	MsgDelayTimeMin = 0
	// MsgDelayTimeMax is the max value of the well known attribute for message delay time.
	MsgDelayTimeMax = 1000
)

const (
	// MsgStartDelayTimeName is the name of the well known attribute for message start delay time.
	MsgStartDelayTimeName = "GenMsgStartDelayTime"
	// MsgStartDelayTimeMin is the min value of the well known attribute for message start delay time.
	MsgStartDelayTimeMin = 0
	// MsgStartDelayTimeMax is the max value of the well known attribute for message start delay time.
	MsgStartDelayTimeMax = 100000
)

var (
	// MsgSendTypeName is the name of the well known attribute for message send type.
	MsgSendTypeName = "GenMsgSendType"
	// MsgSendTypeValues are the value of the well known attribute for message send type.
	MsgSendTypeValues = []string{
		"NoMsgSendType",
		"Cyclic",
		"CyclicIfActive",
		"CyclicAndTriggered",
		"CyclicIfActiveAndTriggered",
	}
)

const (
	// SigStartValueName is the name of the well known attribute for signal start value.
	SigStartValueName = "GenSigStartValue"
	// SigStartValueMin is the min value of the well known attribute for signal start value.
	SigStartValueMin = 0
	// SigStartValueMax is the max value of the well known attribute for signal start value.
	SigStartValueMax = 10000
)

var (
	// SigSendTypeName is the name of the well known attribute for signal send type.
	SigSendTypeName = "GenSigSendType"
	// SigSendTypeValues are the value of the well known attribute for signal send type.
	SigSendTypeValues = []string{
		"NoSigSendType",
		"Cyclic",
		"OnWrite",
		"OnWriteWithRepetition",
		"OnChange",
		"OnChangeWithRepetition",
		"IfActive",
		"IfActiveWithRepetition",
	}
)
