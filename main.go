package main

import (
	"flagger/flagger"
	"os"

	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
)

var src = `
package main

import ("fmt")

func testFunction() {
        if !config.Config.FeatureFlags.GP_users_details_v2__2018_05_26 {
				fmt.Println("test")
        } else {
			fmt.Println("test else ")
		}
}
`

func main() {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, "", src, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}
	flag := &flagger.Flagger{}
	for _, funcName := range f.Decls {
		switch funcD := funcName.(type) {
		case *ast.FuncDecl:
			flag.CheckForFlag(funcD)
		}
	}

	printer.Fprint(os.Stderr, token.NewFileSet(), f)

}