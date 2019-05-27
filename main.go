package main

import (
	"flag"
	"flagger/flagger"
	"github.com/dave/dst/decorator"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	ast "github.com/dave/dst"
	"log"
)

func processFile(path string, flaggerObj *flagger.CommonFlagger, writeToFile bool) {
	f, err := decorator.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	var unFlagged = false
	for _, funcName := range f.Decls {
		switch funcD := funcName.(type) {
		case *ast.FuncDecl:
			localFlag, updatedList := flaggerObj.RemoveFlag(funcD.Body.List)
			if localFlag {
				unFlagged = true
				funcD.Body.List = updatedList
			}
		}
	}

	if unFlagged {
		if writeToFile {
			fWriter, err := os.OpenFile(path, os.O_WRONLY, 0644)
			if err != nil {
				panic(err)
			}
			decorator.Fprint(fWriter, f)
		} else {
			if err := decorator.Print(f); err !=nil {
				panic(err)
			}
		}
	}
}

func validations(flaggerType, dateFormat, flagName, path *string, recursive *bool) {
	if *flaggerType != flagger.NAME && *flaggerType != flagger.DATE {
		panic("Invalid flaggerType specified choose name or date")
	}
	if *flaggerType == flagger.NAME && *flagName == "" {
		panic("Name flagger is selected but flag name was not specified")
	}
	if *path == "" {
		panic("Path is mandatory")
	}
}

func main() {
	// Commandline Flags initialization
	recursive := flag.Bool("recursive", false, "Recursively parse flags")
	flaggerType := flag.String("type", "date", "Flagger Type. Possible values date, name")
	dateFormat := flag.String("date-format", "2006_01_02", "Format of the date embedded in flag")
	flagName := flag.String("name", "",  "Name of the flag to be removed")
	path := flag.String("path", "", "Relative or Absolute Path of the file or folder")
	writeToFile := flag.Bool("write", false, "Enable this flag to update contents to file or it'll be written to stdout")
	flag.Parse()

	validations(flaggerType, dateFormat, flagName, path, recursive)

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

	flaggerObj := flagger.NewFlagger(*flaggerType, *flagName, *dateFormat)

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
		if strings.HasSuffix(file, ".go") {
			processFile(file, flaggerObj, *writeToFile)
		}
	}

}