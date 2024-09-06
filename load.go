package htmlc

import (
	"bytes"
	"os"

	regex "github.com/tkdeng/goregex"
	"github.com/tkdeng/goutil"
)

// loadLayout loads a new layout into the output file
func loadLayout(out *os.File, name string, buf *[]byte, usedRandID *[][]byte) {
	randID := regex.JoinBytes(
		'_',
		regex.CompRE2(`[^\w_]`).RepStrLit(bytes.ReplaceAll([]byte(name), []byte{'/'}, []byte{'_'}), []byte{}),
		'_', goutil.URandBytes(16, usedRandID),
	)

	regex.Comp(`(?ms)^([\s\t]*}[\s\t]*#_MAP_LAYOUT)$`).RepFileStr(out, regex.JoinBytes(
		"\t\t", '"', goutil.HTML.EscapeArgs([]byte(name), '"'), '"',
		` => :`, randID, ',',
		"\n$1",
	), false)

	regex.Comp(`(?ms)^([\s\t]*end[\s\t]*#_LAYOUT)$`).RepFileStr(out, regex.JoinBytes(
		"\t", `def `, randID, `(args, cont) do`, '\n',
		// '\'', goutil.HTML.EscapeArgs(*buf, '\''), '\'',
		// '"', goutil.HTML.EscapeArgs(*buf, '"'), '"',

		"\t\t", '"', *buf, '"',
		// "\t\t", '\'', *buf, '\'',
		// `"#{<<`, *buf, `>>}"`,

		"\n\tend",
		"\n$1",
	), false)
}

// loadWidget loads a new widget into the output file
func loadWidget(out *os.File, name string, buf *[]byte, usedRandID *[][]byte) {
	randID := regex.JoinBytes(
		'_',
		regex.CompRE2(`[^\w_]`).RepStrLit(bytes.ReplaceAll([]byte(name), []byte{'/'}, []byte{'_'}), []byte{}),
		'_', goutil.URandBytes(16, usedRandID),
	)

	regex.Comp(`(?ms)^([\s\t]*}[\s\t]*#_MAP_WIDGET)$`).RepFileStr(out, regex.JoinBytes(
		"\t\t", '"', goutil.HTML.EscapeArgs([]byte(name), '"'), '"',
		` => :`, randID, ',',
		"\n$1",
	), false)

	regex.Comp(`(?ms)^([\s\t]*end[\s\t]*#_WIDGET)$`).RepFileStr(out, regex.JoinBytes(
		"\t", `def `, randID, `(args) do`, '\n',
		// '\'', goutil.HTML.EscapeArgs(*buf, '\''), '\'',
		// '"', goutil.HTML.EscapeArgs(*buf, '"'), '"',

		"\t\t", '"', *buf, '"',
		// "\t\t", '\'', *buf, '\'',
		// `"#{<<`, *buf, `>>}"`,

		"\n\tend",
		"\n$1",
	), false)
}

// loadPage loads a new page into the output file
func loadPage(out *os.File, name string, buf *map[string][]byte, usedRandID *[][]byte) {
	randID := regex.JoinBytes(
		'_',
		regex.CompRE2(`[^\w_]`).RepStrLit(bytes.ReplaceAll([]byte(name), []byte{'/'}, []byte{'_'}), []byte{}),
		'_', goutil.URandBytes(16, usedRandID),
	)

	regex.Comp(`(?ms)^([\s\t]*}[\s\t]*#_MAP_PAGE)$`).RepFileStr(out, regex.JoinBytes(
		"\t\t", '"', goutil.HTML.EscapeArgs([]byte(name), '"'), '"',
		` => :`, randID, ',',
		"\n$1",
	), false)

	resBuf := regex.JoinBytes(
		"\t", `def `, randID, `(layout, args) do`, '\n',
		"\t\t", `App.layout layout, args, %{`, '\n',
	)

	for key, val := range *buf {
		resBuf = regex.JoinBytes(resBuf,
			"\t\t\t", regex.CompRE2(`[^\w_]`).RepStrLit([]byte(key), []byte{}), ": ",
			// '"', goutil.HTML.EscapeArgs(val, '"'), '"', ",\n",
			// '"', goutil.HTML.EscapeArgs(val, '"'), '"', ",\n",

			'"', val, '"', ",\n",
			// '\'', val, '\'', ",\n",
			// `"#{<<`, val, `>>}"`, ",\n",
		)
	}

	resBuf = regex.JoinBytes(resBuf,
		"\t\t}",
		"\n\tend",
		"\n$1",
	)

	regex.Comp(`(?ms)^([\s\t]*end[\s\t]*#_PAGE)$`).RepFileStr(out, resBuf, false)
}
