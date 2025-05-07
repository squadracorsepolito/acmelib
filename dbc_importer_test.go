package acmelib

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const dbcTestFile = "testdata/expected.dbc"

func Test_ImportDBCFile(t *testing.T) {
	assert := assert.New(t)

	dbcFile, err := os.Open(dbcTestFile)
	assert.NoError(err)
	defer dbcFile.Close()

	_, err = ImportDBCFile(dbcTestFile, dbcFile)
	assert.NoError(err)
}
