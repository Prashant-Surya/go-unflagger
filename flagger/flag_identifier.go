package flagger

import (
	"fmt"
	"go/ast"
	"strings"
	"time"
)

type Flagger struct {
}

func(flagger *Flagger) isFlag(condition ast.Expr) bool {
	temp := condition
	breakLoop := false
	var conditions []string
	for !breakLoop {
		switch tempType := temp.(type) {
		case *ast.BinaryExpr:

		case *ast.SelectorExpr:
			conditions = append(conditions, fmt.Sprint(tempType.Sel))
			temp = tempType.X
		default:
			breakLoop = true
			break
		}
	}

	var flag string

	for i, v := range conditions {
		if v == "FeatureFlags" {
			flag = conditions[i-1]
			break
		}
	}

	flagsSplit := strings.Split(flag, "__")
	if len(flagsSplit) <= 1 {
		return false
	}

	date := extractDate(flagsSplit[len(flagsSplit) - 1])
	if date != nil {
		now := time.Now()
		if now.After(*date) {
			return true
		}
	}

	return false
}

func(flagger *Flagger) elseImplementation(updatedList *[]ast.Stmt, elseBlock ast.Stmt) {
	elseStmt := elseBlock.(*ast.BlockStmt)
	for _, item := range elseStmt.List {
		*updatedList = append(*updatedList, item)
	}
}

func(flagger *Flagger) CheckForFlag(function *ast.FuncDecl) *ast.FuncDecl{
	body := function.Body.List
	var updatedList []ast.Stmt
	for _, stmt := range body {
		switch stmtType := stmt.(type) {
		case *ast.IfStmt:
			var flag bool
			var ifImplementation = true
			if unary, ok := stmtType.Cond.(*ast.UnaryExpr); ok {
				flag = flagger.isFlag(unary.X)
				ifImplementation = false
			} else {
				flag = flagger.isFlag(stmtType.Cond)
			}

			if flag {
				if ifImplementation {
					for _, item := range stmtType.Body.List {
						updatedList = append(updatedList, item)
					}
				} else {
					flagger.elseImplementation(&updatedList, stmtType.Else)
				}
				continue
			}
			updatedList = append(updatedList, stmt)
		default:
			updatedList = append(updatedList, stmt)
		}
	}
	function.Body.List = updatedList
	return function
}