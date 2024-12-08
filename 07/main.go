package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type operation func(a, b int64) int64

var opsPartA = map[byte]operation{
	'+': func(a, b int64) int64 { return a + b },
	'*': func(a, b int64) int64 { return a * b },
}

var opsPartB = map[byte]operation{
	'+': func(a, b int64) int64 { return a + b },
	'*': func(a, b int64) int64 { return a * b },
	'|': func(a, b int64) int64 {
		i := strconv.FormatInt(a, 10) + strconv.FormatInt(b, 10)
		v, err := strconv.ParseInt(i, 10, 64)
		if err != nil {
			panic(err)
		}
		return v
	},
}

type calibration struct {
	value int64
	seq   []int64
}

func main() {

	sumA, err := run("input.txt", opsPartA)
	if err != nil {
		fmt.Fprintf(os.Stderr, "puzzle errored: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "A: sum of valid calibrations: %d\n", sumA)

	sumB, err := run("input.txt", opsPartB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "puzzle errored: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "B: sum of valid calibrations: %d\n", sumB)

}

func newCalibration(line string) (calibration, error) {
	v, s, found := strings.Cut(line, ":")
	c := calibration{}
	if !found {
		return c, fmt.Errorf("no ':' found in line: %s", line)
	}
	value, err := strconv.Atoi(v)
	if err != nil {
		return c, fmt.Errorf("invalid test value %s: %w", v, err)
	}
	c.value = int64(value)
	s = strings.TrimSpace(s)
	for _, i := range strings.Split(s, " ") {
		v, err := strconv.Atoi(i)
		if err != nil {
			return c, fmt.Errorf("invalid sequence value %s in %s: %w", i, s, err)
		}
		c.seq = append(c.seq, int64(v))
	}
	return c, nil
}

func (c *calibration) IsValid(ops map[byte]operation) bool {
	if len(c.seq) == 0 {
		return false
	}

	s := c.seq[0]
	for k := range ops {
		if solve(s, c.seq[1:], k, ops, c.value) {
			return true
		}
	}
	return false
}

func solve(sum int64, rest []int64, op byte, ops map[byte]operation, expected int64) bool {
	if len(rest) == 0 {
		return expected == sum
	}

	s := ops[op](sum, rest[0])

	for k := range ops {
		if solve(s, rest[1:], k, ops, expected) {
			return true
		}
	}

	return false
}

func run(inputFile string, ops map[byte]operation) (int64, error) {
	in, err := os.ReadFile(inputFile)
	if err != nil {
		return -1, err
	}

	sum := int64(0)
	for _, line := range strings.Split(string(in), "\n") {
		c, err := newCalibration(line)
		if err != nil {
			return -1, err
		}
		if c.IsValid(ops) {
			sum += c.value
		}
	}

	return sum, nil
}
