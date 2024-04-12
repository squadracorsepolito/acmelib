package main

import (
	"fmt"
	"log"

	"github.com/FerroO2000/acmelib"
)

// !TODO: guardie sulle enum

const dm1MsgName = "DM1"

var uint16bit, _ = acmelib.NewIntegerSignalType("16_bit", 16, false)
var mVUnit = acmelib.NewSignalUnit("milli_volt", acmelib.SignalUnitKindVoltage, "mV")
var flagType = acmelib.NewFlagSignalType("flag")
var uint8bit, _ = acmelib.NewIntegerSignalType("uint_8", 8, false)

func main() {
	bmslv := acmelib.NewNode("BMS_LV", 0)

	battPackGen := acmelib.NewMessage("BMSLV_BatteryPackGeneral", 6)

	currSens, _ := acmelib.NewStandardSignal("Current_Sensor_mV", uint16bit, 0, uint16bit.Max(), 0, 0.076, mVUnit)
	totVolt, _ := acmelib.NewStandardSignal("Total_voltage", uint16bit, 0, uint16bit.Max(), 0, 0.076, mVUnit)

	if err := battPackGen.AppendSignal(currSens); err != nil {
		panic(err)
	}
	if err := battPackGen.AppendSignal(totVolt); err != nil {
		panic(err)
	}

	msgStatus := acmelib.NewMessage("BMSLV_Status", 6)

	isRelOpen, _ := acmelib.NewStandardSignal("is_relary_open", flagType, flagType.Min(), flagType.Max(), 0, 1, nil)

	if err := msgStatus.AppendSignal(isRelOpen); err != nil {
		panic(err)
	}

	bmslv.AddMessage(battPackGen)
	bmslv.AddMessage(msgStatus)

	sinEnum, _ := genSin(bmslv.Messages())

	dm1Msg := acmelib.NewMessage(dm1MsgName, 8)
	sinSig, _ := acmelib.NewEnumSignal("sin", sinEnum)
	fmiSig, _ := acmelib.NewEnumSignal("fmi", initFmi())
	occSig, _ := acmelib.NewStandardSignal("occ_counter", uint8bit, uint8bit.Min(), uint8bit.Max(), 0, 1, nil)

	if err := dm1Msg.InsertSignal(sinSig, 0); err != nil {
		panic(err)
	}
	if err := dm1Msg.InsertSignal(fmiSig, 8); err != nil {
		panic(err)
	}
	if err := dm1Msg.InsertSignal(occSig, 16); err != nil {
		panic(err)
	}

	bmslv.AddMessage(dm1Msg)

	log.Print(bmslv.String())
}

func genSin(messages []*acmelib.Message) (*acmelib.SignalEnum, error) {
	sinEnum := acmelib.NewSignalEnum("sin")

	sigNames := make(map[string]int)
	idx := 0
	for _, msg := range messages {
		for _, sig := range msg.Signals() {
			if sig.Name() == dm1MsgName {
				continue
			}

			if _, ok := sigNames[sig.Name()]; !ok {
				sigNames[sig.Name()] = idx
				idx++
			}
		}
	}

	for enumValname, enumValIdx := range sigNames {
		sigEnumVal := acmelib.NewSignalEnumValue(enumValname, enumValIdx)
		if err := sinEnum.AddValue(sigEnumVal); err != nil {
			panic(err)
		}
	}

	if sinEnum.GetSize() > 8 {
		return nil, fmt.Errorf("too many signals")
	}

	return sinEnum, nil
}

func initFmi() *acmelib.SignalEnum {
	fmi := acmelib.NewSignalEnum("fmi")
	fmi.AddValue(acmelib.NewSignalEnumValue("high_severity", 0))
	fmi.AddValue(acmelib.NewSignalEnumValue("condition_exists", 31))
	return fmi
}
