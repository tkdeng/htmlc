package htmlc

import (
	"os"
)

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
