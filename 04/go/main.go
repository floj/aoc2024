package main

import (
	"bytes"
	"fmt"
	"os"
)

const XMAS = "XMAS"
const XMAS_REVERSE = "SAMX"

func main() {
	cnt, err := countXMAS("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "puzzle errored with %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "occurences: %d\n", cnt)
}

func log(a ...any) {
	if true {
		fmt.Fprintln(os.Stderr, a...)
	}

}

func countXMAS(inputFile string) (int, error) {
	bytes, err := os.ReadFile(inputFile)
	if err != nil {
		return -1, err
	}
	occurences := 0

	m := NewMatrix(bytes)

	for y := range m.NumRows() {
		for x := range m.NumCols() {
			// horizontal
			seq, ok := m.Seq(len(XMAS), x, y, 1, 0)
			log("h", x, y, ok, seq)
			if ok && isXMAS(seq) {
				occurences++
			}

			// vertical
			seq, ok = m.Seq(len(XMAS), x, y, 0, 1)
			log("v", x, y, ok, seq)
			if ok && isXMAS(seq) {
				occurences++
			}

			// diagonal form top left to bottom right
			seq, ok = m.Seq(len(XMAS), x, y, 1, 1)
			log("d1", x, y, ok, seq)
			if ok && isXMAS(seq) {
				occurences++
			}

			// diagonal form top right to bottom left
			log("d2", x, y, ok, seq)
			seq, ok = m.Seq(len(XMAS), x, y, -1, 1)
			if ok && isXMAS(seq) {
				occurences++
			}
			log("----")
		}
		log("####")
	}

	return occurences, nil
}

func isXMAS(s string) bool {
	return s == XMAS || s == XMAS_REVERSE
}

func NewMatrix(s []byte) matrix {
	data := bytes.Split(s, []byte{'\n'})
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
