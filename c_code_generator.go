package acmelib

import (
	"io"
	"math"

	"fmt"
	"strings"
	"text/template"
)

const tmpTemplatesFolder = "../templates"

type Segment struct {
	Index        int
	Shift        int
	ShiftDir     string
	Mask         uint8
}

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
	fileName := "test"

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
		"getLenByte": func(len int) int {
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
		"formatRange": formatRange,
		"sub": func(a, b int) int {
			return a - b
		},
		"isSignedType": isSignedType,
		"isEnumSigned": isEnumSigned,
		"getIntType":   getIntType,
		"div": func (a, b int) int {
			if b == 0 {
				panic("divided by zero")
			}
			return a / b
		},
		"add": func (a, b int) int {
			return a + b
		},
		"mod": func (a, b int) int {
			if b == 0 {
				panic("mod by zero")
			}
			return a % b
		},
		"hexMap": func(mask interface{}) string {
    		return fmt.Sprintf("0x%02x", mask) + "u"
		},
		"getMask": getMask,
		"getByteIndex": func(startBit int) int {
			return startBit / 8
		},
		"segments": segments,
	}	

	hTmpl, err := template.New("c_header").Funcs(funcMap).ParseGlob(tmpTemplatesFolder + "/*.gtpl")
	if err != nil {
		return err
	}

	cTmpl, err := template.New("c_source").Funcs(funcMap).ParseGlob(tmpTemplatesFolder + "/*.gtpl")
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"Bus": bus,
		"dbName": dbName,
		"fileName": fileName,
	}

	if err := hTmpl.ExecuteTemplate(g.hFile, "bus_h", data); err != nil {
		return err
	}

	if err := cTmpl.ExecuteTemplate(g.cFile, "bus_c", data); err != nil {
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

func isSignedType(isSigned bool) string {
	if isSigned {
		return "int"
	}
	return "uint"
}

func isEnumSigned(enumValues []*SignalEnumValue) string {
	if len(enumValues) == 0 {
		return "uint"
	}
	firstIndex := enumValues[0].Index()
	if firstIndex < 0 {
		return "int"
	}
	return "uint"
}

func getIntType(kind string, isSigned bool, enumValues []*SignalEnumValue) string {
	if kind == "enum" {
		return isEnumSigned(enumValues)
	} else if kind == "standard" {
		return isSignedType(isSigned)
	}
	return "uint"
}

func segments(startBit, length int) []Segment {
	remaining := length
    index := startBit / 8
    pos := startBit % 8
    var result []Segment

    for remaining > 0 {
        var segment Segment
        segment.Index = index

        bitsInByte := min(8-pos, remaining)
        segment.Mask = ((1 << bitsInByte) - 1) << pos

        segment.Shift = (index - startBit/8) * 8
        segment.ShiftDir = "left"

        result = append(result, segment)

        remaining -= bitsInByte
        index++
        pos = 0
    }

	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getMask(signalSize int) int {
	return ((1 << signalSize) - 1)
}