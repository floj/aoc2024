package main

import (
	"testing"
)

func TestRunA(t *testing.T) {
	table := []struct {
		file     string
		expected int
	}{
		{file: "test-1.txt", expected: 1},
		{file: "test-2.txt", expected: 2},
		{file: "test-3.txt", expected: 4},
		{file: "test-4.txt", expected: 3},
		{file: "test-5.txt", expected: 36},
	}
	for _, td := range table {
		t.Run(td.file, func(t *testing.T) {
			score, _, err := run(td.file)
			if err != nil {
				t.Fatalf("input %s failed with error: %v", td.file, err)
			}
			if score != td.expected {
				t.Fatalf("input %s produced wrong score. expected %d, got %d", td.file, td.expected, score)
			}
		})
	}
}
