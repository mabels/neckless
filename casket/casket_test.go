package casket

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/mabels/neckless/key"
	"github.com/mabels/neckless/member"
)

func TestCreate(t *testing.T) {
	now := time.Now()
	ks, _, err := Create(CreateArg{
		MemberArg: member.MemberArg{
			Name:   "Test",
			Device: "TestDevice",
			Type:   member.Person,
		},
		DryRun: true,
	})
	if err != nil {
		t.Error("no error expected", err)
	}
	if strings.Compare(path.Base(*ks.CasketFname), "casket") != 0 {
		t.Error("Fname", *ks.CasketFname)
	}
	if ks.Created.Before(now) {
		t.Error("no error expected", err)
	}
	if ks.Updated.Before(now) {
		t.Error("no error expected", err)
	}
	if len(ks.Members) != 1 {
		t.Error("not ok", err)
	}
	for k := range ks.Members {
		if strings.Compare(k, ks.Members[k].Id) != 0 {
			t.Error("key error", k)
		}
		if strings.Compare("Test", ks.Members[k].Name) != 0 {
			t.Error("Name error", k)
		}
		if strings.Compare("TestDevice", ks.Members[k].Device) != 0 {
			t.Error("Device error", k)
		}
		if strings.Compare("Person", string(ks.Members[k].Type)) != 0 {
			t.Error("Person error", k)
		}
	}
}

func TestLs(t *testing.T) {
	id := os.Process{}.Pid
	fname := fmt.Sprintf("./test-%d/casket-test.%d.json", id, id)
	os.RemoveAll(fname)
	os.RemoveAll(path.Dir(fname))
	_, t1, err := Create(CreateArg{
		MemberArg: member.MemberArg{
			Name:   "Test1",
			Device: "TestDevice",
			Type:   member.Person,
		},
		Fname: &fname,
	})
	if err != nil {
		t.Error("create failed", err)
	}
	_, t2, err := Create(CreateArg{
		MemberArg: member.MemberArg{
			Name:   "Test2",
			Device: "TestDevice",
			Type:   member.Person,
		},
		Fname: &fname,
	})
	if err != nil {
		t.Error("create failed", err)
	}
	ks, err := Ls(fname)
	if err != nil {
		t.Error("no error expected", err)
	}
	if len(ks.Members) != 2 {
		t.Error("len members not expected", len(ks.Members))
	}
	cnt := 0
	for k := range ks.Members {
		val := ks.Members[k]
		if !(strings.Compare(val.Name, t1.Name) == 0 ||
			strings.Compare(val.Name, t2.Name) == 0) {
			t.Error("no right name", val.Name)
		}
		if strings.Compare(string(val.PrivateKey.Key.Style), string(key.Private)) != 0 {
			t.Error("is not private Key")
		}
		if !(bytes.Compare(val.PrivateKey.Key.Raw[:], t1.PrivateKey.Key.Raw[:]) == 0 ||
			bytes.Compare(val.PrivateKey.Key.Raw[:], t2.PrivateKey.Key.Raw[:]) == 0) {
			t.Error("no nice key")
		}
		cnt++
	}
	if cnt != 2 {
		t.Error("not cnt ok:", cnt)
	}
	// os.RemoveAll(fname)
	// os.RemoveAll(path.Dir(fname))
}

func TestLsEnvNecklessPrivkey(t *testing.T) {
	os.Setenv("NECKLESS_PRIVKEY", "xkxxkxk")
	ks, err := Ls()
	if err != nil {
		t.Error("no error expected", err)
	}
	if len(ks.Members) != 0 {
		t.Error("len members not expected", len(ks.Members))
	}
	if *ks.CasketAttribute.CasketFname != "ENV:NECKLESS_PRIVKEY" {
		t.Error("fname should ENV:NECKLESS_PRIVKEY==", *ks.CasketAttribute.CasketFname)
	}
}

func TestRm(t *testing.T) {
	id := os.Process{}.Pid
	fname := fmt.Sprintf("./test-%d/casket-test.%d.json", id, id)
	os.RemoveAll(fname)
	os.RemoveAll(path.Dir(fname))
	_, t1, err := Create(CreateArg{
		MemberArg: member.MemberArg{
			Name:   "Test1",
			Device: "TestDevice",
			Type:   member.Person,
		},
		Fname: &fname,
	})
	if err != nil {
		t.Error("create failed", err)
	}
	_, t2, err := Create(CreateArg{
		MemberArg: member.MemberArg{
			Name:   "Test2",
			Device: "TestDevice",
			Type:   member.Person,
		},
		Fname: &fname,
	})
	casket, pks, err := Rm(RmArg{
		Ids:    []string{t1.Id, "doof"},
		Fname:  &fname,
		DryRun: true,
	})
	if err != nil {
		t.Error("expect no error", err)
	}
	if len(pks) != 1 {
		t.Error("not expected pks")
	}
	if len(casket.Members) != 1 {
		t.Error("not expect len")
	}
	if _, found := casket.Members[t2.Id]; !found {
		t.Error("not found ")
	}
	_, _, err = Rm(RmArg{
		Ids:   []string{t1.Id, "bloed"},
		Fname: &fname,
	})
	casket, err = Ls(fname)
	if err != nil {
		t.Error("expect no error", err)
	}
	if len(casket.Members) != 1 {
		t.Error("not expect len")
	}
	if _, found := casket.Members[t2.Id]; !found {
		t.Error("not found ")
	}
	// os.RemoveAll(fname)
	// os.RemoveAll(path.Dir(fname))
}
