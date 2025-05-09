package acmelib

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExportDBCBus(t *testing.T) {
	assert := assert.New(t)

	tdNet := initNetwork(assert)

	dbcRes := new(strings.Builder)
	ExportDBCBus(dbcRes, tdNet.bus)
	compareDBCFiles(assert, dbcRes)
}

func compareDBCFiles(assert *assert.Assertions, actual *strings.Builder) {
	actualLines := make(map[string]bool)
	for line := range strings.SplitSeq(actual.String(), "\n") {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}

		actualLines[trimmedLine] = false
	}

	expectedFile, err := os.Open(dbcTestFile)
	assert.NoError(err)
	defer expectedFile.Close()

	buff := bufio.NewScanner(expectedFile)
	for buff.Scan() {
		line := strings.TrimSpace(buff.Text())
		if line == "" {
			continue
		}

		if _, ok := actualLines[line]; !ok {
			assert.Fail(fmt.Sprintf("missing line: %s", line))
			continue
		}

		actualLines[line] = true
	}

	for line, ok := range actualLines {
		if !ok {
			assert.Fail(fmt.Sprintf("unexpected extra line: %s", line))
		}
	}
}
