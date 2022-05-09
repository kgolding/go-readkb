package readkb

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"sync"
)

/*	struct input_event {
		struct timeval time;
		unsigned short type;
		unsigned short code;
		unsigned int value;
}; */

type InputEvent struct {
	Timeval Timeval
	Type    uint16
	Code    uint16
	V1      uint16
	V2      uint16
}

type Event struct {
	Char     rune
	Scancode uint16
}

type Keyboard struct {
	r   io.Reader
	C   chan *Event
	Err error
	sync.Once
}

func NewFromPath(dev string) (*Keyboard, error) {
	r, err := os.OpenFile(dev, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	return New(r), nil
}

func New(r io.Reader) *Keyboard {
	k := &Keyboard{
		r: r,
		C: make(chan *Event, 20),
	}
	go k.run()
	return k
}

func (k *Keyboard) Close() {
	if k == nil {
		return
	}
	k.Do(func() {
		close(k.C)
		k = nil
	})
}

type Scancode struct {
	Key        rune
	ShiftedKey rune
}

// http://www.quadibloc.com/comp/scan.htm
var CodeMap = map[uint16]Scancode{
	// Main keyboard
	0x1e: {'1', '!'},
	0x1f: {'2', '"'}, // US: @
	0x20: {'3', '£'}, // US: #
	0x21: {'4', '$'},
	0x22: {'5', '%'},
	0x23: {'6', '^'},
	0x24: {'7', '&'},
	0x25: {'8', '*'},
	0x26: {'9', '('},
	0x27: {'0', ')'},
	0x2d: {'-', '_'},
	0x2e: {'=', '+'},
	0x2a: {'\b', '\b'}, // Backspace
	0x2b: {'\t', '\t'}, // Tab
	0x14: {'q', 'Q'},
	0x1a: {'w', 'W'},
	0x08: {'e', 'E'},
	0x15: {'r', 'R'},
	0x17: {'t', 'T'},
	0x1c: {'y', 'Y'},
	0x18: {'u', 'U'},
	0x0c: {'i', 'I'},
	0x12: {'o', 'O'},
	0x13: {'p', 'P'},
	0x2f: {'[', '{'},
	0x30: {']', '}'},
	0x28: {'\n', '\n'},
	0x04: {'a', 'A'},
	0x16: {'s', 'S'},
	0x07: {'d', 'D'},
	0x09: {'f', 'F'},
	0x0a: {'g', 'G'},
	0x0b: {'h', 'H'},
	0x0d: {'j', 'J'},
	0x0e: {'k', 'K'},
	0x0f: {'l', 'L'},
	0x33: {';', ':'},
	0x34: {'\'', '@'},
	0x35: {'#', '~'},
	0x31: {'\\', '|'},
	0x1d: {'z', 'Z'},
	0x1b: {'x', 'X'},
	0x06: {'c', 'C'},
	0x19: {'v', 'V'},
	0x05: {'b', 'B'},
	0x11: {'n', 'N'},
	0x10: {'m', 'M'},
	0x36: {',', '<'},
	0x37: {'.', '>'},
	0x38: {'/', '?'},
	0x39: {' ', ' '},
	// Numpad
	0x5f: {'7', '\u21F1'}, // Home
	0x60: {'8', '\u2191'}, // Up
	0x61: {'9', '\u21de'}, // Page up
	0x5c: {'4', '\u2190'}, // Left
	0x97: {'5', '5'},      // @TODO Numpad 5 isn't working?
	0x5e: {'6', '\u2192'}, // Right
	0x57: {'+', '+'},
	0x59: {'1', '\u21f2'}, // End
	0x5a: {'2', '\u2193'}, // Down
	0x5b: {'3', '\u21df'}, // Page down
	0x62: {'0', '0'},      // Insert
	0x63: {'.', '\u2326'}, // Delete
	0x54: {'/', '/'},
	0x55: {'*', '*'},
	0x56: {'-', '-'},
	0x58: {'\n', '\n'},
}

func (k *Keyboard) run() {
	// InputEvent size varies depending on the CPU due to the Linux TimeVal definition
	inputEventSize := binary.Size(InputEvent{})

	// Custom scanner that splits on inputEventSize byte chunks
	scanner := bufio.NewScanner(k.r)
	scanner.Split(func(data []byte, atEOF bool) (int, []byte, error) {
		if len(data) < inputEventSize {
			return 0, nil, nil
		}
		return inputEventSize, data[:inputEventSize], nil
	})

	// It takes more than one InputEvent to detect a key being pressed, so
	// we keep track of the last InputEvent for reference
	var lastTimeval Timeval
	lastScancode := uint16(0)

	for scanner.Scan() {
		b := scanner.Bytes()
		var ie InputEvent
		err := binary.Read(bytes.NewBuffer(b), binary.LittleEndian, &ie)
		if err != nil {
			continue
		}
		if ie.Type == 4 && ie.Code == 4 { // Key down/up
			lastTimeval = ie.Timeval
			lastScancode = ie.V1
		} else if ie.Type == 1 && ie.V1 == 1 && lastTimeval.Equals(ie.Timeval) {
			sc, ok := CodeMap[lastScancode]
			if ok {
				e := &Event{
					Char:     sc.Key,
					Scancode: lastScancode,
				}
				k.C <- e
				// } else {
				// 	println("Unknown scancode", lastScancode)
			}
			// } else {
			// 	fmt.Printf("Read: %#v [% X]\n", ie, b)
		}
	}
	k.Close()
}
