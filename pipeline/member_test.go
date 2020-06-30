package pipeline

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestMakePrivateMemberEmpty(t *testing.T) {
	_, err := MakePrivateMember(&PrivateMemberArg{
		Member: MemberArg{},
	})
	if err == nil {
		t.Error("Err should set")
	}
}

func TestMakePrivateMemberInitial(t *testing.T) {
	pm, err := MakePrivateMember(&PrivateMemberArg{
		Member: MemberArg{
			Type:   Person,
			Name:   "testName",
			Device: "testDevice",
		},
	})
	if err != nil {
		t.Error("there show no error")
	}
	if strings.Compare(string(pm.PrivateKey.Key.Style), string(Private)) != 0 {
		t.Error("there show no error")
	}

}
func compareMember(m1 *Member, m2 *Member, t *testing.T) {

	if strings.Compare(m1.Id, m2.Id) != 0 {
		t.Error("Id error")
	}
	if strings.Compare(string(m1.Type), string(m2.Type)) != 0 {
		t.Error("Type error")
	}
	if strings.Compare(m1.Name, m2.Name) != 0 {
		t.Error("Name error")
	}
	if strings.Compare(m1.Device, m2.Device) != 0 {
		t.Error("Device error")
	}
	// t.Error(pmi.Member.ValidUntil.String(), pm.Member.ValidUntil.String())
	if !time.Time.Equal(m1.ValidUntil, m2.ValidUntil) {
		t.Error("ValidUntil error")
	}
	if !time.Time.Equal(m1.Created, m2.Created) {
		t.Error("Created error")
	}
	if !time.Time.Equal(m1.Updated, m2.Updated) {
		t.Error("Updated error")
	}
}
func TestMakePrivateMemberPreset(t *testing.T) {
	pmi, _ := MakePrivateMember(&PrivateMemberArg{
		Member: MemberArg{
			Type:   Person,
			Name:   "testName",
			Device: "testDevice",
		},
	})
	pm, _ := MakePrivateMember(&PrivateMemberArg{
		Member: MemberArg{
			Id:         pmi.Id,
			Type:       pmi.Type,
			Name:       pmi.Name,
			Device:     pmi.Device,
			ValidUntil: &pmi.ValidUntil,
			Updated:    &pmi.Updated,
			Created:    &pmi.Created,
		},
		PrivateKey: &pmi.PrivateKey,
	})
	compareMember(&pmi.Member, &pm.Member, t)
	if bytes.Compare(pmi.PrivateKey.Key.Raw[:], pm.PrivateKey.Key.Raw[:]) != 0 {
		t.Error("Key are not the same")
	}
}

func TestMakePublicMemberInitial(t *testing.T) {
}

func TestMakePublicFromPrivate(t *testing.T) {
	pkm, _ := MakePrivateMember(&PrivateMemberArg{
		Member: MemberArg{
			Type:   Person,
			Name:   "testName",
			Device: "testDevice",
		},
	})
	pub := pkm.Public()
	compareMember(&pkm.Member, &pub.Member, t)
	if strings.Compare(string(pub.PublicKey.Key.Style), string(Public)) != 0 {
		t.Error("not public")
	}
	if bytes.Compare(pub.PublicKey.Key.Raw[:], pkm.PrivateKey.Public().Key.Raw[:]) != 0 {
		t.Error("no match public")
	}
}

func TestPrivateMemberJson(t *testing.T) {
	pkm, _ := MakePrivateMember(&PrivateMemberArg{
		Member: MemberArg{
			Type:   Person,
			Name:   "testName",
			Device: "testDevice",
		},
	})
	json, _ := pkm.AsJson().String()
	jpkm, _, err := FromJson(json)
	if err != nil {
		t.Error(err)
	}
	compareMember(&pkm.Member, &jpkm.Member, t)
	if bytes.Compare(pkm.PrivateKey.Key.Raw[:], jpkm.PrivateKey.Key.Raw[:]) != 0 {
		t.Error("keys not matching")
	}
}

func TestPublicMemberJson(t *testing.T) {
	pkm, _ := MakePrivateMember(&PrivateMemberArg{
		Member: MemberArg{
			Type:   Person,
			Name:   "testName",
			Device: "testDevice",
		},
	})
	pubk := pkm.Public()
	json, _ := pubk.AsJson().String()
	_, jpkm, err := FromJson(json)
	if err != nil {
		t.Error(err)
	}
	compareMember(&pkm.Member, &jpkm.Member, t)
	if bytes.Compare(pubk.PublicKey.Key.Raw[:], jpkm.PublicKey.Key.Raw[:]) != 0 {
		t.Error("keys not matching")
	}
}
