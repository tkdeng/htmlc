package htmlc

import (
	regex "github.com/tkdeng/goregex"
	"github.com/tkdeng/goutil"
)

func isCharAlphaNumeric(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

var regEscExsArgs *regex.Regexp = regex.Comp(`([\\]*)([\\#])`)

func EscapeExsArgs(html []byte, quote ...byte) []byte {
	html = goutil.HTML.EscapeArgs(html, quote...)

	return regEscExsArgs.RepFunc(html, func(data func(int) []byte) []byte {
		if len(data(1))%2 == 0 && data(2)[0] == '#' {
			return regex.JoinBytes(data(1), '\\', data(2))
		}
		if data(2)[0] == '#' {
			return append(data(1), data(2)...)
		}
		return data(0)
	})
}
