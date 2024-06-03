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
	modifySignalTypeName(dspaceInt, "DSPACE_rtdACK", "RTD_FSM_STATE", "rtd_fsm_t")

	// adding xpc tx/rx

	expMsgID := acmelib.MessageID(10)
	for _, nodeInt := range mcb.NodeInterfaces() {
		nodeName := nodeInt.Node().Name()

		msgName := fmt.Sprintf("%s_xcp", nodeName)
		tmpMsg := acmelib.NewMessage(msgName, expMsgID, 8)
		checkErr(nodeInt.AddMessage(tmpMsg))

		msgDesc := ""
		if nodeInt.Node().ID() == 0 {
			for idx, rec := range mcb.NodeInterfaces() {
				if idx == 0 {
					continue
				}
				tmpMsg.AddReceiver(rec)
			}

			msgDesc = "The message used to flash a board."
		} else {
			tmpMsg.AddReceiver(mcb.NodeInterfaces()[0])
			msgDesc = "The message used to notify a board is flashed."
		}

		tmpMsg.SetDesc(msgDesc)
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
