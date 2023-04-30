package comb

import (
	"fmt"
	"github.com/RunThem/comb/ast"
	"regexp"
	"strings"

	vec "github.com/RunThem/u.go/vector"
)

var table map[string]*Comb = make(map[string]*Comb, 10)

type Ctag int

const (
	_ Ctag = iota
	C_ANY
	C_MATCH
	C_AND
	C_OR
	C_MAYBE
	C_MAYB1
	C_MAYB0
	C_COMB
)

func (c Ctag) String() string {
	switch c {
	case C_ANY:
		return "C_ANY"
	case C_MATCH:
		return "C_MATCh"
	case C_AND:
		return "C_AND"
	case C_OR:
		return "C_OR"
	case C_MAYBE:
		return "C_MAYBE"
	case C_MAYB1:
		return "C_MAYB_1"
	case C_MAYB0:
		return "C_MAYB_0"
	case C_COMB:
		return "C_COMB"
	}

	return ""
}

type Comb struct {
	Tag     Ctag
	Comment string

	Match   string
	Froward *vec.Vector[*Comb]
}

func New(tag Ctag, comment string) *Comb {
	return &Comb{Tag: tag, Comment: comment, Froward: vec.New[*Comb](8)}
}

func Match(match string) *Comb {
	return &Comb{Tag: C_MATCH, Match: match}
}

func Any(match string) *Comb {
	return &Comb{Tag: C_ANY, Match: match}
}

func (c *Comb) Add(comb ...*Comb) {
	for i := 0; i < len(comb); i++ {
		c.Froward.Push(comb[i])
	}
}

func (c *Comb) Dump() {
	history := make(map[string]bool, 10)

	var dump func(*Comb, int)
	dump = func(c *Comb, level int) {
		fmt.Printf("%s", strings.Repeat("    ", level))

		if len(c.Comment) != 0 {
			if _, ok := history[c.Comment]; ok == true {
				fmt.Printf("{%s:%p}\n", c.Comment, c)
				return
			}

			fmt.Printf("%s:%p ", c.Comment, c)
			history[c.Comment] = true
		}

		fmt.Printf("{%s '%s'}", c.Tag, c.Match)

		if c.Froward == nil || c.Froward.Size() == 0 {
			fmt.Println("")
			return
		}

		fmt.Printf(" size %d\n", c.Froward.Size())
		for i := 0; i < c.Froward.Size(); i++ {
			dump(c.Froward.At(i), level+1)
		}
	}

	dump(c, 0)
}

func match(c *Comb, in *Input) *ast.Ast {
	fmt.Printf("%v\t\tmatch \"%s\"\n", string(in.input[in.idx:]), c.Match)
	idx := in.GetIdx()

	ch := in.Read(len(c.Match))
	if c.Match == ch {
		return ast.NewAst(ch)
	}

	in.SetIdx(idx)
	return nil
}

func any(c *Comb, in *Input) *ast.Ast {
	fmt.Printf("%v\t\tany \"%s\"\n", string(in.input[in.idx:]), c.Match)
	idx := in.GetIdx()
	re, _ := regexp.Compile("^" + c.Match)

	result := re.Find(in.input[in.idx:])
	if len(result) == 0 {
		return nil
	}

	in.SetIdx(idx + len(result))
	return ast.NewAst(string(result))
}

func and(c *Comb, in *Input) *ast.Ast {
	idx := in.GetIdx()
	a := ast.NewAst("")

	for i := 0; i < c.Froward.Size() && in.idx != len(in.input); i++ {
		fmt.Printf("%v\t\tand %d\n", string(in.input[in.idx:]), i)
		comb := c.Froward.At(i)
		ast := comb.Parse(in)

		if ast == nil {
			if comb.Tag == C_MAYB0 || comb.Tag == C_MAYB1 {
				continue
			}

			in.SetIdx(idx)
			return nil
		}

		if comb.Tag == C_MAYBE || comb.Tag == C_MAYB1 {
			for i := 0; i < ast.Size(); i++ {
				a.Add(ast.At(i))
			}
		} else {
			a.Add(ast)
		}
	}

	return a
}

func or(c *Comb, in *Input) *ast.Ast {
	idx := in.GetIdx()

	for i := 0; i < c.Froward.Size() && in.idx != len(in.input); i++ {
		fmt.Printf("%v\t\tor %d\n", string(in.input[in.idx:]), i)
		ast := c.Froward.At(i).Parse(in)

		if ast != nil {
			return ast
		}

		in.SetIdx(idx)
	}

	return nil
}

func maybe(c *Comb, in *Input) *ast.Ast {
	buf := make([]*ast.Ast, 0, 8)
	a := ast.NewAst("")

	for in.idx != len(in.input) {
		fmt.Printf("%v\t\tmaybe\n", string(in.input[in.idx:]))
		idx := in.GetIdx()
		i := 0
		for ; i < c.Froward.Size(); i++ {
			fmt.Printf("%v\t\tmaybe %d\n", string(in.input[in.idx:]), i)
			ast := c.Froward.At(i).Parse(in)
			if ast == nil {
				break
			}

			buf = append(buf, ast)
		}

		if i == c.Froward.Size() {
			for j := 0; j < len(buf); j++ {
				a.Add(buf[j])
			}

			buf = buf[:0]
		} else {
			in.SetIdx(idx)
			break
		}
	}

	return a
}

func mayb1(c *Comb, in *Input) *ast.Ast {
	buf := make([]*ast.Ast, 0, 8)
	a := ast.NewAst("")
	one := false

	for in.idx != len(in.input) {
		fmt.Printf("%v\t\tmayb1\n", string(in.input[in.idx:]))
		idx := in.GetIdx()
		i := 0
		for ; i < c.Froward.Size(); i++ {
			fmt.Printf("%v\t\tmayb1 %d\n", string(in.input[in.idx:]), i)
			ast := c.Froward.At(i).Parse(in)
			if ast == nil {
				break
			}

			buf = append(buf, ast)
		}

		if i == c.Froward.Size() {
			one = true
			for j := 0; j < len(buf); j++ {
				a.Add(buf[j])
			}

			buf = buf[:0]
		} else {
			in.SetIdx(idx)
			break
		}
	}

	if one {
		return a
	}

	return nil
}

func (c *Comb) Parse(in *Input) *ast.Ast {
	switch c.Tag {
	case C_ANY:
		return any(c, in)
	case C_MATCH:
		return match(c, in)
	case C_AND:
		return and(c, in)
	case C_OR:
		return or(c, in)
	case C_MAYBE:
		return maybe(c, in)
	case C_MAYB1:
		return mayb1(c, in)
	}

	return nil
}

var env = `
<factor> := /[0-9]*/ | '(' <expr> ')';
<term>   := <factor> (('*' | '/') <factor>)*;
<expr>   := <term> (('+' | '-') <term>)*;

( -> ?
[ -> *
{ -> +

factor := /[0-9]*/ | '(' <expr> ')';
term   := $factor [<'*' | '/'> $factor];
expr   := $term [<'+' | '-'> $term];
`

func CombL(text string) *Comb {
	table["factor"] = New(C_AND, "factor")
	table["term"] = New(C_AND, "term")
	table["expr"] = New(C_AND, "expr")

	return nil
}
