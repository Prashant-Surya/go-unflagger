package flagger

import (
	"fmt"
	"go/ast"
)

type NameFlagger struct {
	Name string
}

func(flagger *NameFlagger) isFlag(condition ast.Expr) bool {
	temp := condition
	breakLoop := false
	var conditions []string
	for !breakLoop {
		switch tempType := temp.(type) {
		case *ast.SelectorExpr:
			conditions = append(conditions, fmt.Sprint(tempType.Sel))
			temp = tempType.X
		default:
			breakLoop = true
			break
		}
	}

	for i, v := range conditions {
		if v == "FeatureFlags" && conditions[i-1] == flagger.Name {
			return true
			break
		}
	}

	return false
}

func(flagger *NameFlagger) elseImplementation(updatedList *[]ast.Stmt, elseBlock ast.Stmt) {
	elseStmt := elseBlock.(*ast.BlockStmt)
	for _, item := range elseStmt.List {
		*updatedList = append(*updatedList, item)
	}
}

func(flagger *NameFlagger) CheckForFlag(function *ast.FuncDecl) (unFlagged bool){
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
				unFlagged = true
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

	return
}