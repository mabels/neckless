package pearl

import (
	"encoding/json"
	"strings"
	"testing"

	"neckless.adviser.com/keys"
	"neckless.adviser.com/member"
)

func TestPearl(t *testing.T) {
	o1, _ := member.MakePrivateMember(&member.PrivateMemberArg{
		Member: member.MemberArg{
			Type: member.Person,
			Name: "o1",
		},
	})
	o2, _ := member.MakePrivateMember(&member.PrivateMemberArg{
		Member: member.MemberArg{
			Type: member.Person,
			Name: "o2",
		},
	})
	owners := []keys.PublicKey{}
	owners = append(owners, o1.Public().PublicKey)
	owners = append(owners, o2.Public().PublicKey)
	// fmt.Printf("%x:%x\n", o1.Public().PublicKey.Key.Raw, o2.Public().PublicKey.Key.Raw)
	//signer.Key.Raw, pub.Key.Raw, privPubKey)

	pearl, err := Close(&CloseRequestPearl{
		Type:    "test",
		Payload: []byte("Meno"),
		Signer:  o1.PrivateKey,
		Owners:  owners,
	})
	if err != nil {
		t.Error("unexpected errro:", err)
	}
	// data, err := json.Marshal(pearl)
	// t.Error(string(data))
	opc1, err := Open(&o1.PrivateKey, pearl)
	if err != nil {
		t.Error("unexpected errro:", err)
	}
	// data, err = json.Marshal(opc1)
	if strings.Compare("Meno", string(opc1.Payload)) != 0 {
		t.Error("not expected string:", string(opc1.Payload))
	}
	// t.Error(string(opc1.Payload))

	opc2, err := Open(&o2.PrivateKey, pearl)
	if err != nil {
		t.Error("unexpected errro:", err)
	}
	// data, err = json.Marshal(opc2)
	// if err != nil {
	// 	t.Error("unexpected errro:", err)
	// }
	if strings.Compare("Meno", string(opc2.Payload)) != 0 {
		t.Error("not expected string:", string(opc2.Payload))
	}
	json1, _ := json.Marshal(opc1.Claim)
	// json2, _ := json.Marshal(opc2.Claim)
	claim := PearlClaim{}
	json.Unmarshal(json1, &claim)
	if strings.Compare(opc1.Claim.EncryptedPayloadPassword, claim.EncryptedPayloadPassword) != 0 {
		t.Error("should be the same")
	}
	if strings.Compare(opc1.Claim.PayloadChecksum, claim.PayloadChecksum) != 0 {
		t.Error("should be the same")
	}

}
