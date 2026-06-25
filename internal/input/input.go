package input

import (
	"os"
	"overflow/internal/entity"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

// Input handles raw terminal input.
type Input struct {
	keyCh    chan entity.Key
	done     chan struct{}
	mu       sync.Mutex
	oldState *syscall.Termios
	lastKey  entity.Key
}

// New creates a new Input handler.
func New() *Input {
	return &Input{
		keyCh: make(chan entity.Key, 64),
		done:  make(chan struct{}),
	}
}

// EnableRawMode puts the terminal into raw mode (no echo, no line buffering).
func (in *Input) EnableRawMode() error {
	var oldState syscall.Termios
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, os.Stdin.Fd(), syscall.TCGETS, uintptr(unsafe.Pointer(&oldState)), 0, 0, 0); err != 0 {
		return err
	}
	in.oldState = &oldState

	newState := oldState
	newState.Iflag &^= syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK | syscall.ISTRIP | syscall.INLCR | syscall.IGNCR | syscall.ICRNL | syscall.IXON
	newState.Oflag &^= syscall.OPOST
	newState.Lflag &^= syscall.ECHO | syscall.ECHONL | syscall.ICANON | syscall.ISIG | syscall.IEXTEN
	newState.Cflag &^= syscall.CSIZE | syscall.PARENB
	newState.Cflag |= syscall.CS8
	newState.Cc[syscall.VMIN] = 0
	newState.Cc[syscall.VTIME] = 1

	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, os.Stdin.Fd(), syscall.TCSETS, uintptr(unsafe.Pointer(&newState)), 0, 0, 0); err != 0 {
		return err
	}
	return nil
}

// Restore restores the terminal to its original mode.
func (in *Input) Restore() {
	if in.oldState != nil {
		syscall.Syscall6(syscall.SYS_IOCTL, os.Stdin.Fd(), syscall.TCSETS, uintptr(unsafe.Pointer(in.oldState)), 0, 0, 0)
	}
}

// Start begins reading input in a background goroutine.
func (in *Input) Start() {
	go in.readLoop()
}

// Stop stops the input reading goroutine.
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

	// Escape sequences for arrow keys: \033[A, \033[B, \033[C, \033[D
	if data[0] == '\033' && len(data) >= 3 && data[1] == '[' {
		switch data[2] {
		case 'A':
			return entity.KeyUp
		case 'B':
			return entity.KeyDown
		case 'C':
			return entity.KeyRight
		case 'D':
			return entity.KeyLeft
		}
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
	case ch == '\033':
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
