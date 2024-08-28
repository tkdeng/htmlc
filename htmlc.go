package htmlc

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	_ "embed"

	"github.com/AspieSoft/go-regex/v8"
	"github.com/tkdeng/goutil"
)

//todo: allow dynamic compiling and listening for file changes

//go:embed template.exs
var template []byte

func Compile(src string, out string) error {
	if stat, err := os.Stat(src); err != nil || !stat.IsDir() {
		return errors.Join(err, os.ErrNotExist, errors.New("src \""+src+"\": directory not found"))
	}

	if err := os.MkdirAll(filepath.Dir(out), 0755); err != nil {
		return err
	}

	outfile, err := os.OpenFile(out, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer outfile.Close()

	outfile.Write(template)
	// os.WriteFile(out, template, 0755)

	compileDir(outfile, src, "", 'r')

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

func compileDir(out *os.File, dir string, name string, dirType byte) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			if path, err := goutil.JoinPath(dir, file.Name()); err == nil {
				if dirType == 'r' {
					switch file.Name() {
					case "layouts", "layout":
						compileDir(out, path, "", 'l')
					case "widgets", "widget":
						compileDir(out, path, "", 'w')
					case "pages", "page":
						compileDir(out, path, "", 'p')
					case "errors", "error":
						compileDir(out, path, "", 'e')
					default:
						compileDir(out, path, file.Name(), 'd')
					}
				} else {
					compileDir(out, path, name+"/"+file.Name(), dirType)
				}
			}
		} else if strings.HasSuffix(file.Name(), ".html") {
			if path, err := goutil.JoinPath(dir, file.Name()); err == nil {
				if buf, err := os.ReadFile(path); err == nil {
					n := name
					if n != "" {
						n += "/"
					}
					n += strings.TrimSuffix(file.Name(), ".html")

					if dirType == 'l' {
						loadLayout(out, n, buf)
					} else if dirType == 'w' {
						loadWidget(out, n, buf)
					} else if dirType == 'p' || dirType == 'e' {
						buf_page := map[string][]byte{}

						//todo: pull `<% elixir %>` and `<md>` (markdown) tags as reference to prevent regex conflicts
						// also pull `<script>` and `<style>` tags

						regex.Comp(`(?s)<(_?@)([\w_-]+)>(.*)</\1\2>\r?\n?`).RepFunc(buf, func(data func(int) []byte) []byte {
							name := string(data(2))
							if _, ok := buf_page[name]; !ok {
								buf_page[name] = data(3)
							} else {
								buf_page[name] = append(buf_page[name], data(3)...)
							}
							return []byte{}
						}, true)

						if len(buf_page) == 0 {
							buf_page["body"] = buf
						}

						loadPage(out, n, buf_page)
					} else {
						var buf_layout []byte
						buf_page := map[string][]byte{}

						//todo: pull `<% elixir %>` and `<md>` (markdown) tags as reference to prevent regex conflicts
						// also pull `<script>` and `<style>` tags

						buf = regex.Comp(`(?s)(<!DOCTYPE(?:\s+[^>]*|)>|<html(?:\s+[^>]*|)>(.*)</html>)\r?\n?`).RepFunc(buf, func(data func(int) []byte) []byte {
							buf_layout = append(buf_layout, data(0)...)
							return []byte{}
						})

						buf = regex.Comp(`(?s)<(_?@)([\w_-]+)>(.*)</\1\2>\r?\n?`).RepFunc(buf, func(data func(int) []byte) []byte {
							id := string(data(2))
							if _, ok := buf_page[id]; !ok {
								buf_page[id] = data(3)
							} else {
								buf_page[id] = append(buf_page[id], data(3)...)
							}
							return []byte{}
						})

						if len(buf_layout) != 0 {
							loadLayout(out, n, buf_layout)
						}

						if len(buf_page) != 0 {
							loadPage(out, n, buf_page)
						}

						if len(buf) != 0 && len(regex.Comp(`(?s)[\s\r\n\t ]+`).RepStrLit(buf, []byte{})) != 0 {
							loadWidget(out, n, buf)
						}
					}
				}
			}
		}
	}

	return nil
}

func loadLayout(out *os.File, name string, buf []byte) {
	//todo: compile layout
	// fmt.Println("layout:", name)
}

func loadWidget(out *os.File, name string, buf []byte) {
	//todo: compile widget
	// fmt.Println("widget:", name)
}

func loadPage(out *os.File, name string, buf map[string][]byte) {
	//todo: compile page
	// fmt.Println("page:", name)
}
