package htmlc

import (
	"bytes"

	regex "github.com/tkdeng/goregex"
	"github.com/tkdeng/goutil"
	"github.com/tkdeng/htmlc/common"
	"github.com/tkdeng/htmlc/plugin"
)

// compileExs compiles elixir in the html buffer
//
// call compile to start main compiler
type compileExs struct {
	buf *[]byte
}

func (comp *compileExs) clean() {
	goutil.Clean(*comp.buf)
	bytes.ReplaceAll(*comp.buf, []byte{0}, []byte{})
}

// compileExs compiles elixir in the html buffer
func (comp *compileExs) compile(compMode string) error {
	comp.clean()

	compHtml := true
	if compMode != "" {
		if html, ok := plugin.CompilerHTML[compMode]; ok {
			compHtml = html
		}
	}

	// compile html template
	if compHtml {
		goutil.StepBytes(comp.buf, func(i *int, b func(int) byte, m goutil.StepBytesMethod) bool {
			// skip <!-- comments -->
			if bytes.Equal(m.GetBuf(4), []byte("<!--")) {
				m.Inc(4)

				m.Loop(func() bool { return !bytes.Equal(m.GetBuf(3), []byte("-->")) }, func() bool {
					m.Inc(1)
					return true
				})
				m.Inc(2)

				return true
			}

			// get "string"
			if b(0) == '"' || b(0) == '\'' {
				q := b(0)
				str := []byte{}
				m.Inc(1)
				ind := [2]int{*i}

				m.Loop(func() bool { return b(0) != q }, func() bool {
					if b(0) == '\\' {
						if b(1) == q || common.IsCharAlphaNumeric(b(1)) {
							str = append(str, b(0))
						}
						m.Inc(1)
					}

					str = append(str, b(0))
					m.Inc(1)
					return true
				})
				ind[1] = *i

				comp.compStr(q, &str)
				m.Replace(&ind, &str)

				return true
			}

			// get {var object}
			if b(0) == '{' {
				obj := []byte{}
				ind := [2]int{*i}
				m.Inc(1)

				m.Loop(func() bool { return b(0) != '}' }, func() bool {
					obj = append(obj, b(0))
					m.Inc(1)
					return true
				})
				ind[1] = *i + 1

				comp.compObj(&obj, false)
				m.Replace(&ind, &obj)

				return true
			}

			//todo: consider adding shortcode support (to call custom elixir functions imported by possible plugins)
			// call with the apply() method with PluginName as the module, and its functions as hooks

			// get <#widget>
			if b(0) == '<' && (b(1) == '#' || (b(1) == '_' && b(2) == '#')) {
				ind := [2]int{*i}

				var close []byte
				if b(1) == '_' && b(2) == '#' {
					m.Inc(3)
					close = []byte("</_#")
				} else {
					m.Inc(2)
					close = []byte("</#")
				}

				name := []byte{}
				m.Loop(func() bool { return b(0) != '>' && b(0) != ' ' }, func() bool {
					name = append(name, b(0))
					m.Inc(1)
					return true
				})
				close = append(close, name...)
				close = append(close, '>')

				if b(0) != '>' {
					m.Inc(1)
				}

				args := map[string][]byte{}
				argName := []byte{}
				m.Loop(func() bool { return b(0) != '>' }, func() bool {
					if b(0) == '=' {
						m.Inc(1)

						str := []byte{}
						if b(0) == '"' || b(0) == '\'' {
							q := b(0)
							m.Inc(1)

							m.Loop(func() bool { return b(0) != q }, func() bool {
								if b(0) == '\\' {
									if b(1) == q || common.IsCharAlphaNumeric(b(1)) {
										str = append(str, b(0))
									}
									m.Inc(1)
								}

								str = append(str, b(0))
								m.Inc(1)
								return true
							})
							ind[1] = *i

							comp.compStr(q, &str)
						} else {
							m.Loop(func() bool { return b(0) != ' ' && b(0) != '>' }, func() bool {
								if b(0) == '\\' {
									if common.IsCharAlphaNumeric(b(1)) {
										str = append(str, b(0))
									}
									m.Inc(1)
								}

								str = append(str, b(0))
								m.Inc(1)
								return true
							})
						}

						if len(argName) != 0 {
							args[string(argName)] = str
						} else {
							args[string(str)] = []byte("true")
						}
						argName = []byte{}

						m.Inc(1)
						return true
					} else if b(0) == ' ' && len(argName) != 0 {
						args[string(argName)] = []byte("true")
						argName = []byte{}
						m.Inc(1)
						return true
					}

					if b(0) != ' ' {
						argName = append(argName, b(0))
					}
					m.Inc(1)

					return true
				})

				if len(argName) != 0 {
					args[string(argName)] = []byte("true")
				}

				if b(-1) == '/' {
					ind[1] = *i + 1
					m.Replace(&ind, comp.embedWedget(name, args))
					return true
				}

				// get widget content

				m.Inc(1)
				cont := []byte{}
				m.Loop(func() bool { return !bytes.Equal(m.GetBuf(len(close)), close) }, func() bool {
					cont = append(cont, b(0))
					m.Inc(1)
					return true
				})
				m.Inc(len(close) - 1)
				ind[1] = *i + 1

				compCont := compileExs{buf: &cont}
				compCont.compile("")

				args["body"] = cont

				m.Replace(&ind, comp.embedWedget(name, args))

				return true
			}

			// get <% elixir %>
			if b(0) == '<' && b(1) == '%' {
				exs := []byte{}
				m.Inc(2)
				ind := [2]int{*i}

				m.Loop(func() bool { return b(0) != '%' && b(1) != '>' }, func() bool {
					exs = append(exs, b(0))
					m.Inc(1)
					return true
				})
				ind[1] = *i
				m.Inc(1)

				if resBuf, err := plugin.RunCompiler("exs", &exs, IexMode); err == nil {
					m.Replace(&ind, resBuf)
				}

				return true
			}

			// get plugin tags
			if b(0) == '<' {
				for tag, ext := range plugin.CompilerTags {
					if bytes.Equal(m.GetBuf(len(tag)+1), []byte("<"+tag)) {
						closeTag := []byte("</" + tag + ">")
						m.Inc(len(tag) + 1)
						if b(0) != '>' && b(0) != ' ' {
							return true
						}

						m.Loop(func() bool { return b(0) != '>' }, func() bool {
							if b(0) == '"' || b(0) == '\'' {
								q := b(0)
								str := []byte{}
								m.Inc(1)
								ind := [2]int{*i}

								m.Loop(func() bool { return b(0) != q }, func() bool {
									if b(0) == '\\' {
										if b(1) == q || common.IsCharAlphaNumeric(b(1)) {
											str = append(str, b(0))
										}
										m.Inc(1)
									}

									str = append(str, b(0))
									m.Inc(1)
									return true
								})
								ind[1] = *i

								comp.compStr(q, &str)
								m.Replace(&ind, &str)
							}

							m.Inc(1)
							return true
						})

						tagBuf := []byte{}
						m.Inc(1)
						ind := [2]int{*i}
						m.Loop(func() bool { return !bytes.Equal(m.GetBuf(len(closeTag)), closeTag) }, func() bool {
							tagBuf = append(tagBuf, b(0))
							m.Inc(1)
							return true
						})
						ind[1] = *i

						if resBuf, err := plugin.RunCompiler(ext, &tagBuf, IexMode); err == nil {
							m.Replace(&ind, resBuf)
						}

						return true
					}
				}
			}

			// get <markdown>, <script>, and <style> tags
			/* if b(0) == '<' && (bytes.Equal(m.GetBuf(3), []byte("<md")) || bytes.Equal(m.GetBuf(9), []byte("<markdown"))) {
				var closeTag []byte
				if bytes.Equal(m.GetBuf(3), []byte("<md")) {
					m.Inc(3)
					closeTag = []byte("</md>")
				} else {
					m.Inc(9)
					closeTag = []byte("</markdown>")
				}

				if b(0) != '>' && b(0) != ' ' {
					return true
				}

				m.Loop(func() bool { return b(0) != '>' }, func() bool {
					if b(0) == '"' || b(0) == '\'' {
						q := b(0)
						str := []byte{}
						m.Inc(1)
						ind := [2]int{*i}

						m.Loop(func() bool { return b(0) != q }, func() bool {
							if b(0) == '\\' {
								if b(1) == q || common.IsCharAlphaNumeric(b(1)) {
									str = append(str, b(0))
								}
								m.Inc(1)
							}

							str = append(str, b(0))
							m.Inc(1)
							return true
						})
						ind[1] = *i

						comp.compStr(q, &str)
						m.Replace(&ind, &str)
					}

					m.Inc(1)
					return true
				})

				md := []byte{}
				m.Inc(1)
				ind := [2]int{*i}
				m.Loop(func() bool { return !bytes.Equal(m.GetBuf(len(closeTag)), closeTag) }, func() bool {
					md = append(md, b(0))
					m.Inc(1)
					return true
				})
				ind[1] = *i

				mdComp := compileExs{buf: &md}
				mdComp.compMD()
				m.Replace(&ind, mdComp.buf)

				return true
			} else if b(0) == '<' && bytes.Equal(m.GetBuf(7), []byte("<script")) {
				m.Inc(7)
				if b(0) != '>' && b(0) != ' ' {
					return true
				}

				m.Loop(func() bool { return b(0) != '>' }, func() bool {
					if b(0) == '"' || b(0) == '\'' {
						q := b(0)
						str := []byte{}
						m.Inc(1)
						ind := [2]int{*i}

						m.Loop(func() bool { return b(0) != q }, func() bool {
							if b(0) == '\\' {
								if b(1) == q || common.IsCharAlphaNumeric(b(1)) {
									str = append(str, b(0))
								}
								m.Inc(1)
							}

							str = append(str, b(0))
							m.Inc(1)
							return true
						})
						ind[1] = *i

						comp.compStr(q, &str)
						m.Replace(&ind, &str)
					}

					m.Inc(1)
					return true
				})

				js := []byte{}
				m.Inc(1)
				ind := [2]int{*i}
				m.Loop(func() bool { return !bytes.Equal(m.GetBuf(9), []byte("</script>")) }, func() bool {
					js = append(js, b(0))
					m.Inc(1)
					return true
				})
				ind[1] = *i

				jsComp := compileExs{buf: &js}
				jsComp.compJS()
				m.Replace(&ind, jsComp.buf)

				return true
			} else if b(0) == '<' && bytes.Equal(m.GetBuf(6), []byte("<style")) {
				m.Inc(6)
				if b(0) != '>' && b(0) != ' ' {
					return true
				}

				m.Loop(func() bool { return b(0) != '>' }, func() bool {
					if b(0) == '"' || b(0) == '\'' {
						q := b(0)
						str := []byte{}
						m.Inc(1)
						ind := [2]int{*i}

						m.Loop(func() bool { return b(0) != q }, func() bool {
							if b(0) == '\\' {
								if b(1) == q || common.IsCharAlphaNumeric(b(1)) {
									str = append(str, b(0))
								}
								m.Inc(1)
							}

							str = append(str, b(0))
							m.Inc(1)
							return true
						})
						ind[1] = *i

						comp.compStr(q, &str)
						m.Replace(&ind, &str)
					}

					m.Inc(1)
					return true
				})

				css := []byte{}
				m.Inc(1)
				ind := [2]int{*i}
				m.Loop(func() bool { return !bytes.Equal(m.GetBuf(8), []byte("</style>")) }, func() bool {
					css = append(css, b(0))
					m.Inc(1)
					return true
				})
				ind[1] = *i

				cssComp := compileExs{buf: &css}
				cssComp.compCSS()
				m.Replace(&ind, cssComp.buf)

				return true
			} */

			return true
		})
	}

	// compile plugin files
	if compMode != "" {
		plugin.RunCompiler(compMode, comp.buf, IexMode)
	}

	// encode to <<base64>>
	ind := [2]int{0, 0}
	goutil.StepBytes(comp.buf, func(i *int, b func(int) byte, m goutil.StepBytesMethod) bool {
		if b(0) == '#' && b(1) == '{' {
			ind[1] = *i

			if ind[1]-ind[0] > 0 {
				bit := common.EncodeBit(&ind, comp.buf, IexMode)
				// bit := regex.JoinBytes(`#{`, common.EncodeBit(&ind, comp.buf), `}`)
				m.Replace(&ind, &bit)
			}

			m.Inc(2)

			ins := 0
			m.Loop(func() bool { return b(0) != '}' || ins != 0 }, func() bool {
				if b(0) == '{' {
					ins++
				} else if b(0) == '}' {
					ins--
				}

				m.Inc(1)
				return true
			})

			ind[0] = *i + 1
			return true
		}

		if *i == len(*comp.buf)-1 {
			ind[1] = *i

			if ind[1]-ind[0] > 0 {
				bit := common.EncodeBit(&ind, comp.buf, IexMode)
				// bit := regex.JoinBytes(`#{`, common.EncodeBit(&ind, comp.buf), `}`)
				m.Replace(&ind, &bit)
			}
		}

		return true
	})

	if len(*comp.buf) != 0 && (*comp.buf)[len(*comp.buf)-1] == '\n' {
		*comp.buf = (*comp.buf)[:len(*comp.buf)-1]
	}

	return nil
}

func (comp *compileExs) compStr(q byte, str *[]byte) {
	goutil.StepBytes(str, func(i *int, b func(int) byte, m goutil.StepBytesMethod) bool {
		// get {var object}
		if b(0) == '{' {
			obj := []byte{}
			ind := [2]int{*i}
			m.Inc(1)

			m.Loop(func() bool { return b(0) != '}' }, func() bool {
				if b(0) == '\\' {
					if b(1) == q || common.IsCharAlphaNumeric(b(1)) {
						obj = append(obj, b(0))
					}
					m.Inc(1)
				}

				obj = append(obj, b(0))
				m.Inc(1)
				return true
			})
			ind[1] = *i + 1

			comp.compObj(&obj, true)
			m.Replace(&ind, &obj)

			return true
		}

		return true
	})
}

func (comp *compileExs) compObj(obj *[]byte, inStr bool) {
	v := []byte("#{")

	if !inStr && (*obj)[0] == '@' {
		v = append(v, []byte("cont[:")...)
		if len(*obj) == 1 {
			v = append(v, []byte("body]}")...)
			*obj = v
			return
		}
		*obj = (*obj)[1:]

		v = append(v, regex.CompRE2(`[^\w_]`).RepStrLit(*obj, []byte{})...)
		v = append(v, ']', '}')
		*obj = v
		return
	}

	if (*obj)[0] == '#' {
		v = append(v, []byte("this")...)
		if len(*obj) == 1 {
			v = append(v, '}')
			*obj = v
			return
		}
		*obj = (*obj)[1:]
	} else if (*obj)[0] == '&' {
		v = append(v, []byte("args")...)
		*obj = (*obj)[1:]
	} else {
		if inStr {
			v = append(v, []byte(`App.escapeArg `)...)
		} else {
			v = append(v, []byte(`App.escapeHTML `)...)
		}

		v = append(v, []byte("args")...)
	}

	arg := bytes.Split(*obj, []byte{'.'})
	for _, a := range arg {
		a = regex.CompRE2(`[^\w_]`).RepStrLit(a, []byte{})
		if len(a) != 0 {
			v = append(v, '[', ':')
			v = append(v, a...)
			v = append(v, ']')
		}
	}

	v = append(v, '}')
	*obj = v
}

func (comp *compileExs) embedWedget(name []byte, args map[string][]byte) *[]byte {
	name = bytes.ReplaceAll(name, []byte{'.'}, []byte{'/'})

	buf := regex.JoinBytes(`#{App.widget "`, name, `", Map.merge(args, %{`, '\n')

	for key, val := range args {
		buf = regex.JoinBytes(buf, '\t', regex.CompRE2(`[^\w_]`).RepStrLit([]byte(key), []byte{}), `: `)

		if regex.CompRE2(`^([0-9]+(\.[0-9]+|)|true|false|nil|null)$`).Match(val) {
			if bytes.Equal(val, []byte("null")) {
				buf = append(buf, []byte("nil")...)
			} else {
				buf = append(buf, val...)
			}
		} else {
			buf = regex.JoinBytes(buf, '"', common.EscapeExsArgs(val, '"'), '"')
		}

		buf = append(buf, ',', '\n')
	}

	buf = append(buf, '}', ')', '}')

	return &buf
}

//* lang compilers
// handle like a normal compiler

/* func (comp *compileExs) compExs() {
	//todo: compile exs scripts
}

func (comp *compileExs) compMD() {
	//todo: compile markdown
}

func (comp *compileExs) compJS() {
	//todo: compile js script
}

func (comp *compileExs) compCSS() {
	//todo: compile css style
} */
