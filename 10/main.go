package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func (g *Grid) ResetVisited() {
	for i := range g.visited {
		g.visited[i] = 0
	}
}

type Grid struct {
	field   []byte
	cols    int
	visited []int
}

func NewGrid(in []byte) *Grid {
	cols := bytes.IndexByte(in, '\n')
	if cols < 0 {
		panic("no line breal in input")
	}
	f := bytes.ReplaceAll(in, []byte{'\n'}, []byte{})
	g := Grid{
		field:   f,
		cols:    cols,
		visited: make([]int, len(f)),
	}
	return &g
}

func (g *Grid) String() string {
	b := bytes.Buffer{}
	numLen := len(strconv.Itoa(max(g.cols, len(g.field)/g.cols)))
	for j := range numLen {
		b.WriteString(strings.Repeat(" ", numLen+2))
		for i := 0; i < g.cols; i++ {
			c := strconv.Itoa(i)
			c = strings.Repeat(" ", numLen-len(c)) + c
			b.WriteByte(c[j])
		}
		b.WriteString("\n")
	}
	b.WriteString(strings.Repeat(" ", numLen+2))
	b.WriteString(strings.Repeat("↓", g.cols))
	b.WriteByte('\n')
	b.WriteString(strings.Repeat(" ", numLen+1))
	b.WriteString("┌")
	b.WriteString(strings.Repeat("─", g.cols))
	b.WriteString("┐")
	for i := 0; i < len(g.field); i = i + g.cols {
		b.WriteByte('\n')
		r := strconv.Itoa(i / g.cols)
		r = strings.Repeat(" ", numLen-len(r)) + r + "→"
		b.WriteString(r)
		b.WriteString("│")
		//b.Write(bytes.ReplaceAll(g.field[i:i+g.cols], []byte{'.'}, []byte{' '}))
		b.Write(g.field[i : i+g.cols])
		b.WriteString("│")
	}
	b.WriteByte('\n')
	b.WriteString(strings.Repeat(" ", numLen+1))
	b.WriteString("└")
	b.WriteString(strings.Repeat("─", g.cols))
	b.WriteString("┘")
	return b.String()
}

func (g *Grid) p2i(c coord) (int, bool) {
	if c.x < 0 || c.y < 0 {
		return -1, false
	}
	if c.x >= g.cols {
		return -1, false
	}
	idx := c.y*g.cols + c.x
	if idx >= len(g.field) {
		return -1, false
	}
	return idx, true
}

func (g *Grid) MustP2i(c coord) int {
	if idx, ok := g.p2i(c); ok {
		return idx
	}
	panic("could not convert coordinate to index: " + c.String())
}

func (g *Grid) i2p(idx int) (coord, bool) {
	if idx < 0 || idx >= len(g.field) {
		return coord{}, false
	}
	return coord{x: idx % g.cols, y: idx / g.cols}, true
}

func (g *Grid) MustI2p(idx int) coord {
	if c, ok := g.i2p(idx); ok {
		return c
	}
	panic("could not convert index to coordinate: " + strconv.Itoa(idx))
}

func (g *Grid) Get(c coord) (byte, bool) {
	if idx, ok := g.p2i(c); ok {
		return g.field[idx], true
	}
	return 0, false
}

func (g *Grid) MustGet(c coord) byte {
	if v, ok := g.Get(c); ok {
		return v
	}
	panic("could not get " + c.String())
}

func (g *Grid) Set(c coord, v byte) (byte, bool) {
	if idx, ok := g.p2i(c); ok {
		o := g.field[idx]
		g.field[idx] = v
		return o, true
	}
	return 0, false
}

func (g *Grid) MustSet(c coord, v byte) byte {
	if p, ok := g.Set(c, v); ok {
		return p
	}
	panic("could not set " + c.String())
}

type coord struct {
	x, y int
}

func (c coord) String() string {
	return fmt.Sprintf("(%d,%d)", c.x, c.y)
}

func (c coord) Add(o coord) coord {
	return coord{x: c.x + o.x, y: c.y + o.y}
}

var offsets = []coord{
	{x: 0, y: -1},
	{x: 0, y: +1},
	{x: -1, y: 0},
	{x: +1, y: 0},
}

func (g *Grid) Climb(c coord, next byte) {
	v, ok := g.Get(c)
	if !ok {
		return
	}

	if v != next {
		return
	}

	// end of trail reached
	if v == '9' {
		pos := g.MustP2i(c)
		g.visited[pos]++
		fmt.Fprintln(debugW, "skipping ", c, pos, v-'0', g.visited[pos])
		return
	}

	// explore if there is a path to the next number
	for _, offset := range offsets {
		nc := c.Add(offset)
		g.Climb(nc, next+1)
	}
}

func run(file string) (int, int, error) {
	in, err := os.ReadFile(file)
	if err != nil {
		return -1, -1, err
	}

	g := NewGrid(in)

	scoreA, scoreB := 0, 0

	for pos, v := range g.field {
		if v != '0' {
			continue
		}
		g.ResetVisited()
		c := g.MustI2p(pos)
		fmt.Fprintln(debugW, "checking", c, pos)
		g.Climb(c, '0')
		for _, v := range g.visited {
			if v > 0 {
				scoreA++
			}
			scoreB += v
		}
	}

	fmt.Println("scoreA", scoreA)
	fmt.Println("scoreB", scoreB)

	return scoreA, scoreB, nil
}

var debugW = os.Stderr

func main() {
	if _, _, err := run("input.txt"); err != nil {
		fmt.Fprintf(os.Stderr, "puzzle errored with %v\n", err)
		os.Exit(1)
	}
}
