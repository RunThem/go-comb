package parser

import (
	"fmt"
	"github.com/RunThem/comb/ast"
	"github.com/RunThem/comb/comb"
	lexer "github.com/RunThem/comb/lexer"
	"strconv"
	"testing"
)

func par_eval(ast *ast.Ast) int {
	if ast.Size() == 0 {
		n, _ := strconv.Atoi(ast.Code)
		return n
	}

	num := par_eval(ast.At(0))
	fmt.Println(num)
	for i := 1; i < ast.Size(); i++ {
		a := ast.At(i)
		switch a.Code {
		case "*":
			i += 1
			num *= par_eval(ast.At(i))
			fmt.Println("*", num)
		case "/":
			i += 1
			num /= par_eval(ast.At(i))
			fmt.Println("/", num)
		case "+":
			i += 1
			num += par_eval(ast.At(i))
			fmt.Println("+", num)
		case "-":
			i += 1
			num -= par_eval(ast.At(i))
			fmt.Println("-", num)
		case "(":
			i += 1
			return par_eval(ast.At(i))
		case ")":
			i += 1
		default:
			return par_eval(ast.At(i))
		}
	}

	return num
}

func TestNewParser(t *testing.T) {
	//	code := `
	//$factor = /[0-9]*/ | '(' <expr> ')';
	//$term   = $factor [<'*' | '/'> $factor];
	//$expr   = $term [<'+' | '-'> $term];
	//`

	code := `
$factor = /[0-9]*/ | '(' $expr ')';
$term   = $factor [<'*' | '/'> $factor];
$expr   = $term [<'+' | '-'> $term];
`

	lex := lexer.NewLexer(code)

	for i := 0; i < lex.Lex.Size(); i++ {
		lst := lex.Lex.At(i)

		for it := lst.Front(); it != nil; it = it.Next() {
			fmt.Println(it.Value)
		}

		fmt.Println()
	}

	expr := NewParser(lex)

	expr.Dump()

	in := comb.NewInput("4234+9*(3*3-3)")
	//in := comb.NewInput("4+3-(4+2+(2-3))")
	//in := comb.NewInput("4+4")
	ast := expr.Parse(in)

	ast.Dump()

	//fmt.Println(par_eval(ast))
}

func TestJson(t *testing.T) {
	code := `
$number 		= /[\d+.\d]*/;
$string			= //
$item_pair 		= $number ':' $value;
$object_items 	= $item_pair [',' $item_pair];
$object 		= '{' $object_items '}';
$array_items 	= $value [',' $value];
$array 			= '[' $array_items ']';
$value 			= $object | $array | $number | 'true' | 'false' | 'null';
$json 			= $value;
`

	json := NewParser(lexer.NewLexer(code))
	json.Dump()

	in := comb.NewInput(`{1:[11,12],2:{22:true,23:false,24:[244,245]},3:3}`)

	ast := json.Parse(in)

	ast.Dump()
}
