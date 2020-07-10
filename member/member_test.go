package member

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"neckless.adviser.com/key"
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
	if strings.Compare(string(pm.PrivateKey.Key.Style), string(key.Private)) != 0 {
		t.Error("there show no error")
	}

}
func compareMember(m1 *Member, m2 *Member, t *testing.T) {

	if strings.Compare(m1.Id, m2.Id) != 0 {
		t.Error("Id error", m1.Id, m2.Id)
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
	if strings.Compare(string(pub.PublicKey.Key.Style), string(key.Public)) != 0 {
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
	// fmt.Printf("TestPrivateMemberJson:%s\n", json)
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

func TestFilterById(t *testing.T) {
	lst := []*PrivateMember{}
	for i := 0; i < 10; i++ {
		person, _ := MakePrivateMember(&PrivateMemberArg{
			Member: MemberArg{
				Type: Person,
				Name: fmt.Sprintf("testPerson%d", i),
			},
		})
		lst = append(lst, person)
		device, _ := MakePrivateMember(&PrivateMemberArg{
			Member: MemberArg{
				Type: Device,
				Name: fmt.Sprintf("testDevice%d", i),
			},
		})
		lst = append(lst, device)
	}
	if len(FilterById(lst)) != len(lst) {
		t.Error("Filter by nothing")
	}
	if len(FilterById(lst, []string{}...)) != len(lst) {
		t.Error("Filter by nothing")
	}
	if len(FilterById(lst)) != len(lst) {
		t.Error("Filter by nothing")
	}
	if len(FilterById(lst, []string{"darf nix passieren"}...)) != 0 {
		t.Error("Filter by nothing")
	}
	if len(FilterById(lst, lst[0].Id, lst[3].Id, "garnix")) != 2 {
		t.Error("Filter by nothing")
	}
}
func TestFilterByType(t *testing.T) {
	lst := []*PrivateMember{}
	for i := 0; i < 10; i++ {
		person, _ := MakePrivateMember(&PrivateMemberArg{
			Member: MemberArg{
				Type: Person,
				Name: fmt.Sprintf("testPerson%d", i),
			},
		})
		lst = append(lst, person)
		device, _ := MakePrivateMember(&PrivateMemberArg{
			Member: MemberArg{
				Type: Device,
				Name: fmt.Sprintf("testDevice%d", i),
			},
		})
		lst = append(lst, device)
	}
	if len(FilterByType(lst)) != len(lst) {
		t.Error("Filter by nothing")
	}
	mtyps := make([]MemberType, 0)
	if len(FilterByType(lst, mtyps...)) != len(lst) {
		t.Error("Filter by nothing")
	}
	if len(FilterByType(lst, Person)) != len(lst)/2 {
		t.Error("Filter by nothing")
	}
	if len(FilterByType(lst, Device)) != len(lst)/2 {
		t.Error("Filter by nothing")
	}
	if len(FilterByType(lst, Device, Person)) != len(lst) {
		t.Error("Filter by nothing")
	}
	if len(FilterByType(lst, Person, Device)) != len(lst) {
		t.Error("Filter by nothing")
	}
}
