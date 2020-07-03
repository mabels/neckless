package symmetric

import (
	"bytes"
	"encoding/base64"
	"errors"

	"crypto/cipher"
	"crypto/sha256"

	"golang.org/x/crypto/chacha20poly1305"
	"neckless.adviser.com/keys"
)

type KeyAndNonce struct {
	Key   []byte
	Nouce []byte
}

type JsonKeyAndNonce struct {
	Key   string
	Nouce string
}

func (my *KeyAndNonce) AsJson() *JsonKeyAndNonce {
	return &JsonKeyAndNonce{
		Key:   base64.StdEncoding.EncodeToString(my.Key),
		Nouce: base64.StdEncoding.EncodeToString(my.Nouce),
	}
}

func (my *JsonKeyAndNonce) AsKeyAndNonce() (*KeyAndNonce, error) {
	key, err := base64.StdEncoding.DecodeString(my.Key)
	if err != nil {
		return nil, err
	}
	nonce, err := base64.StdEncoding.DecodeString(my.Nouce)
	if err != nil {
		return nil, err
	}
	return &KeyAndNonce{
		Key:   key,
		Nouce: nonce,
	}, nil
}

type SealedContainer struct {
	Checksum []byte
	Payload  []byte
}

func keyFromSeed(seed [][]byte) ([]byte, []byte) {
	randkey, _ := keys.CreateRandomKey()
	randnonce, _ := keys.CreateRandomKey()
	key := sha256.Sum256(bytes.Join(append(seed, randkey[:]), []byte{}))
	nonce := sha256.Sum256(bytes.Join(append(seed, randnonce[:]), []byte{}))
	return key[:], nonce[:]
}

//	// key, nouce := keyFromSeed(seed)

type SealRequest struct {
	Key      keys.RawKey
	Payload  []byte
	Checksum []byte
}

func Checksum(sr *SealRequest) *SealRequest {
	data := append(sr.Payload, sr.Key[:]...)
	my := sha256.Sum256(data)
	sr.Checksum = my[:]
	// fmt.Printf("Seal:%x=>%x\n", sr.Checksum, data)
	return sr
}

func nonce(aead cipher.AEAD, n []byte) []byte {
	csum := make([]byte, aead.NonceSize())
	if len(n) < aead.NonceSize() {
		copy(csum, n)
	} else {
		copy(csum, n[:aead.NonceSize()])
	}
	return csum
}

func Seal(sr *SealRequest) (*SealedContainer, error) {
	aead, err := chacha20poly1305.NewX(sr.Key[:])
	if err != nil {
		return nil, err
	}
	nonce := nonce(aead, sr.Checksum)
	sealed := aead.Seal(nil, nonce, sr.Payload, nil)
	// fmt.Printf("Seal:%d:%x=>%x:%x\n", len(nonce), nonce, sr.Key, sealed)
	return &SealedContainer{
		Checksum: sr.Checksum,
		Payload:  sealed,
	}, nil
}

type OpenContainer struct {
	Checksum []byte
	Payload  []byte
}

func Verify(csum []byte, key *keys.RawKey, open *[]byte) bool {
	data := append(*open, key[:]...)
	// fmt.Printf("Verify:%x=>%x", csum, data)
	tmp := sha256.Sum256(data)
	return bytes.Equal(tmp[:], csum)
}

func SkipVerify(csum []byte, key *keys.RawKey, open *[]byte) bool {
	return true
}

func Open(key *keys.RawKey, sc *SealedContainer, verify func([]byte, *keys.RawKey, *[]byte) bool) (*OpenContainer, error) {
	aead, err := chacha20poly1305.NewX(key[:])
	if err != nil {
		return nil, err
	}
	nonce := nonce(aead, sc.Checksum)
	// fmt.Printf("Open:%d:%x=>%x:%x\n", len(nonce), nonce, key, sc.Payload)
	open, err := aead.Open(nil, nonce, sc.Payload, nil)
	// fmt.Printf("Open:%x:%x\n", key[:], sc.Checksum[:aead.NonceSize()])
	if err != nil {
		return nil, err
	}
	if !verify(sc.Checksum, key, &open) {
		return nil, errors.New("checksum error")
	}
	return &OpenContainer{
		Checksum: sc.Checksum,
		Payload:  open,
	}, nil
}

type Base64SealContainer struct {
	KeyAndNonce
	Checksum []byte
	Payload  string
}

func OpenBase64(key *keys.RawKey, pp *Base64SealContainer, verify func([]byte, *keys.RawKey, *[]byte) bool) (*OpenContainer, error) {
	plain, err := base64.RawStdEncoding.DecodeString(pp.Payload)
	if err != nil {
		return nil, err
	}
	return Open(key, &SealedContainer{
		Checksum: pp.Checksum,
		Payload:  plain,
	}, verify)
}
