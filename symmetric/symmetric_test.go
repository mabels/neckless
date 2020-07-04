package symmetric

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"testing"

	"neckless.adviser.com/key"
)

func TestSealOpen(t *testing.T) {
	val := []byte("So was von Geheim")
	key, _ := key.CreateRandomKey()
	sc, err := Seal(Checksum(&SealRequest{
		Key:     *key,
		Payload: val,
	}))
	if err != nil {
		t.Error("is not expected")
	}
	op, err := Open(key, sc, Verify)
	if err != nil {
		t.Error("is not expected", err)
	}
	if !bytes.Equal(val, op.Payload) {
		t.Error("no equal")
	}
	tmp := sha256.Sum256(append(val, key[:]...))
	if !bytes.Equal(tmp[:], op.Checksum) {
		t.Error("no equal")
	}
	op, err = OpenBase64(key, &Base64SealContainer{
		Checksum: sc.Checksum,
		Payload:  base64.StdEncoding.EncodeToString(sc.Payload),
	}, Verify)
	if err != nil {
		t.Error("is not expected")
	}
	if !bytes.Equal(val, op.Payload) {
		t.Error("no equal")
	}
}

func TestExternChecksum(t *testing.T) {
	val := []byte("So was von Geheim")
	testkey, _ := key.CreateRandomKey()
	csum := []byte("So was von Checksum und 24Byte lang")
	sc, err := Seal(&SealRequest{
		Key:      *testkey,
		Payload:  val,
		Checksum: csum,
	})
	if err != nil {
		t.Error("is not expected")
	}
	op, err := Open(testkey, sc, func(csum1 []byte, key1 *key.RawKey, open *[]byte) bool {
		return bytes.Equal(csum1, csum) &&
			bytes.Equal(key1[:], testkey[:]) &&
			bytes.Equal(*open, val)
	})
	if err != nil {
		t.Error("is not expected", err)
	}
	if !bytes.Equal(val, op.Payload) {
		t.Error("no equal data")
	}
	if !bytes.Equal(csum, op.Checksum) {
		t.Error("no equal csum")
	}
}
