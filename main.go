package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/mgutz/ansi"

	"golang.org/x/tools/go/ast/astutil"
)

func main() {
	app := cli.NewApp()
	app.Name = "gomove"
	app.Usage = "Move Golang packages to a new path."
	app.Version = "0.0.1"
	app.ArgsUsage = "[old path] [new path]"
	app.Author = "Kaushal Subedi <kaushal@subedi.co>"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "dir, d",
			Value: "./",
			Usage: "directory to scan",
		},
		cli.StringFlag{
			Name:  "file, f",
			Usage: "only move imports in a file",
		},
	}

	app.Action = func(c *cli.Context) {
		file := c.String("file")
		dir := c.String("dir")
		from := c.Args().Get(0)
		to := c.Args().Get(1)

		if file != "" {
			ProcessFileNative(file, from, to)
		} else {
			RunApp(dir, from, to, c)
		}

	}

	app.Run(os.Args)
}

func RunApp(dir string, from string, to string, c *cli.Context) {

	if from != "" && to != "" {
		filepath.Walk(dir, func(filePath string, info os.FileInfo, err error) error {
			ProcessFileNative(filePath, from, to)
			return nil
		})

	} else {
		cli.ShowAppHelp(c)
	}

}

func ProcessFileNative(filePath string, from string, to string) {

	// Open file to read
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
	}

	// Scan file line by line
	scanner := bufio.NewScanner(bytes.NewReader(fileContent))

	// Track line that is being scanned
	scanLine := 0
	// Track number of changes in file
	numChages := 0

	// Control variables
	isImportLine := false

	// Store final output text
	output := ""

	// Scan through the lines of go file
	for scanner.Scan() {

		scanLine++
		line := scanner.Text()
		bareLine := strings.Replace(line, " ", "", -1)

		// If it is a single import statement, replace the path in that line
		if strings.Contains(bareLine, "import\"") {
			fmt.Println("Found Import On Line", scanLine)
			newImport := strings.Replace(line, from, to, -1)
			output += newImport + "\n"
			numChages++
			continue
		}

		// Change isImportLine accordingly if import statements are detected
		if strings.Contains(bareLine, "import(") {
			fmt.Println("Found Multiple Imports Starting On Line", scanLine)
			isImportLine = true
		} else if isImportLine && strings.Contains(bareLine, ")") {
			fmt.Println("Imports Finish On Line", scanLine)
			isImportLine = false
		}

		// If it is a import line, replace the import
		if isImportLine {
			newImport := strings.Replace(line, from, to, -1)
			fmt.Println("Replacing", line, "to", newImport, "on line", scanLine)
			output += newImport + "\n"
			numChages++
			continue
		}

		// Just copy the rest of the lines to the output
		output += line + "\n"

	}

	// Only write if changes were made
	if numChages > 0 {
		ioutil.WriteFile(filePath, []byte(output), os.ModePerm)
	}
}

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
