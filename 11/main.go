package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

func main() {
	if err := runA("input.txt", 25); err != nil {
		fmt.Fprintf(os.Stderr, "puzzle errored: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("#################################")
	if err := runA("input.txt", 75); err != nil {
		fmt.Fprintf(os.Stderr, "puzzle errored: %v\n", err)
		os.Exit(1)
	}

}

var seqCache = &sync.Map{}

func getCache(remaining, s int) (uint64, bool) {
	b := bytes.Buffer{}
	fmt.Fprint(&b, remaining, '|', s)
	if v, ok := seqCache.Load(b.String()); ok {
		return v.(uint64), true
	}
	return 0, false
}

func setCache(remaining, s int, v uint64) uint64 {
	b := bytes.Buffer{}
	fmt.Fprint(&b, remaining, '|', s)
	seqCache.LoadOrStore(b.String(), v)
	return v
}

func CountStones(s, remaining int) uint64 {
	if remaining <= 0 {
		return 1
	}

	// check sequence cache
	if v, ok := getCache(remaining, s); ok {
		return v
	}

	// If the stone is engraved with the number 0, it is replaced by a stone engraved with the number 1.
	if s == 0 {
		v := CountStones(1, remaining-1)
		return setCache(remaining, s, v)
	}

	// If the stone is engraved with a number that has an even number of digits, it is replaced by two stones.
	// The left half of the digits are engraved on the new left stone, and the right half of the digits are
	// engraved on the new right stone.
	// poor mans check for length using strings
	// but quick and dirty is faster to type
	num := strconv.Itoa(s)
	if len(num)%2 == 0 {
		front, _ := strconv.Atoi(num[0 : len(num)/2])
		s1 := CountStones(front, remaining-1)
		back, _ := strconv.Atoi(num[len(num)/2:])
		s2 := CountStones(back, remaining-1)
		v := s1 + s2
		return setCache(remaining, s, v)
	}

	// If none of the other rules apply, the stone is replaced by a new stone;
	// the old stone's number multiplied by 2024 is engraved on the new stone.
	v := CountStones(s*2024, remaining-1)
	return setCache(remaining, s, v)
}

func runA(file string, blinks int) error {
	in, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	stones := []int{}
	for _, v := range strings.Split(string(in), " ") {
		num, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("could not parse '%s' as number: %w", v, err)
		}
		stones = append(stones, num)
	}

	wg := &sync.WaitGroup{}
	sum := &atomic.Uint64{}

	for i, num := range stones {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v := CountStones(num, blinks)
			fmt.Printf("stone %d with value %d evolves info %d stones\n", i, num, v)
			sum.Add(v)
		}()
	}
	wg.Wait()

	fmt.Fprintf(os.Stderr, "number of stones after %d blinks: %d\n", blinks, sum.Load())
	return nil
}
