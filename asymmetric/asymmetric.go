package asymmetric

import (
	"golang.org/x/crypto/curve25519"
	"neckless.adviser.com/keys"
)

func CreateShared(priv *keys.RawKey, pub *keys.RawKey) keys.RawKey {
	shared := [32]byte{}
	curve25519.ScalarMult(&shared, priv.As32Byte(), pub.As32Byte())
	return shared
}

// func EncryptFor(keyAlice *keys.PrivateKey, pubKeyBob *keys.PublicKey, msg []byte) []byte {
// 	shared := CreateShared(keyAlice, pubKeyBob)
// 	aead, _ := chacha20poly1305.NewX(shared[:])
// 	nonce := make([]byte, chacha20poly1305.NonceSizeX)
// 	return aead.Seal(nil, nonce, []byte(msg), nil)
// }

// func DecryptWith(keyBob *keys.PrivateKey, pubKeyAlice *keys.PublicKey, msg []byte) ([]byte, error) {
// 	shared := CreateShared(keyBob, pubKeyAlice)
// 	aead, _ := chacha20poly1305.NewX(shared[:])
// 	nonce := make([]byte, chacha20poly1305.NonceSizeX)
// 	return aead.Open(nil, nonce, []byte(msg), nil)
// }
