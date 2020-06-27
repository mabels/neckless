package pipeline

import (
	"bytes"
	"crypto/rand"
	"crypto/subtle"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/curve25519"
)

// PipelineArgs Command Args
type PipelineArgs struct {
	name string
	// valid period.Period
}

// KeySize is Default curve25519 key size
const KeySize = 32

type PublicKey struct {
}
type PrivateKey struct {
	// PrivateKey is curve25519 key.
	raw [KeySize]byte
	// type PrivateKey [KeySize]byte
}

// KeyPair is the pub and privkey container
// KeyPair is the pub and privkey container
type KeyPair struct {
	Priv PublicKey
	Publ PrivateKey
}

func (*KeyPair) PublAsKey() []byte {
	return []byte{}
}

func (*KeyPair) PrivAsKey() []byte {
	return []byte{}
}

// Pipeline is The Created Entry for the PipeLine
type Pipeline struct {
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	KeyPair    KeyPair   `json:"keyPair"`
	ValidUntil time.Time `json:"valid"`
	Updated    time.Time `json:"updated"`
	Created    time.Time `json:"created"`
}

// Key is curve25519 key.
type Key [KeySize]byte

// NewPresharedKey generates a new random key.
func NewPresharedKey() (*Key, error) {
	var k [KeySize]byte
	_, err := rand.Read(k[:])
	if err != nil {
		return nil, err
	}
	return (*Key)(&k), nil
}

// NewPrivateKey generates a new curve25519 secret key.
// It conforms to the format described on https://cr.yp.to/ecdh.html.
func NewPrivateKey() (PrivateKey, error) {
	k, err := NewPresharedKey()
	if err != nil {
		return PrivateKey{}, err
	}
	k[0] &= 248
	k[31] = (k[31] & 127) | 64
	return (PrivateKey)(*k), nil
}

// MarshalText create a String of the Private Key
func (k PrivateKey) MarshalText() ([]byte, error) {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, `privkey:%x`, k[:])
	return buf.Bytes(), nil
}

func (k *Key) MarshalText() ([]byte, error) {
	if k == nil {
		return []byte("null"), nil
	}
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%x", k[:])
	return buf.Bytes(), nil
}

func (k *Key) IsZero() bool {
	if k == nil {
		return true
	}
	var zeros Key
	return subtle.ConstantTimeCompare(zeros[:], k[:]) == 1
}

func (k *PrivateKey) IsZero() bool {
	pk := Key(*k)
	return pk.IsZero()
}

// Public computes the public key matching this curve25519 secret key.
func (k *PrivateKey) Public() Key {
	pk := Key(*k)
	if pk.IsZero() {
		panic("Tried to generate emptyPrivateKey.Public()")
	}
	var p [KeySize]byte
	curve25519.ScalarBaseMult(&p, (*[KeySize]byte)(k))
	return (Key)(p)
}

// Create is used to Create a Pipeline
func Create(arg PipelineArgs) *Pipeline {
	now := time.Now()
	pk, err := NewPrivateKey()
	if err != nil {
		log.Fatal("can not create new private key")
	}
	pkTxt, err := pk.MarshalText()
	if err != nil {
		log.Fatal("can not marshal private to text")
	}
	pubKey := pk.Public()
	pbTxt, err := pubKey.MarshalText()
	if err != nil {
		log.Fatal("can not marshal public to text")
	}
	return &Pipeline{
		Id:   uuid.New().String(),
		Name: arg.name,
		KeyPair: KeyPair{
			Priv: pkTxt,
			Publ: pbTxt,
		},
		ValidUntil: now.AddDate(5, 0, 0),
		Updated:    now,
		Created:    now,
	}
}

type MemberClaim struct {
	SignerName   string `json:"signerName"`
	SignerPubkey string `json:"signerPubKey"`
	jwt.StandardClaims
}

func SignMember(signer *Pipeline, member *Pipeline) string {
	claims := &MemberClaim{
		SignerName:   member.Name,
		SignerPubkey: string(member.KeyPair.Publ),
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: signer.ValidUntil.Unix(),
			Id:        uuid.New().String(),
			IssuedAt:  time.Now().Unix(),
			Subject:   "SignerClaim",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(signer.KeyPair.Priv)
	if err != nil {
		log.Fatal("Siging failed")
	}
	return tokenString
}

func VerifyAndClaim(tknStr string, pl *Pipeline) (*MemberClaim, *jwt.Token, error) {
	claims := MemberClaim{}
	token, err := jwt.ParseWithClaims(tknStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return pl.KeyPair.Priv, nil
	})
	return &claims, token, err
}

func toByte32(p []byte) [32]byte {
	var ret [32]byte
	for i := range ret[:] {
		ret[i] = p[i]
	}
	return ret
}

func EncryptFor(t *testing.T, keyAlice *KeyPair, pubKeyBob []byte, msg string) []byte {
	prA := toByte32(keyAlice.PrivAsKey())
	puB := toByte32(pubKeyBob)
	var shared [32]byte
	curve25519.ScalarMult(&shared, &prA, &puB)
	t.Error(shared)
	aead, _ := chacha20poly1305.NewX(shared[:])
	nonce := make([]byte, chacha20poly1305.NonceSizeX)
	return aead.Seal(nil, nonce, []byte(msg), nil)
}

func DecryptFor(t *testing.T, keyBob *KeyPair, pubKeyAlice []byte, msg []byte) ([]byte, error) {
	prB := toByte32(keyBob.PrivAsKey())
	puA := toByte32(pubKeyAlice)
	var shared [32]byte
	curve25519.ScalarMult(&shared, &prB, &puA)
	t.Error(shared)
	aead, _ := chacha20poly1305.NewX(shared[:])
	nonce := make([]byte, chacha20poly1305.NonceSizeX)
	return aead.Open(nil, nonce, []byte(msg), nil)
}
