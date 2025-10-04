package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestWithBenchmark(t *testing.T) {
	// This test passes but shows we can handle different test types
	result := strings.Repeat("a", 1000)
	if len(result) != 1000 {
		t.Errorf("Expected length 1000, got %d", len(result))
	}
}

func BenchmarkStringConcat(b *testing.B) {
	i := 0
	for b.Loop() {
		_ = fmt.Sprintf("test%d", i)
		i++
	}
}

func BenchmarkStringBuilder(b *testing.B) {
	i := 0
	for b.Loop() {
		var sb strings.Builder
		sb.WriteString("test")
		sb.WriteString(fmt.Sprintf("%d", i))
		_ = sb.String()
		i++
	}
}
