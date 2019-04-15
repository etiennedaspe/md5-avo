package md5

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"testing"
	"testing/quick"
)

//go:generate go run asm.go -out md5.s -stubs stub.go

func TestVectors(t *testing.T) {
	cases := []struct {
		Data      string
		HexDigest string
	}{
		{"", "d41d8cd98f00b204e9800998ecf8427e"},
		{"The quick brown fox jumps over the lazy dog", "9e107d9d372bb6826bd81d3542a419d6"},
		{"The quick brown fox jumps over the lazy cog", "1055d3e698d289f2af8663725127bd4b"},
	}
	for _, c := range cases {
		digest := Sum([]byte(c.Data))
		got := hex.EncodeToString(digest[:])
		if got != c.HexDigest {
			t.Errorf("Sum(%#v) = %s; expect %s", c.Data, got, c.HexDigest)
		}
	}
}

func TestCmp(t *testing.T) {
	if err := quick.CheckEqual(Sum, md5.Sum, nil); err != nil {
		t.Fatal(err)
	}
}

func TestLengths(t *testing.T) {
	data := make([]byte, BlockSize)
	for n := 0; n <= BlockSize; n++ {
		got := Sum(data[:n])
		expect := md5.Sum(data[:n])
		if !bytes.Equal(got[:], expect[:]) {
			t.Errorf("failed on length %d", n)
		}
	}
}
