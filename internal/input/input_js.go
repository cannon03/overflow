//go:build js

package input

import (
	"os"
	"time"

	"overflow/internal/entity"
)

// Input for WASM — reads from os.Stdin in a readLoop, same as the Unix version.
// No syscall, no raw mode — xterm.js handles terminal input via fs.readSync.
type Input struct {
	keyCh   chan entity.Key
	done    chan struct{}
	lastKey entity.Key
}

// New creates a new Input handler for WASM.
func New() *Input {
	return &Input{
		keyCh: make(chan entity.Key, 64),
		done:  make(chan struct{}),
	}
}

// EnableRawMode is a no-op in WASM (xterm.js handles raw mode).
func (in *Input) EnableRawMode() error { return nil }

// Restore is a no-op in WASM.
func (in *Input) Restore() {}

// Start begins reading input from os.Stdin in a background goroutine.
func (in *Input) Start() {
	go in.readLoop()
}

// Stop stops the input handler.
func (in *Input) Stop() {
	close(in.done)
}

// GetKey returns a key press (non-blocking). Returns KeyNone if no key.
func (in *Input) GetKey() entity.Key {
	select {
	case k := <-in.keyCh:
		in.lastKey = k
		return k
	default:
		return entity.KeyNone
	}
}

// LastKey returns the last key pressed without consuming it.
func (in *Input) LastKey() entity.Key {
	return in.lastKey
}

func (in *Input) readLoop() {
	buf := make([]byte, 128)
	for {
		select {
		case <-in.done:
			return
		default:
		}

		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			time.Sleep(1 * time.Millisecond)
			continue
		}

		key := parseKey(buf[:n])
		select {
		case in.keyCh <- key:
		default:
		}
	}
}

func parseKey(data []byte) entity.Key {
	if len(data) == 0 {
		return entity.KeyNone
	}

	ch := rune(data[0])
	switch {
	case ch == 'w' || ch == 'W':
		return entity.KeyW
	case ch == 'a' || ch == 'A':
		return entity.KeyA
	case ch == 's' || ch == 'S':
		return entity.KeyS
	case ch == 'd' || ch == 'D':
		return entity.KeyD
	case ch == ' ':
		return entity.KeySpace
	case ch == '\x1b' || ch == 27:
		return entity.KeyEsc
	case ch == 'p' || ch == 'P':
		return entity.KeyP
	case ch == 'r' || ch == 'R':
		return entity.KeyR
	case ch == 'q' || ch == 'Q':
		return entity.KeyQ
	}
	return entity.KeyNone
}
