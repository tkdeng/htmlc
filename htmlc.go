package htmlc

import (
	"errors"
	"os"
	"path/filepath"
)

//todo: allow dynamic compiling and listening for file changes

func Compile(src string, out string) error {
	if stat, err := os.Stat(src); err != nil || !stat.IsDir() {
		return errors.Join(err, os.ErrNotExist, errors.New("src \""+src+"\": directory not found"))
	}

	if err := os.MkdirAll(filepath.Dir(out), 0755); err != nil {
		return err
	}

	//temp
	os.WriteFile(out, []byte{}, 0755)

	//todo: compile templates to elixir (dist)
	// may detect contents for diff between pages/layouts/widgets
	//
	// `<!DOCTYPE> | <html></html>` tags for layouts
	//
	// `<_@body></_@body>` @embed tags for pages
	//
	// outside html tags and @embed tags for widgets
	//
	// for error type, detect files names `error.html` and `404.html`
	// automatic 404 errors should try to reference relative to the directory the file was called from
	// error type can be treated in a similar way to page type

	//todo: may consider returning a struct for render method
	// may also add a file watcher, and have a stop method in the struct

	return nil
}
