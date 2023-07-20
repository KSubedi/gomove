package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mgutz/ansi"
)

var (
	red          = ansi.ColorCode("red+bh")
	white        = ansi.ColorCode("white+bh")
	green        = ansi.ColorCode("green+bh")
	yellow       = ansi.ColorCode("yellow+bh")
	blackOnWhite = ansi.ColorCode("black+b:white+h")
	//Reset the color
	reset = ansi.ColorCode("reset")
)

// ProcessFileNative processes files using native string search instead of ast parsing
func ProcessFileNative(filePath string, from string, to string) {
	//Colors to be used on the console

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
	numChanges := 0

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

			newImport := strings.Replace(line, from, to, -1)
			output += newImport + "\n"
			if line != newImport {
				numChanges++

				fmt.Println(red+"Updating "+
					reset+white+
					strings.TrimSpace(line)+
					reset+red+
					" to "+
					reset+white+
					strings.TrimSpace(newImport)+
					reset+red+
					"on line", scanLine, reset)
			}

			continue
		}

		// Change isImportLine accordingly if import statements are detected
		if strings.Contains(bareLine, "import(") {
			//	fmt.Println(green+"Found Multiple Imports Starting On Line", scanLine, reset)
			isImportLine = true
		} else if isImportLine && strings.Contains(bareLine, ")") {
			//	fmt.Println(green+"Imports Finish On Line", scanLine, reset)
			isImportLine = false
		}

		// If it is a import line, replace the import
		if isImportLine {
			newImport := strings.Replace(line, from, to, -1)
			output += newImport + "\n"
			if line != newImport {
				numChanges++
				fmt.Println(red+"Updating text "+
					reset+white+
					strings.TrimSpace(line)+
					reset+red+
					" to "+
					reset+white+
					strings.TrimSpace(newImport)+
					reset+red+
					" on line", scanLine, reset)

			}
			continue
		}

		// Just copy the rest of the lines to the output
		output += line + "\n"

	}

	// Only write if changes were made
	if numChanges > 0 {
		fmt.Print(yellow+
			"File",
			filePath,
			" saved after ",
			numChanges,
			" changes",
			reset, "\n")
		ioutil.WriteFile(filePath, []byte(output), os.ModePerm)
	} /* else {
		fmt.Print(yellow+
			"No changes to write on this file.",
			reset, "\n\n\")
	} */
}

// replaceFile goes through a non go file and does a direct replacement of the lookup[from][to] values provided.
// it returns a list of keys that were replaced.
func replaceFile(path string, deleteLines bool, lookup map[string]string) ([]string, error) {
	foundKeys := make(map[string]struct{})
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("open file %v: %w", path, err)
	}
	scanner := bufio.NewScanner(bytes.NewReader(file))
	output := ""
	numChanges := 0
loop:
	for scanLine := 0; scanner.Scan(); scanLine++ {
		line := scanner.Text()
		for k, v := range lookup {
			if strings.Contains(line, k) {
				numChanges++
				foundKeys[k] = struct{}{}
				if deleteLines {
					line = ""
					continue loop
				}
				line = strings.ReplaceAll(line, k, v)
			}

		}
		output += line + "\n"
	}
	err = os.WriteFile(path, []byte(output), os.ModePerm)

	replaceValues := make([]string, 0, len(foundKeys))
	for k, _ := range foundKeys {
		replaceValues = append(replaceValues, k)
	}
	return replaceValues, err
}
