package main

import (
	"os"
	"path"
	"path/filepath"

	"github.com/codegangsta/cli"
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
			ScanDir(dir, from, to, c)
		}

	}

	app.Run(os.Args)
}

// ScanDir scans a directory for go files and
func ScanDir(dir string, from string, to string, c *cli.Context) {

	// If from and to are not empty scan all files
	if from != "" && to != "" {
		// Scan directory for files
		filepath.Walk(dir, func(filePath string, info os.FileInfo, err error) error {
			// Only process go files
			if path.Ext(filePath) == ".go" {
				ProcessFileNative(filePath, from, to)
			}
			return nil
		})

	} else {
		cli.ShowAppHelp(c)
	}

}
