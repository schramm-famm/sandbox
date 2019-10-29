package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"text/template"
	"time"
)

/* TESTS */

func TestAbs(t *testing.T) {
	got := Abs(-1)
	if got != 1 {
		t.Errorf("Abs(-1) = %d; want 1", got)
	}
}

func TestTimeConsuming(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	time.Sleep(1000 * time.Millisecond)
}

/* SUBTESTS */

func TestFoo(t *testing.T) {
	// <setup code>
	t.Run("A=1", func(t *testing.T) {})
	t.Run("A=2", func(t *testing.T) {})
	t.Run("B=1", func(t *testing.T) {})
	// <tear-down code>
}

type Test struct {
	Name string
}

var tests = [5]Test{{"1"}, {"2"}, {"3"}, {"4"}, {"5"}}

func TestGroupedParallel(t *testing.T) {
	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
		})
	}
}

func TestTeardownParallel(t *testing.T) {
	// This Run will not return until the parallel tests finish.
	t.Run("group", func(t *testing.T) {
		t.Run("Test1", func(t *testing.T) {
			t.Parallel()
		})
		t.Run("Test2", func(t *testing.T) {
			t.Parallel()
		})
		t.Run("Test3", func(t *testing.T) {
			t.Parallel()
		})
	})
	// <tear-down code>
}

/* BENCHMARKS */

func BenchmarkHello(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("hello")
	}
}

func BenchmarkBigLen(b *testing.B) {
	big := NewBig()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		big.Len()
	}
}

func BenchmarkTemplateParallel(b *testing.B) {
	templ := template.Must(template.New("test").Parse("Hello, {{.}}!"))
	b.RunParallel(func(pb *testing.PB) {
		var buf bytes.Buffer
		for pb.Next() {
			buf.Reset()
			templ.Execute(&buf, "World")
		}
	})
}

/* EXAMPLES */

func ExampleHello() {
	fmt.Println("hello")
	// Output: hello
}

func ExampleSalutations() {
	fmt.Println("hello, and")
	fmt.Println("goodbye")
	// Output:
	// hello, and
	// goodbye
}

func ExamplePerm() {
	for _, value := range rand.Perm(5) {
		fmt.Println(value)
	}
	// Unordered output: 4
	// 2
	// 1
	// 3
	// 0
}

/* MAIN (optional) */

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}
