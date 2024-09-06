package htmlc

import (
	"bytes"
	"encoding/base64"

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

func encodeBit(ind *[2]int, buf *[]byte) []byte {
	if IexMode {
		return EscapeExsArgs((*buf)[ind[0]:ind[1]], '\'')
	}

	/* bit := []byte(fmt.Sprintf("%#v", (*buf)[ind[0]:ind[1]]))
	bit = bytes.TrimPrefix(bit, []byte(`[]byte{`))
	bit = bytes.TrimSuffix(bit, []byte(`}`)) */

	bit := make([]byte, (ind[1]-ind[0])*2)
	base64.StdEncoding.Encode(bit, (*buf)[ind[0]:ind[1]])
	bit = bytes.ReplaceAll(bit, []byte{0}, []byte{})

	return regex.JoinBytes(`<<`, bit, `>>`)
}
