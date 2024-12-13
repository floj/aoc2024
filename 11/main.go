package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if err := runA("input.txt", 25); err != nil {
		fmt.Fprintf(os.Stderr, "puzzle errored: %v\n", err)
		os.Exit(1)
	}
}

type stoneRow []stone

type stone int

func (r stoneRow) String() string {
	buf := bytes.Buffer{}
	for _, v := range r {
		fmt.Fprintf(&buf, "%d ", v)
	}
	return buf.String()
}

func (r stoneRow) Blink() stoneRow {
	nextRow := make(stoneRow, 0, len(r))

	for _, s := range r {
		// If the stone is engraved with the number 0, it is replaced by a stone engraved with the number 1.
		if s == 0 {
			nextRow = append(nextRow, 1)
			continue
		}
		// If the stone is engraved with a number that has an even number of digits, it is replaced by two stones.
		// The left half of the digits are engraved on the new left stone, and the right half of the digits are
		// engraved on the new right stone.
		// poor mans check for length using strings
		// but quick and dirty is faster to type
		v := strconv.Itoa(int(s))
		if len(v)%2 == 0 {
			front, _ := strconv.Atoi(v[0 : len(v)/2])
			back, _ := strconv.Atoi(v[len(v)/2:])

			nextRow = append(nextRow, stone(front), stone(back))
			continue
		}

		// If none of the other rules apply, the stone is replaced by a new stone;
		// the old stone's number multiplied by 2024 is engraved on the new stone.
		nextRow = append(nextRow, stone(s*2024))
	}

	return nextRow
}

func runA(file string, blinks int) error {
	in, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	row := stoneRow{}

	for _, num := range strings.Split(string(in), " ") {
		num = strings.TrimSpace(num)
		if num == "" {
			continue
		}
		v, err := strconv.Atoi(num)
		if err != nil {
			return fmt.Errorf("could not parse '%s' as number: %w", num, err)
		}
		row = append(row, stone(v))
	}

	fmt.Fprintf(os.Stderr, "%s\n", row)
	for i := range blinks {
		row = row.Blink()
		fmt.Fprintf(os.Stderr, "round %d: %d\n", i+1, len(row))
		// fmt.Fprintf(os.Stderr, "%s\n", row)
	}

	return nil
}
