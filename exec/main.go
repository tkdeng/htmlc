package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	regex "github.com/tkdeng/goregex"
	"github.com/tkdeng/goutil"
	"github.com/tkdeng/htmlc"
)

func main() {
	args := goutil.ReadArgs()

	src := args.Get("./src", "src", "root", "")
	out := args.Get("", "out", "output", "o", "dist")
	port := args.Get("", "port", "")
	noCompile := args.Get("false", "no-compile", "n")

	if !regex.Comp(`^[0-9]{4,5}$`).Match([]byte(port)) && regex.Comp(`^[0-9]{4,5}$`).Match([]byte(src)) {
		goutil.Swap(&src, &port)
	}

	if out == "" {
		goutil.JoinPath(src, filepath.Base(src)+".exs")
	}

	if out == "" {
		out = "html.exs"
	}

	if !strings.HasSuffix(out, ".exs") {
		out += ".exs"
	}

	if port != "" {
		p, err := strconv.ParseUint(port, 10, 16)
		if err != nil {
			panic(errors.Join(err, errors.New("port \""+port+"\": must be between 3000 and 65535")))
		}
		if p < 3000 || p > 65535 {
			panic(errors.Join(strconv.ErrRange, errors.New("port \""+port+"\": must be between 3000 and 65535")))
		}
	}

	if noCompile == "false" {
		htmlc.Compile(src, out)
	}

	if port != "" {
		runServer(out, port)
	}
}

func runServer(file string, port string) {
	exs, err := htmlc.Engine(file)
	if err != nil {
		panic(err)
	}

	exs.Render("index", htmlc.Map{
		"title": "Title",
		"desc":  "Desc",
		"n":     2,
		"map": htmlc.Map{
			"test": "map",
		},
		"arr": []any{
			"test arr",
		},
	}, "layout")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := goutil.Clean(r.URL.Path)
		if url == "/" {
			url = "index"
		} else {
			url = url[1:]
		}

		args := htmlc.Map{}

		e := []byte(r.URL.Query().Encode())
		if len(e) != 0 {
			query := bytes.Split(e, []byte{'&'})
			for _, arg := range query {
				a := bytes.SplitN(arg, []byte{'='}, 2)

				if len(a) == 1 || bytes.Equal(a[1], []byte("true")) {
					args[string(a[0])] = true
				} else if bytes.Equal(a[1], []byte("false")) {
					args[string(a[0])] = false
				} else if bytes.Equal(a[1], []byte("nil")) || bytes.Equal(a[1], []byte("null")) {
					args[string(a[0])] = nil
				} else if regex.CompRE2(`^[0-9]+(\.[0-9]+|)$`).Match(a[1]) {
					if i, e := strconv.Atoi(string(a[1])); e == nil {
						args[string(a[0])] = i
					} else {
						args[string(a[0])] = a[1]
					}
				} else {
					args[string(a[0])] = a[1]
				}
			}
		}

		// w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		buf, err := exs.Render(url, args)
		if err != nil {
			if err == htmlc.ErrorTimeout {
				w.WriteHeader(408)
				w.Write([]byte("<h1>Error 408</h1><h2>Request Timed Out!</h2>"))
				return
			} else if err == htmlc.ErrorStopped {
				w.WriteHeader(500)
				w.Write([]byte("<h1>Error 500</h1><h2>Internal Server Error!</h2>"))
				return
			}

			w.WriteHeader(400)
			w.Write([]byte("<h1>Error 400</h1><h2>Bad Request!</h2>"))
			return
		}

		w.WriteHeader(200)
		w.Write(buf)
	})

	fmt.Println("Running HTTP Server On Port:", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
