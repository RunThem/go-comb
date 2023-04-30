package lexer

import (
	"fmt"
	"testing"
)

func TestNewLexer(t *testing.T) {
	//	code := `
	//<factor> := /[0-9]*/ | '\'(' <expr> ')';
	//<term>   := <factor> (('*' | '/') <factor>)*;
	//<expr>   := <term> (('+' | '-') <term>)*;
	//
	//( -> ?
	//[ -> *
	//{ -> +
	//
	//factor = /[0-9]*/ | '(' <expr> ')';
	//term   = $factor [<'*' | '/'> $factor];
	//expr   = $term [<'+' | '-'> $term];
	//`

	code := `
$factor = /[0-9]*/ | '(' <expr> ')';
$term   = $factor [<'*' | '/'> $factor];
$expr   = $term [<'+' | '-'> $term];
`

	lexer := NewLexer(code)

	for i := 0; i < lexer.Lex.Size(); i++ {
		lst := lexer.Lex.At(i)

		for it := lst.Front(); it != nil; it = it.Next() {
			fmt.Println(it.Value)
		}

		fmt.Println()
	}
}
