package htmlc

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "embed"

	regex "github.com/tkdeng/goregex"
	"github.com/tkdeng/goutil"
	"github.com/tkdeng/htmlc/plugin"
)

const Ext = "htmlc"

//go:embed template.exs
var template []byte

var IexMode bool = false

// Compile will compile a directory into an exs template
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

	if IexMode {
		rmList := []string{
			// bug: strange error where `\[\]` is not recognized correctly
			// remove json import
			`Mix\.install\(.:jason.\)[\r\n]?`,

			// remove App.listen function
			`def\s+listen\(\)\s*do([\r\n]|.)*?^[\s\t]*end\s*#_LISTEN`,

			// remove App.listen function call
			`App\.listen\(\)`,
		}

		regex.Comp(`(?m)^[\s\t]*(`+strings.Join(rmList, "|")+`)[\r\n]?`).RepFileStr(outfile, []byte{}, true, 1024*1024)
		outfile.Sync()
	}

	usedRandID := [][]byte{}

	if err = compileDir(outfile, src, "", 'r', &usedRandID); err != nil {
		return err
	}

	return nil
}

func compileDir(out *os.File, dir string, name string, dirType byte, usedRandID *[][]byte) error {
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
						if err = compileDir(out, path, "", 'l', usedRandID); err != nil {
							return err
						}
					case "widgets", "widget":
						if err = compileDir(out, path, "", 'w', usedRandID); err != nil {
							return err
						}
					case "pages", "page":
						if err = compileDir(out, path, "", 'p', usedRandID); err != nil {
							return err
						}
					case "errors", "error":
						if err = compileDir(out, path, "", 'e', usedRandID); err != nil {
							return err
						}
					default:
						if err = compileDir(out, path, file.Name(), 'd', usedRandID); err != nil {
							return err
						}
					}
				} else {
					if err = compileDir(out, path, name+"/"+file.Name(), dirType, usedRandID); err != nil {
						return err
					}
				}
			}
		} else if strings.HasSuffix(file.Name(), ".html") || strings.HasSuffix(file.Name(), "."+Ext) {
			if path, err := goutil.JoinPath(dir, file.Name()); err == nil {
				if buf, err := os.ReadFile(path); err == nil {
					n := name
					if n != "" {
						n += "/"
					}
					n += strings.TrimSuffix(strings.TrimSuffix(file.Name(), ".html"), "."+Ext)

					if dirType == 'l' {
						comp := compileExs{buf: &buf}
						if err := comp.compile(""); err != nil {
							return err
						}
						loadLayout(out, n, &buf, usedRandID)
					} else if dirType == 'w' {
						comp := compileExs{buf: &buf}
						if err := comp.compile(""); err != nil {
							return err
						}
						loadWidget(out, n, &buf, usedRandID)
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
							comp := compileExs{buf: &b}
							if err := comp.compile(""); err != nil {
								return err
							}
							buf_page[k] = b
						}
						loadPage(out, n, &buf_page, usedRandID)
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
							comp := compileExs{buf: &buf_layout}
							if err := comp.compile(""); err != nil {
								return err
							}
							loadLayout(out, n, &buf_layout, usedRandID)
						}

						if len(buf_page) != 0 {
							for k, b := range buf_page {
								comp := compileExs{buf: &b}
								if err := comp.compile(""); err != nil {
									return err
								}
								buf_page[k] = b
							}
							loadPage(out, n, &buf_page, usedRandID)
						}

						if len(buf) != 0 && len(regex.Comp(`(?s)[\s\r\n\t ]+`).RepStrLit(buf, []byte{})) != 0 {
							comp := compileExs{buf: &buf}
							if err := comp.compile(""); err != nil {
								return err
							}
							loadWidget(out, n, &buf, usedRandID)
						}
					}
				}
			}
		} else if strings.HasSuffix(file.Name(), ".exs") {
			if path, err := goutil.JoinPath(dir, file.Name()); err == nil {
				if buf, err := os.ReadFile(path); err == nil {
					n := name
					if n != "" {
						n += "/"
					}
					n += strings.TrimSuffix(strings.TrimSuffix(file.Name(), ".html"), "."+Ext)

					comp := compileExs{buf: &buf}
					comp.compExs()

					loadExs(out, n, comp.buf, usedRandID)
				}
			}
		} else if regex.Comp(`\.([\w_-]+)$`).Match([]byte(file.Name())) {
			var ext string
			regex.Comp(`\.([\w_-]+)$`).RepFunc([]byte(file.Name()), func(data func(int) []byte) []byte {
				ext = string(data(1))
				return nil
			})

			if path, err := goutil.JoinPath(dir, file.Name()); err == nil {
				if buf, err := os.ReadFile(path); err == nil {
					n := name
					if n != "" {
						n += "/"
					}
					n += strings.TrimSuffix(file.Name(), "."+ext)

					comp := compileExs{buf: &buf}
					if err := comp.compile(ext); err != nil {
						return err
					}

					loadPluginFile(out, n, ext, comp.buf, usedRandID)
				}
			}
		}
	}

	return nil
}

// LiveEngine will Compile a directory and creates a live updating Engine
func LiveEngine(src string, out string) (*ExsEngine, error) {
	if err := Compile(src, out); err != nil {
		return nil, err
	}

	engine, err := Engine(out)
	if err != nil {
		return nil, err
	}

	var mu sync.Mutex
	lastRecompile := time.Now().UnixMilli()
	recompile := func() {
		now := time.Now().UnixMilli()
		if now-lastRecompile < 800 {
			return
		}
		lastRecompile = now

		fmt.Print("\033[33m htmlc engine restarting...\033[0m   \r")

		time.Sleep(1000 * time.Millisecond)

		if !mu.TryLock() {
			return
		}
		defer mu.Unlock()

		engine.Restart(src)

		fmt.Print("                              \r")
	}

	watcher := goutil.FileWatcher()

	watcher.OnAny = func(path, op string) {
		if strings.HasSuffix(path, ".html") || strings.HasSuffix(path, "."+Ext) {
			go recompile()
			return
		}

		for name := range plugin.Compiler {
			if strings.HasSuffix(path, "."+name) {
				go recompile()
				return
			}
		}
	}

	if err := watcher.WatchDir(src); err != nil {
		return engine, err
	}

	return engine, nil
}
