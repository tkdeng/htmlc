package main

import (
	"fmt"

	"github.com/tkdeng/goutil"
	"golang.org/x/term"
)

func printHelp() {
	tSize := 80
	if w, _, err := term.GetSize(0); err == nil {
		if w-10 > 20 {
			tSize = w - 10
		} else {
			tSize = 20
		}
	}

	fmt.Println(goutil.ToColumns([][]string{
		{"\nUsage: htmlc [src] ?[port] [options...]\n"},
		{"--src, --root", "source dicectory."},
		{"--port", "port number to start http server on localhost. (defaults to no server)."},
		{"--out, --o", "output file path. (defaults to the name of the source directory)."},
		{"--no-compile, -n", "skip compiling, and just run the server [out] file on localhost [port]."},
		{"--iex-mode, --iex", "compile the file for the 'iex' command instead of for the 'elixir' command."},
		{"--help, -h", "print this list."},
	}, tSize, "    ", "\n\n"))
}
