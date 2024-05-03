package dm1

import (
	"fmt"
	"strings"

	"github.com/FerroO2000/acmelib"
)

const dm1MsgNamePrefix = "dm1"

var fmiValues = map[string]int{
	"high_severity":                  0,
	"low_severity":                   1,
	"erratic":                        2,
	"v_above_normal":                 3,
	"v_below_normal":                 4,
	"i_above_normal":                 5,
	"i_below_normal":                 6,
	"system_not_responding_properly": 7,
	"abnormal_frequency":             8,
	"abnormal_update_rate":           9,
	"abnormal_rate_of_change":        10,
	"other_failure_mode":             11,
	"failure":                        12,
	"out_of_calibration":             13,
	"special_instruction":            14,
	"data_valid_above_normal_range0": 15,
	"data_valid_above_normal_range1": 16,
	"data_valid_below_normal_range0": 17,
	"data_valid_below_normal_range1": 18,
	"received_network_data_error":    19,
	"data_drifted_high":              20,
	"data_drifted_low":               21,
	"condition_exists":               31,
}

func generateSIN(messages []*acmelib.Message) *acmelib.EnumSignal {
	sinEnum := acmelib.NewSignalEnum("sin_enum")
	sinEnum.SetMinSize(8)

	valIdx := 0
	signalNames := make(map[string]bool)
	for _, msg := range messages {
		tmpMsgName := msg.Name()
		splNames := strings.Split(tmpMsgName, "_")

		if splNames[0] == dm1MsgNamePrefix {
			continue
		}

		for _, sig := range msg.Signals() {
			valName := strings.ReplaceAll(sig.Name(), " ", "_")
			if _, ok := signalNames[valName]; ok {
				continue
			}
			signalNames[valName] = true

			tmpVal := acmelib.NewSignalEnumValue(valName, valIdx)
			if err := sinEnum.AddValue(tmpVal); err != nil {
				panic(err)
			}
			valIdx++
		}
	}

	sin, err := acmelib.NewEnumSignal("sin", sinEnum)
	if err != nil {
		panic(err)
	}

	return sin
}

func generateFMI() *acmelib.EnumSignal {
	fmiEnum := acmelib.NewSignalEnum("fmi_enum")

	for valName, valIdx := range fmiValues {
		tmpVal := acmelib.NewSignalEnumValue(valName, valIdx)
		if err := fmiEnum.AddValue(tmpVal); err != nil {
			panic(err)
		}
	}

	fmi, err := acmelib.NewEnumSignal("fmi", fmiEnum)
	if err != nil {
		panic(err)
	}

	return fmi
}

func generateOccCounter() *acmelib.StandardSignal {
	occSigType, err := acmelib.NewIntegerSignalType("uint_8", 8, false)
	if err != nil {
		panic(err)
	}

	occ, err := acmelib.NewStandardSignal("occ_counter", occSigType)
	if err != nil {
		panic(err)
	}

	return occ
}

func GenerateDM1Messages(bus *acmelib.Bus) (*acmelib.Bus, error) {
	messages := []*acmelib.Message{}
	for _, node := range bus.Nodes() {
		messages = append(messages, node.Messages()...)
	}

	for _, node := range bus.Nodes() {
		msgName := fmt.Sprintf("%s_%s", dm1MsgNamePrefix, node.Name())

		dm1 := acmelib.NewMessage(msgName, 8)
		if err := dm1.InsertSignal(generateSIN(messages), 0); err != nil {
			return nil, err
		}
		if err := dm1.InsertSignal(generateFMI(), 8); err != nil {
			return nil, err
		}
		if err := dm1.InsertSignal(generateOccCounter(), 16); err != nil {
			return nil, err
		}

		if err := node.AddMessage(dm1); err != nil {
			return nil, err
		}
	}

	return bus, nil
}
