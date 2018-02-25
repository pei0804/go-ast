package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"log"

	"github.com/pei0804/go-ast/eval"
)

// https://qiita.com/tenntenn/items/13340f2845316532b55a

/*
PackageClauseOnly: packageの部分までパースする
ImportsOnly: importの部分までパースする
ParseComments: コメントをパースして、ASTに追加する
Trace: パースの過程をトレースする
DeclarationErrors: 定義エラーを見つける（多重定義とか）
SpuriousErrors: 互換のためにある。AllErrorsと同じ
AllErrors: すべてのエラーを出す。そうじゃないと10個まで。
*/

// main コメント
func main() {
	parseFile()
	parseDir()
	parseExpr()
	parseError()
}

func parseFile() {
	log.Println("---parseFile---")
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, "testdata/test1.go", nil, parser.Trace)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%#v", astFile)
}

func parseDir() {
	log.Println("---parseDir---")
	fset := token.NewFileSet()
	astFiles, err := parser.ParseDir(fset, "./testdata", nil, 0)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range astFiles {
		log.Println(k)
		log.Println(v)
	}
}

func parseExpr() {
	log.Println("---parseExpr---")
	expr, err := parser.ParseExpr("1+2")
	if err != nil {
		panic(err)
	}

	log.Println("evalPlus")
	if be, ok := expr.(*ast.BinaryExpr); ok {
		if v, err := eval.Plus(be); err == nil {
			fmt.Println(v)
		} else {
			fmt.Println("Error:", err)
		}
	}

	log.Println("evalBinaryExpr")
	if be, ok := expr.(*ast.BinaryExpr); ok {
		if v, err := eval.BinaryExpr(be); err == nil {
			fmt.Println(v)
		} else {
			fmt.Println("Error:", err)
		}
	}
}

func parseError() {
	log.Println("---parseError---")
	fset := token.NewFileSet()
	_, err := parser.ParseFile(fset, "sample.go", nil, parser.DeclarationErrors|parser.AllErrors)
	if err != nil {
		switch err := err.(type) {
		case scanner.ErrorList:
			for _, e := range err {
				fmt.Printf("%s:%d:%d %s\n", e.Pos.Filename, e.Pos.Line, e.Pos.Column, e.Msg)
			}
		default:
			fmt.Println(err)
		}
		return
	}
}
