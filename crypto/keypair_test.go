package crypto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGeneratePrivateKey(t *testing.T) {
	privKey := GeneratePrivateKey()
	assert.Equal(t, privKeyLen, len(privKey.Bytes()))
	pubKey := privKey.PublicKey()
	assert.Equal(t, pubKeyLen, len(pubKey.Bytes()))

	addr := pubKey.Address()
	assert.Equal(t, addressLen, len(addr.Bytes()))
}
