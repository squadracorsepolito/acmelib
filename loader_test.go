package acmelib

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadNetwork(t *testing.T) {
	assert := assert.New(t)

	tdNet := initNetwork(assert)

	net := tdNet.net
	buf := new(bytes.Buffer)
	assert.NoError(SaveNetwork(net, &SaveNetworkOptions{
		WireWriter: buf,
	}))

	log.Print("----------------------------------------")

	loadNet, err := LoadNetwork(buf, SaveEncodingWire)
	assert.NoError(err)

	assert.Equal(net.Name(), loadNet.Name())
	assert.Len(loadNet.Buses(), 1)

	dbcRes := new(strings.Builder)
	ExportDBCBus(dbcRes, loadNet.Buses()[0])
	compareDBCFiles(assert, dbcRes)
}
