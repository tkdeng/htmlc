package plugin

import (
	"errors"
)

var Compiler = map[string]func(buf *[]byte, iexMode bool){}
var CompilerHTML = map[string]bool{}
var CompilerTags = map[string]string{}

// AddCompiler adds a lang compiler method
//
// @name: .ext
//
// @html: include initial html compiler
//
// @tags: html <tags> to compile content </tags>
//
// @cb: method to run when compiling
func AddCompiler(name string, html bool, tags []string, cb func(buf *[]byte, iexMode bool)) error {
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
