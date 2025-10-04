package main

import (
	"fmt"
	"strings"
	"testing"
)

func BenchmarkSubtests(b *testing.B) {
	b.Run("group", func(b *testing.B) {
		b.Run("Test1", func(b *testing.B) {
			i := 0
			for b.Loop() {
				_ = fmt.Sprintf("test%d", i)
				i++
			}
		})
		b.Run("Test2", func(b *testing.B) {
			i := 0
			for b.Loop() {
				var sb strings.Builder
				sb.WriteString("test")
				sb.WriteString(fmt.Sprintf("%d", i))
				_ = sb.String()
				i++
			}
		})
	})
}
