package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gomove"
	app.Usage = "Move Golang packages to a new path."
	app.Version = "0.2.17"
	app.ArgsUsage = "[old path] [new path]"
	app.Author = "Kaushal Subedi <kaushal@subedi.co>"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "dir, d",
			Value: "./",
			Usage: "directory to scan. files under vendor/ are ignored",
		},
		cli.StringFlag{
			Name:  "conf, c",
			Value: "",
			Usage: "json lookup config for updating multiple paths at the same time",
		},
		cli.StringFlag{
			Name:  "file, f",
			Usage: "only move imports in a file",
		},
		cli.StringFlag{
			Name:  "safe-mode, s",
			Value: "false",
			Usage: "run program in safe mode (comments will be wiped)",
		},
	}

	app.Action = func(c *cli.Context) {
		lookup := make(map[string]string)
		file := c.String("file")
		dir := c.String("dir")
		conf := c.String("conf")

		from := c.Args().Get(0)
		to := c.Args().Get(1)

		if conf != "" {
			b, err := os.ReadFile(conf)
			if err != nil {
				log.Fatalf("error reading file %v: %v", conf, err)
			}
			if err := json.Unmarshal(b, &lookup); err != nil {
				log.Fatalf("error parsing file %v: %v", conf, err)
			}

		} else if from == "" || to == "" {
			cli.ShowAppHelp(c)
			log.Println("is this needed?")
			return

		} else {
			lookup = map[string]string{from: to}
		}

		if file != "" {
			ProcessFile(file, from, to, c)
		} else {
			ParseMod(dir, lookup)
			ScanDir(dir, lookup, c)
		}

	}

	app.Run(os.Args)
}

// ScanDir scans a directory for go files and
func ScanDir(dir string, lookup map[string]string, c *cli.Context) {
	// construct the path of the vendor dir of the project for prefix matching
	vendorDir := path.Join(dir, "vendor")
	// Scan directory for files
	filepath.Walk(dir, func(filePath string, info os.FileInfo, err error) error {
		// ignore vendor path
		if matched := strings.HasPrefix(filePath, vendorDir); matched {
			return nil
		}
		// Only process go files
		if path.Ext(filePath) == ".go" {
			fmt.Println(blackOnWhite+"Processing file", filePath, reset)
			for from, to := range lookup {
				ProcessFile(filePath, from, to, c)
			}
		}

		return nil
	})
}

// ParseMod goes through the go.mod and go.sum file
// it will optimize the lookup map by removing values not found in the
// go.mod file.
func ParseMod(dir string, lookup map[string]string) {
	foundKeys := make(map[string]string)
	modPath := strings.TrimRight(dir, "/") + "/go.mod"
	modFile, err := os.ReadFile(modPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	sumPath := strings.TrimRight(dir, "/") + "/go.sum"
	sumFile, err := os.ReadFile(sumPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	numChanges := 0
	output := ""
	scanner := bufio.NewScanner(bytes.NewReader(modFile))
modScan:
	for scanLine := 0; scanner.Scan(); scanLine++ {
		line := scanner.Text()

		for key, value := range lookup {
			if strings.Contains(line, key) {
				numChanges++
				foundKeys[key] = value
				output += strings.Replace(line, key, value, 1) + "\n"
				continue modScan
			}
		}
		output += line + "\n"
	}
	// remove unneeded keys
	for k := range lookup {
		if _, found := foundKeys[k]; !found {
			fmt.Printf("Remove unused key: %v\n", k)
			delete(lookup, k)
		}
	}

	if numChanges > 0 {
		os.WriteFile(modPath, []byte(output), os.ModePerm)
	}

	// update go.sum file
	numChanges = 0
	output = ""
	scanner = bufio.NewScanner(bytes.NewReader(sumFile))
sumScan:
	for scanLine := 0; scanner.Scan(); scanLine++ {
		line := scanner.Text()

		for key, _ := range lookup {
			if strings.Contains(line, key) {
				numChanges++
				continue sumScan
			}
		}
		output += line + "\n"
	}
	if numChanges > 0 {
		os.WriteFile(sumPath, []byte(output), os.ModePerm)
	}
}

// ProcessFile processes file according to what mode is chosen
func ProcessFile(filePath string, from string, to string, c *cli.Context) {
	if c.String("safe-mode") == "true" {
		ProcessFileAST(filePath, from, to)
	} else {
		ProcessFileNative(filePath, from, to)
	}
}
