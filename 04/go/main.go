package main

import (
	"bytes"
	"fmt"
	"os"
)

const XMAS = "XMAS"
const XMAS_REVERSE = "SAMX"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "puzzle errored with %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		return err
	}

	// part 1
	xmas, err := countXMAS(bytes[:])
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "occurences part 1: %d\n", xmas)

	// part 2
	masX, err := countMASX(bytes[:])
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "occurences part 2: %d\n", masX)
	return nil
}

func countMASX(input []byte) (int, error) {
	patterns := []rect{
		toRect(mas1),
		toRect(mas2),
		toRect(mas3),
		toRect(mas4),
	}

	m := NewMatrix(input)
	occurences := 0

	for y := range m.NumRows() {
		for x := range m.NumCols() {
			r, ok := m.Rect(x, y, 3, 3)
			if !ok {
				continue
			}
			for _, p := range patterns {
				if p.Matches(r) {
					occurences++
					break
				}
			}
		}
	}

	return occurences, nil

}

func countXMAS(input []byte) (int, error) {

	m := NewMatrix(input)
	occurences := 0

	for y := range m.NumRows() {
		for x := range m.NumCols() {
			// horizontal
			if seq, ok := m.Seq(len(XMAS), x, y, 1, 0); ok && isXMAS(seq) {
				occurences++
			}

			// vertical
			if seq, ok := m.Seq(len(XMAS), x, y, 0, 1); ok && isXMAS(seq) {
				occurences++
			}

			// diagonal form top left to bottom right
			if seq, ok := m.Seq(len(XMAS), x, y, 1, 1); ok && isXMAS(seq) {
				occurences++
			}

			// diagonal form top right to bottom left
			if seq, ok := m.Seq(len(XMAS), x, y, -1, 1); ok && isXMAS(seq) {
				occurences++
			}
		}
	}

	return occurences, nil
}

func isXMAS(s string) bool {
	return s == XMAS || s == XMAS_REVERSE
}

func NewMatrix(s []byte) matrix {
	data := bytes.Split(s, []byte("\n"))
	return matrix(data)
}

type matrix [][]byte

func (m matrix) NumCols() int {
	return len(m[0])
}

func (m matrix) NumRows() int {
	return len(m)
}

func (m matrix) Get(x, y int) (byte, bool) {
	// bounds checks
	if x < 0 || y < 0 {
		return 0, false
	}
	if y >= len(m) {
		return 0, false
	}
	row := m[y]
	if x >= len(row) {
		return 0, false
	}
	return row[x], true
}

type rect [][]byte

func (r rect) Matches(other [][]byte) bool {
	if len(r) != len(other) {
		return false
	}
	for i := range other {
		or := other[i]
		rr := r[i]
		if len(rr) != len(or) {
			return false
		}
		for j := range rr {
			if rr[j] == '.' {
				continue
			}
			if or[j] != rr[j] {
				return false
			}
		}
	}
	return true
}

func (r rect) String() string {
	buf := bytes.Buffer{}
	for i, row := range r {
		if i > 0 {
			buf.WriteByte('\n')
		}
		for _, b := range row {
			buf.WriteByte(b)
		}
	}
	return buf.String()
}

func (m *matrix) Rect(x, y, w, h int) (rect, bool) {
	r := make([][]byte, 0, h)
	for ih := range h {
		row := []byte{}
		for iw := range w {
			v, ok := m.Get(x+iw, y+ih)
			if !ok {
				return nil, false
			}
			row = append(row, v)
		}
		r = append(r, row)
	}
	return rect(r), true
}

func (m *matrix) Seq(len, startX, startY, offsetX, offsetY int) (string, bool) {
	buf := bytes.Buffer{}
	x, y := startX, startY
	for range len {
		v, ok := m.Get(x, y)
		if !ok {
			return "", false
		}
		buf.WriteByte(v)
		x += offsetX
		y += offsetY

	}
	return buf.String(), true
}

const mas1 = `
M.S
.A.
M.S
`

const mas2 = `
S.M
.A.
S.M
`

const mas3 = `
M.M
.A.
S.S
`

const mas4 = `
S.S
.A.
M.M
`

func toRect(s string) rect {
	b := []byte(s)
	b = bytes.TrimSpace(b)
	bb := bytes.Split(b, []byte{'\n'})
	return rect(bb)
}
