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
	numLen := 3
	for j := range numLen {
		b.WriteString(strings.Repeat(" ", numLen+1))
		for i := 0; i < g.cols; i++ {
			c := strconv.Itoa(i)
			c = strings.Repeat(" ", numLen-len(c)) + c
			b.WriteByte(c[j])
		}
		b.WriteString("\n")
	}
	b.WriteString(strings.Repeat(" ", numLen+1))
	b.WriteString(strings.Repeat("|", g.cols))
	for i := 0; i < len(g.field); i = i + g.cols {
		b.WriteByte('\n')
		r := strconv.Itoa(i / g.cols)
		r = strings.Repeat(" ", numLen-len(r)) + r + "-"
		b.WriteString(r)

		b.Write(bytes.ReplaceAll(g.field[i:i+g.cols], []byte{'.'}, []byte{' '}))
	}
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

type turnDef struct {
	direction byte
	c         coord
}

var turnDefs = []turnDef{
	{c: coord{x: 0, y: -1}, direction: '^'},
	{c: coord{x: 1, y: 0}, direction: '>'},
	{c: coord{x: 0, y: 1}, direction: 'v'},
	{c: coord{x: -1, y: 0}, direction: '<'},
}

func turns(direction byte) []turnDef {
	cur := slices.IndexFunc(turnDefs, func(e turnDef) bool {
		return e.direction == direction
	})
	if cur < 0 {
		panic("invalid direction " + string(direction))
	}
	l := len(turnDefs)
	return []turnDef{
		turnDefs[cur],
		turnDefs[(l+cur-1)%l],
		turnDefs[(l+cur+1)%l],
	}
}

type Node struct {
	c         coord
	parent    *Node
	direction byte
	score     int
	cost      int
}

func (n *Node) ID() string {
	return n.c.String() + string(n.direction)
}

func (n *Node) String() string {
	return fmt.Sprintf("%s %s (%d/%d)", n.c, string(n.direction), n.cost, n.score)
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
func (g *Grid) Solve(startC, endC coord) ([]*Node, bool) {
	openList := []*Node{{c: startC, direction: '>'}}
	closedList := map[string]bool{}

	// add all walls to the closed list (no need to investigate)
	for idx, v := range g.field {
		if v == '#' {
			c := g.MustI2p(idx)
			closedList[c.String()+"^"] = true
			closedList[c.String()+">"] = true
			closedList[c.String()+"v"] = true
			closedList[c.String()+"<"] = true
		}
	}

	bestPaths := []*Node{}

	for len(openList) > 0 {
		// find smallest f value in open list
		slices.SortFunc(openList, func(l, r *Node) int {
			return l.score - r.score
		})

		currentNode := openList[0]
		debug("checking %s", currentNode)
		openList = openList[1:]
		closedList[currentNode.ID()] = true

		// shortest path found
		if currentNode.c == endC {
			if len(bestPaths) == 0 || currentNode.score <= bestPaths[0].score {
				bestPaths = append(bestPaths, currentNode)
			} else {
				break
			}
		}

		for _, t := range turns(currentNode.direction) {
			// node is in closed list, skip

			cost := 1
			if currentNode.direction != t.direction {
				cost += 1000
			}
			neighborN := &Node{
				c:         currentNode.c.Add(t.c),
				direction: t.direction,
				parent:    currentNode,
				cost:      cost,
				score:     currentNode.score + cost,
			}

			if closedList[neighborN.ID()] {
				debug("skipping %s: in closed list", neighborN)
				continue
			}

			openList = append(openList, neighborN)
		}
	}

	return bestPaths, len(bestPaths) > 0
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

func run(file string) (int, error) {
	in, err := os.ReadFile(file)
	if err != nil {
		return -1, err
	}
	g := NewGrid(in)

	// get start and end
	startI := bytes.IndexByte(g.field, 'S')
	if startI < 0 {
		panic("start not found")
	}
	endI := bytes.IndexByte(g.field, 'E')
	if endI < 0 {
		panic("end not found")
	}

	// The Reindeer start on the Start Tile (marked S)
	startC := g.MustI2p(startI)

	// and need to reach the End Tile (marked E)
	endC := g.MustI2p(endI)

	debug("initial")
	debug(g.String())

	paths, found := g.Solve(startC, endC)
	if !found {
		return -1, fmt.Errorf("could not find path")
	}
	if len(paths) == 0 {
		return -1, fmt.Errorf("could not find path (len=0)")
	}

	for _, p := range paths {
		for _, n := range p.GetPath() {
			g.MustSet(n.c, 'O')
		}
	}
	debug(g.String())

	fmt.Println("total score:", paths[0].score)
	fmt.Printf("found %d paths\n", len(paths))
	fmt.Println("best tiles:", bytes.Count(g.field, []byte{'O'}))

	return paths[0].score, nil
}

const debugEnabled = true

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
	if _, err := run("input.txt"); err != nil {
		fmt.Fprintf(os.Stderr, "puzzle errored: %v\n", err)
		os.Exit(1)
	}
}
