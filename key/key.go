package key

import (
	"bytes"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/curve25519"
)

// KeySize is Default curve25519 key size
// const RawKeyLen = 32

type RawKey [32]byte

func AsRawKey(my []byte) *RawKey {
	ret := RawKey{}
	if len(my) < len(ret) {
		for i := range my {
			ret[i] = my[i]
		}
	} else {
		for i := range ret {
			ret[i] = my[i]
		}
	}
	return &ret
}

// Uncool to copie a key in local memory
func (my *RawKey) As32Byte() *[32]byte {
	out := [32]byte{}
	for i := range *my {
		out[i] = my[i]
	}
	return &out
}

type KeyStyle string

const (
	Private = KeyStyle("Private")
	Public  = KeyStyle("Public")
)

type KeyType struct {
	Id    string
	Style KeyStyle
	Raw   RawKey
}

type JsonKeyType struct {
	Id    string
	Style KeyStyle
	Raw   string
}

func toKey(i []byte) RawKey {
	ret := RawKey{}
	if len(i) < len(ret) {
		for k := range i {
			ret[k] = i[k]
		}
	} else {
		for k := range ret {
			ret[k] = i[k]
		}
	}
	return ret
}

func (k *KeyType) AsJson() *JsonKeyType {
	return &JsonKeyType{
		Id:    k.Id,
		Style: k.Style,
		Raw:   base64.StdEncoding.EncodeToString([]byte(k.Raw[:])),
	}
}

func ToKeyType(jsk JsonKeyType) (*KeyType, error) {
	raw, err := base64.StdEncoding.DecodeString(jsk.Raw)
	if err != nil {
		return nil, err
	}
	return &KeyType{
		Id:    jsk.Id,
		Style: jsk.Style,
		Raw:   toKey(raw),
	}, nil
}

type PublicKey struct {
	Key KeyType
}

type PrivateKey struct {
	Key KeyType
}

func MakePublicKey(key *RawKey, id string) *PublicKey {
	return &PublicKey{
		KeyType{
			Id:    id,
			Style: Public,
			Raw:   *key,
		},
	}
}
func MakePrivateKey(key *RawKey, ids ...string) *PrivateKey {
	var id string
	if len(ids) == 0 || (len(ids) > 0 && len(ids[0]) == 0) {
		id = uuid.New().String()
	} else {
		id = ids[0]
	}
	return &PrivateKey{
		KeyType{
			Id:    id,
			Style: Private,
			Raw:   *key,
		},
	}
}

// CreateRandomKey generates a new random key.
func CreateRandomKey() (*RawKey, error) {
	k := RawKey{}
	_, err := rand.Read(k[:])
	if err != nil {
		return nil, err
	}
	return &k, nil
}

// NewPrivateKey generates a new curve25519 secret key.
// It conforms to the format described on https://cr.yp.to/ecdh.html.
func NewPrivateKey(pk ...*PrivateKey) (*PrivateKey, error) {
	var k RawKey
	if len(pk) > 0 && pk[0] != nil &&
		pk[0].Key.Style == Private &&
		len(pk[0].Key.Raw) == len(RawKey{}) {
		k = pk[0].Key.Raw
	} else {
		m, err := CreateRandomKey()
		if err != nil {
			return &PrivateKey{}, err
		}
		k = *m
	}
	k[0] &= 248
	k[31] = (k[31] & 127) | 64
	return MakePrivateKey(&k), nil
}

func isZero(k RawKey) bool {
	zeros := RawKey{}
	return subtle.ConstantTimeCompare(zeros[:], k[:]) == 1
}

// Public computes the public key matching this curve25519 secret key.
func (k *PrivateKey) Public() *PublicKey {
	if isZero(k.Key.Raw) {
		panic("Tried to generate emptyPrivateKey.Public()")
	}
	pub := [len(RawKey{})]byte{}
	curve25519.ScalarBaseMult(&pub, k.Key.Raw.As32Byte())
	return MakePublicKey(AsRawKey(pub[:]), k.Key.Id)
}

// MarshalText create a String of the Private Key
func (k *PrivateKey) Marshal() string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, `privkey:%s`, k.Key.Raw.Marshal())
	return buf.String()
}

func (k *RawKey) Marshal() string {
	return base64.StdEncoding.EncodeToString(k[:])
}

func (k *PublicKey) Marshal() string {
	return k.Key.Raw.Marshal()
}

func fromText(pkstr string) (*RawKey, error) {
	ret := RawKey{}
	dc, err := base64.StdEncoding.DecodeString(pkstr)
	if err != nil {
		return nil, err
	}
	if len(dc) != len(RawKey{}) {
		return nil, errors.New("this is not a rawkey")
	}
	for i := range dc {
		ret[i] = dc[i]
	}
	return &ret, nil
}

func FromText(pkstr string, id string) (*PrivateKey, *PublicKey, error) {
	if strings.HasPrefix(pkstr, "privkey:") {
		keyBytes, err := fromText(strings.TrimPrefix(pkstr, "privkey:"))
		if err != nil {
			return nil, nil, err
		}
		return MakePrivateKey(keyBytes, id), nil, nil
	}
	keyBytes, err := fromText(pkstr)
	if err != nil {
		return nil, nil, err
	}
	return nil, MakePublicKey(keyBytes, id), nil
}
