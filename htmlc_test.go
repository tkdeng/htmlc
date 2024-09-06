package htmlc

import (
	"testing"
)

//!test: ./htmlc --no-compile -o="./test/template.exs" 3000

func Test(t *testing.T) {
	err := Compile("./test/templates", "./test/template.exs")
	if err != nil {
		t.Error(err)
	}
}
