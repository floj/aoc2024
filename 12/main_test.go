package main

import (
	"testing"
)

func TestRunA(t *testing.T) {
	table := []struct {
		file     string
		expected int
	}{
		{file: "input-test-1.txt", expected: 140},
		{file: "input-test-2.txt", expected: 772},
		{file: "input-test-3.txt", expected: 1930},
		{file: "input.txt", expected: 1433460},
	}
	for _, td := range table {
		t.Run(td.file, func(t *testing.T) {
			sumA, _, err := runA(td.file)
			if err != nil {
				t.Fatalf("input %s failed with error: %v", td.file, err)
			}
			if sumA != td.expected {
				t.Fatalf("input %s produced wrong score for A. expected %d, got %d", td.file, td.expected, sumA)
			}
		})
	}
}

func TestRunB(t *testing.T) {
	table := []struct {
		file     string
		expected int
	}{
		{file: "input-test-1.txt", expected: 80},
		{file: "input-test-2.txt", expected: 436},
		{file: "input-test-b2.txt", expected: 236},
		{file: "input-test-b3.txt", expected: 368},
		{file: "input-test-3.txt", expected: 1206},
		// {file: "input.txt", expected: 1433460},
	}
	for _, td := range table {
		t.Run(td.file, func(t *testing.T) {
			_, sumB, err := runA(td.file)
			if err != nil {
				t.Fatalf("input %s failed with error: %v", td.file, err)
			}
			if sumB != td.expected {
				t.Fatalf("input %s produced wrong score for B. expected %d, got %d", td.file, td.expected, sumB)
			}
		})
	}
}
