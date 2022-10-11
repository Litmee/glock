package treaty

import (
	"bytes"
	"encoding/binary"
)

// Encode message encoding
func Encode(key byte, d byte) []byte {

	// define an empty bytes buffer
	b := new(bytes.Buffer)

	// write the message key
	_ = binary.Write(b, binary.LittleEndian, key)

	// write message entity
	_ = binary.Write(b, binary.LittleEndian, d)

	// Return
	return b.Bytes()
}
