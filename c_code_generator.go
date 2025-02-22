package acmelib

import (
	"io"
	"math"

	"fmt"
	"strings"
	"text/template"
)

const tmpTemplatesFolder = "../templates"

// type Signal struct {
// 	StartBit  int
// 	Size      int
// 	ByteOrder string
// }

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
	dbName := "expected"

	funcMap := template.FuncMap{
		"toUpper": strings.ToUpper,
		"toLower": strings.ToLower,
		"camelToSnake": func(s string) string {
			var res string
			for idx, ch := range s {
				if idx > 0 && ch >= 'A' && ch <= 'Z' {
					// Check if the previous character is lowercase OR the next character is lowercase
					if (s[idx-1] >= 'a' && s[idx-1] <= 'z') || (idx+1 < len(s) && s[idx+1] >= 'a' && s[idx+1] <= 'z') {
						res += "_"
					}
				}
				// Convert uppercase to lowercase
				if ch >= 'A' && ch <= 'Z' {
					res += string(ch + 'a' - 'A')
				} else {
					res += string(ch)
				}
			}
			return string(res)
		},
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
		"getLen": func(len int) int {
			if len <= 8 {
				return 8
			} else if len <= 16 {
				return 16
			} else if len <= 32 {
				return 32
			} else {
				return 64
			}
		},
		"segments": segments,
		"formatRange": formatRange,
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

	if err := cTmpl.ExecuteTemplate(g.cFile, "bus_c", bus); err != nil {
		return err
	}	

	return nil
}

func formatRange(min interface{}, max interface{}, offset interface {}, scale interface{}) string {
	var minStr, maxStr string
	var minimum, maximum interface{}
	var isFloat bool

	switch min.(type) {
	case float64:
		minimum = min.(float64)
		maximum = max.(float64)
		isFloat = true
	case int:
		minimum = min.(int)
		maximum = max.(int)
		isFloat = false
	default:
		return "-"
	}

	physToRaw := func(x interface{}, isFloat bool) interface{} {
		if isFloat {
			return (x.(float64)-offset.(float64))/scale.(float64)
		}
		return math.Round(float64(x.(int)-offset.(int))/float64(scale.(int)))
	}

	if minimum == 0.0 && maximum == 0.0 {
		return "-"
	} else if minimum != nil && maximum != nil {
		minStr = fmt.Sprintf("%v", physToRaw(minimum, isFloat))
		maxStr = fmt.Sprintf("%v", physToRaw(maximum, isFloat))
		if isFloat {
			return fmt.Sprintf("%s..%s (%.5f..%.5f -)", minStr, maxStr, minimum, maximum)
		}
		return fmt.Sprintf("%s..%s (%d..%d -)", minStr, maxStr, minimum, maximum)
	} else if minimum != nil {
		minStr = fmt.Sprintf("%v", physToRaw(minimum, isFloat))
		if isFloat {
			return fmt.Sprintf("%s.. (%.5f.. -)", minStr, minimum)
		}
		return fmt.Sprintf("%s.. (%d.. -)", minStr, minimum)
	} else if maximum != nil {
		maxStr = fmt.Sprintf("%v", physToRaw(maximum, isFloat))
		if isFloat {
			return fmt.Sprintf("..%s (..%.5f -)", maxStr, maximum)
		}
		return fmt.Sprintf("..%s (..%d -)", maxStr, maximum)
	} else {
		return "-"
	}
}

func segments(signal Signal, invertShift bool) []struct {
	Index         int
	Shift         int
	ShiftDir      string
	Mask          int
} {
	var result []struct {
		Index    int
		Shift    int
		ShiftDir string
		Mask     int
	}

	index, pos := signal.GetStartBit()/8, signal.GetStartBit()%8
	left := signal.GetSize()

	for left > 0 {
		var length, shift, mask int
		if signal.ParentMessage().ByteOrder().String() == "big_endian" {
			if left >= pos+1 {
				length = pos + 1
				pos = 7
				shift = -(left - length)
				mask = (1 << length) - 1
			} else {
				length = left
				shift = pos - length + 1
				mask = ((1 << length) - 1) << (pos - length + 1)
			}
		} else {
			shift = left - signal.GetSize() + pos
			if left >= 8-pos {
				length = 8 - pos
				mask = ((1 << length) - 1) << pos
				pos = 0
			} else {
				length = left
				mask = ((1 << length) - 1) << pos
			}
		}

		shiftDirection := "left"
		if invertShift {
			if shift < 0 {
				shift = -shift
				shiftDirection = "left"
			} else {
				shiftDirection = "right"
			}
		} else {
			if shift < 0 {
				shift = -shift
				shiftDirection = "right"
			} else {
				shiftDirection = "left"
			}
		}

		result = append(result, struct {
			Index    int
			Shift    int
			ShiftDir string
			Mask     int
		}{
			Index:    index,
			Shift:    shift,
			ShiftDir: shiftDirection,
			Mask:     mask,
		})

		left -= length
		index++
	}

	return result
}

