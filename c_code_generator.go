package acmelib

import (
	"io"

	"text/template"

	"github.com/squadracorsepolito/acmelib/templates"
)

func GenerateCCode(bus *Bus, hFile io.Writer, cFile io.Writer) error {
	csGen := newCSourceGenerator(hFile, cFile)
	return csGen.generateBus(bus)
}

type cCodeGenerator struct {
	hFile io.Writer
	cFile io.Writer
}

func newCSourceGenerator(hFile io.Writer, cFile io.Writer) *cCodeGenerator {
	return &cCodeGenerator{
		hFile: hFile,
		cFile: cFile,
	}
}

func (g *cCodeGenerator) generateBus(bus *Bus) error {
	hTmpl, err := template.New("c_header").Parse(templates.BusHeader)
	if err != nil {
		return err
	}

	cTmpl, err := template.New("c_source").Parse(templates.BusSource)
	if err != nil {
		return err
	}

	if err := hTmpl.Execute(g.hFile, bus); err != nil {
		return err
	}

	if err := cTmpl.Execute(g.cFile, bus); err != nil {
		return err
	}

	return nil
}
