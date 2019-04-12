// +build ignore

package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

func main() {
	TEXT("block", 0, "func(h *[4]uint32, m []byte)")
	Doc("block MD5 hashes the 64-byte message m into the running state h.")
	h := Mem{Base: Load(Param("h"), GP64())}
	m := Mem{Base: Load(Param("m").Base(), GP64())}

	// Store message values on the stack.
	w := AllocLocal(64)
	W := func(i int) Mem { return w.Offset((i % 16) * 4) }

	Comment("Load initial hash.")
	hash := [4]Register{GP32(), GP32(), GP32(), GP32()}
	for i, r := range hash {
		MOVL(h.Offset(4*i), r)
	}

	Comment("Initialize registers.")
	a, b, c, d := GP32(), GP32(), GP32(), GP32()
	for i, r := range []Register{a, b, c, d} {
		MOVL(hash[i], r)
	}

	var (
		f func(x, y, z Register) Register
		r = [64]uint8{
			7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22,
			5, 9, 14, 20, 5, 9, 14, 20, 5, 9, 14, 20, 5, 9, 14, 20,
			4, 11, 16, 23, 4, 11, 16, 23, 4, 11, 16, 23, 4, 11, 16, 23,
			6, 10, 15, 21, 6, 10, 15, 21, 6, 10, 15, 21, 6, 10, 15, 21,
		}
		k = [64]uint32{
			0xd76aa478, 0xe8c7b756, 0x242070db, 0xc1bdceee, 0xf57c0faf, 0x4787c62a, 0xa8304613, 0xfd469501, 0x698098d8, 0x8b44f7af, 0xffff5bb1, 0x895cd7be, 0x6b901122, 0xfd987193, 0xa679438e, 0x49b40821,
			0xf61e2562, 0xc040b340, 0x265e5a51, 0xe9b6c7aa, 0xd62f105d, 0x2441453, 0xd8a1e681, 0xe7d3fbc8, 0x21e1cde6, 0xc33707d6, 0xf4d50d87, 0x455a14ed, 0xa9e3e905, 0xfcefa3f8, 0x676f02d9, 0x8d2a4c8a,
			0xfffa3942, 0x8771f681, 0x6d9d6122, 0xfde5380c, 0xa4beea44, 0x4bdecfa9, 0xf6bb4b60, 0xbebfbc70, 0x289b7ec6, 0xeaa127fa, 0xd4ef3085, 0x4881d05, 0xd9d4d039, 0xe6db99e5, 0x1fa27cf8, 0xc4ac5665,
			0xf4292244, 0x432aff97, 0xab9423a7, 0xfc93a039, 0x655b59c3, 0x8f0ccc92, 0xffeff47d, 0x85845dd1, 0x6fa87e4f, 0xfe2ce6e0, 0xa3014314, 0x4e0811a1, 0xf7537e82, 0xbd3af235, 0x2ad7d2bb, 0xeb86d391,
		}
		not Register
	)
	for i := 0; i < 64; i++ {
		Commentf("Round %d. Operation %d.", i/16+1, i+1)
		u := GP32()

		if i < 16 {
			MOVL(m.Offset((i%16)*4), u)
			f = round1
		} else if i < 32 {
			MOVL(W(5*i+1), u)

			// precomputing NOT(z) for round 4
			not = GP32()
			MOVL(d, not)
			NOTL(not)

			f = func(x, y, z Register) Register {
				r := GP32()
				MOVL(x, r)
				ANDL(z, r)
				r2 := GP32()
				MOVL(not, r2)
				ANDL(y, r2)
				ORL(r2, r)
				return r
			}
		} else if i < 48 {
			MOVL(W(3*i+5), u)
			f = round3
		} else {
			MOVL(W(7*i), u)
			f = func(x, y, z Register) Register {
				r := GP32()
				MOVL(not, r)
				ORL(x, r)
				XORL(y, r)
				return r
			}
		}

		t := GP32()
		MOVL(a, t)
		ADDL(f(b, c, d), t)
		ADDL(U32(k[i]), t)
		ADDL(u, t)
		ROLL(U8(r[i]), t)
		ADDL(b, t)
		a, d, c, b = d, c, b, t
	}

	Comment("Final add.")
	for i, r := range []Register{a, b, c, d} {
		ADDL(r, hash[i])
	}

	Comment("Store results back.")
	for i, r := range hash {
		MOVL(r, h.Offset(4*i))
	}
	RET()

	Generate()
}

func round1(x, y, z Register) Register {
	r := GP32()
	MOVL(z, r)
	XORL(y, r)
	ANDL(x, r)
	XORL(z, r)
	return r
}

func round3(x, y, z Register) Register {
	r := GP32()
	MOVL(x, r)
	XORL(y, r)
	XORL(z, r)
	return r
}
