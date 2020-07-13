package gem

import (
	"encoding/json"
	"strings"
	"testing"

	"neckless.adviser.com/member"
)

func TestCreate(t *testing.T) {
	gem := Create()
	if len(gem.PubKeys) != 0 {
		t.Error("not expected")
	}
	pk1, _ := member.MakePrivateMember(&member.PrivateMemberArg{
		Member: member.MemberArg{
			Type:   member.Person,
			Name:   "testName",
			Device: "testDevice",
		},
	})
	pk2, _ := member.MakePrivateMember(&member.PrivateMemberArg{
		Member: member.MemberArg{
			Type:   member.Person,
			Name:   "testName",
			Device: "testDevice",
		},
	})
	gem = Create(pk1.Public(), pk2.Public())
	if len(gem.PubKeys) != 2 {
		t.Error("not expected")
	}
}

func TestAdd(t *testing.T) {
	gem := Create()
	if len(gem.PubKeys) != 0 {
		t.Error("not expected")
	}
	pk1, _ := member.MakePrivateMember(&member.PrivateMemberArg{
		Member: member.MemberArg{
			Type:   member.Person,
			Name:   "testName",
			Device: "testDevice",
		},
	})
	gem.Add(pk1.Public())
	gem.Add(pk1.Public())
	pk2, _ := member.MakePrivateMember(&member.PrivateMemberArg{
		Member: member.MemberArg{
			Type:   member.Person,
			Name:   "testName",
			Device: "testDevice",
		},
	})
	gem.Add(pk2.Public())
	if len(gem.PubKeys) != 2 {
		t.Error("not expected len")
	}
}

func TestRm(t *testing.T) {
	pk1, _ := member.MakePrivateMember(&member.PrivateMemberArg{
		Member: member.MemberArg{
			Type:   member.Person,
			Name:   "testName",
			Device: "testDevice",
		},
	})
	pk2, _ := member.MakePrivateMember(&member.PrivateMemberArg{
		Member: member.MemberArg{
			Type:   member.Person,
			Name:   "testName",
			Device: "testDevice",
		},
	})
	gem := Create(pk1.Public(), pk2.Public())
	if len(gem.PubKeys) != 2 {
		t.Error("not expected")
	}
	gem.Rm()
	if len(gem.PubKeys) != 2 {
		t.Error("not expected")
	}
	gem.Rm("doof", "dumm")
	if len(gem.PubKeys) != 2 {
		t.Error("not expected")
	}
	gem.Rm(pk1.Id, pk2.Id)
	if len(gem.PubKeys) != 0 {
		t.Error("not expected")
	}
}

func TestLs(t *testing.T) {
	pk1, _ := member.MakePrivateMember(&member.PrivateMemberArg{
		Member: member.MemberArg{
			Type:   member.Person,
			Name:   "testName",
			Device: "testDevice",
		},
	})
	pk2, _ := member.MakePrivateMember(&member.PrivateMemberArg{
		Member: member.MemberArg{
			Type:   member.Person,
			Name:   "testName",
			Device: "testDevice",
		},
	})
	gem := Create(pk1.Public(), pk2.Public())
	if len(gem.Ls()) != 2 {
		t.Error("not expected")
	}
	if len(gem.Ls("doof")) != 0 {
		t.Error("not expected")
	}
	if len(gem.Ls("doof", pk1.Id)) != 1 {
		t.Error("not expected")
	}
	if len(gem.Ls("doof", pk1.Id, pk2.Id, "offk")) != 2 {
		t.Error("not expected")
	}
}

func TestLsType(t *testing.T) {
	pk1, _ := member.MakePrivateMember(&member.PrivateMemberArg{
		Member: member.MemberArg{
			Type:   member.Person,
			Name:   "testName",
			Device: "testDevice",
		},
	})
	pk2, _ := member.MakePrivateMember(&member.PrivateMemberArg{
		Member: member.MemberArg{
			Type:   member.Person,
			Name:   "testName",
			Device: "testDevice",
		},
	})
	gem := Create(pk1.Public(), pk2.Public())
	if len(gem.LsByType(member.Device)) != 0 {
		t.Error("not expected")
	}
	if len(gem.LsByType(member.Person)) != 2 {
		t.Error("not expected")
	}
	pk3, _ := member.MakePrivateMember(&member.PrivateMemberArg{
		Member: member.MemberArg{
			Type:   member.Device,
			Name:   "testName",
			Device: "testDevice",
		},
	})
	gem.Add(pk3.Public())
	if len(gem.LsByType(member.Device)) != 1 {
		t.Error("not expected")
	}
	if len(gem.LsByType(member.Person)) != 2 {
		t.Error("not expected")
	}
}

func TestJson(t *testing.T) {
	pk1, _ := member.MakePrivateMember(&member.PrivateMemberArg{
		Member: member.MemberArg{
			Type:   member.Person,
			Name:   "testName",
			Device: "testDevice",
		},
	})
	pk2, _ := member.MakePrivateMember(&member.PrivateMemberArg{
		Member: member.MemberArg{
			Type:   member.Person,
			Name:   "testName",
			Device: "testDevice",
		},
	})
	gem := Create(pk1.Public(), pk2.Public())
	str, _ := json.Marshal(gem.AsJSON())
	back, err := FromJSON(str)
	if err != nil {
		t.Error("not expected")
	}
	if len(back.PubKeys) != 2 {
		t.Error("not expected")
	}
	if strings.Compare(gem.PubKeys[pk1.Id].Id, pk1.Id) != 0 {
		t.Error("not expected")
	}
	if strings.Compare(gem.PubKeys[pk2.Id].Id, pk2.Id) != 0 {
		t.Error("not expected")
	}
}
