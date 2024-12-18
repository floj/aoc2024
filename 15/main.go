package main

import (
	"bytes"
	"fmt"
	"os"
	"slices"
	"strconv"
)

type Grid struct {
	field []byte
	cols  int
}

func NewGrid(in []byte) *Grid {
	cols := bytes.IndexByte(in, '\n')
	if cols < 0 {
		panic("no linebreak found")
	}
	field := bytes.ReplaceAll(in, []byte{'\n'}, []byte{})

	return &Grid{
		field: field,
		cols:  cols,
	}
}

func (g *Grid) Clone() *Grid {
	return &Grid{
		field: slices.Clone(g.field),
		cols:  g.cols,
	}
}

func (g *Grid) String() string {
	b := bytes.Buffer{}
	b.WriteString("  ")
	for i := 0; i < g.cols; i++ {
		c := strconv.Itoa(i)
		if len(c) == 1 {
			c = " " + c
		}
		b.WriteByte(c[0])
	}
	b.WriteString("\n  ")
	for i := 0; i < g.cols; i++ {
		c := strconv.Itoa(i)
		if len(c) == 1 {
			c = " " + c
		}
		b.WriteByte(c[1])
	}

	for i := 0; i < len(g.field); i = i + g.cols {
		b.WriteByte('\n')
		r := strconv.Itoa(i / g.cols)
		if len(r) == 1 {
			r = " " + r
		}
		b.WriteString(r)
		b.Write(g.field[i : i+g.cols])
	}
	return b.String()
}

func (g *Grid) p2i(c coord) (int, bool) {
	if c.x < 0 || c.y < 0 {
		return -1, false
	}
	idx := c.y*g.cols + c.x
	if idx >= len(g.field) {
		return -1, false
	}
	return idx, true
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
		fmt.Printf("drawing %s with %s\n", c, string(v))
		o := g.field[idx]
		g.field[idx] = v
		return o, true
	}
	return 0, false
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

func (c coord) AddX(offset int) coord {
	return coord{x: c.x + offset, y: c.y}
}

func (c coord) AddY(offset int) coord {
	return coord{x: c.x, y: c.y + offset}
}

func (c coord) AddXY(offsetX, offsetY int) coord {
	return coord{x: c.x + offsetX, y: c.y + offsetY}
}

func (c coord) GPS() int {
	return c.y*100 + c.x
}

var offsets = map[byte]coord{
	'^': {x: 0, y: -1},
	'>': {x: 1, y: 0},
	'v': {x: 0, y: 1},
	'<': {x: -1, y: 0},
}

func (g *Grid) Robot() (coord, bool) {
	rPos := bytes.IndexByte(g.field, '@')
	if rPos < 0 {
		return coord{}, false
	}
	rC, ok := g.i2p(rPos)
	return rC, ok
}

func (g *Grid) MoveWarehouse1(srcC coord, direction byte) bool {
	off, found := offsets[direction]
	if !found {
		panic("invalid direction")
	}

	srcV, ok := g.Get(srcC)
	if !ok {
		panic("invalid srcC  " + srcC.String())
	}
	destC := srcC.Add(off)
	destV, ok := g.Get(destC)
	if !ok {
		panic("invalid destC  " + srcC.String())
	}
	fmt.Printf("checking src=%s (%s) dest=%s (%s)\n", srcC, string(srcV), destC, string(destV))
	switch destV {
	case '#':
		// walls can't be pushed
		return false
	case '.':
		// empty field
		g.Set(destC, srcV)
		g.Set(srcC, '.')
		return true
	case 'O':
		// box
		if g.MoveWarehouse1(destC, direction) {
			g.Set(destC, srcV)
			g.Set(srcC, '.')
			return true
		}
		return false
	default:
		panic("invalid dest field value " + string(destV))
	}
}

func (g *Grid) MoveWarehouse2(indent string, srcC coord, direction byte) bool {
	off, found := offsets[direction]
	if !found {
		panic("invalid direction")
	}

	srcV := g.MustGet(srcC)

	fmt.Printf("%schecking src=%s (%s)\n", indent, srcC, string(srcV))
	// if src is a box, apply special handling

	switch srcV {
	case '#':
		return false
	case '.':
		return true
	case '[':
		fallthrough
	case ']':
		return g.MoveBox(indent+"  ", srcC, off, direction)
	}

	destC := srcC.Add(off)
	destV := g.MustGet(destC)

	switch destV {
	case '.':
		fmt.Printf("%s-> empty space, moving %s (%s) to %s\n", indent, srcC, string(srcV), destC)
		g.Set(destC, srcV)
		g.Set(srcC, '.')
		return true
	case '[':
		fallthrough
	case ']':
		if !g.MoveBox(indent+"  ", destC, off, direction) {
			return false
		}
		g.Set(destC, srcV)
		g.Set(srcC, '.')
		return true
	case '#':
		fmt.Printf("%s-> wall\n", indent)
		return false
	default:
		panic("unknown symbol: " + string(destV))
	}
}

func (g *Grid) MoveBox(indent string, srcC, off coord, direction byte) bool {
	// get top left corner of box
	srcV, ok := g.Get(srcC)
	if !ok {
		panic("invalid srcC  " + srcC.String())
	}
	if srcV == ']' {
		return g.MoveBox(indent, srcC.AddX(-1), off, direction)
	}
	if srcV != '[' {
		panic("not a box " + srcC.String())
	}

	destC := srcC.Add(off)
	fmt.Printf("%smoving box src=%s dest=%s srcV=%s direction=%s\n", indent, srcC, destC, string(srcV), string(direction))

	// when moving a box left or right, just check the one new tile that will be taken
	switch direction {
	case '^':
		fallthrough
	case 'v':
		if !g.MoveWarehouse2(indent+"  ", destC, direction) {
			fmt.Printf(" -> can't move boxL %s\n", destC)
			return false
		}
		if !g.MoveWarehouse2(indent+"  ", destC.AddX(1), direction) {
			fmt.Printf(" -> can't move boxR %s\n", destC.AddX(1))
			return false
		}
		g.Set(destC, '[')
		g.Set(destC.AddX(1), ']')
		g.Set(srcC, '.')
		g.Set(srcC.AddX(1), '.')
		return true
	case '<':
		if !g.MoveWarehouse2(indent+"  ", destC, direction) {
			fmt.Printf("%s-> can't move box %s\n", indent, destC)
			return false
		}
		g.Set(destC, '[')
		g.Set(destC.AddX(1), ']')
		g.Set(srcC.AddX(1), '.')
		return true
	case '>':
		if !g.MoveWarehouse2(indent+"  ", destC.AddX(1), direction) {
			fmt.Printf("%s-> can't move %s\n", indent, destC.AddX(1))
			return false
		}
		g.Set(destC, '[')
		g.Set(destC.AddX(1), ']')
		g.Set(srcC, '.')
		return true
	default:
		panic("invalid direction: " + string(direction))
	}
}

type MoveFn func(g *Grid)

func runA(file string) (int, error) {
	in, err := os.ReadFile(file)
	if err != nil {
		return -1, err
	}

	warehouse, momements, found := bytes.Cut(in, []byte{'\n', '\n'})
	if !found {
		return -1, fmt.Errorf("can't split input")
	}

	g := NewGrid(warehouse)
	momements = bytes.ReplaceAll(momements, []byte{'\n'}, []byte{})

	fmt.Println("initial")
	fmt.Println(g)
	for i, m := range momements {
		rC, ok := g.Robot()
		if !ok {
			return -1, fmt.Errorf("invalid robot position")
		}
		g.MoveWarehouse1(rC, m)
		fmt.Printf("move %d %s %s\n", i+1, rC, string(m))
		// fmt.Println(g)
	}

	sumA := 0

	for i, v := range g.field {
		if v != 'O' {
			continue
		}
		if c, ok := g.i2p(i); ok {
			sumA += c.GPS()
		}

	}

	fmt.Println("sum A", sumA)
	return sumA, nil
}

func ResizeForB(in []byte) []byte {
	resized := make([]byte, 0, len(in)*2)
	for i, b := range in {
		switch b {
		case '\n':
			resized = append(resized, '\n')
		case '#':
			resized = append(resized, '#', '#')
		case 'O':
			resized = append(resized, '[', ']')
		case '.':
			resized = append(resized, '.', '.')
		case '@':
			resized = append(resized, '@', '.')
		default:
			panic(fmt.Sprintf("invalid symbol on map at index %d: '%s'", i, string(b)))
		}
	}
	return resized
}

func runB(file string) (int, error) {
	in, err := os.ReadFile(file)
	if err != nil {
		return -1, err
	}

	warehouse, momements, found := bytes.Cut(in, []byte{'\n', '\n'})
	if !found {
		return -1, fmt.Errorf("can't split input")
	}

	g := NewGrid(ResizeForB(warehouse))
	momements = bytes.ReplaceAll(momements, []byte{'\n'}, []byte{})

	fmt.Println("initial")
	fmt.Println(g)
	for i, m := range momements {
		rC, ok := g.Robot()
		if !ok {
			return -1, fmt.Errorf("invalid robot position")
		}
		fmt.Printf("move %d %s %s\n", i+1, rC, string(m))

		newG := g.Clone()
		ok = newG.MoveWarehouse2("", rC, m)
		if !ok {
			// do nothing
			continue
		}
		g = newG
		fmt.Println(g)
		// check if field is broken
		if idx := bytes.Index(g.field, []byte(".]")); idx >= 0 {
			panic("split box " + g.MustI2p(idx).String())
		}
		if idx := bytes.Index(g.field, []byte("[.")); idx >= 0 {
			panic("split box " + g.MustI2p(idx).String())
		}

	}

	sumA := 0

	for i, v := range g.field {
		if v == 'O' || v == '[' {
			if c, ok := g.i2p(i); ok {
				sumA += c.GPS()
			}
		}
	}

	fmt.Println("sum A", sumA)
	return sumA, nil
}

func main() {
	if _, err := runA("input.txt"); err != nil {
		fmt.Fprintf(os.Stderr, "puzzle errored: %v\n", err)
		os.Exit(1)
	}
	if _, err := runB("input.txt"); err != nil {
		fmt.Fprintf(os.Stderr, "puzzle errored: %v\n", err)
		os.Exit(1)
	}
}
