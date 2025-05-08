package acmelib

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SaveNetwork(t *testing.T) {
	assert := assert.New(t)

	tdNet := initNetwork(assert)

	resBuf := new(bytes.Buffer)

	net := tdNet.net
	assert.NoError(SaveNetwork(net, &SaveNetworkOptions{
		JSONWriter: resBuf,
	}))
}
