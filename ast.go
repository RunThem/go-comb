package comb

import (
	"fmt"
	"strings"

	vec "github.com/RunThem/u.go/vector"
)

type Ast struct {
	code    string
	forward *vec.Vector[*Ast]
}

func NewAst(code string) *Ast {
	return &Ast{code: code, forward: vec.New[*Ast](8)}
}

func (a *Ast) Add(ast *Ast) {
	a.forward.Push(ast)
}

func (a *Ast) Dump() {
	if a == nil {
		fmt.Println("ast is nil")
		return
	}

	var dump func(a *Ast, level int)
	dump = func(a *Ast, level int) {
		if len(a.code) != 0 {
			fmt.Printf("%s{'%s'}\n", strings.Repeat("    ", level), a.code)
			return
		}

		for i := 0; i < a.forward.Size(); i++ {
			dump(a.forward.At(i), level+1)
		}
	}

	dump(a, 0)
}
