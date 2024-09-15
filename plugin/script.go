package plugin

func init() {
	AddCompiler("js", false, []string{"script"}, func(buf *[]byte, iexMode bool) {
		//todo: compile script {vars} with js `args[key]` variables
		// may add special `js` map for public js args

		//todo: minify js

		// note: this compiler will also apply to external js template files
	})
}
