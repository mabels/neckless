package pipeline

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestDeEnCryptPubKey(t *testing.T) {
	alice, _ := NewPrivateKey()
	bob, _ := NewPrivateKey()
	buf := new(bytes.Buffer)
	fmt.Fprint(buf, "Hello World")
	crypt := EncryptFor(alice, bob.Public(), buf.Bytes())
	if bytes.Contains(crypt, buf.Bytes()) {
		t.Error("no encrypted text", crypt, buf.Bytes())
	}
	mangled, err := DecryptWith(bob, alice.Public(), crypt)
	if err != nil {
		t.Error("we expect to be wirhout error")
	}
	if bytes.Compare(buf.Bytes(), mangled) != 0 {
		t.Error("should be equal")
	}
	if strings.Compare(string(mangled), "Hello World") != 0 {
		t.Error("Mangled is not Hello world:", buf.String(), mangled)
	}
}
