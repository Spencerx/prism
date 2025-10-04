package pkg1

import (
	"fmt"
	"strings"
	"testing"
)

func BenchmarkSimple(b *testing.B) {
	i := 0
	for b.Loop() {
		_ = fmt.Sprintf("test%d", i)
		i++
	}
}

func BenchmarkSimpleBuilder(b *testing.B) {
	i := 0
	for b.Loop() {
		var sb strings.Builder
		sb.WriteString("test")
		sb.WriteString(fmt.Sprintf("%d", i))
		_ = sb.String()
		i++
	}
}
