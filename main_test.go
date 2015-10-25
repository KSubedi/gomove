package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestApp(t *testing.T) {
	fileContent := "package testing\n\nimport \"fmt\"\n\nfunc HelloWorld() {\nfmt.Println(\"Hello World!\")\n}\n"

	ioutil.WriteFile("hello.go", []byte(fileContent), os.ModePerm)

	ProcessFile("hello.go", "fmt", "replacedImport")

	result, err := ioutil.ReadFile("hello.go")
	if err != nil {
		t.Error("Failed to read written file.")
	}

	if !strings.Contains(string(result), "replacedImport") {
		t.Error("Got different results")
	}
	os.Remove("hello.go")
}
