package ast

import (
	"fmt"
	vec "github.com/RunThem/u.go/vector"
	"strings"
)

type Ast struct {
	Code    string
	Forward *vec.Vector[*Ast]
}

func NewAst(code string) *Ast {
	return &Ast{Code: code, Forward: vec.New[*Ast](8)}
}

func (a *Ast) Add(ast ...*Ast) {
	a.Forward.Push(ast...)
}

func (a *Ast) Size() int {
	return a.Forward.Size()
}

func (a *Ast) At(idx int) *Ast {
	return a.Forward.At(idx)
}

func (a *Ast) Dump() {
	if a == nil {
		fmt.Println("ast is nil")
		return
	}

	var dump func(a *Ast, level int)
	dump = func(a *Ast, level int) {
		if a.Size() == 0 {
			fmt.Printf("%s{'%s'}\n", strings.Repeat("    ", level), a.Code)
			return
		} else {
			fmt.Printf("%s{size %d}\n", strings.Repeat("    ", level), a.Forward.Size())
		}

		for i := 0; i < a.Forward.Size(); i++ {
			dump(a.Forward.At(i), level+1)
		}
	}

	dump(a, 0)
}
