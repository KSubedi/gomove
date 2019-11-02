package main

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"

	"github.com/mgutz/ansi"
	"golang.org/x/tools/go/ast/astutil"
)

// ProcessFileAST processes the files using golang's AST parser
func ProcessFileAST(filePath string, from string, to string) {

	//Colors to be used on the console
	red := ansi.ColorCode("red+bh")
	white := ansi.ColorCode("white+bh")
	yellow := ansi.ColorCode("yellow+bh")
	blackOnWhite := ansi.ColorCode("black+b:white+h")
	//Reset the color
	reset := ansi.ColorCode("reset")

	fmt.Println(blackOnWhite+"Processing file", filePath, "in SAFE MODE", reset)

	// New FileSet to parse the go file to
	fSet := token.NewFileSet()

	// Parse the file
	file, err := parser.ParseFile(fSet, filePath, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Keep track of number of changes
	changed := false
	if changed = astutil.RewriteImport(fSet, file, from, to); changed {
		fmt.Println(red +
			"Updating import " +
			reset + white +
			from +
			reset + red +
			" to " +
			reset + white +
			to +
			reset)
	}

	// If the number of changes are more than 0, write file
	if changed {
		// Print the new AST tree to a new output buffer. These Config settings intended to match gofmt.
		printerMode := printer.TabIndent | printer.UseSpaces
		printConfig := &printer.Config{Mode: printerMode, Tabwidth: 8}

		var outputBuffer bytes.Buffer
		err := printConfig.Fprint(&outputBuffer, fSet, file)
		if err != nil {
			fmt.Println(err)
			return
		}

		ioutil.WriteFile(filePath, outputBuffer.Bytes(), os.ModePerm)
		fmt.Println(yellow+
			"File",
			filePath,
			"saved",
			reset, "\n\n")
	} else {
		fmt.Println(yellow+
			"No changes to write on this file.",
			reset, "\n\n")
	}
}
