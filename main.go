package main

import (
	"flag"
	"flagger/flagger"
	"go/printer"
	"os"
	"path/filepath"
	"strings"

	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func processFile(path string, dateFormat string) {
	if !strings.HasSuffix(path, ".go") {
		return
	}
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, path, nil, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}
	flaggerObj := &flagger.Flagger{}
	var unFlagged = false
	for _, funcName := range f.Decls {
		switch funcD := funcName.(type) {
		case *ast.FuncDecl:
			unFlagged = flaggerObj.CheckForFlag(funcD)
		}
	}

	if unFlagged {
		//var buf bytes.Buffer
		//printer.Fprint(&buf, token.NewFileSet(), f)
		//fileErr := ioutil.WriteFile(path, buf.Bytes(), 0644)
		//if fileErr != nil {
		//	panic(fileErr)
		//}
		printer.Fprint(os.Stdout, token.NewFileSet(), f)
	}
}

func main() {
	// Commandline Flags initialization
	recursive := flag.Bool("recursive", false, "Recursively parse flags")
	dateFormat := flag.String("date-format", "2008_01_01", "Format of the date embedded in flag")

	path := flag.String("path", "", "Relative or Absolute Path of the file or folder")

	flag.Parse()

	if *path == "" {
		panic("Path is mandatory")
	}

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	filesPath := *path
	if !strings.HasPrefix(filesPath, "/") {
		filesPath = currentDir + "/" + filesPath
	}

	fi, err := os.Stat(filesPath)
	if os.IsNotExist(err) {
		panic("Path provided does not exist " + filesPath)
	}

	var filesList []string
	if fi.IsDir()  {
		if !*recursive {
			panic("Given folder as path but recursive flag is not enabled")
		}
		traversalErr := filepath.Walk(filesPath,
			func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir(){
				filesList = append(filesList, path)
			}
			return nil
		})

		if traversalErr != nil {
			panic(traversalErr)
		}

	} else {
		filesList = append(filesList, filesPath)
	}
	for _, file := range filesList {
		processFile(file, *dateFormat)
	}

}