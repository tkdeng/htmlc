package plugin

func init() {
	addCompiler("exs", false, []string{}, func(comp *compileExs, iexMode bool) {
		//todo: compile exs scripts
	})
}
