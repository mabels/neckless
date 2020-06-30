package pipeline

import (
	"bytes"
	"crypto/subtle"
	"strings"
	"testing"
)

func TestMakePublicKey(t *testing.T) {
	pk := MakePublicKey([KeySize]byte{1})
	if pk.Key.Style != KeyStyle(Public) {
		t.Error("Has to be Public Style")
	}
	if pk.Key.Raw[0] != 1 {
		t.Error("No the Right Raw Value")
	}
}

func TestMakePrivateKey(t *testing.T) {
	pk := MakePrivateKey([KeySize]byte{1})
	if pk.Key.Style != KeyStyle(Private) {
		t.Error("Has to be Private Style")
	}
	if pk.Key.Raw[0] != 1 {
		t.Error("No the Right Raw Value")
	}
}

func TestNewPrivateKey(t *testing.T) {
	pk, err := NewPrivateKey()
	if err != nil {
		t.Error("NewPrivateKey failed", err)
	}
	if pk.Key.Style != KeyStyle(Private) {
		t.Error("Has to be Private Style")
	}
	zeros := [KeySize]byte{}
	if subtle.ConstantTimeCompare(zeros[:], pk.Key.Raw[:]) == 1 {
		t.Error("No the Right Raw Value")
	}
}

func TestPublicKey(t *testing.T) {
	pk, err := NewPrivateKey()
	if err != nil {
		t.Error("NewPrivateKey failed", err)
	}
	pubkey := pk.Public()
	if pubkey.Key.Style != KeyStyle(Public) {
		t.Error("Has to be Private Style")
	}
	zeros := [KeySize]byte{}
	if subtle.ConstantTimeCompare(zeros[:], pubkey.Key.Raw[:]) == 1 {
		t.Error("No the Right Raw Value")
	}
}

func TestMarshalPrivate(t *testing.T) {
	pk, err := NewPrivateKey()
	if err != nil {
		t.Error("NewPrivateKey failed", err)
	}
	pkstr := pk.Marshal()
	if !strings.HasPrefix(pkstr, "privkey:") {
		t.Error("Private wrong prefix")
	}
}

func TestMarshalPublic(t *testing.T) {
	pk, err := NewPrivateKey()
	if err != nil {
		t.Error("NewPrivateKey failed", err)
	}
	pkstr := pk.Public().Marshal()
	if strings.HasPrefix(pkstr, "privkey:") {
		t.Error("Public wrong prefix")
	}
}

func TestFromTextKaputtEmpty(t *testing.T) {
	_, _, err := FromText("")
	if err == nil {
		t.Error("expecting error")
	}
}
func TestFromTextKaputtPublicIllegal(t *testing.T) {
	_, _, err := FromText("blabla")
	if err == nil {
		t.Error("expecting error")
	}
}

func TestFromTextKaputtPublicLong(t *testing.T) {
	val := [KeySize * 2]byte{}
	hex := []byte{'0', '1', '2', '3',
		'4', '5', '6', '7',
		'8', '9', 'a', 'b',
		'c', 'd', 'e', 'f'}
	for i := range val[:] {
		val[i] = hex[i%len(hex)]
	}
	val[29] = 'k'
	// t.Error(string(val[:]))
	_, _, err := FromText(string(val[:]))
	if err == nil {
		t.Error("expecting error")
	}
}

func TestFromPrivateSerialize(t *testing.T) {
	inpk, _ := NewPrivateKey()
	pk, pb, err := FromText(inpk.Marshal())
	if err != nil {
		t.Error("expecting error")
	}
	if pb != nil {
		t.Error("expecting error")
	}
	if bytes.Compare(pk.Key.Raw[:], inpk.Key.Raw[:]) != 0 {
		t.Error("expecting error")
	}
}

func TestFromPublicSerialize(t *testing.T) {
	inpk, _ := NewPrivateKey()
	pk, pb, err := FromText(inpk.Public().Marshal())
	if err != nil {
		t.Error("expecting error")
	}
	if pk != nil {
		t.Error("expecting error")
	}
	if bytes.Compare(pb.Key.Raw[:], inpk.Public().Key.Raw[:]) != 0 {
		t.Error("expecting error")
	}
}
