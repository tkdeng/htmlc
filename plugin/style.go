package plugin

func init() {
	AddCompiler("css", false, []string{"style"}, func(buf *[]byte, iexMode bool) {
		//todo: compile style {vars} with css `--args` variables
		// may add special `css` map for public css vars

		//todo: minify css

		// note: this compiler will also apply to external css template files
	})
}
