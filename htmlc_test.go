package htmlc

import (
	"testing"
)

func Test(t *testing.T) {
	err := Compile("./test/templates", "./test/template.exs")
	if err != nil {
		t.Error(err)
	}
}
