package token

const (
	TOK_IDENT = "IDENT"
	TOK_MATCH = "MATCH"
	TOK_REGEX = "REGEX"

	TOK_DEFINE = "="

	TOK_OR = "|"

	TOK_L_ANGLE_BRACK = "<"
	TOK_R_ANGLE_BRACK = ">"

	TOK_L_PARENT     = "("
	TOK_R_PARENT     = ")"
	TOK_L_BRACK      = "["
	TOK_R_BRACK      = "]"
	TOK_L_BIG_PARANT = "{"
	TOK_R_BIG_PARANT = "}"
)

type Token struct {
	Type    string
	Literal string
}

func NewToken(typ string, literal string) *Token {
	return &Token{
		Type:    typ,
		Literal: literal,
	}
}
