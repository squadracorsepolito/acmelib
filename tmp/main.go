package main

import (
	"bytes"
	"os"
	"path"

	"github.com/squadracorsepolito/acmelib"
)

func main() {
	// opening the test network file
	netWireFile, err := os.Open("../testdata/expected.binpb")
	must(err)
	defer netWireFile.Close()

	// loading the test network
	net, err := acmelib.LoadNetwork(netWireFile, acmelib.SaveEncodingWire)
	must(err)

	// extract the first bus from the network
	buses := net.Buses()
	if len(buses) == 0 {
		panic("no buses")
	}
	bus := buses[0]

	// creating the buffers for the files
	hFileBuf := new(bytes.Buffer)
	wFileBuf := new(bytes.Buffer)

	// generating the c code
	must(acmelib.GenerateCCode(bus, hFileBuf, wFileBuf))

	// writing the files
	must(os.WriteFile(path.Join("res", "test.h"), hFileBuf.Bytes(), 0666))
	must(os.WriteFile(path.Join("res", "test.c"), wFileBuf.Bytes(), 0666))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
