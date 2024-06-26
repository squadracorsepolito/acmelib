package dbc

type tokenKind uint

const (
	tokenError tokenKind = iota
	tokenEOF
	tokenSpace

	tokenIdent
	tokenNumber
	tokenNumberRange
	tokenMuxIndicator
	tokenString
	tokenKeyword
	tokenPunct
)

var tokenNames = map[tokenKind]string{
	tokenError: "error",
	tokenEOF:   "eof",
	tokenSpace: "space",

	tokenIdent:        "ident",
	tokenNumber:       "number",
	tokenNumberRange:  "number_range",
	tokenMuxIndicator: "mux_indicator",
	tokenString:       "string",
	tokenKeyword:      "keyword",
	tokenPunct:        "punct",
}

type token struct {
	kind      tokenKind
	kindName  string
	value     string
	startLine int
	startCol  int
	endLine   int
	endCol    int
}

func (t *token) isEOF() bool {
	return t.kind == tokenEOF
}

func (t *token) isError() bool {
	return t.kind == tokenError
}

func (t *token) isSpace() bool {
	return t.kind == tokenSpace
}

func (t *token) isNumber() bool {
	return t.kind == tokenNumber
}

func (t *token) isNumberRange() bool {
	return t.kind == tokenNumberRange
}

func (t *token) isMuxIndicator() bool {
	return t.kind == tokenMuxIndicator
}

func (t *token) isIdent() bool {
	return t.kind == tokenIdent
}

func (t *token) isString() bool {
	return t.kind == tokenString
}

func (t *token) isKeyword(k keywordKind) bool {
	if t.kind != tokenKeyword {
		return false
	}
	return getKeywordKind(t.value) == k
}

func (t *token) isPunct(s punctKind) bool {
	if t.kind != tokenPunct {
		return false
	}
	return getPunctKind(t.value) == s
}
