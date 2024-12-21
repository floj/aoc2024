package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
)

func main() {
	if _, _, err := runA("input.txt"); err != nil {
		fmt.Fprintf(os.Stderr, "puzzle errored: %v\n", err)
		os.Exit(1)
	}
}

type Grid struct {
	field []byte
	cols  int

	visited []byte
}

func NewGrid(file string) (*Grid, error) {
	in, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	cols := bytes.Index(in, []byte{'\n'})
	field := bytes.ReplaceAll(in, []byte{'\n'}, []byte{})

	return &Grid{
		field:   field,
		cols:    cols,
		visited: make([]byte, len(field)),
	}, nil
}

func (g *Grid) String() string {
	b := bytes.Buffer{}
	for i := 0; i < len(g.field); i = i + g.cols {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.Write(g.field[i : i+g.cols])
	}
	return b.String()
}

func (g *Grid) p2i(x, y int) (int, bool) {
	if x < 0 || y < 0 {
		return -1, false
	}
	if c.x >= g.cols {
		return -1, false
	}
	idx := y*g.cols + x
	if idx >= len(g.field) {
		return -1, false
	}
	return idx, true
}

func (g *Grid) i2p(idx int) (int, int, bool) {
	if idx < 0 || idx >= len(g.field) {
		return -1, -1, false
	}
	return idx % g.cols, idx / g.cols, true
}

func (g *Grid) Get(x, y int) (byte, bool) {
	if idx, ok := g.p2i(x, y); ok {
		return g.field[idx], true
	}
	return 0, false
}

func (g *Grid) GetOffset(idx, offX, offY int) (byte, bool) {
	i := idx + offX + (offY * g.cols)
	if i < 0 || i >= len(g.field) {
		return 0, false
	}
	return g.field[i], true
}

func (g *Grid) Set(x, y int, v byte) (byte, bool) {
	if idx, ok := g.p2i(x, y); ok {
		o := g.field[idx]
		g.field[idx] = v
		return o, true
	}
	return 0, false
}

type Neighbor struct {
	x       int
	y       int
	bitmask byte
}

var neighbors = map[string]Neighbor{
	"top":    {x: 0, y: -1, bitmask: 0b0100},
	"bottom": {x: 0, y: +1, bitmask: 0b0001},
	"right":  {x: +1, y: 0, bitmask: 0b0010},
	"left":   {x: -1, y: 0, bitmask: 0b1000},
}

func (g *Grid) FloodFill(idx int, v byte) []int {
	// out of bounds
	if idx < 0 || idx >= len(g.field) {
		return nil
	}
	// already visited
	if g.visited[idx] > 0 {
		return nil
	}
	// incorrect value
	if g.field[idx] != v {
		return nil
	}
	// matches, fill field
	pos := []int{idx}
	g.visited[idx]++
	for _, n := range neighbors {
		pos = append(pos, g.FloodFill(idx+n.x+(g.cols*n.y), v)...)
	}
	sort.Ints(pos)
	return pos
}

func (g *Grid) RectFrom(pos []int, v byte) *Rect {
	if len(pos) == 0 {
		return nil
	}

	fields := make([]byte, len(g.field))
	for _, p := range pos {
		fields[p] = v
	}

	minX, minY := pos[0]%g.cols, pos[0]/g.cols
	maxX, maxY := minX, minY

	for _, p := range pos[1:] {
		x, y := p%g.cols, p/g.cols
		minX, maxX = min(minX, x), max(maxX, x)
		minY, maxY = min(minY, y), max(maxY, y)
	}

	// fmt.Println("rect", "min", minX, minY, "max", maxX, maxY)

	rect := Rect{
		v:    v,
		cols: maxX - minX + 1,
	}

	for r := minY; r <= maxY; r++ {
		posFrom := r*g.cols + minX
		posTo := r*g.cols + maxX + 1
		rect.field = append(rect.field, fields[posFrom:posTo]...)
	}

	for i := range rect.field {
		if rect.field[i] != v {
			rect.field[i] = '.'
		}
	}

	return &rect
}

type Rect struct {
	v      byte
	field  []byte
	cols   int
	sidesI int
	sidesO int
	peri   int
}

func (r *Rect) Perimeter() int {
	if r.peri > 0 {
		return r.peri
	}

	for i, v := range r.field {
		if v != r.v {
			continue
		}
		ix, iy := i%r.cols, i/r.cols
		for _, n := range neighbors {
			nx, ny := ix+n.x, iy+n.y
			npos := ny*r.cols + nx
			if nx < 0 || nx >= r.cols {
				r.peri++
				continue
			}
			if npos < 0 || npos >= len(r.field) {
				r.peri++
				continue
			}

			// other area, use fence
			if r.field[npos] != v {
				r.peri++
				continue
			}
			// else same area, no fence
		}
	}

	return r.peri
}

func (r *Rect) Sides() int {
	return r.SidesInner() + r.SidesOuter()
}

type touchpoint struct {
	offX int
	offY int
	bm   byte
}

func (r *Rect) SidesOuter() int {
	if r.sidesO > 0 {
		return r.sidesO
	}
	sides := make([]byte, len(r.field))

	// first set flags for all sides
	for i, v := range r.field {
		if v != r.v {
			continue
		}
		ix, iy := i%r.cols, i/r.cols
		for _, n := range neighbors {
			nx, ny := ix+n.x, iy+n.y
			npos := ny*r.cols + nx
			if nx < 0 || nx >= r.cols {
				sides[i] = sides[i] | n.bitmask
				continue
			}
			if npos < 0 || npos >= len(r.field) {
				sides[i] = sides[i] | n.bitmask
				continue
			}

			// other area, use fence
			if r.field[npos] != v {
				sides[i] = sides[i] | n.bitmask
				continue
			}
			// else same area, no fence
		}
	}

	corners := [][]string{
		// outer corners
		{"top", "left"},
		{"top", "right"},
		{"bottom", "right"},
		{"bottom", "left"},
	}
	touch := map[string]touchpoint{
		"top-left":     {offX: -1, offY: -1, bm: neighbors["bottom"].bitmask | neighbors["right"].bitmask},
		"top-right":    {offX: +1, offY: -1, bm: neighbors["bottom"].bitmask | neighbors["left"].bitmask},
		"bottom-left":  {offX: -1, offY: +1, bm: neighbors["top"].bitmask | neighbors["right"].bitmask},
		"bottom-right": {offX: +1, offY: +1, bm: neighbors["top"].bitmask | neighbors["left"].bitmask},
	}

	// count corners, this is equivalent to sides
	for i, v := range sides {
		ix, iy := i%r.cols, i/r.cols

		for _, c := range corners {
			bm := neighbors[c[0]].bitmask | neighbors[c[1]].bitmask
			res := v & bm
			if res == bm {
				r.sidesO++
				fmt.Printf("outer %d %d %s-%s mask=%04b v=%04b r=%04b %t %d\n", ix, iy, c[0], c[1], bm, v, res, res == bm, r.sidesO)

				// hack! see if another corner is touching
				if t, ok := touch[c[0]+"-"+c[1]]; ok {
					tx, ty := ix+t.offX, iy+t.offY

					tpos := ty*r.cols + tx
					if tx < 0 || tx >= r.cols || tpos < 0 || tpos >= len(r.field) {
						continue
					}

					if sides[tpos]&t.bm == t.bm {
						r.sidesO--
						fmt.Printf("outer touch %d %d %d %d %d\n", ix, iy, tx, ty, r.sidesO)
					}
				}
			}
		}
	}

	return r.sidesO
}

func (r *Rect) SidesInner() int {
	if r.sidesI > 0 {
		return r.sidesI
	}

	sides := make([]byte, len(r.field))

	// first set flags for all sides
	for i, v := range r.field {
		if v == r.v {
			continue
		}
		ix, iy := i%r.cols, i/r.cols
		for _, n := range neighbors {
			nx, ny := ix+n.x, iy+n.y
			npos := ny*r.cols + nx
			if nx < 0 || nx >= r.cols {
				continue
			}
			if npos < 0 || npos >= len(r.field) {
				continue
			}
			if r.field[npos] == v {
				continue
			}
			sides[i] = sides[i] | n.bitmask
			// else same area, no fence
		}
	}

	corners := [][]string{
		// outer corners
		{"top", "left"},
		{"top", "right"},
		{"bottom", "right"},
		{"bottom", "left"},
	}

	// count corners, this is equivalent to sides
	for i, v := range sides {
		ix, iy := i%r.cols, i/r.cols

		for _, c := range corners {
			bm := neighbors[c[0]].bitmask | neighbors[c[1]].bitmask
			res := v & bm
			if res == bm {
				r.sidesI++
				fmt.Printf("inner %d %d %s-%s mask=%04b v=%04b r=%04b %t %d\n", ix, iy, c[0], c[1], bm, v, res, res == bm, r.sidesI)
			}
		}
	}

	return r.sidesI
}

func (r *Rect) Area() int {
	sum := 0
	for _, v := range r.field {
		if v == r.v {
			sum++
		}
	}
	return sum
}

func (r *Rect) String() string {
	b := bytes.Buffer{}
	fmt.Fprintf(&b, "area %s cols=%d size=%d peri=%d sides=%d\n", []byte{r.v}, r.cols, r.Area(), r.Perimeter(), r.Sides())
	for i := 0; i < len(r.field); i = i + r.cols {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.Write(r.field[i : i+r.cols])
	}
	return b.String()
}

func runA(file string) (int, int, error) {
	g, err := NewGrid(file)
	if err != nil {
		return -1, -1, err
	}

	areas := []*Rect{}

	for {
		found := false
		for i, v := range g.field {
			if g.visited[i] > 0 {
				continue
			}
			area := g.FloodFill(i, v)
			// fmt.Printf("filled at %d: %+v\n", i, area)
			areas = append(areas, g.RectFrom(area, v))
			found = true
			break
		}
		if !found {
			break
		}
	}

	sumA, sumB := 0, 0
	// fmt.Println("areas")
	for _, a := range areas {
		sumA += a.Area() * a.Perimeter()
		sumB += a.Area() * a.Sides()
		fmt.Printf("%s\n", a)
		// fmt.Println("###########")
	}

	fmt.Println("sum A", sumA)
	fmt.Println("sum B", sumB)
	return sumA, sumB, nil
}
