package htmlc

import (
	"errors"
	"os"
	"path/filepath"
)

//todo: make htmlc compiler a separate git repo/project
// also allow repo to be imported

//todo: allow dynamic compiling and listening for file changes

func Compile(src string, dist string) error {
	if stat, err := os.Stat(src); err != nil || !stat.IsDir() {
		return errors.Join(err, os.ErrNotExist, errors.New("src \""+src+"\": directory not found"))
	}

	if err := os.MkdirAll(filepath.Dir(dist), 0755); err != nil {
		return err
	}

	//note: dist is the dist file "not the directory"
	//temp
	os.WriteFile(dist, []byte{}, 0755)

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
