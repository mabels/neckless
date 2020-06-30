package pipeline

import (
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/curve25519"
)

func EncryptFor(keyAlice *PrivateKey, pubKeyBob *PublicKey, msg []byte) []byte {
	prA := keyAlice.Key.Raw
	puB := pubKeyBob.Key.Raw
	var shared [32]byte
	curve25519.ScalarMult(&shared, &prA, &puB)
	aead, _ := chacha20poly1305.NewX(shared[:])
	nonce := make([]byte, chacha20poly1305.NonceSizeX)
	return aead.Seal(nil, nonce, []byte(msg), nil)
}

func DecryptWith(keyBob *PrivateKey, pubKeyAlice *PublicKey, msg []byte) ([]byte, error) {
	prB := keyBob.Key.Raw
	puA := pubKeyAlice.Key.Raw
	var shared [32]byte
	curve25519.ScalarMult(&shared, &prB, &puA)
	aead, _ := chacha20poly1305.NewX(shared[:])
	nonce := make([]byte, chacha20poly1305.NonceSizeX)
	return aead.Open(nil, nonce, []byte(msg), nil)
}
