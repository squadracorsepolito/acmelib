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

	res := new(strings.Builder)
	ExportDBCBus(res, tdNet.bus)

	resLines := make(map[string]bool)
	for line := range strings.SplitSeq(res.String(), "\n") {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}

		resLines[trimmedLine] = false
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

		if _, ok := resLines[line]; !ok {
			assert.Fail(fmt.Sprintf("missing line: %s", line))
			continue
		}

		resLines[line] = true
	}

	for line, ok := range resLines {
		if !ok {
			assert.Fail(fmt.Sprintf("unexpected extra line: %s", line))
		}
	}
}
