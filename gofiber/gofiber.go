package gofiber

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v3"
	core "github.com/gofiber/template"
	"github.com/tkdeng/htmlc"
)

type Engine struct {
	core.Engine

	htmlcEngine  *htmlc.ExsEngine
	templatePath string
}

// New returns a HTML render engine for Fiber
func New(directory string) (*Engine, error) {
	return newEngine(directory, nil)
}

// NewFileSystem returns a HTML render engine for Fiber with file system
func NewFileSystem(fs http.FileSystem) (*Engine, error) {
	return newEngine("/", fs)
}

// newEngine creates a new Engine instance with common initialization logic.
func newEngine(directory string, fs http.FileSystem) (*Engine, error) {
	out := filepath.Clean(directory)
	if out == "" || out == "." {
		out = "temp/template"
		os.MkdirAll("temp", 0755)
	}
	out += ".exs"

	htmlcEngine, err := htmlc.LiveEngine(directory, out)
	if err != nil {
		return nil, err
	}

	engine := &Engine{
		Engine: core.Engine{
			Directory:  directory,
			FileSystem: fs,
			Extension:  htmlc.Ext,
			Funcmap:    make(map[string]interface{}),
		},

		htmlcEngine:  htmlcEngine,
		templatePath: out,
	}

	return engine, nil
}

// Load parses the templates to the engine.
func (e *Engine) Load() error {
	// nothing to load :)
	return nil
}

// Render will execute the template name along with the given values.
func (e *Engine) Render(out io.Writer, name string, binding interface{}, layout ...string) error {
	args := map[string]any{}
	if b, ok := binding.(map[string]any); ok {
		args = b
	} else if b, ok := binding.(fiber.Map); ok {
		args = b
	}

	buf, err := e.htmlcEngine.Render(name, args, layout...)
	if err != nil {
		return err
	}

	if _, err = out.Write(buf); err != nil {
		return fmt.Errorf("render: %w", err)
	}
	return nil
}

// Layout will Render a layout.
func (e *Engine) Layout(out io.Writer, name string, binding interface{}, content ...map[string]string) error {
	args := map[string]any{}
	if b, ok := binding.(map[string]any); ok {
		args = b
	} else if b, ok := binding.(fiber.Map); ok {
		args = b
	}

	buf, err := e.htmlcEngine.Layout(name, args, content...)
	if err != nil {
		return err
	}

	if _, err = out.Write(buf); err != nil {
		return fmt.Errorf("render: %w", err)
	}
	return nil
}

// Widget will Render a widget.
func (e *Engine) Widget(out io.Writer, name string, binding interface{}) error {
	args := map[string]any{}
	if b, ok := binding.(map[string]any); ok {
		args = b
	} else if b, ok := binding.(fiber.Map); ok {
		args = b
	}

	buf, err := e.htmlcEngine.Widget(name, args)
	if err != nil {
		return err
	}

	if _, err = out.Write(buf); err != nil {
		return fmt.Errorf("render: %w", err)
	}
	return nil
}
