package asymmetric

import (
	"golang.org/x/crypto/curve25519"
	"github.com/mabels/neckless/key"
)

func CreateShared(priv *key.RawKey, pub *key.RawKey) key.RawKey {
	shared := [32]byte{}
	curve25519.ScalarMult(&shared, priv.As32Byte(), pub.As32Byte())
	return shared
}
