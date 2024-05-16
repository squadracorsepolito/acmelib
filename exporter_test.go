package acmelib

import (
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const expectedDBCFilename = "testdata/expected.dbc"

func Test_ExportBus(t *testing.T) {
	assert := assert.New(t)

	net := NewNetwork("net")

	bus0 := NewBus("bus_0")
	bus0.SetDesc("bus0 description")
	assert.NoError(net.AddBus(bus0))

	node0 := NewNode("node_0", 0, 1)
	node0.SetDesc("node0 description")
	node0Int := node0.Interfaces()[0]
	assert.NoError(bus0.AddNodeInterface(node0Int))

	recNode0Int := NewNode("rec_node_0", 1, 1).Interfaces()[0]
	assert.NoError(bus0.AddNodeInterface(recNode0Int))

	// msg0 has a signal in big endian order and a multiplexer signal
	msg0 := NewMessage("msg_0", 1, 8)
	assert.NoError(node0Int.AddMessage(msg0))

	size4Type, err := NewIntegerSignalType("4_bits", 4, false)
	assert.NoError(err)

	stdSig0, err := NewStandardSignal("std_sig_0", size4Type)
	assert.NoError(err)
	stdSig0.SetUnit(NewSignalUnit("deg_celsius", SignalUnitKindTemperature, "degC"))
	// stdSig0.SetByteOrder(SignalByteOrderBigEndian)
	assert.NoError(msg0.AppendSignal(stdSig0))

	muxSig0, err := NewMultiplexerSignal("mux_sig_0", 4, 16)
	assert.NoError(err)
	assert.NoError(msg0.AppendSignal(muxSig0))

	fixedSig, err := NewStandardSignal("fixed_sig", size4Type)
	assert.NoError(err)
	assert.NoError(muxSig0.InsertSignal(fixedSig, 0))

	multiGroupSig0, err := NewStandardSignal("multi_group_sig_0", size4Type)
	assert.NoError(err)
	assert.NoError(muxSig0.InsertSignal(multiGroupSig0, 4, 0, 2, 3))

	oneGroupSig0, err := NewStandardSignal("one_group_sig_0", size4Type)
	assert.NoError(err)
	assert.NoError(muxSig0.InsertSignal(oneGroupSig0, 4, 1))

	// msg1 has a 2 level nested multiplexer signals
	msg1 := NewMessage("msg_1", 2, 8)
	assert.NoError(node0Int.AddMessage(msg1))
	msg1.SetDesc("msg1 description")

	muxSig1, err := NewMultiplexerSignal("mux_sig_1", 4, 16)
	assert.NoError(err)
	assert.NoError(msg1.AppendSignal(muxSig1))
	muxSig1.SetDesc("mux1 description")

	oneGroupSig1, err := NewStandardSignal("one_group_sig_1", size4Type)
	assert.NoError(err)
	assert.NoError(muxSig1.InsertSignal(oneGroupSig1, 0, 0))

	nestedMuxSig1, err := NewMultiplexerSignal("nested_mux sig_1", 2, 8)
	assert.NoError(err)
	assert.NoError(muxSig1.InsertSignal(nestedMuxSig1, 0, 1))

	oneGroupSig2, err := NewStandardSignal("one_group_sig_2", size4Type)
	assert.NoError(err)
	assert.NoError(nestedMuxSig1.InsertSignal(oneGroupSig2, 0, 0))

	multiGroupSig1, err := NewStandardSignal("multi_group sig_1", size4Type)
	assert.NoError(err)
	assert.NoError(nestedMuxSig1.InsertSignal(multiGroupSig1, 4))

	// msg2 has an enum signal
	msg2 := NewMessage("msg_2", 3, 8)
	assert.NoError(node0Int.AddMessage(msg2))
	msg2.AddReceiver(recNode0Int)

	enum := NewSignalEnum("enum")
	assert.NoError(enum.AddValue(NewSignalEnumValue("VALUE_0", 0)))
	assert.NoError(enum.AddValue(NewSignalEnumValue("VALUE_1", 1)))
	assert.NoError(enum.AddValue(NewSignalEnumValue("VALUE_15", 15)))

	enumSig0, err := NewEnumSignal("enum_sig_0", enum)
	assert.NoError(err)
	assert.NoError(msg2.AppendSignal(enumSig0))

	// msg3 is big endian
	msg3 := NewMessage("msg_3", 4, 1)
	msg3.SetByteOrder(MessageByteOrderBigEndian)
	assert.NoError(node0Int.AddMessage(msg3))
	stdSig1, err := NewStandardSignal("std_sig_1", size4Type)
	assert.NoError(err)
	stdSig2, err := NewStandardSignal("std_sig_2", size4Type)
	assert.NoError(err)
	assert.NoError(msg3.AppendSignal(stdSig1))
	assert.NoError(msg3.AppendSignal(stdSig2))

	// attributes testing
	strAtt := NewStringAttribute("str_att", "")

	intAtt, err := NewIntegerAttribute("int_att", 0, 0, 10000)
	assert.NoError(err)

	hexAtt, err := NewIntegerAttribute("hex_att", 0, 0, 10000)
	assert.NoError(err)
	hexAtt.SetFormatHex()

	floatAtt, err := NewFloatAttribute("float_att", 1.5, 0, 100.5)
	assert.NoError(err)

	enumAtt, err := NewEnumAttribute("enum_att", "VALUE_0", "VALUE_1", "VALUE_2", "VALUE_3")
	assert.NoError(err)

	bus0.AddAttributeValue(strAtt, "bus0_value")
	node0Int.node.AddAttributeValue(intAtt, 1)
	msg0.AddAttributeValue(hexAtt, 1)
	stdSig0.AddAttributeValue(enumAtt, "VALUE_1")
	muxSig0.AddAttributeValue(floatAtt, 50.75)

	// special attributes testing
	msg2.SetCycleTime(10)
	msg2.SetDelayTime(20)
	msg2.SetStartDelayTime(30)
	msg2.SetSendType(MessageSendTypeCyclicIfActiveAndTriggered)

	enumSig0.SetSendType(SignalSendTypeOnChangeWithRepetition)

	// exporting the bus
	fileBuf := &strings.Builder{}
	ExportBus(fileBuf, bus0)

	testFile, err := os.Open(expectedDBCFilename)
	assert.NoError(err)

	testFileBuf := &strings.Builder{}
	_, err = io.Copy(testFileBuf, testFile)
	assert.NoError(err)
	testFile.Close()

	// thanks to Windows that puts \r after \n
	re := regexp.MustCompile(`\r?\n`)
	expectedFileStr := re.ReplaceAllString(testFileBuf.String(), "")

	fileStr := strings.ReplaceAll(fileBuf.String(), "\n", "")

	assert.Equal(expectedFileStr, fileStr)
}
