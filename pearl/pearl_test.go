package pearl

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"neckless.adviser.com/key"
	"neckless.adviser.com/member"
)

func createTestKeyO1() *key.RawKey {
	x := key.RawKey{}
	for i := range x {
		x[i] = byte(i)
	}
	return &x
}

func createTestKeyO2() *key.RawKey {
	x := key.RawKey{}
	for i := range x {
		x[i] = byte(i + 32)
	}
	return &x
}
func TestPearl(t *testing.T) {

	o1, _ := member.MakePrivateMember(&member.PrivateMemberArg{
		Member: member.MemberArg{
			Type: member.Person,
			Name: "o1",
		},
		PrivateKey: key.MakePrivateKey(createTestKeyO1()),
	})
	o2, _ := member.MakePrivateMember(&member.PrivateMemberArg{
		Member: member.MemberArg{
			Type: member.Person,
			Name: "o2",
		},
		PrivateKey: key.MakePrivateKey(createTestKeyO2()),
	})
	owners := []*key.PublicKey{}
	owners = append(owners, &o1.Public().PublicKey)
	owners = append(owners, &o2.Public().PublicKey)
	// fmt.Printf("%x:%x\n", o1.Public().PublicKey.Key.Raw, o2.Public().PublicKey.Key.Raw)
	//signer.Key.Raw, pub.Key.Raw, privPubKey)

	pearl, err := Close(&CloseRequestPearl{
		Type:    "test",
		Payload: []byte("Meno"),
		Owners: PearlOwner{
			Signer: &o1.PrivateKey,
			Owners: owners,
		},
	})
	if err != nil {
		t.Error("unexpected errro:", err)
	}
	if bytes.Compare(pearl.FingerPrint, []byte{
		187, 49, 78, 117, 205, 12, 114, 70, 150, 235, 216, 139, 127, 53, 58, 179, 96,
		158, 114, 155, 231, 225, 185, 124, 140, 239, 33, 66, 237, 30, 68, 104,
	}) != 0 {
		t.Error("unexpected errro:", pearl.FingerPrint)
	}
	// data, err := json.Marshal(pearl)
	// t.Error(string(data))
	opc1, err := Open([]*key.PrivateKey{&o1.PrivateKey}, pearl)
	if err != nil {
		t.Error("unexpected errro:", err)
	}
	// data, err = json.Marshal(opc1)
	if strings.Compare("Meno", string(opc1.Payload)) != 0 {
		t.Error("not expected string:", string(opc1.Payload))
	}
	// t.Error(string(opc1.Payload))

	opc2, err := OpenOne(&o2.PrivateKey, pearl)
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

func TestJsonPearl(t *testing.T) {
	o1, _ := member.MakePrivateMember(&member.PrivateMemberArg{
		Member: member.MemberArg{
			Type: member.Person,
			Name: "o1",
		},
	})
	owners := []*key.PublicKey{&o1.Public().PublicKey}
	pearl, err := Close(&CloseRequestPearl{
		Type:    "test",
		Payload: []byte("Meno"),
		Owners: PearlOwner{
			Signer: &o1.PrivateKey,
			Owners: owners,
		},
	})
	if err != nil {
		t.Error("unexpected errro:", err)
	}

	str, err := json.MarshalIndent(pearl.AsJSON(), "", "  ")
	if err != nil {
		t.Error("json should work", err)
	}

	jsbackPearl := &JSONPearl{}
	json.Unmarshal(str, &jsbackPearl)
	// t.Error(string(str))
	// t.Error(jsbackPearl)
	// data, err := json.Marshal(pearl)
	// t.Error(string(data))
	backPearl, err := jsbackPearl.FromJSON()
	if err != nil {
		t.Error("fromjson should not error", err)
	}
	opc1, err := Open([]*key.PrivateKey{&o1.PrivateKey}, backPearl)
	if err != nil {
		t.Error("unexpected errro:", err)
	}
	// t.Error(opc1)
	// if bytes.Compare(owners[0].Key.Raw, opc1.
	// t.Error("not expected string:", string(opc1.Payload))
	// }
	// data, err = json.Marshal(opc1)
	if strings.Compare("Meno", string(opc1.Payload)) != 0 {
		t.Error("not expected string:", string(opc1.Payload))
	}

	if bytes.Compare(opc1.Closed.FingerPrint, pearl.FingerPrint) != 0 {
		t.Error("unexpected errro:", pearl.FingerPrint)
	}

}
