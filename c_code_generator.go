package acmelib

import (
	"io"
	"math"

	"fmt"
	"strings"
	"strconv"
	"text/template"
	"bytes"
)

const tmpTemplatesFolder = "../templates"

type Segment struct {
	Index        int
	Shift        int
	ShiftDir     string
	Mask         uint8
}

// define DB name
var dbName = "new_simple"
var fileName = "new_simple_out"

const packLeftShiftFmt = `
static inline uint8_t pack_left_shift_u{{ .Length }}(
    {{ .VarType }} value,
    uint8_t shift,
    uint8_t mask)
{
    return (uint8_t)((uint8_t)(value << shift) & mask);
}
`

const packRightShiftFmt = `
static inline uint8_t pack_right_shift_u{{ .Length }}(
    {{ .VarType }} value,
    uint8_t shift,
    uint8_t mask)
{
    return (uint8_t)((uint8_t)(value >> shift) & mask);
}
`

const unpackLeftShiftFmt = `
static inline {{ .VarType }} unpack_left_shift_u{{ .Length }}(
    uint8_t value,
    uint8_t shift,
    uint8_t mask)
{
    return ({{ .VarType }})(({{ .VarType }})(value & mask) << shift);
}
`

const unpackRightShiftFmt = `
static inline {{ .VarType }} unpack_right_shift_u{{ .Length }}(
    uint8_t value,
    uint8_t shift,
    uint8_t mask)
{
    return ({{ .VarType }})(({{ .VarType }})(value & mask) >> shift);
}
`

type HelperKind struct {
	Length  int
	VarType string
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

	kinds := []HelperKind{
		{Length: 8, VarType: "uint8_t"},
		{Length: 16, VarType: "uint16_t"},
	}

	packHelpers, unpackHelpers := generatePackUnpackHelpers(kinds)

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
		"getLenByte": getLenByte,
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
		"getMask": func (signalSize int) int {
			return ((1 << signalSize) - 1)
		},
		"getByteIndex": func(startBit int) int {
			return startBit / 8
		},
		"segments": segments,
		"isInRange": isInRange,
		"generateEncoding": GenerateEncoding,
		"generateDecoding": GenerateDecoding,
		"generatePackHelpers": func() []string {
			return packHelpers
		},
		"generateUnpackHelpers": func() []string {
			return unpackHelpers
		},
		"ExtractSignalsFromMux": ExtractSignalsFromMux,
		"GenerateSignalStruct": GenerateSignalStruct,
		"GenerateEncodingDeclaration": GenerateEncodingDeclaration,
		"GenerateDecodingDeclaration": GenerateDecodingDeclaration,
		"GenerateIsInRangeDeclaration": GenerateIsInRangeDeclaration,
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
		"packHelpers":   packHelpers,
		"unpackHelpers": unpackHelpers, 
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

func getLenByte(len int) int {
	if len <= 8 {
		return 8
	} else if len <= 16 {
		return 16
	} else if len <= 32 {
		return 32
	} else {
		return 64
	}
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

func isInRange(min interface{}, max interface{}, size int, isSigned int) string {
	var minStr, maxStr string
	var minimum, maximum float64
	var isFloat bool

	switch min.(type) {
	case float64:
		minimum = min.(float64)
		maximum = max.(float64)
		isFloat = true
	case int:
		minimum = float64(min.(int))
		maximum = float64(max.(int))
		isFloat = false
	default:
		return "true"
	}

	check := []string{}

	if minimum == 0 && maximum == 0 {
		return "true"
	}

	if !math.IsNaN(minimum) && minimum != 0 {
		if !isFloat {
			minimum = math.Round(minimum)
		}
		minStr = strconv.FormatFloat(minimum, 'f', -1, 64)
		check = append(check, fmt.Sprintf("(value >= %su)", minStr))
	}

	if !math.IsNaN(maximum) && maximum != 0 {
		if !isFloat {
			maximum = math.Round(maximum)
		}
		maxStr = strconv.FormatFloat(maximum, 'f', -1, 64)
		check = append(check, fmt.Sprintf("(value <= %su)", maxStr))
	}

	if len(check) == 0 {
		return "true"
	}

	if len(check) == 1 {
		return check[0][1 : len(check[0])-1]
	}

	return strings.Join(check, " && ")
}

func GenerateEncoding(scale, offset float64, useFloat bool) string {
	scaleLiteral := formatLiteral(scale, useFloat)
	offsetLiteral := formatLiteral(offset, useFloat)

	if offset == 0 && scale == 1 {
		return "value"
	} else if offset != 0 && scale != 1 {
		return fmt.Sprintf("(value - %s) / %s", offsetLiteral, scaleLiteral)
	} else if offset != 0 {
		return fmt.Sprintf("value - %s", offsetLiteral)
	} else {
		return fmt.Sprintf("value / %s", scaleLiteral)
	}
}

func formatLiteral(value float64, useFloat bool) string {
	strValue := strconv.FormatFloat(value, 'f', -1, 64)
	if useFloat {
		if !strings.Contains(strValue, ".") {
			strValue += ".0"
		}
	}
	return strValue
}

func GenerateDecoding(scale, offset float64, useFloat bool) string {
	floatingPointType := getFloatingPointType(useFloat)

	scaleLiteral := formatLiteral(scale, useFloat)
	offsetLiteral := formatLiteral(offset, useFloat)

	if offset == 0 && scale == 1 {
		return fmt.Sprintf("(%s)value", floatingPointType)
	} else if offset != 0 && scale != 1 {
		return fmt.Sprintf("((%s)value * %s) + %s", floatingPointType, scaleLiteral, offsetLiteral)
	} else if offset != 0 {
		return fmt.Sprintf("(%s)value + %s", floatingPointType, offsetLiteral)
	} else {
		return fmt.Sprintf("(%s)value * %s", floatingPointType, scaleLiteral)
	}
}

func getFloatingPointType(useFloat bool) string {
	if useFloat {
		return "float"
	}
	return "double"
}

func generatePackUnpackHelpers(kinds []HelperKind) ([]string, []string) {
	packHelpers := generateHelpers(kinds, packLeftShiftFmt, packRightShiftFmt)
	unpackHelpers := generateHelpers(kinds, unpackLeftShiftFmt, unpackRightShiftFmt)
	return packHelpers, unpackHelpers
}

func generateHelpers(kinds []HelperKind, leftFormat, rightFormat string) []string {
	var helpers []string
	tmplLeft := template.Must(template.New("left").Parse(leftFormat))
	tmplRight := template.Must(template.New("right").Parse(rightFormat))

	for _, kind := range kinds {
		var helper string
		var buf bytes.Buffer

		tmplLeft.Execute(&buf, kind)
		tmplRight.Execute(&buf, kind)

		helper = buf.String()
		helpers = append(helpers, helper)
	}
	return helpers
}

func ExtractSignalsFromMux(signalGroups [][]Signal) []Signal {
	var res []Signal
	for _, group := range signalGroups {
		for _, s := range group {
			if s.Kind() == SignalKindMultiplexer {
				mux, err := s.ToMultiplexer()
				if err != nil {
					return nil
				}
				res = append(res, ExtractSignalsFromMux(mux.GetSignalGroups())...)
			} else {
				res = append(res, s)
			}
		}
	}
	return res
}

func GenerateSignalStruct(signal Signal) string {
	var res, rangeStr, scaleStr, offsetStr string

	switch signal.Kind() {
	case SignalKindStandard:
		standardSignal, err := signal.ToStandard()
		if err != nil {
			return ""
		}
		rangeStr = formatRange(standardSignal.Type().Min(), standardSignal.Type().Max(), standardSignal.Type().Offset(), standardSignal.Type().Scale())
		scaleStr = fmt.Sprintf("%v", standardSignal.Type().Scale())
		offsetStr = fmt.Sprintf("%v", standardSignal.Type().Offset())
		res = fmt.Sprintf("\t/**\n\t * Range: %s\n\t * Scale: %s\n\t * Offset: %s\n\t */\n\t%s%d_t %s;", rangeStr, scaleStr, offsetStr, isSignedType(standardSignal.Type().Signed()), getLenByte(standardSignal.Type().Size()), strings.ToLower(standardSignal.Name()))
	case SignalKindEnum: 
		enumSignal, err := signal.ToEnum()
		if err != nil {
			return ""
		}
		rangeStr = formatRange(enumSignal.Enum().Values()[0].Index(), enumSignal.Enum().Values()[len(enumSignal.Enum().Values())-1].Index(), 0, 1)
		scaleStr = "1"
		offsetStr = "0"
		res = fmt.Sprintf("\t/**\n\t * Range: %s\n\t * Scale: %s\n\t * Offset: %s\n\t */\n\t%s%d_t %s;", rangeStr, scaleStr, offsetStr, isEnumSigned(enumSignal.Enum().Values()), getLenByte(enumSignal.GetSize()), strings.ToLower(enumSignal.Name()))
	case SignalKindMultiplexer:
		muxSignal, err := signal.ToMultiplexer()
		if err != nil {
			return ""
		}
		rangeStr = "-"
		scaleStr = "1"
		offsetStr = "0"
		res = fmt.Sprintf("\t/**\n\t * Range: %s\n\t * Scale: %s\n\t * Offset: %s\n\t */\n\tuint8_t %s;", rangeStr, scaleStr, offsetStr, strings.ToLower(muxSignal.Name()))
		for _, s := range ExtractSignalsFromMux(muxSignal.GetSignalGroups()) {
			res += "\n" + GenerateSignalStruct(s)
		}
	}

	return res
}

func GenerateEncodingDeclaration(signal Signal, messageName string) string {
	var res = "/**\n * Encode given signal by applying scaling and offset.\n *\n * @param[in] value Signal to encode.\n *\n * @return Encoded signal.\n*/\n"
	
	switch signal.Kind() {
	case SignalKindStandard:
		standardSignal, err := signal.ToStandard()
		if err != nil {
			return ""
		}
		res += fmt.Sprintf("%s%d_t %s_%s_%s_encode(double value);\n", isSignedType(standardSignal.Type().Signed()), getLenByte(standardSignal.GetSize()), dbName, messageName, strings.ToLower(standardSignal.Name()))
	case SignalKindEnum:
		enumSignal, err := signal.ToEnum()
		if err != nil {
			return ""
		}
		res += fmt.Sprintf("%s%d_t %s_%s_%s_encode(double value);\n", isEnumSigned(enumSignal.Enum().Values()), getLenByte(enumSignal.GetSize()), dbName, messageName, strings.ToLower(enumSignal.Name()))
	case SignalKindMultiplexer:
		muxSignal, err := signal.ToMultiplexer()
		if err != nil {
			return ""
		}
		res += fmt.Sprintf("uint8_t %s_%s_%s_encode(double value);\n", dbName, messageName, strings.ToLower(muxSignal.Name()))
		for _, s := range ExtractSignalsFromMux(muxSignal.GetSignalGroups()) {
			res += "\n" + GenerateEncodingDeclaration(s, messageName)
		}
	}
	
	return res
}

func GenerateDecodingDeclaration(signal Signal, messageName string) string {
	var res = "/**\n * Decode given signal by applying scaling and offset.\n *\n * @param[in] value Signal to decode.\n *\n * @return Decoded signal.\n*/\n"
	
	switch signal.Kind() {
	case SignalKindStandard:
		standardSignal, err := signal.ToStandard()
		if err != nil {
			return ""
		}
		res += fmt.Sprintf("double %s_%s_%s_decode(%s%d_t value);\n", dbName, messageName, strings.ToLower(standardSignal.Name()), isSignedType(standardSignal.Type().Signed()), getLenByte(standardSignal.GetSize()))
	case SignalKindEnum:
		enumSignal, err := signal.ToEnum()
		if err != nil {
			return ""
		}
		res += fmt.Sprintf("double %s_%s_%s_decode(%s%d_t value);\n", dbName, messageName, strings.ToLower(enumSignal.Name()), isEnumSigned(enumSignal.Enum().Values()), getLenByte(enumSignal.GetSize()))
	case SignalKindMultiplexer:
		muxSignal, err := signal.ToMultiplexer()
		if err != nil {
			return ""
		}
		res += fmt.Sprintf("double %s_%s_%s_decode(uint8_t value);\n", dbName, messageName, strings.ToLower(muxSignal.Name()))
		for _, s := range ExtractSignalsFromMux(muxSignal.GetSignalGroups()) {
			res += "\n" + GenerateDecodingDeclaration(s, messageName)
		}
	}
	
	return res
}

func GenerateIsInRangeDeclaration(signal Signal, messageName string) string {
	var res = "/**\n * Check that given signal is in allowed range.\n *\n * @param[in] value Signal to check.\n *\n * @return true if in range, false otherwise.\n*/\n"
	
	switch signal.Kind() {
	case SignalKindStandard:
		standardSignal, err := signal.ToStandard()
		if err != nil {
			return ""
		}
		res += fmt.Sprintf("bool %s_%s_%s_is_in_range(%s%d_t value);\n", dbName, messageName, strings.ToLower(standardSignal.Name()), isSignedType(standardSignal.Type().Signed()), getLenByte(standardSignal.GetSize()))
	case SignalKindEnum:
		enumSignal, err := signal.ToEnum()
		if err != nil {
			return ""
		}
		res += fmt.Sprintf("bool %s_%s_%s_is_in_range(%s%d_t value);\n", dbName, messageName, strings.ToLower(enumSignal.Name()), isEnumSigned(enumSignal.Enum().Values()), getLenByte(enumSignal.GetSize()))
	case SignalKindMultiplexer:
		muxSignal, err := signal.ToMultiplexer()
		if err != nil {
			return ""
		}
		res += fmt.Sprintf("bool %s_%s_%s_is_in_range(uint8_t value);\n", dbName, messageName, strings.ToLower(muxSignal.Name()))
		for _, s := range ExtractSignalsFromMux(muxSignal.GetSignalGroups()) {
			res += "\n" + GenerateIsInRangeDeclaration(s, messageName)
		}
	}
	
	return res
}