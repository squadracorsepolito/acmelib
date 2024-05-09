package main

import (
	"fmt"

	"github.com/squadracorsepolito/acmelib"
)

const bmsLVNodeName = "BMS_LV"
const bmsLVNodeID acmelib.NodeID = 0

var bmsLV = acmelib.NewNode(bmsLVNodeName, bmsLVNodeID)

func bmsLVInit() error {
	h := bmsLVHandler{}

	cellVolt1 := acmelib.NewMessage(h.genMsgName("cell_voltage_1"), 8)
	cellVolt2 := acmelib.NewMessage(h.genMsgName("cell_voltage_2"), 6)
	for i := range 7 {
		tmpSig, err := h.genSigCellVolt(i)
		if err != nil {
			return err
		}

		if i < 4 {
			if err := cellVolt1.AppendSignal(tmpSig); err != nil {
				return err
			}
			continue
		}

		if err := cellVolt2.AppendSignal(tmpSig); err != nil {
			return err
		}
	}

	return nil
}

type bmsLVHandler struct{}

func (h *bmsLVHandler) genMsgName(msgName string) string {
	return fmt.Sprintf("%s_%s", bmsLVNodeName, msgName)
}

func (h *bmsLVHandler) genSigCellVolt(idx int) (*acmelib.StandardSignal, error) {
	sigName := fmt.Sprintf("CELL_%d_Voltage", idx)
	sigDesc := fmt.Sprintf("The voltage of cell %d expressed in mV.", idx)

	sig, err := acmelib.NewStandardSignal(sigName, float16SigType)
	if err != nil {
		return nil, err
	}

	sig.SetDesc(sigDesc)

	if err := sig.SetPhysicalValues(0, 2490.33, 2000, 0.038); err != nil {
		return nil, err
	}

	sig.SetUnit(mVSigUnit)

	return sig, nil
}

func (h *bmsLVHandler) genMsgStatus() (*acmelib.Message, error) {
	msg := acmelib.NewMessage(h.genMsgName("status"), 4)

	relOpenSig, err := acmelib.NewStandardSignal("RELAY_isOpen", flagSigType)
	if err != nil {
		return nil, err
	}
	relOpenSig.SetDesc("States whether the LV relay is open.")
	if err := msg.AppendSignal(relOpenSig); err != nil {
		return nil, err
	}

	for i := range 7 {
		ovSig, err := acmelib.NewStandardSignal(fmt.Sprintf("CELL_%d_isInOV", i), flagSigType)
		if err != nil {
			return nil, err
		}
		ovSig.SetDesc(fmt.Sprintf("States whether the cell %d is in overvoltage.", i))
		if err := msg.AppendSignal(ovSig); err != nil {
			return nil, err
		}

		uvSig, err := acmelib.NewStandardSignal(fmt.Sprintf("CELL_%d_isInUV", i), flagSigType)
		if err != nil {
			return nil, err
		}
		uvSig.SetDesc(fmt.Sprintf("States whether the cell %d is in undervoltage.", i))
		if err := msg.AppendSignal(uvSig); err != nil {
			return nil, err
		}
	}

	for i := range 12 {
		otSig, err := acmelib.NewStandardSignal(fmt.Sprintf("TEMP_%d_isInOT", i), flagSigType)
		if err != nil {
			return nil, err
		}
		otSig.SetDesc(fmt.Sprintf("States whether the cell temperature sensor %d detects an over temperature.", i))
		if err := msg.AppendSignal(otSig); err != nil {
			return nil, err
		}
	}

	return msg, nil
}

func (h *bmsLVHandler) genMsgBattGeneral() (*acmelib.Message, error) {
	msg := acmelib.NewMessage(h.genMsgName("BatteryGeneral"), 6)

	currSensSig, err := acmelib.NewStandardSignal("LV_BATT_CurrentSensorVoltage", float16SigType)
	if err != nil {
		return nil, err
	}
	currSensSig.SetDesc("The voltage returned by the current sensor of the battery pack.")
	if err := currSensSig.SetPhysicalValues(0, 4980.66, 0, 0.076); err != nil {
		return nil, err
	}

	battVoltSig, err := acmelib.NewStandardSignal("LV_BATT_Voltage", float16SigType)
	if err != nil {
		return nil, err
	}
	battVoltSig.SetDesc("The total read voltage of the battery pack.")
	if err := battVoltSig.SetPhysicalValues(0, 17497.845, 14000, 0.267); err != nil {
		return nil, err
	}
	battVoltSig.SetUnit(mVSigUnit)
	if err := msg.AppendSignal(battVoltSig); err != nil {
		return nil, err
	}

	battVoltSumSig, err := acmelib.NewStandardSignal("LV_BATT_VoltageSummed", float16SigType)
	if err != nil {
		return nil, err
	}
	battVoltSumSig.SetDesc("The total voltage of the battery pack calculated by summing each cell voltage.")
	if err := battVoltSumSig.SetPhysicalValues(0, 17497.845, 14000, 0.267); err != nil {
		return nil, err
	}
	battVoltSumSig.SetUnit(mVSigUnit)
	if err := msg.AppendSignal(battVoltSumSig); err != nil {
		return nil, err
	}

	return msg, nil
}
