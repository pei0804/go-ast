package main

import (
	"bufio"
	"errors"
	"fmt"
	"go/ast"
	"go/constant"
	"go/parser"
	"go/token"
	"io"
	"os"
)

// https://qiita.com/tenntenn/items/a312d2c5381e36cf4cd3

func eval(e string) (string, error) {
	expr, err := parser.ParseExpr(e)
	if err != nil {
		return "", err
	}

	v, err := evalExpr(expr)
	if err != nil {
		return "", err
	}

	return v.String(), nil
}

func evalExpr(expr ast.Expr) (v constant.Value, rerr error) {
	defer func() {
		if r := recover(); r != nil {
			v, rerr = constant.MakeUnknown(), fmt.Errorf("%v", r)
		}
	}()

	switch e := expr.(type) {
	case *ast.ParenExpr:
		return evalExpr(e.X)
	case *ast.BinaryExpr:
		return evalBinaryExpr(e)
	case *ast.UnaryExpr:
		return evalUnaryExpr(e)
	case *ast.BasicLit:
		return constant.MakeFromLiteral(e.Value, e.Kind, 0), nil
	case *ast.Ident:
		return evalIdent(e)
	}

	return constant.MakeUnknown(), errors.New("unkown node")
}

func evalBinaryExpr(expr *ast.BinaryExpr) (constant.Value, error) {
	x, err := evalExpr(expr.X)
	if err != nil {
		return constant.MakeUnknown(), err
	}

	y, err := evalExpr(expr.Y)
	if err != nil {
		return constant.MakeUnknown(), err
	}

	switch expr.Op {
	case token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ:
		return constant.MakeBool(constant.Compare(x, expr.Op, y)), nil
	}

	return constant.BinaryOp(x, expr.Op, y), nil
}

func evalUnaryExpr(expr *ast.UnaryExpr) (constant.Value, error) {
	x, err := evalExpr(expr.X)
	if err != nil {
		return constant.MakeUnknown(), err
	}

	return constant.UnaryOp(expr.Op, x, 0), nil
}

func evalIdent(expr *ast.Ident) (constant.Value, error) {
	switch {
	case expr.Name == "true":
		return constant.MakeBool(true), nil
	case expr.Name == "false":
		return constant.MakeBool(false), nil
	}

	return constant.MakeUnknown(), errors.New("unkown ident")
}

func repl(r io.Reader) {
	s := bufio.NewScanner(r)
	for {
		fmt.Print(">")
		if !s.Scan() {
			break
		}

		l := s.Text()
		switch {
		case l == "bye":
			return
		default:
			r, err := eval(l)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			fmt.Println(r)
		}
	}

	if err := s.Err(); err != nil {
		fmt.Println("Error:", err)
	}
}

func main() {
	repl(os.Stdin)
}
