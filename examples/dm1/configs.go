package main

import "github.com/squadracorsepolito/acmelib"

var flagSigType = acmelib.NewFlagSignalType("flag")
var float16SigType, _ = acmelib.NewFloatSignalType("uint16", 16)

var mVSigUnit = acmelib.NewSignalUnit("milli_volt", acmelib.SignalUnitKindElectrical, "mV")
