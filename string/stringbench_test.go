package string_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func BenchmarkConnectStringWithOperator(b *testing.B) {
	hello := "hello"
	world := "world"
	for i := 0; i < b.N; i++ {
		hello = hello + world
	}
}

func BenchmarkConnectStringWithSprintf(b *testing.B) {
	hello := "hello"
	world := "world"
	for i := 0; i < b.N; i++ {
		hello = fmt.Sprintf("%s,%s", hello, world)
	}
}

func BenchmarkConnectStringWithJoin(b *testing.B) {
	hello := "hello"
	world := "world"
	for i := 0; i < b.N; i++ {
		hello = strings.Join([]string{hello}, world)
	}
}

func BenchmarkConnectStringWithBuffer(b *testing.B) {
	hello := "hello"
	world := "world"
	var buffer bytes.Buffer
	buffer.WriteString(hello)
	for i := 0; i < b.N; i++ {
		buffer.WriteString(world)
	}
	_ = buffer.String()
}

func BenchmarkConnectStringWithBuilder(b *testing.B) {
	hello := "hello"
	world := "world"
	var builder strings.Builder
	builder.WriteString(hello)
	for i := 0; i < b.N; i++ {
		builder.WriteString(world)
	}
	_ = builder.String()
}
