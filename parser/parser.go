package parser

import (
	"container/list"
	"fmt"
	"github.com/RunThem/comb/comb"
	"github.com/RunThem/comb/lexer"
	"github.com/RunThem/comb/token"
)

var table map[string]*comb.Comb = make(map[string]*comb.Comb, 10)

func NewParser(lex *lexer.Lexer) *comb.Comb {
	for i := 0; i < lex.Lex.Size(); i++ {
		tok := lex.Lex.At(i).Front().Value.(*token.Token)
		if tok.Type != token.TOK_IDENT {
			panic("front item is not ident")
		}

		table[tok.Literal] = comb.New(comb.C_AND, tok.Literal)
	}

	for i := 0; i < lex.Lex.Size(); i++ {
		parser(lex.Lex.At(i))
	}

	return table[lex.Lex.Back().Front().Value.(*token.Token).Literal]
}

func parser(lst *list.List) {
	cur_tok := lst.Front()
	root := table[typeof(cur_tok).Literal]

	cur_tok = cur_tok.Next()
	define := typeof(cur_tok)
	if define.Type != token.TOK_DEFINE {
		return
	}

	cur_tok = cur_tok.Next()

	stmt, cur_tok := parse_stmt(cur_tok)
	if stmt.Tag == comb.C_OR {
		root.Tag = comb.C_OR
	}
	//root.Add(stmt)
	root.Froward = stmt.Froward
}

func parse_stmt(element *list.Element) (*comb.Comb, *list.Element) {
	var stmt *comb.Comb
	st := comb.New(comb.C_AND, "")

	for element != nil {
		tok := typeof(element)
		fmt.Println(tok.Type, tok.Literal)

		switch tok.Type {
		case token.TOK_MATCH:
			st.Add(comb.Match(tok.Literal))
		case token.TOK_REGEX:
			st.Add(comb.Any(tok.Literal))
		case token.TOK_IDENT:
			st.Add(table[tok.Literal])
		case token.TOK_OR:
			if stmt == nil {
				stmt = comb.New(comb.C_OR, "")
			}

			if st.Froward.Size() == 1 {
				stmt.Add(st.Froward.At(0))
			} else {
				stmt.Add(st)
			}

			st = comb.New(comb.C_AND, "")
		case token.TOK_L_ANGLE_BRACK:
			tmp, elem := parse_stmt(element.Next())
			st.Add(tmp)

			if elem == nil {
				goto loop
			}

			element = elem
		case token.TOK_L_PARENT:
			tmp, elem := parse_stmt(element.Next())
			t := comb.New(comb.C_MAYB0, "")

			if tmp.Tag == comb.C_AND {
				t.Froward = tmp.Froward
			} else {
				t.Add(tmp)
			}

			st.Add(t)

			if elem == nil {
				goto loop
			}

			element = elem
		case token.TOK_L_BRACK:
			tmp, elem := parse_stmt(element.Next())
			t := comb.New(comb.C_MAYB1, "")
			if tmp.Tag == comb.C_AND {
				t.Froward = tmp.Froward
			} else {
				t.Add(tmp)
			}
			st.Add(t)

			if elem == nil {
				goto loop
			}

			element = elem
		case token.TOK_L_BIG_PARANT:
			tmp, elem := parse_stmt(element.Next())
			t := comb.New(comb.C_MAYBE, "")

			if tmp.Tag == comb.C_AND {
				t.Froward = tmp.Froward
			} else {
				t.Add(tmp)
			}

			st.Add(t)

			if elem == nil {
				goto loop
			}

			element = elem
		case token.TOK_R_ANGLE_BRACK, token.TOK_R_PARENT, token.TOK_R_BRACK, token.TOK_R_BIG_PARANT:
			goto loop
		}

		element = element.Next()
	}
loop:

	if stmt == nil {
		return st, element
	}

	if st.Froward.Size() == 1 {
		stmt.Add(st.Froward.At(0))
	} else {
		stmt.Add(st)
	}

	return stmt, element
}

func typeof(element *list.Element) *token.Token {
	return element.Value.(*token.Token)
}
