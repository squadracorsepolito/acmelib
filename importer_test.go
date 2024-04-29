package acmelib

import (
	"os"
	"testing"

	"github.com/FerroO2000/acmelib/dbc"
	"github.com/stretchr/testify/assert"
)

func Test_ImportDBCFile(t *testing.T) {
	assert := assert.New(t)

	file, err := os.ReadFile(cmpExportBusFilename)
	assert.NoError(err)

	parser := dbc.NewParser("test_ExportBus", file)
	dbcFile, err := parser.Parse()
	assert.NoError(err)

	// t.Log(dbcFile)

	bus, err := ImportDBCFile(dbcFile)
	assert.NoError(err)

	t.Log(bus)

	// resFile, err := os.Create("testdata/res.dbc")
	// assert.NoError(err)

	// ExportBus(resFile, bus)
	// resFile.Close()
}
