package htmlc

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
	"sync"
	"time"

	regex "github.com/tkdeng/goregex"
	"github.com/tkdeng/goutil"
)

type ExsEngine struct {
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.ReadCloser
	out     <-chan []byte
	running *bool
	mu      sync.Mutex
}

type Map map[string]any

var ErrorStopped error = errors.New("elixir cmd stopped running")
var ErrorTimeout error = errors.New("request timed out")

func Engine(file string) (*ExsEngine, error) {
	cmd := exec.Command(`iex`, file)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return &ExsEngine{}, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		stdin.Close()
		return &ExsEngine{}, err
	}

	out := make(chan []byte)
	running := true
	ready := false

	go func() {
		for {
			if !running {
				break
			}

			buf := make([]byte, 1024*1024)
			n, err := stdout.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				continue
			}

			if !ready {
				continue
			}

			buf = buf[:n]

			if regex.CompRE2(`(?im)^\s*iex\s*\(.*?\)\s*>`).Match(buf) || regex.CompRE2(`(?im)^\s*nil`).Match(buf) {
				continue
			}

			hasOut := false
			regex.CompRE2(`~c"((?:\\"|.)*?)"`).RepFunc(buf, func(data func(int) []byte) []byte {
				if hasOut {
					return nil
				}
				out <- bytes.ReplaceAll(data(1), []byte("\\n"), []byte("\n"))
				hasOut = true
				return nil
			})

			if !hasOut {
				out <- []byte("<h1>Error 500</h1><h2>Elixir Template Error!</h2>")
			}
		}

		running = false
		stdin.Close()
		stdout.Close()
	}()

	err = cmd.Start()
	if err != nil {
		running = false
		stdin.Close()
		stdout.Close()
		return &ExsEngine{}, err
	}

	time.Sleep(1 * time.Second)

	ready = true

	return &ExsEngine{
		cmd:     cmd,
		stdin:   stdin,
		stdout:  stdout,
		out:     out,
		running: &running,
	}, nil
}

func (exs *ExsEngine) Render(name string, args Map, layout ...string) ([]byte, error) {
	if !*exs.running {
		return []byte{}, ErrorStopped
	}

	lay := []byte("layout")
	if len(layout) != 0 {
		lay = EscapeExsArgs([]byte(layout[0]), '"')
	}

	in := regex.JoinBytes(
		`App.render "`, EscapeExsArgs([]byte(name), '"'), '"',
		`, "`, lay, `", `, renderExsMap(args), '\n',
	)

	in = append(in, []byte("\n")...)

	var out []byte

	go func() {
		exs.mu.Lock()
		defer exs.mu.Unlock()
		exs.stdin.Write(in)
		out = <-exs.out
		time.Sleep(10 * time.Nanosecond)
	}()

	start := time.Now().UnixMilli()
	for len(out) == 0 && time.Now().UnixMilli()-start < 10*1000 {
		time.Sleep(10 * time.Nanosecond)
	}

	if len(out) == 0 {
		return []byte{}, ErrorTimeout
	}

	return out, nil
}

func renderExsMap(args Map) []byte {
	buf := []byte("%{")

	for key, val := range args {
		buf = append(buf, []byte(key+": ")...)
		if v, ok := val.(string); ok {
			buf = regex.JoinBytes(buf, '"', EscapeExsArgs([]byte(v), '"'), '"', ',')
		} else if v, ok := val.([]byte); ok {
			buf = regex.JoinBytes(buf, '"', EscapeExsArgs(v, '"'), '"', ',')
		} else if v, ok := val.(Map); ok {
			buf = regex.JoinBytes(buf, renderExsMap(v), ',')
		} else if v, ok := val.([]any); ok {
			buf = regex.JoinBytes(buf, renderExsArray(v), ',')
		} else {
			buf = regex.JoinBytes(buf, goutil.ToType[[]byte](val), ',')
		}
	}

	if buf[len(buf)-1] == ',' {
		buf = buf[:len(buf)-1]
	}

	buf = append(buf, '}')
	return buf
}

func renderExsArray(args []any) []byte {
	buf := []byte("[")

	for _, val := range args {
		if v, ok := val.(string); ok {
			buf = regex.JoinBytes(buf, '"', EscapeExsArgs([]byte(v), '"'), '"', ',')
		} else if v, ok := val.([]byte); ok {
			buf = regex.JoinBytes(buf, '"', EscapeExsArgs(v, '"'), '"', ',')
		} else if v, ok := val.(Map); ok {
			buf = regex.JoinBytes(buf, renderExsMap(v), ',')
		} else if v, ok := val.([]any); ok {
			buf = regex.JoinBytes(buf, renderExsArray(v), ',')
		} else {
			buf = regex.JoinBytes(buf, goutil.ToType[[]byte](v), ',')
		}
	}

	if buf[len(buf)-1] == ',' {
		buf = buf[:len(buf)-1]
	}

	buf = append(buf, ']')
	return buf
}
