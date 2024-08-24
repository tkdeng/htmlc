package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/AspieSoft/go-regex-re2/v2"
	"github.com/tkdeng/goutil"
	"github.com/tkdeng/htmlc"
)

func main() {
	args := goutil.ReadArgs()

	src := args.Get("./src", "src", "root", "")
	dist := args.Get("", "dist", "output", "out", "o")
	port := args.Get("", "port", "")

	if !regex.Comp(`^[0-9]{4,5}$`).Match([]byte(port)) && regex.Comp(`^[0-9]{4,5}$`).Match([]byte(src)) {
		goutil.Swap(&src, &port)
	}

	if dist == "" {
		goutil.JoinPath(src, filepath.Base(src)+".exs")
	}

	if !strings.HasSuffix(dist, ".exs") {
		dist += ".exs"
	}

	if port != "" {
		p, err := strconv.ParseUint(port, 10, 16)
		if err != nil {
			panic(errors.Join(err, errors.New("port \""+port+"\": must be between 3000 and 65535")))
		}
		if p < 3000 || p > 65535 {
			panic(errors.Join(strconv.ErrRange, errors.New("port \""+port+"\": must be between 3000 and 65535")))
		}
		fmt.Println(p)
	}

	htmlc.Compile(src, dist)

	//todo: add optional render method if port != ""
}
