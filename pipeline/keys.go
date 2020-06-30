package pipeline

import (
	"bytes"
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/curve25519"
)

// KeySize is Default curve25519 key size
const KeySize = 32

type KeyStyle string

const (
	Private = KeyStyle("Private")
	Public  = KeyStyle("Public")
)

type KeyType struct {
	Style KeyStyle
	Raw   [KeySize]byte
}

type PublicKey struct {
	Key KeyType
}

type PrivateKey struct {
	Key KeyType
}

func MakePublicKey(key [KeySize]byte) *PublicKey {
	return &PublicKey{
		KeyType{
			Style: Public,
			Raw:   key,
		},
	}
}
func MakePrivateKey(key [KeySize]byte) *PrivateKey {
	return &PrivateKey{
		KeyType{
			Style: Private,
			Raw:   key,
		},
	}
}

// NewPresharedKey generates a new random key.
func newPresharedKey() ([KeySize]byte, error) {
	var k [KeySize]byte
	_, err := rand.Read(k[:])
	if err != nil {
		return k, err
	}
	return k, nil
}

// NewPrivateKey generates a new curve25519 secret key.
// It conforms to the format described on https://cr.yp.to/ecdh.html.
func NewPrivateKey(pk ...*PrivateKey) (*PrivateKey, error) {
	var k [KeySize]byte
	if len(pk) > 0 && pk[0] != nil &&
		pk[0].Key.Style == Private &&
		len(pk[0].Key.Raw) == KeySize {
		k = pk[0].Key.Raw
	} else {
		m, err := newPresharedKey()
		if err != nil {
			return &PrivateKey{}, err
		}
		k = m
	}
	k[0] &= 248
	k[31] = (k[31] & 127) | 64
	return MakePrivateKey(k), nil
}

func isZero(k [KeySize]byte) bool {
	var zeros [KeySize]byte
	return subtle.ConstantTimeCompare(zeros[:], k[:]) == 1
}

// Public computes the public key matching this curve25519 secret key.
func (k *PrivateKey) Public() *PublicKey {
	if isZero(k.Key.Raw) {
		panic("Tried to generate emptyPrivateKey.Public()")
	}
	var pub [KeySize]byte
	curve25519.ScalarBaseMult(&pub, &k.Key.Raw)
	return MakePublicKey(pub)
}

// MarshalText create a String of the Private Key
func (k *PrivateKey) Marshal() string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, `privkey:%x`, k.Key.Raw)
	return buf.String()
}

func (k *PublicKey) Marshal() string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, `%x`, k.Key.Raw)
	return buf.String()
}

func fromText(pkstr string) ([KeySize]byte, error) {
	ret := [KeySize]byte{}
	if len(pkstr) < KeySize*2 {
		return ret, errors.New("To Short")
	}
	for i := range ret[:] {
		val, err := strconv.ParseUint(string(pkstr[(i*2):(i*2)+2]), 16, 8)
		if err != nil {
			return ret, err
		}
		ret[i] = byte(val)

	}
	return ret, nil
}

func FromText(pkstr string) (*PrivateKey, *PublicKey, error) {
	if strings.HasPrefix(pkstr, "privkey:") {
		keyBytes, err := fromText(strings.TrimPrefix(pkstr, "privkey:"))
		if err != nil {
			return nil, nil, err
		}
		return MakePrivateKey(keyBytes), nil, nil
	}
	keyBytes, err := fromText(pkstr)
	if err != nil {
		return nil, nil, err
	}
	return nil, MakePublicKey(keyBytes), nil
}
