package main

import (
	"log"
	"os"

	"github.com/squadracorsepolito/acmelib"
)

func main() {
	sc24 := acmelib.NewNetwork("SC24")

	mcbFile, err := os.Open("MCB.dbc")
	if err != nil {
		panic(err)
	}
	defer mcbFile.Close()

	mcb, err := acmelib.ImportDBCFile("mcb", mcbFile)
	if err != nil {
		panic(err)
	}

	if err := mcb.UpdateName("Main CAN Bus"); err != nil {
		panic(err)
	}

	if err := sc24.AddBus(mcb); err != nil {
		panic(err)
	}

	// hvcbFile, err := os.Open("HVCB.dbc")
	// if err != nil {
	// 	panic(err)
	// }
	// defer hvcbFile.Close()

	// hvcb, err := acmelib.ImportDBCFile("hvcb", hvcbFile)
	// if err != nil {
	// 	panic(err)
	// }

	// if err := sc24.AddBus(hvcb); err != nil {
	// 	panic(err)
	// }

	mcb.SetBaudrate(1_000_000)

	busLoad, err := acmelib.CalculateBusLoad(mcb, 1000)
	if err != nil {
		panic(err)
	}
	log.Print("BUS LOAD: ", busLoad)

	mdFile, err := os.Create("SC24.md")
	if err != nil {
		panic(err)
	}
	defer mdFile.Close()

	if err := acmelib.ExportToMarkdown(sc24, mdFile); err != nil {
		panic(err)
	}
}
