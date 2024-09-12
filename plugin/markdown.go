package plugin

import (
	"bytes"
	"os"

	regex "github.com/tkdeng/goregex"
	"github.com/tkdeng/goutil"
	"github.com/tkdeng/htmlc/common"
)

func init() {
	addCompiler("md", false, []string{"md", "markdwon"}, func(comp *compileExs, iexMode bool) {
		//todo: compile markdown
	})

	addLoader("md", func(out *os.File, name string, buf *[]byte, usedRandID *[][]byte, iexMode bool) {
		// temp
		return

		randID := regex.JoinBytes(
			'_',
			regex.CompRE2(`[^\w_]`).RepStrLit(bytes.ReplaceAll([]byte(name), []byte{'/'}, []byte{'_'}), []byte{}),
			'_', goutil.URandBytes(16, usedRandID),
		)

		regex.Comp(`(?ms)^([\s\t]*}[\s\t]*#_MAP_MARKDOWN)$`).RepFileStr(out, regex.JoinBytes(
			"\t\t", '"', common.EscapeExsArgs([]byte(name), '"'), '"',
			` => :`, randID, ',',
			"\n$1",
		), false)

		q := '"'
		t := "\t\t"
		if iexMode {
			q = '\''
			t = ""
		}

		regex.Comp(`(?ms)^([\s\t]*end[\s\t]*#_MARKDOWN)$`).RepFileStr(out, regex.JoinBytes(
			"\t", `def `, randID, `(args) do`, '\n',
			t, q, *buf, q,
			"\n\tend",
			"\n$1",
		), false)

		out.Sync()
	})
}
