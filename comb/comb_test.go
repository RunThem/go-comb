package comb

import (
	"fmt"
	"github.com/RunThem/comb/ast"
	"strconv"
	"testing"
)

func eval(ast *ast.Ast) int {
	if ast.Size() == 0 {
		n, _ := strconv.Atoi(ast.Code)
		return n
	}

	num := eval(ast.At(0))
	fmt.Println(num)
	for i := 1; i < ast.Size(); i++ {
		a := ast.At(i)
		switch a.Code {
		case "*":
			i += 1
			num *= eval(ast.At(i))
			fmt.Println(num)
		case "/":
			i += 1
			num /= eval(ast.At(i))
			fmt.Println(num)
		case "+":
			i += 1
			num += eval(ast.At(i))
			fmt.Println(num)
		case "-":
			i += 1
			num -= eval(ast.At(i))
			fmt.Println(num)
		case "(":
			i += 1
			return eval(ast.At(i))
		case ")":
			i += 1
		default:
			return eval(ast.At(i))
		}
	}

	return num
}
func TestMatch(t *testing.T) {

	/*
	 * Factor = '[0-9]*' | '(' <Expr> ')'
	 * Term   = <Factor> (('*' | '/') <Factor>)*
	 * Expr   = <Term> (('+' | '-') <Term>)*
	 */

	Expr := New(C_AND, "Expr")
	Term := New(C_AND, "Term")
	Factor := New(C_OR, "Factor")

	factor := New(C_AND, "factor")
	factor.Add(Match("("), Expr, Match(")"))

	Factor.Add(Any("[1-9]+"), factor)

	term := New(C_MAYBE, "term")

	op := New(C_OR, "term_op")
	op.Add(Match("*"), Match("/"))

	term.Add(op, Factor)

	Term.Add(Factor, term)

	expr := New(C_MAYBE, "expr")

	op = New(C_OR, "expr_op")
	op.Add(Match("+"), Match("-"))

	expr.Add(op, Term)

	Expr.Add(Term, expr)

	//in := NewInput("4234+9*(3*3-3)")
	//in := NewInput("(4234)")
	in := NewInput("4+4")
	ast := Expr.Parse(in)

	Expr.Dump()

	ast.Dump()

	fmt.Println(eval(ast))
}

const (
	_ int = iota
	LOWEST
	SUM
	PRODUCT
	PREFIX
	CALL
)

func TestPrattParsing(t *testing.T) {

}
