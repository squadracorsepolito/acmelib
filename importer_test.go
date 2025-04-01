package acmelib

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ImportDBCFile(t *testing.T) {
	assert := assert.New(t)

	inputFile, err := os.Open(expectedDBCFilename)
	assert.NoError(err)

	bus, err := ImportDBCFile(expectedDBCFilename, inputFile)
	assert.NoError(err)
	inputFile.Close()

	// // exporting the bus
	// fileBuf := &strings.Builder{}
	// ExportBus(fileBuf, bus)

	// testFile, err := os.Open(expectedDBCFilename)
	// assert.NoError(err)

	// testFileBuf := &strings.Builder{}
	// _, err = io.Copy(testFileBuf, testFile)
	// assert.NoError(err)
	// testFile.Close()

	// // thanks to Windows that puts \r after \n
	// re := regexp.MustCompile(`\r?\n`)
	// expectedFileStr := re.ReplaceAllString(testFileBuf.String(), "")

	// fileStr := strings.ReplaceAll(fileBuf.String(), "\n", "")

	// assert.Equal(expectedFileStr, fileStr)

	// testing signal types aggregation
	sigTypeIDs := make(map[EntityID]bool)
	for _, tmpNodeInt := range bus.NodeInterfaces() {
		for _, tmpMsg := range tmpNodeInt.SentMessages() {
			for _, sig := range tmpMsg.signals.getValues() {
				if sig.Kind() != SignalKindStandard {
					continue
				}

				stdSig, err := sig.ToStandard()
				assert.NoError(err)
				sigTypeIDs[stdSig.Type().EntityID()] = true
			}
		}
	}
	assert.Len(sigTypeIDs, 1)
}

const muxSignalsDBCFilename = "testdata/mux_signals.dbc"

func Test_ImportDBCFile_MuxSignals(t *testing.T) {
	assert := assert.New(t)

	inputFile, err := os.Open(muxSignalsDBCFilename)
	assert.NoError(err)

	bus, err := ImportDBCFile(muxSignalsDBCFilename, inputFile)
	assert.NoError(err)
	inputFile.Close()

	type muxedSig struct {
		name     string
		startPos int
		size     int
	}

	expected := [][]muxedSig{
		{
			{"DADD", 2, 6},
			{"CADD", 8, 3},
			{"ovrd_bank", 11, 2},
			{"write_RADD", 16, 13},
			{"value", 32, 16},
		},
		{
			{"DADD", 2, 6},
			{"CADD", 8, 3},
			{"ovrd_bank", 11, 2},
			{"ovrd_bitset_idx", 13, 3},
			{"ovrd_bitset", 16, 48},
		},
		{
			{"DADD", 2, 6},
			{"CADD", 8, 3},
			{"ovrd_bank", 11, 2},
			{"ovrd_bank_enabled", 13, 1},
		},
		{
			{"DADD", 2, 6},
			{"CADD", 8, 3},
			{"ovrd_bank", 11, 2},
			{"read_RADD", 16, 13},
		},
	}

	node0Int, err := bus.GetNodeInterfaceByNodeName("node_0")
	assert.NoError(err)

	msg0, err := node0Int.GetSentMessageByName("msg_0")
	assert.NoError(err)

	sig, err := msg0.GetSignalByName("ovrd_cmd")
	assert.NoError(err)

	muxSig, err := sig.ToMultiplexer()
	assert.NoError(err)

	assert.Equal(4, muxSig.GroupCount())
	assert.Equal(62, muxSig.GroupSize())

	for groupID, group := range muxSig.GetSignalGroups() {
		expectedGroup := expected[groupID]
		for idx, tmpSig := range group {
			expectedSig := expectedGroup[idx]
			assert.Equal(expectedSig.name, tmpSig.Name())
			assert.Equal(expectedSig.startPos, tmpSig.GetStartBit())
			assert.Equal(expectedSig.size, tmpSig.GetSize())
		}
	}
}
