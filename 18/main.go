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

func NewGrid(w, h int) *Grid {
	g := Grid{
		field: make([]byte, w*h),
		cols:  w,
	}
	for i := range g.field {
		g.field[i] = '.'
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

var turnDefs = []coord{
	{x: 0, y: -1},
	{x: 1, y: 0},
	{x: 0, y: 1},
	{x: -1, y: 0},
}

func (g *Grid) neighbors(c coord) []coord {
	next := []coord{}
	for _, t := range turnDefs {
		n := c.Add(t)
		if v, ok := g.Get(n); ok && v != '#' {
			debug("adding neighbors %s", n)
			next = append(next, n)
		}
	}
	return next
}

type Node struct {
	c      coord
	parent *Node
	score  int
}

func (n *Node) ID() string {
	return fmt.Sprintf("%s %d", n.c, n.score)
}

func (n *Node) String() string {
	return fmt.Sprintf("%s (%d)", n.c, n.score)
}

func (n *Node) GetPath() []*Node {
	path := []*Node{n}
	for n != nil {
		path = append(path, n)
		n = n.parent
	}
	return path
}

// use a* to find shortest path
// implementation from https://de.wikipedia.org/wiki/A*-Algorithmus
func (g *Grid) Solve(startC, endC coord) (*Node, bool) {
	openList := []*Node{{c: startC}}
	closedList := map[coord]bool{}

	for len(openList) > 0 {
		// find smallest f value in open list
		slices.SortFunc(openList, func(l, r *Node) int {
			return l.score - r.score
		})

		currentNode := openList[0]
		debug("checking %s", currentNode)
		openList = openList[1:]
		closedList[currentNode.c] = true

		// shortest path found
		if currentNode.c == endC {
			return currentNode, true
		}

		for _, t := range g.neighbors(currentNode.c) {
			// node is in closed list, skip
			neighborN := &Node{
				c:      t,
				parent: currentNode,
				score:  currentNode.score + 1,
			}

			if closedList[neighborN.c] {
				debug("skipping %s: in closed list", neighborN)
				continue
			}

			idx := slices.IndexFunc(openList, func(n *Node) bool {
				return n.c == t
			})
			if idx < 0 {
				openList = append(openList, neighborN)
			} else if neighborN.score < openList[idx].score {
				openList[idx] = neighborN
			}
		}
	}

	return nil, false
}

func ParseCoord(s string) (coord, error) {
	x, y, ok := strings.Cut(s, ",")
	if !ok {
		return coord{}, fmt.Errorf("could not parse %s as coordinate", s)
	}
	ix, err := strconv.Atoi(x)
	if err != nil {
		return coord{}, fmt.Errorf("could not parse %s as coordinate: %w", s, err)
	}
	iy, err := strconv.Atoi(y)
	if err != nil {
		return coord{}, fmt.Errorf("could not parse %s as coordinate: %w", s, err)
	}
	return coord{x: ix, y: iy}, nil
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

func runA(file string, w, h int, dropBytes int) (int, error) {
	in, err := os.ReadFile(file)
	if err != nil {
		return -1, err
	}

	g := NewGrid(w, h)

	for _, line := range strings.Split(string(in), "\n") {
		dropBytes--
		if dropBytes < 0 {
			break
		}
		c, err := ParseCoord(line)
		if err != nil {
			return -1, err
		}
		g.MustSet(c, '#')
	}

	// get start and end
	startI := 0              // top left
	endI := len(g.field) - 1 // bottom right

	startC := g.MustI2p(startI)
	endC := g.MustI2p(endI)

	g.MustSet(startC, 'S')
	g.MustSet(endC, 'E')

	debug("initial")
	debug(g.String())

	path, found := g.Solve(startC, endC)
	if !found {
		return -1, fmt.Errorf("could not find path")
	}

	for _, n := range path.GetPath() {
		fmt.Println(n.c)
		g.MustSet(n.c, 'O')
	}
	debug(g.String())

	pathLen := bytes.Count(g.field, []byte{'O'})
	fmt.Println("total score:", path.score)
	fmt.Println("path len:", pathLen)

	return path.score, nil
}

func runB(file string, w, h int) (int, error) {
	in, err := os.ReadFile(file)
	if err != nil {
		return -1, err
	}

	g := NewGrid(w, h)

	drops := strings.Split(string(in), "\n")

	// get start and end
	startI := 0              // top left
	endI := len(g.field) - 1 // bottom right

	startC := g.MustI2p(startI)
	endC := g.MustI2p(endI)

	g.MustSet(startC, 'S')
	g.MustSet(endC, 'E')

	debug("initial")
	debug(g.String())

	for i, drop := range drops {
		dropC, err := ParseCoord(drop)
		if err != nil {
			return -1, err
		}
		fmt.Println("dropping", i, dropC)

		g.MustSet(dropC, '#')
		_, found := g.Solve(startC, endC)
		if !found {
			fmt.Println("no path after dropping", dropC)
			break
		}
	}

	return 0, nil
}

const debugEnabled = false

func debug(msg string, args ...any) {
	if !debugEnabled {
		return
	}
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, msg)
		return
	}
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
}

func main() {
	// if _, err := run("input-test.txt", 7, 7, 12); err != nil {
	if _, err := runA("input.txt", 71, 71, 1024); err != nil {
		fmt.Fprintf(os.Stderr, "puzzle errored: %v\n", err)
		os.Exit(1)
	}
	if _, err := runB("input.txt", 71, 71); err != nil {
		fmt.Fprintf(os.Stderr, "puzzle errored: %v\n", err)
		os.Exit(1)
	}
}
