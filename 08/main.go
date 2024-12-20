package main

import (
	"bytes"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Grid struct {
	field []byte
	cols  int
}

func NewGrid(in []byte) *Grid {
	cols := bytes.IndexByte(in, '\n')
	if cols < 0 {
		panic("no line breal in input")
	}
	g := Grid{
		field: bytes.ReplaceAll(in, []byte{'\n'}, []byte{}),
		cols:  cols,
	}
	return &g
}

func (g *Grid) Clone() *Grid {
	return &Grid{
		field: slices.Clone(g.field),
		cols:  g.cols,
	}
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

func (c coord) Invert() coord {
	return coord{x: -c.x, y: -c.y}
}

type coordPair struct {
	c1, c2 coord
}

func (c coordPair) Distance() coord {
	return coord{
		x: c.c1.x - c.c2.x,
		y: c.c1.y - c.c2.y,
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func permute(v []coord) []coordPair {
	pairs := []coordPair{}
	for i := range v {
		for j := i + 1; j < len(v); j++ {
			pairs = append(pairs, coordPair{c1: v[i], c2: v[j]})
		}
	}
	return pairs
}

func runA(file string) error {
	in, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	g := NewGrid(in)
	// find antennas
	antennas := map[string][]coord{}
	for i, v := range g.field {
		if v == '.' {
			continue
		}
		antennas[string(v)] = append(antennas[string(v)], g.MustI2p(i))
	}

	fmt.Fprintf(debugWriter, "antennas: %+v\n", antennas)

	for k, vv := range antennas {
		if len(vv) < 2 {
			continue
		}
		pairs := permute(vv)
		fmt.Fprintf(debugWriter, "checking %s: %v\n", string(k), pairs)
		for _, p := range pairs {
			dist := p.Distance()
			g.Set(p.c1.Add(dist), '#')
			g.Set(p.c2.Add(dist.Invert()), '#')
		}
		fmt.Fprintf(debugWriter, "%s\n", g)
	}

	anti := bytes.Count(g.field, []byte{'#'})
	fmt.Println("antinodes", anti)
	return nil
}

var debugWriter = os.Stderr

func main() {
	if err := runA("input.txt"); err != nil {
		fmt.Fprintf(os.Stderr, "puzzle errored: %v\n", err)
		os.Exit(1)
	}

}
