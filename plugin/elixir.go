package plugin

func init() {
	//todo: move elixir plugin to core
	// also include a separate compiler for files

	addCompiler("exs", false, []string{}, func(comp *compileExs, iexMode bool) {
		//todo: compile exs scripts
	})
}
