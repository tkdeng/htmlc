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

				// if basic widget `<#name/>`
				if bytes.HasSuffix(name, []byte{'/'}) || b(-1) == '/' || (b(0) == '/' && b(1) == '>') {
					if b(0) == '/' && b(1) == '>' {
						m.Inc(1)
					}

					name = bytes.TrimSuffix(name, []byte{'/'})
					ind[1] = *i + 1
					m.Replace(&ind, comp.embedWedget(name, map[string][]byte{}))
					return true
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

				if bytes.Equal(argName, []byte{'/'}) {
					argName = []byte{}
				}

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

				compExs := compileExs{buf: &exs}
				compExs.compExs()
				m.Replace(&ind, compExs.buf)

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

						if cb, ok := plugin.Compiler[ext]; ok {
							cb(&tagBuf, IexMode)
							m.Replace(&ind, &tagBuf)
						}

						return true
					}
				}
			}

			// skip <markdown>, <script>, and <style> tags (if missed by plugin)
			if b(0) == '<' && (bytes.Equal(m.GetBuf(3), []byte("<md")) || bytes.Equal(m.GetBuf(9), []byte("<markdown")) || bytes.Equal(m.GetBuf(7), []byte("<script")) || bytes.Equal(m.GetBuf(6), []byte("<style"))) {
				var closeTag []byte
				if bytes.Equal(m.GetBuf(3), []byte("<md")) {
					m.Inc(3)
					closeTag = []byte("</md>")
				} else if bytes.Equal(m.GetBuf(9), []byte("<markdown")) {
					m.Inc(9)
					closeTag = []byte("</markdown>")
				} else if bytes.Equal(m.GetBuf(7), []byte("<script")) {
					m.Inc(7)
					closeTag = []byte("</script>")
				} else if bytes.Equal(m.GetBuf(6), []byte("<style")) {
					m.Inc(7)
					closeTag = []byte("</style>")
				} else {
					return true
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

				m.Inc(1)
				m.Loop(func() bool { return !bytes.Equal(m.GetBuf(len(closeTag)), closeTag) }, func() bool {
					m.Inc(1)
					return true
				})

				return true
			}

			return true
		})
	}

	// compile plugin files
	if compMode != "" {
		if cb, ok := plugin.Compiler[compMode]; ok {
			cb(comp.buf, IexMode)
		}
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

	if len(args) == 0 {
		buf := regex.JoinBytes(`#{App.widget "`, name, `", args}`)
		return &buf
	}

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

func (comp *compileExs) compExs() {
	//todo: compile exs scripts
}
