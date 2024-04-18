package main

import (
	"fmt"

	"github.com/FerroO2000/acmelib"
)

var (
	flagType       = acmelib.NewFlagSignalType("flag")
	float16Type, _ = acmelib.NewFloatSignalType("float16", 16)
)

var (
	voltageUnit    = acmelib.NewSignalUnit("voltage", acmelib.SignalUnitKindElectrical, "V")
	celsiusDegUnit = acmelib.NewSignalUnit("celsius_deg", acmelib.SignalUnitKindTemperature, "degC")
)

func main() {
	hvcb := acmelib.NewBus("hvcb")

	hvbNode := acmelib.NewNode("HVB", 0)
	err := hvcb.AddNode(hvbNode)
	panicErr(err)

	ivtMainNode := acmelib.NewNode("IVTMain", 1)
	err = hvcb.AddNode(ivtMainNode)
	panicErr(err)

	pcNode := acmelib.NewNode("PC", 2)
	err = hvcb.AddNode(pcNode)
	panicErr(err)

	vcuNode := acmelib.NewNode("VCU", 3)
	err = hvcb.AddNode(vcuNode)
	panicErr(err)

	chargerNode := acmelib.NewNode("Charger", 4)
	err = hvcb.AddNode(chargerNode)
	panicErr(err)

	info01DbgVMsg := info01DbgV()
	err = hvbNode.AddMessage(info01DbgVMsg)
	panicErr(err)
	info01DbgVMsg.AddReceiver(pcNode)

	info02DbgTMsg := info02DbgT()
	err = hvbNode.AddMessage(info02DbgTMsg)
	panicErr(err)
	info02DbgTMsg.AddReceiver(pcNode)

	hvbTXVCUCmdMsg := hvbTXVCUCmd()
	err = vcuNode.AddMessage(hvbTXVCUCmdMsg)
	panicErr(err)
	hvbTXVCUCmdMsg.AddReceiver(hvbNode)

	nvbRXDiagnosisMsg := nvbRXDiagnosis()
	err = hvbNode.AddMessage(nvbRXDiagnosisMsg)
	panicErr(err)
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func info01DbgV() *acmelib.Message {
	msg := acmelib.NewMessage("INFO_01_DbgV", 8)
	msg.SetID(288)

	muxSig, err := acmelib.NewMultiplexerSignal("BMS_eDbgV", 56, 8)
	panicErr(err)

	selectValue := -1
	for i := 0; i < 255; i++ {
		if i%3 == 0 {
			selectValue++
		}

		sig, err := acmelib.NewStandardSignal(fmt.Sprintf("BMS_VDbgV%3d", i), float16Type)
		panicErr(err)

		sig.SetPhysicalValues(0, 4.95, 0, 0.001)
		sig.SetUnit(voltageUnit)
		sig.SetDesc(fmt.Sprintf("Cell %d voltage.", i))

		err = muxSig.AppendMuxSignal(selectValue, sig)
		panicErr(err)
	}

	msg.AppendSignal(muxSig)

	return msg
}

func info02DbgT() *acmelib.Message {
	msg := acmelib.NewMessage("INFO_01_DbgV", 8)
	msg.SetID(289)

	muxSig, err := acmelib.NewMultiplexerSignal("BMS_eDbgT", 55, 7)
	panicErr(err)

	selectValue := -1
	for i := 0; i < 127; i++ {
		if i%3 == 0 {
			selectValue++
		}

		sig, err := acmelib.NewStandardSignal(fmt.Sprintf("BMS_TDbgT%3d", i), float16Type)
		panicErr(err)

		sig.SetPhysicalValues(-40, 105, -273.15, 0.01)
		sig.SetUnit(celsiusDegUnit)
		sig.SetDesc(fmt.Sprintf("Thermistor %d temperature.", i))

		err = muxSig.AppendMuxSignal(selectValue, sig)
		panicErr(err)
	}

	return msg
}

func hvbTXVCUCmd() *acmelib.Message {
	msg := acmelib.NewMessage("HVB_TX_VCUCmd", 8)
	msg.SetID(336)

	enum := acmelib.NewSignalEnum("status")
	err := enum.AddValue(acmelib.NewSignalEnumValue("DISABLED", 0))
	panicErr(err)
	err = enum.AddValue(acmelib.NewSignalEnumValue("ENABLED", 1))
	panicErr(err)

	invReqSig, err := acmelib.NewEnumSignal("VCU_bHvbInvReq", enum)
	panicErr(err)
	invReqSig.SetDesc("Requested closing Inverter conductors  by VCU")
	err = msg.InsertSignal(invReqSig, 0)
	panicErr(err)

	clrErrSig, err := acmelib.NewStandardSignal("VCU_ClrErr", flagType)
	panicErr(err)
	err = msg.InsertSignal(clrErrSig, 6)
	panicErr(err)

	balReqSig, err := acmelib.NewEnumSignal("VCU_bBalReq", enum)
	panicErr(err)
	balReqSig.SetDesc("Enables pack balancing.")
	err = msg.InsertSignal(balReqSig, 8)
	panicErr(err)

	allVTReqSig, err := acmelib.NewEnumSignal("VCU_bAllVTReq", enum)
	panicErr(err)
	balReqSig.SetDesc("Enables pack all V & T message.")
	err = msg.InsertSignal(allVTReqSig, 14)
	panicErr(err)

	return msg
}

func nvbRXDiagnosis() *acmelib.Message {
	msg := acmelib.NewMessage("HVB_RX_Diagnosis", 8)
	msg.SetID(512)

	sigNames := []string{
		"HVB_Diag_Flash", "HVB_Diag_eeprom", "HVB_Diag_RAM",
		"HVB_Diag_CAN", "HVB_Diag_UART",
		"HVB_Diag_cell_sna", "HVB_Diag_vcu_can_sna", "HVB_Diag_bat_curr_sna", "HVB_Diag_inv_vlt_sna", "HVB_Diag_bat_vlt_sna",
		"HVB_Diag_cell_ut", "HVB_Diag_cell_ot", "HVB_Diag_cell_uv", "HVB_Diag_cell_ov", "HVB_Diag_bat_uv", "HVB_Diag_imd_sna", "HVB_Diag_imd__low_r",
		"HVB_Diag_bat_curr_oc", "HVB_Diag_inv_vlt_ov",
		"HVB_Recovery_Active",
	}

	sigStartBits := []int{
		0, 1, 2,
		4, 5,
		16, 17, 18, 19, 20,
		32, 33, 34, 35, 36, 37, 38,
		40, 41,
		56,
	}

	for i, name := range sigNames {
		sig, err := acmelib.NewStandardSignal(name, flagType)
		panicErr(err)
		err = msg.InsertSignal(sig, sigStartBits[i])
		panicErr(err)
	}

	return msg
}
