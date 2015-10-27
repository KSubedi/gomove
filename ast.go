package main

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/mgutz/ansi"
	"golang.org/x/tools/go/ast/astutil"
)

func ProcessFileAST(filePath string, from string, to string) {

	//Colors to be used on the console
	red := ansi.ColorCode("red+bh")
	white := ansi.ColorCode("white+bh")
	greenUnderline := ansi.ColorCode("green+buh")
	blackOnWhite := ansi.ColorCode("black+b:white+h")
	//Reset the color
	reset := ansi.ColorCode("reset")

	// If the file is a go file scan it
	if path.Ext(filePath) == ".go" {
		// New FileSet to parse the go file to
		fSet := token.NewFileSet()

		// Parse the file
		file, err := parser.ParseFile(fSet, filePath, nil, 0)
		if err != nil {
			fmt.Println(err)
		}

		// Get the list of imports from the ast
		imports := astutil.Imports(fSet, file)

		// Keep track of number of changes
		numChanges := 0

		// Iterate through the imports array
		for _, mPackage := range imports {
			for _, mImport := range mPackage {
				// Since astutil returns the path string with quotes, remove those
				importString := mImport.Path.Value

				// If the path matches the oldpath, replace it with the new one
				if strings.Contains(importString, from) {
					//If it needs to be replaced, increase numChanges so we can write the file later
					numChanges++

					// Join the path of the import package with the remainder from the old one after removing the old import package
					replacePackage := strings.Replace(importString, from, to, -1)

					fmt.Println(red + "Updating import " + importString + " from file " + reset + white + filePath + reset + red + " to " + reset + white + replacePackage + reset)

					// Remove the old import and replace it with the replacement
					astutil.DeleteImport(fSet, file, importString)
					astutil.AddImport(fSet, file, replacePackage)
				}
			}
		}

		// If the number of changes are more than 0, write file
		if numChanges > 0 {
			// Print the new AST tree to a new output buffer
			var outputBuffer bytes.Buffer
			printer.Fprint(&outputBuffer, fSet, file)

			ioutil.WriteFile(filePath, outputBuffer.Bytes(), os.ModePerm)
			fmt.Printf(blackOnWhite+"File "+filePath+" saved after %d changes."+reset+"\n", numChanges)
		} else {
			fmt.Println(greenUnderline + "No changes needed on file " + filePath + reset)
		}
	}
}
