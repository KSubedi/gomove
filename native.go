package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

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
