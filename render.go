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
	filename   string
	cmd        *exec.Cmd
	stdin      io.WriteCloser
	stdout     io.ReadCloser
	out        <-chan []byte
	running    *bool
	restarting *bool
	mu         sync.Mutex

	// will be used when recompiling
	compMU sync.RWMutex
}

type Map map[string]any

var ErrorStopped error = errors.New("elixir cmd stopped running")
var ErrorTimeout error = errors.New("request timed out")

// Engine starts a new htmlc engine for an exs template
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
	restarting := false
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
		if !restarting {
			fmt.Println("\033[31mhtmlc engine stopped!\033[0m")
		}
	}()

	err = cmd.Start()
	if err != nil {
		running = false
		stdin.Close()
		stdout.Close()
		return nil, err
	}

	time.Sleep(1 * time.Second)
	ready = true

	return &ExsEngine{
		filename:   file,
		cmd:        cmd,
		stdin:      stdin,
		stdout:     stdout,
		out:        out,
		running:    &running,
		restarting: &restarting,
	}, nil
}

// Restart will restart (and optionally recompile) the htmlc Engine
//
// @compSrc: directory to compile (leave blank for no recompile)
func (exs *ExsEngine) Restart(compSrc ...string) error {
	exs.compMU.Lock()
	defer exs.compMU.Unlock()

	fmt.Print("\033[33m htmlc engine restarting...\033[0m   \r")

	time.Sleep(250 * time.Millisecond)

	*exs.restarting = true
	*exs.running = false
	exs.stdin.Write([]byte("stop\n"))

	time.Sleep(250 * time.Millisecond)

	if len(compSrc) != 0 {
		if err := Compile(compSrc[0], exs.filename); err != nil {
			fmt.Println("\033[31mhtmlc engine stopped!\033[0m         ")
			return err
		}
	}

	engine, err := Engine(exs.filename)
	if err != nil {
		fmt.Println("\033[31mhtmlc engine stopped!\033[0m         ")
		return err
	}

	exs.cmd = engine.cmd
	exs.stdin = engine.stdin
	exs.stdout = engine.stdout
	exs.out = engine.out
	exs.running = engine.running
	exs.restarting = engine.restarting

	time.Sleep(250 * time.Millisecond)

	fmt.Print("                              \r")

	return nil
}

// Render renders an html page
func (exs *ExsEngine) Render(name string, args Map, layout ...string) ([]byte, error) {
	exs.compMU.RLock()
	defer exs.compMU.RUnlock()

	if !*exs.running {
		return []byte{}, ErrorStopped
	}

	lay := []byte("layout")
	if len(layout) != 0 {
		lay = regex.CompRE2(`["'\']`).RepStrLit([]byte(layout[0]), []byte{})
	}

	//todo: fix issue with args not rendering properly

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

// Layout renders an html layout
func (exs *ExsEngine) Layout(name string, args Map, cont ...map[string]string) ([]byte, error) {
	exs.compMU.RLock()
	defer exs.compMU.RUnlock()

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

// Widget renders an html widget
func (exs *ExsEngine) Widget(name string, args Map) ([]byte, error) {
	exs.compMU.RLock()
	defer exs.compMU.RUnlock()

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
