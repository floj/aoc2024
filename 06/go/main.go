package main

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"slices"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	if err := runA("input.txt"); err != nil {
		fmt.Fprintf(os.Stderr, "puzzle errored: %v\n", err)
		os.Exit(1)
	}

	if err := runB("input.txt"); err != nil {
		fmt.Fprintf(os.Stderr, "puzzle errored: %v\n", err)
		os.Exit(1)
	}
}

const (
	NORTH byte = '^'
	EAST       = '>'
	SOUTH      = 'v'
	WEST       = '<'

	MOVED        state = 0
	LEFT_FIELD         = iota
	ENTERED_LOOP       = iota
)

type state int

var playerSym = bytes.Runes([]byte{NORTH, EAST, SOUTH, WEST})

func turn(d byte) byte {
	switch d {
	case NORTH:
		return EAST
	case EAST:
		return SOUTH
	case SOUTH:
		return WEST
	case WEST:
		return NORTH
	default:
		panic("unknown direction")
	}
}

type Area struct {
	fields []byte
	visits []byte
	cols   int
	rows   int
	moves  int
}

func (a *Area) Count(v byte) int {
	count := 0
	for _, b := range a.fields {
		if b == v {
			count++
		}
	}
	return count
}

func (a *Area) Get(x, y int) byte {
	idx := y*a.cols + x
	if idx >= len(a.fields) {
		return 0
	}
	return a.fields[idx]
}

func (a *Area) Player() (int, int, byte) {
	pos := bytes.IndexFunc(a.fields, func(r rune) bool {
		return slices.Contains(playerSym, r)
	})
	if pos < 0 {
		return -1, -1, 0
	}

	return pos % a.cols, pos / a.cols, a.fields[pos]
}

func (a *Area) Move() state {
	a.moves++
	oldX, oldY, direction := a.Player()
	// fmt.Fprintf(os.Stderr, "%d (%d,%d)\n", a.moves, oldX, oldY)
	newX, newY := oldX, oldY
	a.Set(oldX, oldY, 'X')
	switch direction {
	case NORTH:
		newY--
	case EAST:
		newX++
	case SOUTH:
		newY++
	case WEST:
		newX--
	default:
		panic("unknown direction")
	}
	if newX < 0 || newX >= a.cols || newY < 0 || newY >= a.rows {
		return LEFT_FIELD
	}

	// obstacle on new field, rotate right
	if a.Get(newX, newY) == '#' {
		direction = turn(direction)
		a.Set(oldX, oldY, direction)
		return MOVED
	}

	a.Set(newX, newY, direction)
	if a.Visited(newX, newY, direction) {
		return ENTERED_LOOP
	}

	return MOVED
}

func (a *Area) Set(x, y int, v byte) {
	idx := y*a.cols + x
	if idx >= len(a.fields) {
		return
	}
	a.fields[idx] = v
}

func (a *Area) Visited(x, y int, direction byte) bool {
	idx := y*a.cols + x
	if idx >= len(a.fields) {
		return true
	}

	flags := a.visits[idx]
	switch direction {
	case NORTH:
		flags = flags | 0b0001
	case EAST:
		flags = flags | 0b0010
	case SOUTH:
		flags = flags | 0b0100
	case WEST:
		flags = flags | 0b1000
	}
	if a.visits[idx] == flags {
		// same field visited in the same direction -> loop
		return true
	}
	a.visits[idx] = flags
	return false
}

func (a *Area) String() string {
	buf := bytes.Buffer{}
	x, y, _ := a.Player()
	fmt.Fprintf(&buf, "(%d,%d)\n", x, y)
	for i := 0; i < len(a.fields); i = i + a.cols {
		if i > 0 {
			buf.WriteByte('\n')
		}
		buf.Write(a.fields[i : i+a.cols])
	}
	return buf.String()
}

func NewArea(in []byte) Area {
	field := slices.Clone(in)
	cols := bytes.Index(field, []byte{'\n'})
	f := bytes.ReplaceAll(field, []byte{'\n'}, []byte{})
	rows := len(f) / cols

	return Area{
		fields: f,
		cols:   cols,
		rows:   rows,
		visits: make([]byte, len(f)),
	}
}

func runA(inputFile string) error {
	in, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	a := NewArea(in)
	for a.Move() == MOVED {
	}

	fmt.Printf("done, visited %d\n", a.Count('X'))
	return nil
}

func runB(inputFile string) error {
	in, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	leftField, enteredLoop := &atomic.Int32{}, &atomic.Int32{}
	field := NewArea(in).fields
	remaining := &atomic.Int32{}
	wg := &sync.WaitGroup{}

	start := time.Now()
	for range runtime.NumCPU() * 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				i := remaining.Add(1)
				if int(i) >= len(field) {
					break
				}
				// if field is already taken (existing obstacle or player), skip
				if field[i] != '.' {
					continue
				}
				// add additional obstacle
				a := NewArea(in)
				a.fields[i] = '#'

				s := MOVED
				for s == MOVED {
					s = a.Move()
					switch s {
					case LEFT_FIELD:
						fmt.Fprintf(os.Stderr, "obstacle at %d of %d: OK\n", i, len(field))
						leftField.Add(1)
					case ENTERED_LOOP:
						fmt.Fprintf(os.Stderr, "obstacle at %d of %d: LOOP\n", i, len(field))
						enteredLoop.Add(1)
					}
				}
			}
		}()
	}

	wg.Wait()
	fmt.Printf("done, looped %d, %v\n", enteredLoop.Load(), time.Since(start))
	return nil
}
