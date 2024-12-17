package crypto

import "encoding/hex"

type Address struct {
	value []byte
}

func (a *Address) Bytes() []byte {
	return a.value
}

// String converts Address to string
func (a *Address) String() string {
	return hex.EncodeToString(a.value)
}
