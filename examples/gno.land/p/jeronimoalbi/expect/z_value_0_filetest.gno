package main

import (
	"strings"

	"gno.land/p/jeronimoalbi/expect"
)

var (
	output strings.Builder
	t      = expect.MockTestingT(&output)
)

func main() {
	expect.Value(t, "foo").ToEqual("bar")
	expect.Value(t, []byte("foo")).ToEqual([]byte("bar"))
	expect.Value(t, true).ToEqual(false)
	expect.Value(t, float32(1)).ToEqual(float32(2))
	expect.Value(t, float64(1.1)).ToEqual(float64(1.2))
	expect.Value(t, uint(1)).ToEqual(uint(2))
	expect.Value(t, uint8(1)).ToEqual(uint8(2))
	expect.Value(t, uint16(1)).ToEqual(uint16(2))
	expect.Value(t, uint32(1)).ToEqual(uint32(2))
	expect.Value(t, uint64(1)).ToEqual(uint64(2))
	expect.Value(t, int(1)).ToEqual(int(2))
	expect.Value(t, int8(1)).ToEqual(int8(2))
	expect.Value(t, int16(1)).ToEqual(int16(2))
	expect.Value(t, int32(1)).ToEqual(int32(2))
	expect.Value(t, int64(1)).ToEqual(int64(2))

	expect.Value(t, 0).ToEqual([]string{})

	println(output.String())
}

// Output:
// Expected values to match
// Got: foo
// Want: bar
// Expected values to match
// Got: foo
// Want: bar
// Expected values to match
// Got: true
// Want: false
// Expected values to match
// Got: 1
// Want: 2
// Expected values to match
// Got: 1.1
// Want: 1.2
// Expected values to match
// Got: 1
// Want: 2
// Expected values to match
// Got: 1
// Want: 2
// Expected values to match
// Got: 1
// Want: 2
// Expected values to match
// Got: 1
// Want: 2
// Expected values to match
// Got: 1
// Want: 2
// Expected values to match
// Got: 1
// Want: 2
// Expected values to match
// Got: 1
// Want: 2
// Expected values to match
// Got: 1
// Want: 2
// Expected values to match
// Got: 1
// Want: 2
// Expected values to match
// Got: 1
// Want: 2
// Unsupported type: unknown
