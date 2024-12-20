package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Input struct {
	Towels   []string
	Patterns []string
}

type Tree map[string]Tree

func readInput(file string) (Input, error) {
	f, err := os.Open(file)
	if err != nil {
		return Input{}, err
	}
	defer f.Close()

	in := Input{}
	s := bufio.NewScanner(f)
	// towels
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" {
			break
		}
		in.Towels = append(in.Towels, strings.Split(line, ",")...)
	}

	// patterns
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		in.Patterns = append(in.Patterns, line)
	}

	if err := s.Err(); err != nil {
		return Input{}, err
	}

	for i := range in.Towels {
		in.Towels[i] = strings.TrimSpace(in.Towels[i])
	}

	// sort longer towles first
	slices.SortFunc(in.Towels, func(l, r string) int {
		// if v := len(r) - len(l); v != 0 {
		if v := len(l) - len(r); v != 0 {
			return v
		}
		return strings.Compare(l, r)
	})

	return in, nil
}

func matchPattern(num int, pattern string, towels []string) int {
	// find candidate patterns

	candidates := []string{}
	for _, t := range towels {
		if strings.Contains(pattern, t) {
			candidates = append(candidates, t)
		}
	}

	fmt.Fprintf(debugOut, "%v\n", candidates)

	return tryMatch([]string{strconv.Itoa(num)}, map[string]int{}, pattern, pattern, candidates)
}

func tryMatch(progress []string, suffixCache map[string]int, org, test string, towels []string) int {
	if len(test) == 0 {
		return 1
	}
	fmt.Fprintf(debugOut, "%2d | %v checking %s\n", len(test), progress, test)
	sum := 0
	for _, t := range towels {
		if !strings.HasPrefix(test, t) {
			continue
		}
		remaining := test[len(t):]
		if cache, hit := suffixCache[remaining]; hit {
			sum += cache
			continue
		}
		s := tryMatch(append(progress[:], t), suffixCache, org, remaining, towels)
		suffixCache[remaining] = s
		sum += s
	}
	return sum
}

func run(file string) error {
	in, err := readInput(file)
	if err != nil {
		return err
	}

	comb, matched := 0, 0
	for i, p := range in.Patterns {
		s := matchPattern(i+1, p, in.Towels)
		if s > 0 {
			matched++
		}
		comb += s
	}
	fmt.Println("matched", matched)
	fmt.Println("combinations", comb)

	return nil
}

var debugOut = io.Discard

func main() {
	if err := run("input.txt"); err != nil {
		fmt.Fprintf(os.Stderr, "puzzle errored: %v", err)
		os.Exit(1)
	}
}
