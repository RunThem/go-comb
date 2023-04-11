package comb

import (
	"fmt"
	"strings"

	vec "github.com/RunThem/u.go/vector"
)

type Ctag int

const (
	_ Ctag = iota
	C_MATCH
	C_AND
	C_OR
	C_MAYBE
	C_MAYB1
	C_COMB
)

type Comb struct {
	tag     Ctag
	comment string

	match   string
	forward *vec.Vector[*Comb]
}

func New(tag Ctag, comment string) *Comb {
	return &Comb{tag: tag, comment: comment, forward: vec.New[*Comb](8)}
}

func Match(match string) *Comb {
	return &Comb{tag: C_MATCH, match: match, forward: vec.New[*Comb](0)}
}

func (c *Comb) Add(comb ...*Comb) {
	for i := 0; i < len(comb); i++ {
		c.forward.Push(comb[i])
	}
}

func (c *Comb) Dump() {
	history := make(map[string]bool, 10)

	var dump func(*Comb, int)
	dump = func(c *Comb, level int) {
		fmt.Printf("%s", strings.Repeat("    ", level))

		if len(c.comment) != 0 {
			if _, ok := history[c.comment]; ok == true {
				fmt.Printf("{%s}\n", c.comment)
				return
			}

			fmt.Printf("%s: ", c.comment)
			history[c.comment] = true
		}

		fmt.Printf("{%d, '%s'}", c.tag, c.match)

		if c.forward.Size() == 0 {
			fmt.Println("")
			return
		}

		fmt.Printf(" size %d\n", c.forward.Size())
		for i := 0; i < c.forward.Size(); i++ {
			dump(c.forward.At(i), level+1)
		}
	}

	dump(c, 0)
}

func match(c *Comb, in *Input) *Ast {
	idx := in.GetIdx()

	ch := in.Read(len(c.match))
	if c.match == ch {
		return NewAst(ch)
	}

	in.SetIdx(idx)
	return nil
}

func and(c *Comb, in *Input) *Ast {
	idx := in.GetIdx()
	a := NewAst("")

	for i := 0; i < c.forward.Size(); i++ {
		ast := c.forward.At(i).Parse(in)

		if ast == nil {
			in.SetIdx(idx)
			return nil
		}

		if c.forward.At(i).tag == C_MAYBE {
			for i := 0; i < ast.forward.Size(); i++ {
				a.Add(ast.forward.At(i))
			}
		} else {
			a.Add(ast)
		}
	}

	return a
}

func or(c *Comb, in *Input) *Ast {
	idx := in.GetIdx()

	for i := 0; i < c.forward.Size(); i++ {
		ast := c.forward.At(i).Parse(in)

		if ast != nil {
			return ast
		}

		in.SetIdx(idx)
	}

	return nil
}

func maybe(c *Comb, in *Input) *Ast {
	buf := make([]*Ast, 0, 8)
	a := NewAst("")

	for true {
		idx := in.GetIdx()
		i := 0
		for ; i < c.forward.Size(); i++ {
			ast := c.forward.At(i).Parse(in)
			if ast == nil {
				break
			}

			buf = append(buf, ast)
		}

		if i == c.forward.Size() {
			for j := 0; j < len(buf); j++ {
				a.forward.Push(buf[j])
			}

			buf = buf[:0]
		} else {
			in.SetIdx(idx)
			break
		}
	}

	return a
}

func (c *Comb) Parse(in *Input) *Ast {
	switch c.tag {
	case C_MATCH:
		return match(c, in)
	case C_AND:
		return and(c, in)
	case C_OR:
		return or(c, in)
	case C_MAYBE:
		return maybe(c, in)
	}

	return nil
}
