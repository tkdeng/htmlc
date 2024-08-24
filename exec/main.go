package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/AspieSoft/go-regex-re2/v2"
	"github.com/tkdeng/goutil"
	"github.com/tkdeng/htmlc"
)

func main() {
	args := goutil.ReadArgs()

	src := args.Get("./src", "src")
	dist := args.Get("./html.exs", "dist") //todo: use folder name as default file name, and keep in same folder
	root := args.Get("./", "root", "")
	port := args.Get("", "port", "")

	if !regex.Comp(`^[3-6][0-9]{3,4}$`).Match([]byte(port)) && regex.Comp(`^[3-6][0-9]{3,4}$`).Match([]byte(root)) {
		goutil.Swap(&root, &port)
	}

	if err := goutil.AppendRoot(root, &src, &dist); err != nil {
		panic(err)
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
