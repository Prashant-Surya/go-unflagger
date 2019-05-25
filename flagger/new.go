package flagger

import "go/ast"

const (
	DATE = "date"
	NAME = "name"
)

type FlaggerInterface interface {
	CheckForFlag(function *ast.FuncDecl) bool
}

func NewFlagger(flaggerType string, name string, dateFormat string) FlaggerInterface {
	switch flaggerType {
	case DATE:
		return &DateFlagger{
			DateFormat: dateFormat,
		}
	case NAME:
		return &NameFlagger{
			Name: name,
		}
	default:
		return nil
	}
}