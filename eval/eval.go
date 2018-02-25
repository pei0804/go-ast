package eval

import (
	"errors"
	"go/ast"
	"go/constant"
	"go/token"
	"strconv"
)

func Plus(expr *ast.BinaryExpr) (int64, error) {
	// 左辺
	xLit, ok := expr.X.(*ast.BasicLit)
	if !ok {
		return 0, errors.New("left operand is not BasicLit")
	}

	// 右辺
	yLit, ok := expr.Y.(*ast.BasicLit)
	if !ok {
		return 0, errors.New("right operand is not BasicLit")
	}

	// 足し算か
	if expr.Op != token.ADD {
		return 0, errors.New("operator is not +")
	}

	// 計算出来るint型に変換する
	x, err := strconv.ParseInt(xLit.Value, 10, 64)
	if err != nil {
		return 0, err
	}

	y, err := strconv.ParseInt(yLit.Value, 10, 64)
	if err != nil {
		return 0, err
	}

	return x + y, nil
}

func BinaryExpr(expr *ast.BinaryExpr) (constant.Value, error) {
	xLit, ok := expr.X.(*ast.BasicLit)
	if !ok {
		return constant.MakeUnknown(), errors.New("left operand is not BasicLit")
	}

	yLit, ok := expr.Y.(*ast.BasicLit)
	if !ok {
		return constant.MakeUnknown(), errors.New("right operand is not BasicLit")
	}

	x := BasicLit(xLit)
	y := BasicLit(yLit)
	return constant.BinaryOp(x, expr.Op, y), nil
}

// いい感じにValueを取ってくる
func BasicLit(expr *ast.BasicLit) constant.Value {
	return constant.MakeFromLiteral(expr.Value, expr.Kind, 0)
}
