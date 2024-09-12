package plugin

import (
	"bytes"
	"errors"
	"os"

	regex "github.com/tkdeng/goregex"
	"github.com/tkdeng/goutil"
	"github.com/tkdeng/htmlc/common"
)

// compileExs compiles elixir in the html buffer
//
// call compile to start main compiler
type compileExs struct {
	buf *[]byte
}

var Compiler = map[string]func(comp *compileExs, iexMode bool){}
var CompilerHTML = map[string]bool{}
var CompilerTags = map[string]string{}
var Loader = map[string]func(out *os.File, name string, buf *[]byte, usedRandID *[][]byte, iexMode bool){}

// addCompiler adds a lang compiler method
//
// @name: .ext
//
// @html: include initial html compiler
//
// @tags: html <tags> to compile content </tags>
//
// @cb: method to run when compiling
func addCompiler(name string, html bool, tags []string, cb func(comp *compileExs, iexMode bool)) error {
	if _, ok := Compiler[name]; !ok {
		Compiler[name] = cb
		CompilerHTML[name] = html

		for _, tag := range tags {
			CompilerTags[tag] = name
		}
		return nil
	}

	return errors.New("compiler name taken")
}

// addLoader adds a new loader to load external files into the script
func addLoader(name string, cb func(out *os.File, name string, buf *[]byte, usedRandID *[][]byte, iexMode bool)) error {
	if _, ok := Loader[name]; !ok {
		Loader[name] = cb
		return nil
	}

	//todo: have loader also load methods into template.exs (add to buf at compile time)

	return errors.New("loader name taken")
}

func RunCompiler(name string, buf *[]byte, iexMode bool) (*[]byte, error) {
	if cb, ok := Compiler[name]; ok {
		comp := compileExs{buf: buf}
		cb(&comp, iexMode)
		return comp.buf, nil
	}

	return &[]byte{}, errors.New("compiler not found")
}

func (comp *compileExs) clean() {
	goutil.Clean(*comp.buf)
	bytes.ReplaceAll(*comp.buf, []byte{0}, []byte{})
}

func (comp *compileExs) compStr(q byte, str *[]byte) {
	goutil.StepBytes(str, func(i *int, b func(int) byte, m goutil.StepBytesMethod) bool {
		// get {var object}
		if b(0) == '{' {
			obj := []byte{}
			ind := [2]int{*i}
			m.Inc(1)

			m.Loop(func() bool { return b(0) != '}' }, func() bool {
				if b(0) == '\\' {
					if b(1) == q || common.IsCharAlphaNumeric(b(1)) {
						obj = append(obj, b(0))
					}
					m.Inc(1)
				}

				obj = append(obj, b(0))
				m.Inc(1)
				return true
			})
			ind[1] = *i + 1

			comp.compObj(&obj, true)
			m.Replace(&ind, &obj)

			return true
		}

		return true
	})
}

func (comp *compileExs) compObj(obj *[]byte, inStr bool) {
	v := []byte("#{")

	if !inStr && (*obj)[0] == '@' {
		v = append(v, []byte("cont[:")...)
		if len(*obj) == 1 {
			v = append(v, []byte("body]}")...)
			*obj = v
			return
		}
		*obj = (*obj)[1:]

		v = append(v, regex.CompRE2(`[^\w_]`).RepStrLit(*obj, []byte{})...)
		v = append(v, ']', '}')
		*obj = v
		return
	}

	if (*obj)[0] == '#' {
		v = append(v, []byte("this")...)
		if len(*obj) == 1 {
			v = append(v, '}')
			*obj = v
			return
		}
		*obj = (*obj)[1:]
	} else if (*obj)[0] == '&' {
		v = append(v, []byte("args")...)
		*obj = (*obj)[1:]
	} else {
		if inStr {
			v = append(v, []byte(`App.escapeArg `)...)
		} else {
			v = append(v, []byte(`App.escapeHTML `)...)
		}

		v = append(v, []byte("args")...)
	}

	arg := bytes.Split(*obj, []byte{'.'})
	for _, a := range arg {
		a = regex.CompRE2(`[^\w_]`).RepStrLit(a, []byte{})
		if len(a) != 0 {
			v = append(v, '[', ':')
			v = append(v, a...)
			v = append(v, ']')
		}
	}

	v = append(v, '}')
	*obj = v
}
