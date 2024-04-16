package dbc

type keywordKind uint

const (
	keywordVersion keywordKind = iota

	keywordNewSymbols

	keywordBitTiming

	keywordNode

	keywordMessage
	keywordMessageTransmitter

	keywordSignal
	keywordSignalValueType

	keywordValueTable
	keywordValueEncoding

	keywordEnvVar
	keywordEnvVarData

	keywordSignalType
	keywordSignalGroup

	keywordComment

	keywordAttribute
	keywordAttributeDefault
	keywordAttributeValue
	keywordAttributeInt
	keywordAttributeHex
	keywordAttributeFloat
	keywordAttributeString
	keywordAttributeEnum

	keywordExtendedMux
)

var keywords = map[string]keywordKind{
	"VERSION": keywordVersion,

	"NS_": keywordNewSymbols,

	"BS_": keywordBitTiming,

	"BU_": keywordNode,

	"BO_":       keywordMessage,
	"BO_TX_BU_": keywordMessageTransmitter,

	"SG_":          keywordSignal,
	"SIG_VALTYPE_": keywordSignalValueType,

	"VAL_TABLE_": keywordValueTable,
	"VAL_":       keywordValueEncoding,

	"EV_":          keywordEnvVar,
	"ENVVAR_DATA_": keywordEnvVarData,

	"SGTYPE_":    keywordSignalType,
	"SIG_GROUP_": keywordSignalGroup,

	"CM_": keywordComment,

	"BA_DEF_":     keywordAttribute,
	"BA_DEF_DEF_": keywordAttributeDefault,
	"BA_":         keywordAttributeValue,
	"INT":         keywordAttributeInt,
	"HEX":         keywordAttributeHex,
	"FLOAT":       keywordAttributeFloat,
	"STRING":      keywordAttributeString,
	"ENUM":        keywordAttributeEnum,

	"SG_MUL_VAL_": keywordExtendedMux,
}

func getKeywordKind(s string) keywordKind {
	return keywords[s]
}

func getKeyword(kind keywordKind) string {
	for str, k := range keywords {
		if k == kind {
			return str
		}
	}
	return ""
}

var newSymbolsValues = []string{
	"NS_DESC_",
	"CM_",
	"BA_DEF_",
	"BA_",
	"VAL_",
	"VAL_TABLE_",
	"CAT_DEF_",
	"CAT_",
	"FILTER",
	"BA_DEF_DEF_",
	"EV_DATA_",
	"ENVVAR_DATA_",
	"SIG_GROUP_",
	"SGTYPE_",
	"SGTYPE_VAL_",
	"BA_DEF_SGTYPE_",
	"BA_SGTYPE_",
	"SIG_TYPE_REF_",
	"SIG_VALTYPE_",
	"SIGTYPE_VALTYPE_",
	"BO_TX_BU_",
	"BA_DEF_REL_",
	"BA_REL_",
	"BA_DEF_DEF_REL_",
	"BU_SG_REL_",
	"BU_EV_REL_",
	"BU_BO_REL_",
	"SG_MUL_VAL_",
}

var envVarAccessTypes = map[string]EnvVarAccessType{
	"DUMMY_NODE_VECTOR0":    EnvVarDummyNodeVector0,
	"DUMMY_NODE_VECTOR1":    EnvVarDummyNodeVector1,
	"DUMMY_NODE_VECTOR2":    EnvVarDummyNodeVector2,
	"DUMMY_NODE_VECTOR3":    EnvVarDummyNodeVector3,
	"DUMMY_NODE_VECTOR8000": EnvVarDummyNodeVector8000,
	"DUMMY_NODE_VECTOR8001": EnvVarDummyNodeVector8001,
	"DUMMY_NODE_VECTOR8002": EnvVarDummyNodeVector8002,
	"DUMMY_NODE_VECTOR8003": EnvVarDummyNodeVector8003,
}
