package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Grid struct {
	field []byte
	cols  int
}

func NewGrid(w, h int) *Grid {
	field := make([]byte, w*h)
	for i := range field {
		field[i] = '.'
	}
	return &Grid{
		field: field,
		cols:  w,
	}
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

func (g *Grid) Inc(x, y int) (byte, bool) {
	if x < 0 || y < 0 {
		return 0, false
	}
	idx := y*g.cols + x
	if idx >= len(g.field) {
		return 0, false
	}
	if g.field[idx] == '.' {
		g.field[idx] = '1'
	} else {
		g.field[idx]++
	}
	return g.field[idx], true
}

func (g *Grid) Reset() {
	for i := range g.field {
		g.field[i] = '.'
	}
}

func (g *Grid) Rect(x, y, w, h int) *Grid {
	rect := Grid{
		cols: w,
	}

	for iy := y; iy < y+h; iy++ {
		pos := iy*g.cols + x
		rect.field = append(rect.field, g.field[pos:pos+w]...)
	}

	return &rect
}

type Robot struct {
	x    int
	y    int
	velX int
	velY int
}

func GetRobots(file string) ([]*Robot, error) {
	in, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	rr := []*Robot{}
	for _, line := range strings.Split(string(in), "\n") {
		pos, vel, ok := strings.Cut(line, " ")
		if !ok {
			return nil, fmt.Errorf("invalid line: %s", line)
		}
		if !strings.HasPrefix(pos, "p=") {
			return nil, fmt.Errorf("invalid line: %s", line)
		}
		posX, posY, ok := strings.Cut(pos[2:], ",")
		if !ok {
			return nil, fmt.Errorf("invalid line: %s", line)
		}

		if !strings.HasPrefix(vel, "v=") {
			return nil, fmt.Errorf("invalid line: %s", line)
		}
		velX, velY, ok := strings.Cut(vel[2:], ",")
		if !ok {
			return nil, fmt.Errorf("invalid line: %s", line)
		}

		r := &Robot{
			x:    MustAtoi(posX),
			y:    MustAtoi(posY),
			velX: MustAtoi(velX),
			velY: MustAtoi(velY),
		}

		rr = append(rr, r)
	}

	return rr, nil
}

func MustAtoi(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return v
}

func runA(file string, secs int) (int, error) {
	width, height := 101, 103
	// width, height := 11, 7
	g := NewGrid(width, height)

	robots, err := GetRobots(file)
	if err != nil {
		return -1, err
	}
	for _, r := range robots {
		fmt.Printf("%+v\n", r)
	}

	for round := 0; round <= secs; round++ {
		g.Reset()

		for _, r := range robots {
			// fmt.Printf("%+v\n", r)
			if _, ok := g.Inc(r.x, r.y); !ok {
				panic("increment failed")
			}
		}

		if bytes.Index(g.field, []byte("1111111111")) >= 0 {
			fmt.Println(round, "###############")
			fmt.Println(g)
		}

		for _, r := range robots {
			r.x += r.velX
			r.y += r.velY
			for r.x < 0 {
				r.x += width
			}
			for r.x >= width {
				r.x -= width
			}
			for r.y < 0 {
				r.y += height
			}
			for r.y >= height {
				r.y -= height
			}
		}
	}

	adjust := width % 2
	fmt.Println("done", "###############")
	fmt.Println(g)

	// get quadrants
	topL := g.Rect(0, 0, width/2, height/2)
	topR := g.Rect(width/2+adjust, 0, width/2, height/2)
	bottomL := g.Rect(0, height/2+adjust, width/2, height/2)
	bottomR := g.Rect(width/2+adjust, height/2+adjust, width/2, height/2)
	// fmt.Println("topL #############")
	// fmt.Println(topL)
	// fmt.Println("topR #############")
	// fmt.Println(topR)
	// fmt.Println("bottomL #############")
	// fmt.Println(bottomL)
	// fmt.Println("bottomR #############")
	// fmt.Println(bottomR)

	safetyFactor := 1

	for _, q := range []*Grid{topL, topR, bottomL, bottomR} {
		sumQ := 0
		for _, v := range q.field {
			if v != '.' {
				sumQ += int(v - '0')
			}
		}
		fmt.Println("sumQ", sumQ)
		safetyFactor *= sumQ
	}

	fmt.Println("safetyFactor", safetyFactor)
	return safetyFactor, nil
}

func runB(file string) (int, error) {
	width, height := 101, 103
	// width, height := 11, 7
	g := NewGrid(width, height)

	robots, err := GetRobots(file)
	if err != nil {
		return -1, err
	}
	for _, r := range robots {
		fmt.Printf("%+v\n", r)
	}

	treeFound := false
	round := 0
	for ; !treeFound; round++ {
		g.Reset()

		for _, r := range robots {
			// fmt.Printf("%+v\n", r)
			if _, ok := g.Inc(r.x, r.y); !ok {
				panic("increment failed")
			}
		}

		if bytes.Index(g.field, []byte("1111111111")) >= 0 {
			fmt.Println(round, "###############")
			fmt.Println(g)
			treeFound = true
			break
		}

		for _, r := range robots {
			r.x += r.velX
			r.y += r.velY
			for r.x < 0 {
				r.x += width
			}
			for r.x >= width {
				r.x -= width
			}
			for r.y < 0 {
				r.y += height
			}
			for r.y >= height {
				r.y -= height
			}
		}
	}

	fmt.Println("tree found in round", round)
	return round, nil
}

func main() {
	if _, err := runA("input.txt", 100); err != nil {
		fmt.Fprintf(os.Stderr, "puzzle errored: %v\n", err)
		os.Exit(1)
	}
	if _, err := runB("input.txt"); err != nil {
		fmt.Fprintf(os.Stderr, "puzzle errored: %v\n", err)
		os.Exit(1)
	}
}
