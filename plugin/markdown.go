package plugin

func init() {
	AddCompiler("md", false, []string{"md", "markdown"}, func(buf *[]byte, iexMode bool) {
		//todo: compile markdown
	})
}
