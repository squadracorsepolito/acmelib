package main

import (
	"os"

	"github.com/squadracorsepolito/acmelib"
)

func main() {
	sc24 := acmelib.NewNetwork("SC24")

	dbcFile, err := os.Open("MCB.dbc")
	if err != nil {
		panic(err)
	}
	defer dbcFile.Close()

	mcb, err := acmelib.ImportDBCFile("mcb", dbcFile)
	if err != nil {
		panic(err)
	}

	if err := mcb.UpdateName("Main CAN Bus"); err != nil {
		panic(err)
	}

	if err := sc24.AddBus(mcb); err != nil {
		panic(err)
	}

	mdFile, err := os.Create("SC24.md")
	if err != nil {
		panic(err)
	}
	defer mdFile.Close()

	if err := acmelib.ExportToMarkdown(sc24, mdFile); err != nil {
		panic(err)
	}
}
