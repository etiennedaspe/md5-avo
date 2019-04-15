package md5

import (
	"encoding/binary"
)

const (
	// The size of an MD5 checksum in bytes.
	Size = 16

	// The block size of MD5 in bytes.
	BlockSize = 64
)

// Sum returns the MD5 checksum of data.
func Sum(data []byte) [Size]byte {
	var (
		n = len(data)
		h = [4]uint32{0x67452301, 0xefcdab89, 0x98badcfe, 0x10325476}
	)

	// Consume full blocks.
	for len(data) >= BlockSize {
		block(&h, data)
		data = data[BlockSize:]
	}

	// Final block.
	tmp := make([]byte, BlockSize)
	copy(tmp, data)
	tmp[len(data)] = 0x80

	if len(data) >= 56 {
		block(&h, tmp)
		for i := 0; i < BlockSize; i++ {
			tmp[i] = 0
		}
	}

	binary.LittleEndian.PutUint64(tmp[56:], uint64(8*n))
	block(&h, tmp)

	// Write into byte array.
	var digest [Size]byte
	for i := 0; i < 4; i++ {
		binary.LittleEndian.PutUint32(digest[4*i:], h[i])
	}

	return digest
}
