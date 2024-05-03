package main

import (
	"os"

	"github.com/FerroO2000/acmelib"
	"github.com/FerroO2000/acmelib/plugins/dm1"
)

func main() {
	mcbFile, err := os.Open("MCB.dbc")
	if err != nil {
		panic(err)
	}

	mcbBus, err := acmelib.ImportDBCFile("MCB", mcbFile)
	if err != nil {
		panic(err)
	}

	mcbFile.Close()

	mcbBusDM1, err := dm1.GenerateDM1Messages(mcbBus)
	if err != nil {
		panic(err)
	}

	mcbDM1File, err := os.Create("MCB_DM1.dbc")
	if err != nil {
		panic(err)
	}

	acmelib.ExportBus(mcbDM1File, mcbBusDM1)

	mcbDM1File.Close()
}
