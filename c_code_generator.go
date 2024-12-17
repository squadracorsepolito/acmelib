package acmelib

import (
	"io"

	"text/template"
)

const tmpTemplatesFolder = "../templates"

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
	hTmpl, err := template.New("c_header").ParseGlob(tmpTemplatesFolder + "/*.tmpl")
	if err != nil {
		return err
	}

	cTmpl, err := template.New("c_source").ParseGlob(tmpTemplatesFolder + "/*.tmpl")
	if err != nil {
		return err
	}

	if err := hTmpl.ExecuteTemplate(g.hFile, "bus_h", bus); err != nil {
		return err
	}

	if err := cTmpl.ExecuteTemplate(g.cFile, "bus_c", bus); err != nil {
		return err
	}

	return nil
}
