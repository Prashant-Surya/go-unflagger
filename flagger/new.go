package flagger

import "go/ast"

const (
	DATE = "date"
	NAME = "name"
)

type FlaggerInterface interface {
	CheckForFlag(function *ast.FuncDecl) bool
	//ValidateCondition(conditions []string) bool
}

type FlagChecker interface {
	IsValidFlag(conditions []string) bool
}

func NewFlagger(flaggerType string, name string, dateFormat string) *CommonFlagger {
	obj := &CommonFlagger{}
	switch flaggerType {
	case DATE:
		obj.FlagCheckerObj = &DateFlagChecker{
			DateFormat: dateFormat,
		}
	case NAME:
		obj.FlagCheckerObj = &NameFlagChecker{
			Name: name,
		}
	default:
		return nil
	}
	return obj
}