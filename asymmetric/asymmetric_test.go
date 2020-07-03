package asymmetric

import (
	"bytes"
	"fmt"
	"testing"

	"neckless.adviser.com/keys"
)

func TestDeEnCryptPubKey(t *testing.T) {
	alice, _ := keys.NewPrivateKey()
	bob, _ := keys.NewPrivateKey()
	buf := new(bytes.Buffer)
	fmt.Fprint(buf, "Hello World")
	aliceBob := CreateShared(&alice.Key.Raw, &bob.Public().Key.Raw)
	bobAlice := CreateShared(&bob.Key.Raw, &alice.Public().Key.Raw)
	// if bytes.Contains(crypt, buf.Bytes()) {
	// 	t.Error("no encrypted text", crypt, buf.Bytes())
	// }
	// mangled, err := DecryptWith(bob, alice.Public(), crypt)
	// if err != nil {
	// 	t.Error("we expect to be wirhout error")
	// }
	if bytes.Compare(aliceBob[:], bobAlice[:]) != 0 {
		t.Error("should be equal")
	}
	// if strings.Compare(string(mangled), "Hello World") != 0 {
	// 	t.Error("Mangled is not Hello world:", buf.String(), mangled)
	// }
}
