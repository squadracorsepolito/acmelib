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

	net := initTestNetwork(assert)

	// exporting the bus
	fileBuf := &strings.Builder{}
	ExportBus(fileBuf, net.Buses()[0])

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
