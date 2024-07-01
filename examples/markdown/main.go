package main

import (
	"fmt"
	"log"
	"os"

	"github.com/squadracorsepolito/acmelib"
)

func main() {
	sc24 := acmelib.NewNetwork("SC24")
	sc24.SetDesc("The CAN network of the squadracorse 2024 formula SAE car.")

	mcbFile, err := os.Open("MCB.dbc")
	checkErr(err)
	defer mcbFile.Close()

	mcb, err := acmelib.ImportDBCFile("mcb", mcbFile)
	checkErr(err)

	if err := mcb.UpdateName("Main CAN Bus"); err != nil {
		panic(err)
	}

	if err := sc24.AddBus(mcb); err != nil {
		panic(err)
	}

	// hvcbFile, err := os.Open("HVCB.dbc")
	// if err != nil {
	// 	panic(err)
	// }
	// defer hvcbFile.Close()

	// hvcb, err := acmelib.ImportDBCFile("hvcb", hvcbFile)
	// if err != nil {
	// 	panic(err)
	// }

	// if err := sc24.AddBus(hvcb); err != nil {
	// 	panic(err)
	// }

	// renaming signal types
	dashInt, err := mcb.GetNodeInterfaceByNodeName("DASH")
	checkErr(err)

	tmpSigType, err := acmelib.NewIntegerSignalType("fan_pwm_t", 4, false)
	checkErr(err)
	tmpSigType.SetMax(10)
	modifySignalType(dashInt, "DASH_peripheralsStatus", "TSAC_FAN_pwmStatus", tmpSigType)

	modifySignalTypeName(dashInt, "DASH_rotarySwitchState", "ROT_SWITCH_0_position", "rotary_switch_pos_t")
	modifySignalTypeName(dashInt, "DASH_appsRangeLimits", "APPS_0_voltageRangeMin", "uint16_t")
	modifySignalTypeName(dashInt, "DASH_lvRelayOverride", "BMS_LV_diagPWD", "bms_lv_password_t")

	bmslvInt, err := mcb.GetNodeInterfaceByNodeName("BMS_LV")
	checkErr(err)

	modifySignalTypeName(bmslvInt, "BMS_LV_hello", "FW_majorVersion", "uint8_t")
	modifySignalTypeName(bmslvInt, "BMS_LV_lvCellVoltage0", "LV_CELL_0_voltage", "lv_cell_voltage_t")
	modifySignalTypeName(bmslvInt, "BMS_LV_lvCellNTCResistance0", "LV_CELL_NTC_00_resistance", "ntc_resistance_t")
	modifySignalTypeName(bmslvInt, "BMS_LV_lvBatGeneral", "LV_BAT_voltage", "lv_bat_voltage_t")
	modifySignalTypeName(bmslvInt, "BMS_LV_lvBatGeneral", "LV_BAT_currentSensVoltage", "lv_bat_current_sens_t")

	dspaceInt, err := mcb.GetNodeInterfaceByNodeName("DSPACE")
	checkErr(err)

	tmpSigType, err = acmelib.NewIntegerSignalType("seconds_t", 6, false)
	checkErr(err)
	tmpSigType.SetMax(59)
	modifySignalType(dspaceInt, "DSPACE_datetime", "DATETIME_seconds", tmpSigType)

	modifySignalTypeName(dspaceInt, "DSPACE_datetime", "DATETIME_month", "month_t")
	modifySignalTypeName(dspaceInt, "DSPACE_datetime", "DATETIME_day", "day_t")
	modifySignalTypeName(dspaceInt, "DSPACE_datetime", "DATETIME_hours", "hours_t")
	modifySignalTypeName(dspaceInt, "DSPACE_datetime", "DATETIME_minutes", "minutes_t")
	modifySignalTypeName(dspaceInt, "DSPACE_status", "DSPACE_FSM_state", "rtd_fsm_t")

	extraNode := acmelib.NewNode("EXTRA_NODE", 8, 1)
	unknownIRMsg := acmelib.NewMessage("unknown_ir", 0x70, 8)
	extraNodeInt := extraNode.Interfaces()[0]
	checkErr(mcb.AddNodeInterface(extraNodeInt))
	checkErr(extraNodeInt.AddMessage(unknownIRMsg))

	scannerInt, err := mcb.GetNodeInterfaceByNodeName("SCANNER")
	checkErr(err)

	// adding tpms
	tpms(mcb, scannerInt, dspaceInt)

	dbcFile, err := os.Create("mcb_parsed.dbc")
	checkErr(err)
	defer dbcFile.Close()
	acmelib.ExportBus(dbcFile, mcb)

	// adding xpc tx/rx
	diagTool := mcb.NodeInterfaces()[0]
	xcpRXMsgID := acmelib.MessageID(10)
	xcpTXCANID := acmelib.CANID(100)
	for _, nodeInt := range mcb.NodeInterfaces() {
		if nodeInt.Node().ID() == 0 {
			continue
		}

		if nodeInt.Node().ID() == 8 {
			break
		}

		nodeName := nodeInt.Node().Name()

		msgRXName := fmt.Sprintf("%s_xcpRX", nodeName)
		tmpRXMsg := acmelib.NewMessage(msgRXName, xcpRXMsgID, 8)
		checkErr(nodeInt.AddMessage(tmpRXMsg))
		tmpRXMsg.AddReceiver(diagTool)
		tmpRXMsg.SetDesc("The message used to notify the diagnostic tool that the board is flashed.")

		msgTXName := fmt.Sprintf("%s_xcpTX", nodeName)
		tmpTXMsg := acmelib.NewMessage(msgTXName, 0, 8)
		tmpTXMsg.SetStaticCANID(xcpTXCANID)
		checkErr(diagTool.AddMessage(tmpTXMsg))
		tmpTXMsg.AddReceiver(nodeInt)
		tmpTXMsg.SetDesc(fmt.Sprintf("The message used to flash the %s.", nodeName))

		xcpTXCANID++
	}

	// calculte bus load
	mcb.SetBaudrate(1_000_000)
	busLoad, err := acmelib.CalculateBusLoad(mcb, 1000)
	checkErr(err)
	log.Print("BUS LOAD: ", busLoad)

	mdFile, err := os.Create("SC24.md")
	checkErr(err)
	defer mdFile.Close()

	if err := acmelib.ExportToMarkdown(sc24, mdFile); err != nil {
		panic(err)
	}

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func modifySignalTypeName(nodeInt *acmelib.NodeInterface, msgName, sigName, newName string) {
	tmpMsg, err := nodeInt.GetMessageByName(msgName)
	checkErr(err)
	tmpSig, err := tmpMsg.GetSignalByName(sigName)
	checkErr(err)
	tmpStdSig, err := tmpSig.ToStandard()
	checkErr(err)
	tmpStdSig.Type().SetName(newName)
}

func modifySignalType(nodeInt *acmelib.NodeInterface, msgName, sigName string, newType *acmelib.SignalType) {
	tmpMsg, err := nodeInt.GetMessageByName(msgName)
	checkErr(err)
	tmpSig, err := tmpMsg.GetSignalByName(sigName)
	checkErr(err)
	tmpStdSig, err := tmpSig.ToStandard()
	checkErr(err)
	tmpStdSig.SetType(newType)
}

func tpms(mcb *acmelib.Bus, scanner, dspace *acmelib.NodeInterface) *acmelib.Node {
	tpms := acmelib.NewNode("TPMS", 9, 1)
	tpms.SetDesc("The tire pressure monitoring system.")
	tpmsInt := tpms.Interfaces()[0]
	checkErr(mcb.AddNodeInterface(tpmsInt))

	idSigType, err := acmelib.NewDecimalSignalType("tire_sens_id_t", 8, false)
	checkErr(err)

	statusSigType, err := acmelib.NewDecimalSignalType("tire_sens_status_t", 8, false)
	checkErr(err)
	statusSigType.SetDesc("Bit #2: 0 if battery voltage > 2.2V, otherwise 1; Bit #3: 0 if wheel spinning, 1 otherwise")

	tempSigType, err := acmelib.NewIntegerSignalType("tire_temp_t", 8, false)
	checkErr(err)
	tempSigType.SetMin(0x0a)
	tempSigType.SetMax(0xaa)

	tempUnit := acmelib.NewSignalUnit("temp_celsius", acmelib.SignalUnitKindTemperature, "degC")

	pressSigType, err := acmelib.NewDecimalSignalType("tire_press_t", 8, false)
	checkErr(err)
	pressSigType.SetMin(0x01)
	pressSigType.SetMax(0xfe)

	pressUnit := acmelib.NewSignalUnit("press_milli_bar", acmelib.SignalUnitKindCustom, "mB")

	frontMsg := acmelib.NewMessage("TPMS_front", 0x718, 8)
	frontMsg.SetStaticCANID(0x718)
	tmpSig, err := acmelib.NewStandardSignal("TIRE_FL_sensID", idSigType)
	checkErr(err)
	tmpSig.SetDesc("The sensor id of the front left tire.")
	checkErr(frontMsg.AppendSignal(tmpSig))
	tmpSig, err = acmelib.NewStandardSignal("TIRE_FL_status", statusSigType)
	checkErr(err)
	tmpSig.SetDesc("The status the front left tire.")
	checkErr(frontMsg.AppendSignal(tmpSig))
	tmpSig, err = acmelib.NewStandardSignal("TIRE_FL_temperature", tempSigType)
	checkErr(err)
	tmpSig.SetUnit(tempUnit)
	tmpSig.SetDesc("The temperature of the front left tire.")
	checkErr(frontMsg.AppendSignal(tmpSig))
	tmpSig, err = acmelib.NewStandardSignal("TIRE_FL_pressure", pressSigType)
	checkErr(err)
	tmpSig.SetUnit(pressUnit)
	tmpSig.SetDesc("The pressure of the front left tire.")
	checkErr(frontMsg.AppendSignal(tmpSig))
	tmpSig, err = acmelib.NewStandardSignal("TIRE_FR_sensID", idSigType)
	checkErr(err)
	tmpSig.SetDesc("The sensor id of the front right tire.")
	checkErr(frontMsg.AppendSignal(tmpSig))
	tmpSig, err = acmelib.NewStandardSignal("TIRE_FR_status", statusSigType)
	checkErr(err)
	tmpSig.SetDesc("The status the front right tire.")
	checkErr(frontMsg.AppendSignal(tmpSig))
	tmpSig, err = acmelib.NewStandardSignal("TIRE_FR_temperature", tempSigType)
	checkErr(err)
	tmpSig.SetUnit(tempUnit)
	tmpSig.SetDesc("The temperature of the front right tire.")
	checkErr(frontMsg.AppendSignal(tmpSig))
	tmpSig, err = acmelib.NewStandardSignal("TIRE_FR_pressure", pressSigType)
	checkErr(err)
	tmpSig.SetUnit(pressUnit)
	tmpSig.SetDesc("The pressure of the front right tire.")
	checkErr(frontMsg.AppendSignal(tmpSig))

	rearMsg := acmelib.NewMessage("TPMS_rear", 0x728, 8)
	rearMsg.SetStaticCANID(0x728)
	tmpSig, err = acmelib.NewStandardSignal("TIRE_RL_sensID", idSigType)
	checkErr(err)
	tmpSig.SetDesc("The sensor id of the rear left tire.")
	checkErr(rearMsg.AppendSignal(tmpSig))
	tmpSig, err = acmelib.NewStandardSignal("TIRE_RL_status", statusSigType)
	checkErr(err)
	tmpSig.SetDesc("The status the rear left tire.")
	checkErr(rearMsg.AppendSignal(tmpSig))
	tmpSig, err = acmelib.NewStandardSignal("TIRE_RL_temperature", tempSigType)
	checkErr(err)
	tmpSig.SetUnit(tempUnit)
	tmpSig.SetDesc("The temperature of the rear left tire.")
	checkErr(rearMsg.AppendSignal(tmpSig))
	tmpSig, err = acmelib.NewStandardSignal("TIRE_RL_pressure", pressSigType)
	checkErr(err)
	tmpSig.SetUnit(pressUnit)
	tmpSig.SetDesc("The pressure of the rear left tire.")
	checkErr(rearMsg.AppendSignal(tmpSig))
	tmpSig, err = acmelib.NewStandardSignal("TIRE_RR_sensID", idSigType)
	checkErr(err)
	tmpSig.SetDesc("The sensor id of the rear right tire.")
	checkErr(rearMsg.AppendSignal(tmpSig))
	tmpSig, err = acmelib.NewStandardSignal("TIRE_RR_status", statusSigType)
	checkErr(err)
	tmpSig.SetDesc("The status the rear right tire.")
	checkErr(rearMsg.AppendSignal(tmpSig))
	tmpSig, err = acmelib.NewStandardSignal("TIRE_RR_temperature", tempSigType)
	checkErr(err)
	tmpSig.SetUnit(tempUnit)
	tmpSig.SetDesc("The temperature of the rear right tire.")
	checkErr(rearMsg.AppendSignal(tmpSig))
	tmpSig, err = acmelib.NewStandardSignal("TIRE_RR_pressure", pressSigType)
	checkErr(err)
	tmpSig.SetUnit(pressUnit)
	tmpSig.SetDesc("The pressure of the rear right tire.")
	checkErr(rearMsg.AppendSignal(tmpSig))

	checkErr(tpmsInt.AddMessage(frontMsg))
	checkErr(tpmsInt.AddMessage(rearMsg))

	frontMsg.AddReceiver(dspace)
	frontMsg.AddReceiver(scanner)

	rearMsg.AddReceiver(dspace)
	rearMsg.AddReceiver(scanner)

	return tpms
}
