package htmlc

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
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
	cmd := exec.Command(`elixir`, file)

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
				fmt.Println(err)
				continue
			}

			if !ready {
				continue
			}
			buf = buf[:n]

			if IexMode {
				if regex.CompRE2(`(?m)^\s*iex\s*\([0-9]+\)\s*>`).Match(buf) {
					continue
				}

				hasData := false
				regex.CompRE2(`(?s)(?:~c|)"(.*)"`).RepFunc(buf, func(data func(int) []byte) []byte {
					buf = bytes.ReplaceAll(data(1), []byte("\\n"), []byte("\n"))
					hasData = true
					return nil
				})

				if !hasData {
					buf = []byte("<h1>Error 500</h1><h2>Internal Server Error!</h2>")
				}
			} else {
				buf = regex.CompRE2(`(?s)<<(.*?)>>`).RepFunc(buf, func(data func(int) []byte) []byte {
					bit := make([]byte, len(data(1)))
					if _, err := base64.StdEncoding.Decode(bit, data(1)); err != nil {
						bit, _ = base64.StdEncoding.DecodeString(string(data(1)))
					}
					bytes.ReplaceAll(bit, []byte{0}, []byte{})
					return bit
				})
			}

			out <- buf
		}

		running = false
		stdin.Close()
		stdout.Close()
		fmt.Println("htmlc engine stopped")
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
		lay = regex.CompRE2(`["'\']`).RepStrLit([]byte(layout[0]), []byte{})
	}

	json, err := goutil.JSON.Stringify(args)
	if err != nil {
		return []byte{}, err
	}

	bit := make([]byte, len(json)*2)
	base64.StdEncoding.Encode(bit, json)
	regex.CompRE2(`["'\']`).RepStrLit([]byte(name), []byte{})

	in := regex.JoinBytes(
		`:render,`,
		regex.CompRE2(`["'\']`).RepStrLit([]byte(name), []byte{}), ',',
		lay, ',', bit,
		'\n',
	)

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

func (exs *ExsEngine) Layout(name string, args Map, cont ...map[string]string) ([]byte, error) {
	if !*exs.running {
		return []byte{}, ErrorStopped
	}

	json, err := goutil.JSON.Stringify(args)
	if err != nil {
		return []byte{}, err
	}

	var contStr []byte
	if len(cont) != 0 {
		contStr, err = goutil.JSON.Stringify(cont[0])
		if err != nil {
			return []byte{}, err
		}
	}

	bit := make([]byte, len(json)*2)
	base64.StdEncoding.Encode(bit, json)
	regex.CompRE2(`["'\']`).RepStrLit([]byte(name), []byte{})

	var in []byte
	if contStr != nil {
		contBit := make([]byte, len(contStr)*2)
		base64.StdEncoding.Encode(contBit, contStr)
		regex.CompRE2(`["'\']`).RepStrLit([]byte(name), []byte{})

		in = regex.JoinBytes(
			`:layout,`,
			regex.CompRE2(`["'\']`).RepStrLit([]byte(name), []byte{}), ',',
			bit, ',', contBit,
			'\n',
		)
	} else {
		in = regex.JoinBytes(
			`:layout,`,
			regex.CompRE2(`["'\']`).RepStrLit([]byte(name), []byte{}), ',',
			bit,
			'\n',
		)
	}

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

func (exs *ExsEngine) Widget(name string, args Map) ([]byte, error) {
	if !*exs.running {
		return []byte{}, ErrorStopped
	}

	json, err := goutil.JSON.Stringify(args)
	if err != nil {
		return []byte{}, err
	}

	bit := make([]byte, len(json)*2)
	base64.StdEncoding.Encode(bit, json)
	regex.CompRE2(`["'\']`).RepStrLit([]byte(name), []byte{})

	in := regex.JoinBytes(
		`:render,`,
		regex.CompRE2(`["'\']`).RepStrLit([]byte(name), []byte{}), ',',
		bit,
		'\n',
	)

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
