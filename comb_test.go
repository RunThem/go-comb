package comb

import (
	"fmt"
	"strconv"
	"testing"
)

func eval(ast *Ast) int {
	if ast.Size() == 0 {
		n, _ := strconv.Atoi(ast.code)
		return n
	}

	num := eval(ast.forward.At(0))
	fmt.Println(num)
	for i := 1; i < ast.Size(); i++ {
		a := ast.forward.At(i)
		switch a.code {
		case "*":
			i += 1
			num *= eval(ast.forward.At(i))
			fmt.Println(num)
		case "/":
			i += 1
			num /= eval(ast.forward.At(i))
			fmt.Println(num)
		case "+":
			i += 1
			num += eval(ast.forward.At(i))
			fmt.Println(num)
		case "-":
			i += 1
			num -= eval(ast.forward.At(i))
			fmt.Println(num)
		case "(":
			i += 1
			return eval(ast.forward.At(i))
		case ")":
			i += 1
		default:
			return eval(ast.forward.At(i))
		}
	}

	return num
}

func TestMatch(t *testing.T) {

	/*
	 * Factor = 'N' | '(' <Expr> ')'
	 * Term   = <Factor> (('*' | '/') <Factor>)*
	 * Expr   = <Term> (('+' | '-') <Term>)*
	 */

	Expr := New(C_AND, "Expr")
	Term := New(C_AND, "Term")
	Factor := New(C_OR, "Factor")

	factor := New(C_AND, "factor")
	factor.Add(Match("("), Expr, Match(")"))

	Factor.Add(Match("3"), factor)

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

	in := NewInput("3+3*(3*3-3)")
	ast := Expr.Parse(in)

	ast.Dump()

	fmt.Println(eval(ast))
}
