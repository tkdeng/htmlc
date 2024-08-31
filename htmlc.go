package htmlc

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "embed"

	"github.com/tkdeng/goregex"
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

	if err = compileDir(outfile, src, "", 'r'); err != nil {
		return err
	}

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
						if err = compileDir(out, path, "", 'l'); err != nil {
							return err
						}
					case "widgets", "widget":
						if err = compileDir(out, path, "", 'w'); err != nil {
							return err
						}
					case "pages", "page":
						if err = compileDir(out, path, "", 'p'); err != nil {
							return err
						}
					case "errors", "error":
						if err = compileDir(out, path, "", 'e'); err != nil {
							return err
						}
					default:
						if err = compileDir(out, path, file.Name(), 'd'); err != nil {
							return err
						}
					}
				} else {
					if err = compileDir(out, path, name+"/"+file.Name(), dirType); err != nil {
						return err
					}
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

					strList := pullStrings(&buf)

					if dirType == 'l' {
						if err := compileExs(&buf, &strList); err != nil {
							return err
						}
						loadLayout(out, n, &buf)
					} else if dirType == 'w' {
						if err := compileExs(&buf, &strList); err != nil {
							return err
						}
						loadWidget(out, n, &buf)
					} else if dirType == 'p' || dirType == 'e' {
						buf_page := map[string][]byte{}

						regex.Comp(`(?s)<(_?@)([\w_-]+)>(.*)</\1\2>\r?\n?`).RepFunc(buf, func(data func(int) []byte) []byte {
							name := string(data(2))
							if _, ok := buf_page[name]; !ok {
								buf_page[name] = data(3)
							} else {
								buf_page[name] = append(buf_page[name], data(3)...)
							}
							return []byte{}
						})

						if len(buf_page) == 0 {
							buf_page["body"] = buf
						}

						for k, b := range buf_page {
							if err := compileExs(&b, &strList); err != nil {
								return err
							}
							buf_page[k] = b
						}
						loadPage(out, n, &buf_page)
					} else {
						var buf_layout []byte
						buf_page := map[string][]byte{}

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
							if err := compileExs(&buf_layout, &strList); err != nil {
								return err
							}
							loadLayout(out, n, &buf_layout)
						}

						if len(buf_page) != 0 {
							for k, b := range buf_page {
								if err := compileExs(&b, &strList); err != nil {
									return err
								}
								buf_page[k] = b
							}
							loadPage(out, n, &buf_page)
						}

						if len(buf) != 0 && len(regex.Comp(`(?s)[\s\r\n\t ]+`).RepStrLit(buf, []byte{})) != 0 {
							if err := compileExs(&buf, &strList); err != nil {
								return err
							}
							loadWidget(out, n, &buf)
						}
					}
				}
			}
		}
	}

	return nil
}

// pullStrings pulls out strings and creates references to them
//
// @restore: alternate functionality to restore pulled strings to references
func pullStrings(buf *[]byte, restore ...*[][]byte) [][]byte {
	if len(restore) != 0 {
		*buf = regex.Comp(`(?s)(["'\'])str:([A-Za-z0-9]+)\1`).RepFunc(*buf, func(data func(int) []byte) []byte {
			if i, err := strconv.ParseUint(string(data(2)), 36, 64); err == nil {
				return regex.JoinBytes(data(1), (*restore[0])[i], data(1))
			}
			return data(0)
		})

		return nil
	}

	strList := [][]byte{}

	*buf = regex.Comp(`(?s)(["'\'])((?:\\(?:\\|\1)|.)*?)\1`).RepFunc(*buf, func(data func(int) []byte) []byte {
		strList = append(strList, data(2))
		return regex.JoinBytes(data(1), `str:`, strconv.FormatUint(uint64(len(strList)-1), 36), data(1))
	})

	return strList
}

// compileExs compiles elixir in the html buffer
func compileExs(buf *[]byte, strList *[][]byte) error {
	var err error

	//todo: may have to consider reading individual bytes in for loop instead of using regex
	// note: may only need in some parts of the exs compiler

	// compile <% elixir %> to #{} for strings
	*buf = regex.Comp(`(?s)<%(.*?)%>`).RepFunc(*buf, func(data func(int) []byte) []byte {
		if err != nil {
			return nil
		}

		// ensure {brackets} are properly opened and closed (to prevent escaping '#{}' in elixir strings)
		count := 0
		regex.Comp(`[{}]`).RepFunc(data(1), func(data func(int) []byte) []byte {
			if err != nil {
				return nil
			}

			if data(0)[0] == '{' {
				count++
			} else {
				count--
				if count < 0 {
					err = errors.New("invalid exs script: '}' was not opened")
					return nil
				}
			}

			return []byte{}
		})

		if err != nil {
			return nil
		}

		if count > 0 {
			err = errors.New("invalid exs script: '{' was not closed")
			return nil
		} else if count < 0 {
			err = errors.New("invalid exs script: '}' was not opened")
			return nil
		}

		return regex.JoinBytes([]byte("#{"), data(1), '}')
	})

	//todo: detect ending #{do} and starting #{end} tags and merge contents into inner string
	// may use previous regex to detect unclosed #{do} keywords

	// compile html components
	*buf = regex.Comp(`(?s)<(_?#)([\w_\-\.]+)(\s+.*?|)(/>|>(.*?)</\1\2>)`).RepFunc(*buf, func(data func(int) []byte) []byte {
		args := map[string][]byte{
			"body": {},
		}

		argList := regex.Comp(`([\w_-]+\s*=\s*(?:".*?"|'.*?'|[\w_-]*))`).Split(data(3))
		for _, arg := range argList {
			if len(bytes.TrimSpace(arg)) != 0 {
				a := bytes.SplitN(arg, []byte{'='}, 2)
				pullStrings(&a[1], strList)

				if len(a[1]) >= 3 {
					if (a[1][0] == '"' && a[1][len(a[1])-1] == '"') || (a[1][0] == '\'' && a[1][len(a[1])-1] == '\'') {
						a[1] = a[1][1 : len(a[1])-1]
					}
				}

				args[string(a[0])] = goutil.HTML.EscapeArgs(a[1])
			}
		}

		if !bytes.Equal(data(4), []byte("/>")) {
			args["body"] = goutil.HTML.EscapeArgs(regex.Comp(`^(?s)>(.*)</_?#[\w_-]+>$`).RepStr(data(4), []byte("$1")))
		}

		b := regex.JoinBytes(`#{App.widget "`, bytes.ReplaceAll(data(2), []byte{'.'}, []byte{'/'}), `", Map.merge(args, %{`)

		for k, v := range args {
			b = append(b, regex.JoinBytes(k, `: "`, v, `",`)...)
		}

		return append(b, []byte(`})}`)...)
	})

	// compile embed variables
	*buf = regex.Comp(`\{@([\w_]+)\}`).RepStr(*buf, []byte(`#{cont[:$1]}`))

	//todo: compile if/each attributes to elixir (include else if and else)
	// example: <header if="args.n == 2"> <else if=""/> <else/> </header>
	// <ul each="menu"> <a href="{#url}">{#name}</a> </ul>
	// <ul each="menu"> <a href="<% this.url %>"><% this.name %></a> </ul>
	// <div each="menu"> {#} </div> | <div each="menu"> <% this %> </div>
	// (note: may have to limit nesting logic)
	// may have to consider pullStrings on inner content of each loops

	pullStrings(buf, strList)

	//todo: use func regex to compile {var.a} variables and more
	// also escape html unless {&arg} to allow html
	// and detect if in an html attr string, to escape htmlattr instead of html
	// (note: may run alternate version before running pullStrings restore)

	// compile arg variables
	//todo: replace with better variable method
	*buf = regex.Comp(`\{([\w_]+)\}`).RepStr(*buf, []byte(`#{args[:$1]}`))

	return err
}

// loadLayout loads a new layout into the output file
func loadLayout(out *os.File, name string, buf *[]byte) {
	//todo: compile layout
	// fmt.Println("layout:", name)
}

// loadWidget loads a new widget into the output file
func loadWidget(out *os.File, name string, buf *[]byte) {
	//todo: compile widget
	// fmt.Println("widget:", name)
}

// loadPage loads a new page into the output file
func loadPage(out *os.File, name string, buf *map[string][]byte) {
	//todo: compile page
	// fmt.Println("page:", name)
}
