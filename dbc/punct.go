package dbc

type punctKind uint

const (
	punctColon punctKind = iota
	punctComma
	punctLeftParen
	punctRightParen
	punctLeftSquareBrace
	punctRightSquareBrace
	punctPipe
	punctSemicolon
	punctAt
	punctPlus
	punctMinus
)

var punctKeywords = map[rune]punctKind{
	':': punctColon,
	',': punctComma,
	'(': punctLeftParen,
	')': punctRightParen,
	'[': punctLeftSquareBrace,
	']': punctRightSquareBrace,
	'|': punctPipe,
	';': punctSemicolon,
	'@': punctAt,
	'+': punctPlus,
	'-': punctMinus,
}

func isPunctKeyword(r rune) bool {
	_, ok := punctKeywords[r]
	return ok
}

func getPunctKind(str string) punctKind {
	return punctKeywords[rune(str[0])]
}

func getPunctRune(kind punctKind) rune {
	var r rune
	for tmpR, k := range punctKeywords {
		if k == kind {
			r = tmpR
			break
		}
	}
	return r
}
