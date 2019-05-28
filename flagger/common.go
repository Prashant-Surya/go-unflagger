package flagger

import (
	"fmt"
	ast "github.com/dave/dst"
)

type CommonFlagger struct {
	FlagCheckerObj FlagChecker
}

func(flagger *CommonFlagger) isFlag(condition ast.Expr) bool {
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

	if flagger.FlagCheckerObj.IsValidFlag(conditions) {
		return true
	}

	return false

}

func(flagger *CommonFlagger) elseImplementation(updatedList *[]ast.Stmt, elseBlock ast.Stmt) {
	if elseBlock == nil {
		return
	}
	elseStmt := elseBlock.(*ast.BlockStmt)
	for _, item := range elseStmt.List {
		*updatedList = append(*updatedList, item)
	}
}

func(flagger *CommonFlagger) RemoveFlag(body []ast.Stmt) (bool, []ast.Stmt){
	var updatedList []ast.Stmt
	var unFlagged = false
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
			} else {
				// If there is no flag check in body if it contains flag
				localFlag, localUpdate := flagger.RemoveFlag(stmtType.Body.List)
				if localFlag {
					unFlagged = true
					stmtType.Body.List = localUpdate
				}
				updatedList = append(updatedList, stmt)
			}
		default:
			updatedList = append(updatedList, stmt)
		}
	}

	return unFlagged, updatedList
}