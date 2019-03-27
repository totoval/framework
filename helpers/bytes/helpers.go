package bytes

import (
	"encoding/binary"
)

func FromUint64(input uint64) []byte {

	theByte := make([]byte, 8)
	binary.BigEndian.PutUint64(theByte, uint64(input))

	return theByte
}
