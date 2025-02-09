package acmelib

import (
	"io"

	"text/template"
	"strings"
	"fmt"
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
	// define DB name
	dbName := "simple"

	funcMap := template.FuncMap{
		"toUpper": strings.ToUpper,
		"toLower": strings.ToLower,
		"toUint": func(i interface{}) string {
			switch v := i.(type) {
			case MessageID:
				return fmt.Sprintf("0x%Xu", uint(v))
			case int:
				return fmt.Sprintf("%du", (v))
			default:
				return "invalid type"
			}
		},
		"isExtended": func(id MessageID) int {
			if (uint(id) & 0x80000000) == 0 {
				return 0
			}
			return 1
		},
	}	

	hTmpl, err := template.New("c_header").Funcs(funcMap).ParseGlob(tmpTemplatesFolder + "/*.tmpl")
	if err != nil {
		return err
	}

	cTmpl, err := template.New("c_source").Funcs(funcMap).ParseGlob(tmpTemplatesFolder + "/*.tmpl")
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"Bus":    bus,
		"dbName": dbName,
	}

	if err := hTmpl.ExecuteTemplate(g.hFile, "bus_h", data); err != nil {
		return err
	}

	// if err := hTmpl.ExecuteTemplate(g.hFile, "bus_h", map[string]interface{}{
	// 	"Bus": bus,
	// 	"dbName": dbName,
	// }); err != nil {
	// 	return err
	// }

	if err := cTmpl.ExecuteTemplate(g.cFile, "bus_c", bus); err != nil {
		return err
	}

	// if err := cTmpl.ExecuteTemplate(g.cFile, "bus_c", map[string]interface{}{
	// 	"Bus":    bus,
	// 	"dbName": dbName,
	// }); err != nil {
	// 	return err
	// }
	

	return nil
}
