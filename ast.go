package comb

import (
	"fmt"
	vec "github.com/RunThem/u.go/vector"
	"strings"
)

type Ast struct {
	code    string
	forward *vec.Vector[*Ast]
}

func NewAst(code string) *Ast {
	return &Ast{code: code, forward: vec.New[*Ast](8)}
}

func (a *Ast) Add(ast ...*Ast) {
	a.forward.Push(ast...)
}

func (a *Ast) Size() int {
	return a.forward.Size()
}

func (a *Ast) Dump() {
	if a == nil {
		fmt.Println("ast is nil")
		return
	}

	var dump func(a *Ast, level int)
	dump = func(a *Ast, level int) {
		if a.Size() == 0 {
			fmt.Printf("%s{'%s'}\n", strings.Repeat("    ", level), a.code)
			return
		} else {
			fmt.Printf("%s{size %d}\n", strings.Repeat("    ", level), a.forward.Size())
		}

		for i := 0; i < a.forward.Size(); i++ {
			dump(a.forward.At(i), level+1)
		}
	}

	dump(a, 0)
}
