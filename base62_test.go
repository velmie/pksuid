package pksuid

import "testing"

type isBase62BytesTest struct {
	in  []byte
	out bool
}

var isBase62BytesTests = []isBase62BytesTest{
	{
		in:  []byte{},
		out: false,
	},
	{
		in:  []byte{0},
		out: false,
	},
	{
		in:  []byte("hello!"),
		out: false,
	},
	{
		in:  []byte("hello "),
		out: false,
	},
	{
		in:  []byte("foo:bar"),
		out: false,
	},
	{
		in:  []byte("hello"),
		out: true,
	},
	{
		in:  []byte("20MgJhy7bR6mHUfKROlb6RfYrBk"),
		out: true,
	},
}

func TestIsBase62Bytes(t *testing.T) {
	for i, test := range isBase62BytesTests {
		got := isBase62Bytes(test.in)
		if got != test.out {
			t.Errorf("#%d isBase62Bytes(%q) returned %t, want %t", i, test.in, got, test.out)
		}
	}
}
