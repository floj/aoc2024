package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "puzzle errored with %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	in, err := loadInput("input.txt")
	if err != nil {
		return err
	}

	correct := in.sumCorrectUpdates()
	fmt.Fprintf(os.Stdout, "correct updates middle page sum: %d\n", correct)

	incorrect, err := in.sumIncorrectUpdates()
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "ordered incorrect updates middle page sum: %d\n", incorrect)

	return nil
}

func (in Input) sumIncorrectUpdates() (int, error) {
	sum := 0
	for _, u := range in.incorrectUpdates() {
		ordered, err := u.orderByRules(in.rules)
		if err != nil {
			return -1, fmt.Errorf("failed to order %v: %w", u, err)
		}
		sum += ordered[len(ordered)/2]
	}
	return sum, nil
}

func (u update) orderByRules(rules []Rule) (update, error) {
	remaining := u[:]
	ordered := update{}

	for len(remaining) > 0 {
		constrMap := buildConstraints(remaining, rules)
		// fmt.Printf("remaining %+v\n", remaining)
		// fmt.Printf("  ordered %+v\n", ordered)
		// for _, p := range remaining {
		// 	fmt.Printf("%d => %+v\n", p, constrMap[p])
		// }

		// find page without contstraints
		found := false
		for k, v := range constrMap {
			if len(v) == 0 {
				// fmt.Println("using constraint", k, v)
				ordered = append(ordered, k)
				remaining = slices.DeleteFunc(remaining, func(e int) bool {
					return e == k
				})
				found = true
				break
			}
		}
		if !found {
			return ordered, fmt.Errorf("no empty constraint found")
		}
	}
	return ordered, nil
}

func buildConstraints(pages []int, rules []Rule) map[int][]int {
	constr := map[int][]int{}
	for _, r := range rules {
		constr[r.behind] = append(constr[r.behind], r.before)
	}

	m := map[int][]int{}
	for _, page := range pages {
		//only keep relevant page orderings
		cc := []int{}
		for _, p := range pages {
			if slices.Contains(constr[page], p) {
				cc = append(cc, p)
			}
		}
		slices.Sort(cc)
		m[page] = cc
	}
	return m
}

type Input struct {
	rules   []Rule
	updates []update
}

func (in Input) sumCorrectUpdates() int {
	sum := 0
	for _, u := range in.updates {
		middlePage, ok := u.Check(in.rules)
		if !ok {
			continue
		}
		sum += middlePage
		//fmt.Fprintf(os.Stderr, "update %v correct, using page %d\n", u, middlePage)
	}
	return sum
}

func (in Input) findBeforeRules(page int) []Rule {
	rules := []Rule{}
	for _, r := range in.rules {
		if r.before == page {
			rules = append(rules, r)
		}
	}
	return rules
}

func (in Input) findBehindRules(page int) []Rule {
	rules := []Rule{}
	for _, r := range in.rules {
		if r.behind == page {
			rules = append(rules, r)
		}
	}
	return rules
}

func (in Input) incorrectUpdates() []update {
	incorrect := []update{}
	for _, u := range in.updates {
		if _, ok := u.Check(in.rules); ok {
			continue
		}
		incorrect = append(incorrect, u)

	}
	return incorrect
}

type update []int

func (u update) Check(rules []Rule) (int, bool) {
	for _, r := range rules {
		beforeIdx := slices.Index(u, r.before)
		behindIdx := slices.Index(u, r.behind)
		if beforeIdx < 0 || behindIdx < 0 {
			// page not in update, ignore
			continue
		}
		if beforeIdx >= behindIdx {
			return -1, false
		}
	}
	return u[len(u)/2], true
}

func newUpdate(s string) (update, error) {
	u := update{}
	for _, p := range strings.Split(s, ",") {
		page, err := strconv.Atoi(p)
		if err != nil {
			return update{}, fmt.Errorf("invalid page number %s: %w", p, err)
		}
		u = append(u, page)
	}
	return u, nil
}

type Rule struct {
	before int
	behind int
}

func newRule(s string) (Rule, error) {
	parts := strings.Split(s, "|")
	if len(parts) != 2 {
		return Rule{}, fmt.Errorf("page ordering rule must be exactly two fields but found %d for %s", len(parts), s)
	}
	page, err := strconv.Atoi(parts[0])
	if err != nil {
		return Rule{}, fmt.Errorf("invalid page number %s: %w", parts[0], err)
	}
	before, err := strconv.Atoi(parts[1])
	if err != nil {
		return Rule{}, fmt.Errorf("invalid page number %s: %w", parts[0], err)
	}
	return Rule{before: page, behind: before}, nil
}

func loadInput(file string) (Input, error) {
	i := Input{}

	f, err := os.Open(file)
	if err != nil {
		return i, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	// process first section
	for scanner.Scan() {
		line := scanner.Text()

		// empty line indicates switch to 2nd section with update instructions
		if line == "" {
			break
		}
		r, err := newRule(line)
		if err != nil {
			return i, err
		}
		i.rules = append(i.rules, r)
	}

	if err := scanner.Err(); err != nil {
		return i, err
	}
	for scanner.Scan() {
		line := scanner.Text()
		u, err := newUpdate(line)
		if err != nil {
			return i, err
		}
		i.updates = append(i.updates, u)
	}

	return i, scanner.Err()
}
