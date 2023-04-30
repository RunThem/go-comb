package lexer

import (
	"bytes"
	"container/list"
	"github.com/RunThem/comb/token"
	u "github.com/RunThem/u.go"
	vec "github.com/RunThem/u.go/vector"
)

type Lexer struct {
	code string
	idx  int
	Lex  *vec.Vector[*list.List]
}

func NewLexer(code string) *Lexer {
	l := &Lexer{code: code, Lex: vec.New[*list.List](10)}

	lst := list.New()
	for l.idx < len(l.code) {
		ch := l.read()

		if u.IsSpace(ch) {
			continue
		}

		switch ch {
		case ';':
			l.Lex.Push(lst)
			lst = list.New()
		case '=':
			lst.PushBack(token.NewToken(token.TOK_DEFINE, "="))
		case '\'':
			lst.PushBack(l.match())
		case '/':
			lst.PushBack(l.regex())
		case '$':
			lst.PushBack(l.ident())
		default:
			lst.PushBack(token.NewToken(string(ch), string(ch)))
		}
	}

	return l
}

func (l *Lexer) match() *token.Token {
	var match bytes.Buffer

	for true {
		ch := l.read()

		if ch == '\\' {
			switch l.peek(1) {
			case '\'':
				match.WriteByte('\'')
			case 'n':
				match.WriteByte('\n')
			case 'r':
				match.WriteByte('\r')
			case 't':
				match.WriteByte('\t')
			default:
				panic("\\?")
			}

			l.read()
		} else if ch == '\'' {
			break
		}

		match.WriteByte(ch)
	}

	return token.NewToken(token.TOK_MATCH, match.String())
}

func (l *Lexer) regex() *token.Token {
	var match bytes.Buffer

	for ch := l.peek(0); ch != '/'; ch = l.peek(0) {
		match.WriteByte(l.read())
	}

	l.read()
	return token.NewToken(token.TOK_REGEX, match.String())
}

func (l *Lexer) ident() *token.Token {
	var match bytes.Buffer

	for ch := l.peek(0); u.IsAlnum(ch) || ch == '_'; ch = l.peek(0) {
		match.WriteByte(l.read())
	}

	return token.NewToken(token.TOK_IDENT, match.String())
}

func (l *Lexer) read() byte {
	l.idx++
	return l.code[l.idx-1]
}

func (l *Lexer) peek(count int) byte {
	return l.code[l.idx+count]
}
