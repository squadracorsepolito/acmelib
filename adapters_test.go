package acmelib

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadDBC(t *testing.T) {
	assert := assert.New(t)

	net, err := LoadDBC("test_network", "./testdata/mcb.dbc")
	assert.NoError(err)

	assert.NoError(os.WriteFile("./testdata/res_mcb.txt", []byte(net.String()), 0644))
}
