package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

type coord struct {
	x, y int
}

func (c coord) String() string {
	return fmt.Sprintf("(%5d,%5d)", c.x, c.y)
}

func (c coord) Add(o coord) coord {
	return coord{x: c.x + o.x, y: c.y + o.y}
}

func (c coord) Dist(o coord) float64 {
	distX := math.Abs(float64(c.x) - float64(o.x))
	distY := math.Abs(float64(c.y) - float64(o.y))
	dist := math.Pow(distX, 2) + math.Pow(distY, 2)
	return dist
}

func (c coord) Mul(i int) coord {
	return coord{x: c.x * i, y: c.y * i}
}

type ClawConf struct {
	A     coord
	B     coord
	Prize coord
}

func CoordFromLine(line string) (coord, error) {
	_, spec, _ := strings.Cut(line, ":")
	xs, ys, ok := strings.Cut(spec, ",")
	if !ok {
		return coord{}, fmt.Errorf("invalid line: %s", line)
	}
	x, err := strconv.Atoi(strings.TrimSpace(xs)[2:])
	if err != nil {
		return coord{}, fmt.Errorf("invalid line '%s': %w", line, err)
	}
	y, err := strconv.Atoi(strings.TrimSpace(ys)[2:])
	if err != nil {
		return coord{}, fmt.Errorf("invalid line '%s': %w", line, err)
	}
	return coord{x: x, y: y}, nil
}

func GetClawConf(file string, prizeShift coord) ([]ClawConf, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	clawConf := []ClawConf{}
	scn := bufio.NewScanner(f)
	c := ClawConf{}
	for scn.Scan() {

		line := scn.Text()
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "Button A:") {
			c.A, err = CoordFromLine(line)
			if err != nil {
				return nil, err
			}
			continue
		}

		if strings.HasPrefix(line, "Button B:") {
			c.B, err = CoordFromLine(line)
			if err != nil {
				return nil, err
			}
			continue
		}

		if strings.HasPrefix(line, "Prize:") {
			prize, err := CoordFromLine(line)
			if err != nil {
				return nil, err
			}
			c.Prize = prize.Add(prizeShift)
			clawConf = append(clawConf, c)
			c = ClawConf{}
			continue
		}
		return nil, fmt.Errorf("unknown directive: %s", line)
	}

	return clawConf, scn.Err()
}

type Node struct {
	c      coord
	parent *Node
	score  int
	cost   int
	btn    string
}

func (n *Node) ID() string {
	return n.c.String()
}

func (n *Node) String() string {
	return fmt.Sprintf("%s (%d)", n.c, n.score)
}

func (n *Node) TotalCost() map[string]int {
	s := map[string]int{}
	c := n
	for {
		if c == nil {
			return s
		}
		s["sum"] += c.cost
		if c.btn != "" {
			s[c.btn]++
		}
		c = c.parent
	}
}

func (cc ClawConf) PrecalcMoves(n *Node, start, end, step int) []*Node {
	nodes := []*Node{}

	add, btn, cost := cc.A, "A", 3
	for i := start; i <= end; i = i + step {
		c := n.c.Add(add.Mul(i))
		if c.x <= cc.Prize.x && c.y <= cc.Prize.y {
			nodes = append(nodes, &Node{c: c, parent: n, score: n.score + (cost * i), cost: cost, btn: btn})
		}
	}

	add, btn, cost = cc.B, "B", 1
	for i := start; i <= end; i = i + step {
		c := n.c.Add(add.Mul(i))
		if c.x <= cc.Prize.x && c.y <= cc.Prize.y {
			nodes = append(nodes, &Node{c: c, parent: n, score: n.score + (cost * i), cost: cost, btn: btn})
		}
	}
	return nodes
}

type TurnFn func(n *Node) []*Node

func (cc ClawConf) Moves(n *Node) []*Node {
	nodes := []*Node{}
	cA := n.c.Add(cc.A)
	if cA.x <= cc.Prize.x && cA.y <= cc.Prize.y {
		// nodes = append(nodes, &Node{c: cA, parent: n, score: n.score + cA.Dist(cc.Prize), cost: 1})
		nodes = append(nodes, &Node{c: cA, parent: n, score: n.score + 3, cost: 3, btn: "A"})
	}
	cB := n.c.Add(cc.B)
	if cB.x <= cc.Prize.x && cB.y <= cc.Prize.y {
		// nodes = append(nodes, &Node{c: cB, parent: n, score: n.score + cB.Dist(cc.Prize)*3, cost: 3})
		nodes = append(nodes, &Node{c: cB, parent: n, score: n.score + 1, cost: 1, btn: "B"})
	}
	return nodes
}

// use a* to find shortest path
// implementation from https://de.wikipedia.org/wiki/A*-Algorithmus
func (cc ClawConf) Solve(tfn TurnFn) *Node {
	openList := []*Node{{c: coord{}}}
	closedList := map[string]bool{}

	for len(openList) > 0 {
		// find smallest f value in open list
		slices.SortFunc(openList, func(l, r *Node) int {
			return r.score - l.score
		})

		currentNode := openList[0]

		openList = openList[1:]
		closedList[currentNode.ID()] = true

		// fmt.Fprintf(os.Stderr, "[%4d | %4d ] checking %v\n", len(openList), len(closedList), currentNode)

		// shortest path found
		if currentNode.c == cc.Prize {
			return currentNode
		}

		for _, t := range tfn(currentNode) {
			// node is in closed list, skip
			if closedList[t.ID()] {
				continue
			}
			openList = append(openList, t)
		}
	}

	return nil
}

func runA(file string) error {
	confs, err := GetClawConf(file, coord{})
	if err != nil {
		return err
	}

	total := 0
	for _, conf := range confs {
		fmt.Fprintf(debugW, "%+v\n", conf)
		n := conf.Solve(conf.Moves)
		if n == nil {
			fmt.Println("no path found")
			continue
		}
		costs := n.TotalCost()
		fmt.Println("costs", costs)
		total += costs["sum"]
	}

	fmt.Println("total", total)
	return nil
}

func Precalc(cc ClawConf) TurnFn {
	return func(n *Node) []*Node {
		moves := cc.Moves(n)
		for i := 100000; i < 10000000000000/1000; i = i * 10 {
			moves = append(moves, cc.PrecalcMoves(n, i, i+i, i/10)...)
		}
		return moves
	}
}

func runB(file string) error {
	confs, err := GetClawConf(file, coord{x: 10000000000000, y: 10000000000000})
	if err != nil {
		return err
	}

	total := 0
	for _, conf := range confs {
		fmt.Fprintf(debugW, "%+v\n", conf)
		n := conf.Solve(Precalc(conf))
		if n == nil {
			fmt.Println("no path found")
			continue
		}
		costs := n.TotalCost()
		fmt.Println("costs", costs)
		total += costs["sum"]
	}

	fmt.Println("total", total)
	return nil
}

var debugW = os.Stderr

func main() {
	if err := runA("input.txt"); err != nil {
		fmt.Fprintf(os.Stderr, "puzzle errored with %v\n", err)
		os.Exit(1)
	}
	// if err := runB("input-test.txt"); err != nil {
	// 	fmt.Fprintf(os.Stderr, "puzzle errored with %v\n", err)
	// 	os.Exit(1)
	// }
}
