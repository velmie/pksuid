package pksuid

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestPrefixString(t *testing.T) {
	const expected = "hello"
	prefix := Prefix{'h', 'e', 'l', 'l', 'o', 0, '!'}
	if prefix.String() != expected {
		t.Errorf("prefix.String() expected %q, returned %q", expected, prefix.String())
	}
}

func TestNew(t *testing.T) {
	id := New(NilPrefix)
	id2 := New(NilPrefix)
	if id.Prefix() != NilPrefix {
		t.Errorf("expected prefix to be nil, got: %s", id.Prefix())
	}
	if id == id2 {
		t.Error("New(...) must generate random values")
	}
}

type parseTest struct {
	in          string
	errContains string
}

var parseTests = []parseTest{
	{
		in:          "",
		errContains: "less than",
	},
	{
		in:          strings.Repeat("1", pksuidByteLength),
		errContains: "failed to parse",
	},
	{
		in:          strings.Repeat("1", pksuidByteLength+1),
		errContains: "greater than",
	},
	{
		in: "208bayYfCoFfqyLD4lZaZ8BvwaF",
	},
	{
		in: "prefix:208bayYfCoFfqyLD4lZaZ8BvwaF",
	},
	{
		in:          fmt.Sprintf("%s%s", strings.Repeat("a", prefixByteLength+1), "208bayYfCoFfqyLD4lZaZ8BvwaF"),
		errContains: "prefix must be less than",
	},
}

func TestParse(t *testing.T) {
	for i, test := range parseTests {
		id, err := Parse(test.in)
		if err != nil && test.errContains == "" {
			t.Errorf("test %d: Unexpected error: %v", i, err)
			continue
		}
		if err != nil && !strings.Contains(err.Error(), test.errContains) {
			t.Errorf(
				"test %d: Parse(%q) returned error %q, want something containing %q",
				i,
				test.in,
				err.Error(),
				test.errContains,
			)
			continue
		}
		if err == nil && id.String() != test.in {
			t.Errorf("test %d: parsed value %q does not match the input %q", i, id.String(), test.in)
		}
	}
}

func TestFromBytes(t *testing.T) {
	short := make([]byte, ksuidByteLength-1)
	if _, err := FromBytes(short); err != errMinSize {
		t.Errorf("expected error %q, got %v", errMinSize, err)
	}
	long := make([]byte, pksuidByteLength+1)
	if _, err := FromBytes(long); err != errMaxSize {
		t.Errorf("expected error %q, got %v", errMaxSize, err)
	}
	if _, err := FromBytes(make([]byte, pksuidByteLength)); err != nil {
		t.Errorf("got unexpected error: %q", err)
	}

}

func TestSetPrefix(t *testing.T) {
	p := Prefix{'k', 'e', 'y', '_', 'l', 'i', 'v', 'e'}
	id := PKSUID{}
	id.SetPrefix(p)
	if id.Prefix() != p {
		t.Errorf("expected prefix %q, got %q", p, id.Prefix())
	}
}

func TestNil(t *testing.T) {
	if !Nil.IsNil() {
		t.Fatal("Nil must be Nil!")
	}
	x, _ := FromBytes(make([]byte, pksuidByteLength))
	if !x.IsNil() {
		t.Fatal("zero-byte array must be Nil!")
	}
	b := make([]byte, pksuidByteLength)
	copy(b[:prefixByteLength], strings.Repeat("a", prefixByteLength))

	w, _ := FromBytes(b)
	if !w.IsNil() {
		t.Fatal("non-zero prefix with following zero bytes must be Nil!")
	}
}

func TestMarshalAndUnmarhsalText(t *testing.T) {
	prefix := Prefix{'f', 'o', 'o'}
	var (
		id1 = New(prefix)
		id2 PKSUID
	)
	if err := id2.UnmarshalText([]byte(id1.String())); err != nil {
		t.Fatal(err)
	}
	if id1 != id2 {
		t.Fatal(id1, "!=", id2)
	}
	if b, err := id2.MarshalText(); err != nil {
		t.Fatal(err)
	} else if s := string(b); s != id1.String() {
		t.Fatal(s)
	}
}

func TestMarshalAndUnmarshalBinary(t *testing.T) {
	prefix := Prefix{'b', 'a', 'r'}
	var (
		id1 = New(prefix)
		id2 PKSUID
	)
	if err := id2.UnmarshalBinary(id1.Bytes()); err != nil {
		t.Fatal(err)
	}
	if id1 != id2 {
		t.Fatal(id1, "!=", id2)
	}
	if b, err := id2.MarshalBinary(); err != nil {
		t.Fatal(err)
	} else if bytes.Compare(b, id1.Bytes()) != 0 {
		t.Fatal("bad binary form:", id2)
	}
}

const encodedWithPrefix = "key_sandbox:208baxy3JWtg6ZyJmGpo2RoEpTX"

func TestSqlValuer(t *testing.T) {
	id, _ := Parse(encodedWithPrefix)

	if v, err := id.Value(); err != nil {
		t.Error(err)
	} else if s, ok := v.(string); !ok {
		t.Error("not a string value")
	} else if s != encodedWithPrefix {
		t.Error("bad string value::", s)
	}
}

func TestSqlScanner(t *testing.T) {
	prefix := Prefix{'b', 'a', 'z'}
	id1 := New(prefix)
	id2 := New(prefix)

	tests := []struct {
		pksuid PKSUID
		value  interface{}
	}{
		{Nil, nil},
		{id1, id1.String()},
		{id2, id2.Bytes()},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%T", test.value), func(t *testing.T) {
			var id PKSUID
			if err := id.Scan(test.value); err != nil {
				t.Error(err)
			}
			if id != test.pksuid {
				t.Error("bad PKSUID:")
				t.Logf("expected %v", test.pksuid)
				t.Logf("found    %v", id)
			}
		})
	}
}

func TestSqlValuerNilValue(t *testing.T) {
	if v, err := Nil.Value(); err != nil {
		t.Error(err)
	} else if v != nil {
		t.Errorf("bad nil value: %v", v)
	}
}

func TestDelegated(t *testing.T) {
	id := New(NilPrefix)
	ksuid := id.KSUID()

	if id.Time() != ksuid.Time() {
		t.Errorf("expected id.Time() to return %q, got %q", ksuid.Time(), id.Time())
	}

	if id.Timestamp() != ksuid.Timestamp() {
		t.Errorf("expected id.Timestamp() to return %q, got %q", ksuid.Time(), id.Time())
	}

	if !bytes.Equal(id.Payload(), ksuid.Payload()) {
		t.Error("expected id.Payload() to be equal to ksuid.Payload()")
	}

	if !bytes.Equal(id.IDBytes(), ksuid.Bytes()) {
		t.Error("expected id.IDBytes() to be equal to ksuid.Bytes()")
	}

	if id.ID() != ksuid.String() {
		t.Errorf("expected id.ID() to return %q, got %q", ksuid.String(), id.ID())
	}
}

func BenchmarkString(b *testing.B) {
	id := New(Prefix{'u', 's', 'e', 'r', '_'})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = id.String()
	}
}

type isNilBenchmark struct {
	in   PKSUID
	name string
}

var isNilBenchMarks = []isNilBenchmark{
	{
		in:   PKSUID{},
		name: "nil without prefix",
	},
	{
		in:   PKSUID{'p', 'r', 'e', 'f', 'i', 'x', ':'},
		name: "nil with prefix",
	},
	{
		in:   New(Prefix{'a', 'n', 'y', ':'}),
		name: "not nil with prefix",
	},
}

func BenchmarkIsNil(b *testing.B) {
	for _, tc := range isNilBenchMarks {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_ = tc.in.IsNil()
			}
		})
	}
}
