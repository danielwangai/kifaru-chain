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

func TestSignature(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.PublicKey()

	// sign
	msg := []byte("hello world")
	sig := privKey.Sign(msg)
	assert.Equal(t, signatureLen, len(sig.Bytes()))

	// verify signature
	isValid := sig.Verify(pubKey, msg)
	assert.True(t, isValid)

	// verify with different keypair
	privKey2 := GeneratePrivateKey()
	pubKey2 := privKey2.PublicKey()
	sig2 := privKey2.Sign(msg)
	assert.False(t, sig2.Verify(pubKey, msg))
	assert.False(t, sig2.Verify(pubKey2, []byte("wrong message")))
}
