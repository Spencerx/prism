package pkg1

import (
	"fmt"
	"strings"
	"testing"
)

// {"Time":"2025-10-03T08:05:50.159422858-05:00","Action":"start","Package":"go.dalton.dog/prism/examples/pkg1"}
// {"Time":"2025-10-03T08:05:50.163939547-05:00","Action":"output","Package":"go.dalton.dog/prism/examples/pkg1","Output":"goos: linux\n"}
// {"Time":"2025-10-03T08:05:50.164055101-05:00","Action":"output","Package":"go.dalton.dog/prism/examples/pkg1","Output":"goarch: amd64\n"}
// {"Time":"2025-10-03T08:05:50.164067724-05:00","Action":"output","Package":"go.dalton.dog/prism/examples/pkg1","Output":"pkg: go.dalton.dog/prism/examples/pkg1\n"}
// {"Time":"2025-10-03T08:05:50.164079518-05:00","Action":"output","Package":"go.dalton.dog/prism/examples/pkg1","Output":"cpu: 13th Gen Intel(R) Core(TM) i7-13700H\n"}
// {"Time":"2025-10-03T08:05:50.164094873-05:00","Action":"run","Package":"go.dalton.dog/prism/examples/pkg1","Test":"BenchmarkSimple"}
// {"Time":"2025-10-03T08:05:50.164103506-05:00","Action":"output","Package":"go.dalton.dog/prism/examples/pkg1","Test":"BenchmarkSimple","Output":"=== RUN   BenchmarkSimple\n"}
// {"Time":"2025-10-03T08:05:50.164112648-05:00","Action":"output","Package":"go.dalton.dog/prism/examples/pkg1","Test":"BenchmarkSimple","Output":"BenchmarkSimple\n"}
// {"Time":"2025-10-03T08:05:52.631809113-05:00","Action":"output","Package":"go.dalton.dog/prism/examples/pkg1","Test":"BenchmarkSimple","Output":"BenchmarkSimple-20           \t 7098393\t       177.9 ns/op\n"}
// {"Time":"2025-10-03T08:05:52.631879584-05:00","Action":"run","Package":"go.dalton.dog/prism/examples/pkg1","Test":"BenchmarkSimpleBuilder"}
// {"Time":"2025-10-03T08:05:52.631886375-05:00","Action":"output","Package":"go.dalton.dog/prism/examples/pkg1","Test":"BenchmarkSimpleBuilder","Output":"=== RUN   BenchmarkSimpleBuilder\n"}
// {"Time":"2025-10-03T08:05:52.631894016-05:00","Action":"output","Package":"go.dalton.dog/prism/examples/pkg1","Test":"BenchmarkSimpleBuilder","Output":"BenchmarkSimpleBuilder\n"}
// {"Time":"2025-10-03T08:05:54.183482818-05:00","Action":"output","Package":"go.dalton.dog/prism/examples/pkg1","Test":"BenchmarkSimpleBuilder","Output":"BenchmarkSimpleBuilder-20    \t 5313279\t       248.4 ns/op\n"}
// {"Time":"2025-10-03T08:05:54.183549086-05:00","Action":"output","Package":"go.dalton.dog/prism/examples/pkg1","Output":"PASS\n"}
// {"Time":"2025-10-03T08:05:54.185182533-05:00","Action":"output","Package":"go.dalton.dog/prism/examples/pkg1","Output":"ok  \tgo.dalton.dog/prism/examples/pkg1\t4.025s\n"}
// {"Time":"2025-10-03T08:05:54.185212667-05:00","Action":"pass","Package":"go.dalton.dog/prism/examples/pkg1","Elapsed":4.026}

func BenchmarkSimple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("test%d", i)
	}
}

func BenchmarkSimpleBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var sb strings.Builder
		sb.WriteString("test")
		sb.WriteString(fmt.Sprintf("%d", i))
		_ = sb.String()
	}
}
