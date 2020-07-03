package member

import (
	"strings"
	"testing"
)

func TestSignedJWT(t *testing.T) {
	pm, _ := MakePrivateMember(&PrivateMemberArg{
		Member: MemberArg{
			Type: "Person",
			Name: "Test",
		},
	})
	jwt, err := MakePublicMemberJWT(&pm.PrivateKey, pm.Public())
	if err != nil {
		t.Error("We don't expect an error", err)
	}
	// t.Error(jwt)
	claim, _, err := VerifyJWT(&pm.PrivateKey, jwt)
	if err != nil {
		t.Error("We don't expect an error", err)
	}
	if strings.Compare(claim.PublicKey, pm.Public().PublicKey.Marshal()) != 0 {
		t.Error("public key missmatch", claim.PublicKey)
	}
}
