package comb

import (
	"testing"
)

func TestExample(t *testing.T) {
	/*
	 * Factor = <number> | '(' <Expr> ')'
	 * Term   = <Factor> (('*' | '/') <Factor>)*
	 * Expr   = <Term> (('+' | '-') <Term>)*
	 */

	Expr := New(C_COMB, "Expr")
	Term := New(C_COMB, "Term")
	Factor := New(C_OR, "Factor")

	op := New(C_OR, "expr_op")
	op.Add(Match("+"), Match("-"))

	expr := New(C_MAYBE, "expr")
	expr.Add(op, Term)

	Expr.Add(Term, expr)

	op = New(C_OR, "term_op")
	op.Add(Match("*"), Match("/"))

	term := New(C_MAYBE, "term")
	term.Add(op, Factor)

	Term.Add(Factor, term)

	factor := New(C_AND, "factor")
	factor.Add(Match("("), Expr, Match(")"))

	Factor.Add(Match("num"), factor)

	Expr.Dump()

	// in := NewInput("3+3")
	// Expr.Parse(in).Dump()
}

func TestMatch(t *testing.T) {
	Term := New(C_AND, "Term")
	Factor := New(C_OR, "Factor")

	factor := New(C_AND, "factor")
	factor.Add(Match("("), Term, Match(")"))

	Factor.Add(Match("N"), factor)

	term := New(C_MAYBE, "term")

	op := New(C_OR, "term_op")
	op.Add(Match("*"), Match("/"))

	term.Add(op, Factor)

	Term.Add(Factor, term)

	in := NewInput("(N*N)*N/(N/N)")
	Term.Parse(in).Dump()
}
