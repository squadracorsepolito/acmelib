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
	tokenSyntax
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
	tokenSyntax:       "syntax",
}

type token struct {
	kind     tokenKind
	kindName string
	value    string
	start    int
	col      int
	line     int
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

func (t *token) isSyntax(s syntaxKind) bool {
	if t.kind != tokenSyntax {
		return false
	}
	return getSyntaxKind(t.value) == s
}
