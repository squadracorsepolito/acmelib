package main

import (
	"log"
	"os"

	"github.com/squadracorsepolito/acmelib"
)

var nodeIDs = map[string]acmelib.NodeID{
	"TLB_BAT":    1,
	"SB_FRONT":   2,
	"SB_REAR":    3,
	"BMS_LV":     4,
	"DASH":       5,
	"DIAG_TOOL":  6,
	"DSPACE":     7,
	"EXTRA_NODE": 8,
	"SCANNER":    9,
	"TPMS":       10,
	"IMU":        11,
	"BRUSA":      12,
}

var messageIDs = map[string]acmelib.MessageID{
	"DSPACE_timeAndDate":          1,
	"TLB_BAT_signalsStatus":       4,
	"SB_FRONT_analogDevice":       5,
	"SB_REAR_analogDevice":        6,
	"SB_REAR_criticalPeripherals": 7,

	"BMS_LV_lvBatGeneral":    20,
	"DASH_hmiDevicesState":   22,
	"DSPACE_peripheralsCTRL": 25,

	"DIAG_TOOL_xcpTxTLB_BAT":    40,
	"DIAG_TOOL_xcpTxSB_FRONT":   41,
	"DIAG_TOOL_xcpTxSB_REAR":    42,
	"DIAG_TOOL_xcpTxBMS_LV":     43,
	"DIAG_TOOL_xcpTxDASH":       44,
	"DIAG_TOOL_xcpTxSCANNER":    45,
	"TLB_BAT_sdcSensingStatus":  46,
	"SB_REAR_sdcSensingStatus":  47,
	"SB_FRONT_sdcSensingStatus": 48,
	"SB_FRONT_potentiometer":    50,
	"SB_REAR_potentiometer":     51,

	"TLB_BAT_xcpTx":            70,
	"SB_FRONT_xcpTx":           70,
	"SB_REAR_xcpTx":            70,
	"BMS_LV_xcpTx":             70,
	"DASH_xcpTx":               70,
	"SCANNER_xcpTx":            70,
	"DSPACE_signals":           73,
	"DSPACE_fsmStates":         74,
	"BMS_LV_cellsStatus":       75,
	"BMS_LV_status":            76,
	"BMS_LV_lvCellVoltage0":    77,
	"BMS_LV_lvCellVoltage1":    78,
	"DASH_peripheralsStatus":   79,
	"TPMS_frontWheelsPressure": 80,
	"TPMS_rearWheelsPressure":  81,

	"TLB_BAT_hello":               100,
	"SB_FRONT_hello":              100,
	"SB_REAR_hello":               100,
	"BMS_LV_hello":                100,
	"DASH_hello":                  100,
	"DSPACE_hello":                100,
	"BMS_LV_lvCellNTCResistance0": 101,
	"BMS_LV_lvCellNTCResistance1": 102,
	"SB_FRONT_ntcResistance":      103,
	"SB_REAR_ntcResistance":       104,
	"DASH_appsRangeLimits":        105,
	"DASH_carCommands":            106,
	"DSPACE_dashLedsColorRGB":     107,
}

func main() {
	sc24 := acmelib.NewNetwork("SC24")
	sc24.SetDesc("The CAN network of the squadracorse 2024 formula SAE car")

	// load mcb
	mcbFile, err := os.Open("MCB.dbc")
	checkErr(err)
	defer mcbFile.Close()

	mcb, err := acmelib.ImportDBCFile("mcb", mcbFile)
	checkErr(err)
	checkErr(mcb.UpdateName("Main CAN Bus"))
	checkErr(sc24.AddBus(mcb))

	// renaming signal types
	dashInt, err := mcb.GetNodeInterfaceByNodeName("DASH")
	checkErr(err)

	tmpSigType, err := acmelib.NewIntegerSignalType("fan_pwm_t", 4, false)
	checkErr(err)
	tmpSigType.SetMax(10)
	modifySignalType(dashInt, "DASH_peripheralsStatus", "TSAC_FAN_pwmDutyCycleStatus", tmpSigType)

	modifySignalTypeName(dashInt, "DASH_hmiDevicesState", "ROT_SW_1_state", "rotary_switch_state_t")
	modifySignalTypeName(dashInt, "DASH_appsRangeLimits", "APPS_0_voltageRangeMin", "uint16_t")
	modifySignalTypeName(dashInt, "DASH_carCommands", "BMS_LV_diagPWD", "bms_lv_password_t")

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
	modifySignalType(dspaceInt, "DSPACE_timeAndDate", "DATETIME_seconds", tmpSigType)

	modifySignalTypeName(dspaceInt, "DSPACE_timeAndDate", "DATETIME_month", "month_t")
	modifySignalTypeName(dspaceInt, "DSPACE_timeAndDate", "DATETIME_day", "day_t")
	modifySignalTypeName(dspaceInt, "DSPACE_timeAndDate", "DATETIME_hours", "hours_t")
	modifySignalTypeName(dspaceInt, "DSPACE_timeAndDate", "DATETIME_minutes", "minutes_t")

	// calculte bus load
	mcb.SetBaudrate(1_000_000)
	busLoad, err := acmelib.CalculateBusLoad(mcb, 1000)
	checkErr(err)
	log.Print("BUS LOAD: ", busLoad)

	// parse IDs
	parseNodeIDs(mcb)
	parseMessageIDs(mcb)

	// save files
	dbcFile, err := os.Create("mcb_parsed.dbc")
	checkErr(err)
	defer dbcFile.Close()
	acmelib.ExportBus(dbcFile, mcb)

	wireFile, err := os.Create("SC24.binpb")
	checkErr(err)
	defer wireFile.Close()
	jsonFile, err := os.Create("SC24.json")
	checkErr(err)
	defer jsonFile.Close()
	checkErr(acmelib.SaveNetwork(sc24, acmelib.SaveEncodingWire|acmelib.SaveEncodingJSON, wireFile, jsonFile, nil))

	mdFile, err := os.Create("SC24.md")
	checkErr(err)
	defer mdFile.Close()

	checkErr(acmelib.ExportToMarkdown(sc24, mdFile))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func modifySignalTypeName(nodeInt *acmelib.NodeInterface, msgName, sigName, newName string) {
	tmpMsg, err := nodeInt.GetSentMessageByName(msgName)
	checkErr(err)
	tmpSig, err := tmpMsg.GetSignalByName(sigName)
	checkErr(err)
	tmpStdSig, err := tmpSig.ToStandard()
	checkErr(err)
	tmpStdSig.Type().SetName(newName)
}

func modifySignalType(nodeInt *acmelib.NodeInterface, msgName, sigName string, newType *acmelib.SignalType) {
	tmpMsg, err := nodeInt.GetSentMessageByName(msgName)
	checkErr(err)
	tmpSig, err := tmpMsg.GetSignalByName(sigName)
	checkErr(err)
	tmpStdSig, err := tmpSig.ToStandard()
	checkErr(err)
	tmpStdSig.SetType(newType)
}

func parseNodeIDs(mcb *acmelib.Bus) {
	interfaces := mcb.NodeInterfaces()
	for i := len(interfaces) - 1; i >= 0; i = i - 1 {
		tmpNodeInt := interfaces[i]
		tmpNode := tmpNodeInt.Node()
		if nodeID, ok := nodeIDs[tmpNode.Name()]; ok {
			checkErr(tmpNode.UpdateID(nodeID))
		}
	}
}

func parseMessageIDs(mcb *acmelib.Bus) {
	for _, tmpNodeInt := range mcb.NodeInterfaces() {
		for _, tmpMsg := range tmpNodeInt.SentMessages() {
			if msgID, ok := messageIDs[tmpMsg.Name()]; ok {
				checkErr(tmpMsg.UpdateID(msgID))
			}
		}
	}
}
