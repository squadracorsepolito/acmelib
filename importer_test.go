package acmelib

import (
	"io"
	"os"
	"regexp"
	"strings"
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

	// exporting the bus
	fileBuf := &strings.Builder{}
	ExportBus(fileBuf, bus)

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

	// testing signal types aggregation
	sigTypeIDs := make(map[EntityID]bool)
	for _, tmpNodeInt := range bus.NodeInterfaces() {
		for _, tmpMsg := range tmpNodeInt.Messages() {
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
