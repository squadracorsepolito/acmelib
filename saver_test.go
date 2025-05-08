package acmelib

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SaveNetwork(t *testing.T) {
	assert := assert.New(t)

	tdNet := initNetwork(assert)

	jsonFile, err := os.Create("./testdata/new_expected.json")
	assert.NoError(err)
	defer jsonFile.Close()

	net := tdNet.net
	assert.NoError(SaveNetwork(net, &SaveNetworkOptions{
		JSONWriter: jsonFile,
	}))
}
