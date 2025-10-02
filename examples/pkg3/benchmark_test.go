package main

import (
	"fmt"
	"strings"
	"testing"
)

func BenchmarkSubtests(b *testing.B) {
	b.Run("group", func(b *testing.B) {
		b.Run("Test1", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = fmt.Sprintf("test%d", i)
			}
		})
		b.Run("Test2", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				var sb strings.Builder
				sb.WriteString("test")
				sb.WriteString(fmt.Sprintf("%d", i))
				_ = sb.String()
			}
		})
	})
}
