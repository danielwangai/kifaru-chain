package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"io"
)

const (
	privKeyLen   = 64 // length of the private key
	seedLen      = 32
	pubKeyLen    = 32 // length of the public key
	addressLen   = 20
	signatureLen = 64
)

type PrivateKey struct {
	key ed25519.PrivateKey
}

func (p *PrivateKey) Bytes() []byte {
	return p.key
}

// Sign uses private key to sign
func (p *PrivateKey) Sign(msg []byte) *Signature {
	return &Signature{
		value: ed25519.Sign(p.key, msg),
	}
}

func GeneratePrivateKey() *PrivateKey {
	seed := make([]byte, seedLen)
	_, err := io.ReadFull(rand.Reader, seed)
	if err != nil {
		panic(err)
	}

	return &PrivateKey{
		key: ed25519.NewKeyFromSeed(seed),
	}
}

// PublicKey returns public key from the private key
func (p *PrivateKey) PublicKey() *PublicKey {
	b := make([]byte, pubKeyLen)
	copy(b, p.key[32:])

	return &PublicKey{
		key: b,
	}
}

type PublicKey struct {
	key ed25519.PublicKey
}

func (p *PublicKey) Bytes() []byte {
	return p.key
}

// Address returns address from public key
// address = last 20 characters on the public address
func (p *PublicKey) Address() Address {
	return Address{
		value: p.key[len(p.key)-addressLen:],
	}
}

type Signature struct {
	value []byte
}

func (s *Signature) Bytes() []byte {
	return s.value
}

func SignatureFromBytes(b []byte) *Signature {
	fmt.Println("Len: ", len(b))
	if len(b) != signatureLen {
		panic("invalid signature length, must be 64")
	}

	return &Signature{value: b}
}

// Verify checks if the signature is valid
func (s *Signature) Verify(pubKey *PublicKey, msg []byte) bool {
	return ed25519.Verify(pubKey.key, msg, s.value)
}
