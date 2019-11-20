gomove
===================


**gomove** is a utility to help you move golang packages by automatically changing the import paths from the old one to new one.

----------

Getting Started
-------------

Getting started with gomove is really easy. If you have a proper `$GOPATH` setup and your path set to `bin` directory in `$GOPATH`, you can do the following to get gomove tool:

    $ go get -u github.com/ksubedi/gomove

You can also download pre-built binary files and debian packages from the [releases page](https://github.com/KSubedi/gomove/releases).

Once you have the gomove installed, using it is really easy. First, move your package directory to the new directory and run gomove tool.

In this example, we are moving package `github.com/ksubedi/go-web-seed` to `github.com/ksubedi/new-project`. First we move the first directory to the second one, then we can do the following to automatically update the imports:

	$ gomove -d $GOPATH/src/github.com/ksubedi/new-project github.com/ksubedi/go-web-seed github.com/ksubedi/new-project
	
You can also `cd` to the directory of `github.com/ksubedi/new-project` and run gomove like this:

	$ gomove github.com/ksubedi/go-web-seed github.com/ksubedi/new-project
	
You can also only replace the contents one file only by using `-f` or `--file` flag.

	$ gomove -f hello.go github.com/bla/bla github.com/foo/bar

You can also run `gomove --help` for help.
	
	$ gomove --help
	NAME:
	   gomove - Move Golang packages to a new path.

	USAGE:
	   gomove [global options] command [command options] [old path] [new path]
	   
	VERSION:
	   0.2.17
	   
	AUTHOR(S):
	   Kaushal Subedi <kaushal@subedi.co> 
	   
	COMMANDS:
	   help, h	Shows a list of commands or help for one command
	   
	GLOBAL OPTIONS:
	   --dir, -d "./"		directory to scan
	   --file, -f 			only move imports in a file
	   --safe-mode, -s "false"	run program in safe mode (comments will be wiped)
	   --help, -h			show help
	   --version, -v		print the version


Safe Mode
-------------
By default this tool uses something called native mode where it scans through the go files, finds import statements and replaces the package names. However if you want to use a more safe way of changing the import statements, you can use Safe Mode which parses the ast tree of the go file and replaces imports from there. The biggest downside to using this is that it will [prettify](https://godoc.org/github.com/golang/go/src/cmd/gofmt) your files. To use it in safe mode set the `--safe-mode` or `-s` flag to true:
    
    gomove --safe-mode true github.com/bla/bla github.com/foo/bar


License
-------------

This software is licensed under the GNU GPL V3 License. Check [LICENSE.md](LICENSE.md) for full license.