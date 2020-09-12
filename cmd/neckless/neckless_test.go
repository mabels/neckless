package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"neckless.adviser.com/gem"
	"neckless.adviser.com/member"
)

/*
NECKLESS="go run neckless.adviser.com/cmd/neckless"
rm -f casket.User1.json
$NECKLESS casket --file casket.User1.json create --person --name "Person.User1"
$NECKLESS casket --file casket.User1.json create --device --name "Device.User1"
rm -f casket.User2.json
$NECKLESS casket --file casket.User2.json create --name "Person.User1"
rm -f neckless.shared.json
$NECKLESS gem add --casketFile casket.User1.json  --file neckless.shared.json
$NECKLESS casket --file casket.User2.json get | $NECKLESS gem --casketFile casket.User1.json -file neckless.shared.json add
$NECKLESS gem --casketFile casket.User1.json -file neckless.shared.json ls
$NECKLESS gem --casketFile casket.User2.json -file neckless.shared.json ls

*/

func cmdNeckless(t *testing.T, strargs string, stdin ...string) (*NecklessIO, error) {
	cargs := []string{"neckless"}
	splitted := strings.Split(strargs, " ")
	cargs = append(cargs, splitted...) //, splitted)
	my := strings.Join(stdin, "\n")
	args := NecklessArgs{
		Nio: NecklessIO{
			// in:  bytes.Buffer,
			in:  bufio.NewReader(strings.NewReader(my)),
			out: NecklessOutputs{nos: []NecklessOutput{{buf: new(bytes.Buffer)}}},
			err: NecklessOutputs{nos: []NecklessOutput{{buf: new(bytes.Buffer)}}},
		},
	}
	// fmt.Println(">>>", cargs)
	_, err := buildArgs(cargs, &args)
	if t != nil && err != nil {
		// pwd, _ := os.Getwd()
		t.Error(fmt.Sprintf("%s=>%s", strargs, err))
	}
	return &args.Nio, err
}

func TestAddUserToGem(t *testing.T) {
	for i := 0; i < 43; i++ {
		// test for unsort maps quirks
		os.Remove("casket.User1.json")
		nio, err := cmdNeckless(t, "casket --file casket.User1.json create --person --name Person.User1 --email test@test.com")
		if err != nil {
			t.Error(err)
		}
		nio, err = cmdNeckless(t, "casket --file casket.User1.json create --device --name Device.User1")
		if err != nil {
			t.Error(err)
		}
		nio, err = cmdNeckless(t, "casket --file casket.User1.json get test@test")
		if err != nil {
			t.Error(err)
		}
		tmp := []member.JsonPublicMember{}
		err = json.Unmarshal(nio.out.first().buf.Bytes(), &tmp)
		if err != nil {
			t.Error(string(nio.out.first().buf.Bytes()))
			t.Error(err)
		}
		if len(tmp) != 1 {
			t.Error(string(nio.out.first().buf.Bytes()))
			t.Error(err)
		}
		if !strings.Contains(tmp[0].Email, "test@test") {
			t.Error(string(nio.out.first().buf.Bytes()))
		}
		nio, err = cmdNeckless(t, "casket --file casket.User1.json get --outFile device1.pub.json --device")
		if err != nil {
			t.Error(err)
		}
		os.Remove("casket.User2.json")
		nio, err = cmdNeckless(t, "casket --file casket.User2.json create --name Person.User2")
		if err != nil {
			t.Error(err)
		}
		nio, err = cmdNeckless(t, "casket --file casket.User2.json get --outFile user2.pub.json")
		if err != nil {
			t.Error(err)
		}
		os.Remove("neckless.shared.json")
		// t.Error("X0.out:", nio.out.String(), "\nX0.err:", nio.err.String())
		nio, err = cmdNeckless(t, "gem --casketFile casket.User1.json --file neckless.shared.json add")
		if err != nil {
			t.Error(err)
		}
		nio, err = cmdNeckless(t, "gem --casketFile casket.User1.json --file neckless.shared.json ls")
		if err != nil {
			t.Error(err)
		}
		// t.Error("X1.out:", nio.out.String(), "\nX1.err:", nio.err.String())
		nio, err = cmdNeckless(t, "gem --casketFile casket.User1.json --file neckless.shared.json add --pubFile user2.pub.json")
		if err != nil {
			t.Error(err)
		}

		nio, err = cmdNeckless(t, "gem --casketFile casket.User1.json --file neckless.shared.json add --pubFile device1.pub.json")
		if err != nil {
			t.Error(err)
		}
		// t.Error("X2", nio.out.String(), nio.err.String())
		// t.Error("X2", nio.out.String(), nio.err.String())
		u1nio, err := cmdNeckless(t, "gem --casketFile casket.User1.json --file neckless.shared.json ls")
		if err != nil {
			t.Error(err)
		}
		gtmp := []gem.JsonGem{}
		err = json.Unmarshal(u1nio.out.first().buf.Bytes(), &gtmp)
		if err != nil {
			pwd, _ := os.Getwd()
			t.Error(pwd, string(u1nio.out.first().buf.Bytes()))
			t.Error(err)
		}
		u2nio, err := cmdNeckless(t, "gem --casketFile casket.User2.json --file neckless.shared.json ls")
		if err != nil {
			t.Error(err)
		}
		gtmp = []gem.JsonGem{}
		err = json.Unmarshal(u2nio.out.first().buf.Bytes(), &gtmp)
		if err != nil {
			t.Error(string(u2nio.out.first().buf.Bytes()))
			t.Error(err)
		}
		if len(gtmp) != 1 {
			t.Error("not expected len", len(gtmp))
		}
		if len(gtmp[0].PubKeys) != 3 {
			t.Error("not expected len pubkeys", len(gtmp[0].PubKeys))
		}
		if bytes.Compare(u1nio.out.first().buf.Bytes(), u2nio.out.first().buf.Bytes()) != 0 {
			t.Error("YYYY", u1nio.out.first().buf.String(), u1nio.err.first().buf.String())
			t.Error("XXXX", u2nio.out.first().buf.String(), u2nio.err.first().buf.String())
		}
		var toDelID string
		for i := range gtmp[0].PubKeys {
			if strings.Compare(string(gtmp[0].PubKeys[i].Type), string(member.Device)) == 0 {
				toDelID = gtmp[0].PubKeys[i].Id
			}
		}
		nio, _ = cmdNeckless(t, fmt.Sprintf("gem --casketFile casket.User2.json --file neckless.shared.json rm %s", toDelID))
		// t.Error("X2", tmp[0].PubKeys[0].Id, nio.out.String(), nio.err.String())

		u1nio, _ = cmdNeckless(t, "gem --casketFile casket.User1.json --file neckless.shared.json ls")
		u2nio, _ = cmdNeckless(t, "gem --casketFile casket.User2.json --file neckless.shared.json ls")
		gtmp = []gem.JsonGem{}
		// t.Error(string(u1nio.out.Bytes()))
		err = json.Unmarshal(u1nio.out.first().buf.Bytes(), &gtmp)
		if err != nil {
			t.Error(err)
		}
		if len(gtmp) != 1 {
			t.Error("YYYY", u1nio.out.first().buf.String(), u1nio.err.first().buf.String())
			t.Error("XXXX", u2nio.out.first().buf.String(), u2nio.err.first().buf.String())
			t.Error("not expected len", len(gtmp))
		}
		if len(gtmp[0].PubKeys) != 2 {
			t.Error("YYYY", u1nio.out.first().buf.String(), u1nio.err.first().buf.String())
			t.Error("XXXX", u2nio.out.first().buf.String(), u2nio.err.first().buf.String())
			t.Error("not expected len pubkeys", len(gtmp[0].PubKeys))
		}
		if bytes.Compare(u1nio.out.first().buf.Bytes(), u2nio.out.first().buf.Bytes()) != 0 {
			t.Error("YYYY", u1nio.out.first().buf.String(), u1nio.err.first().buf.String())
			t.Error("XXXX", u2nio.out.first().buf.String(), u2nio.err.first().buf.String())
		}
		nio.out.first().buf.String() // make the compiler happy
	}
	// t.Error("xx", nio.out.String())

}

func createTestData(t *testing.T) {
	os.Remove("casket.User1.json")
	nio, _ := cmdNeckless(t, "casket --file casket.User1.json create --person --name Person.User1")
	nio, _ = cmdNeckless(t, "casket --file casket.User1.json create --device --name Device.User1")
	nio, _ = cmdNeckless(t, "casket --file casket.User1.json get --outFile device1.pub.json --device")
	os.Remove("neckless.shared.json")
	dev1, _ := ioutil.ReadFile("device1.pub.json")
	nio, _ = cmdNeckless(t, "gem --casketFile casket.User1.json --file neckless.shared.json add", string(dev1))
	fromStdin := nio.out.first().buf.String()
	nio, _ = cmdNeckless(t, "gem --casketFile casket.User1.json --file neckless.shared.json add --pubFile device1.pub.json")
	fromFile := nio.out.first().buf.String()
	if fromStdin != fromFile {
		t.Error("should be the same", fromStdin, fromFile)
	}

	nio, err := cmdNeckless(nil, "kv --casketFile casket.User1.json --file neckless.shared.json ls")
	if err != nil {
		t.Error("there should be an error")
	}
	if len(nio.out.first().buf.Bytes()) != 0 {
		t.Error("should be empty", nio.out.first().buf.String())
	}
	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json add M=1 M=2")
	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json add N=4711 M=3")
	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json add TOTP=4S62BZNFXXSZLCRO")
}

func TestKvs(t *testing.T) {
	createTestData(t)
	nio, _ := cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --json")
	var mya []flatKeyValue
	// inputJS := string(nio.out.first().buf.Bytes())
	// t.Error(inputJS)
	json.Unmarshal(nio.out.first().buf.Bytes(), &mya)
	if len(mya) != 3 {
		t.Error("should be", len(mya))
	}
	if strings.Compare(mya[0].Key, "M") != 0 {
		t.Error("not M")
	}
	if !(mya[0].Value == "3") {
		t.Error("not M")
	}
	if strings.Compare(mya[1].Key, "N") != 0 {
		t.Error("not N")
	}
	if strings.Compare(mya[1].Value, "4711") != 0 {
		t.Error("not N")
	}

	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json add YU=1[T1]")
	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json add YU=2[T2]")
	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --tag T1 YU")
	if strings.Compare(strings.TrimSpace(nio.out.first().buf.String()), "YU=\"1\"") != 0 {
		t.Error("not expected", nio.out.first().buf.String())
	}
	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --tag T2 YU")
	if strings.Compare(strings.TrimSpace(nio.out.first().buf.String()), "YU=\"2\"") != 0 {
		t.Error("not expected", nio.out.first().buf.String())
	}
	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json add YU=22[T2]")
	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --tag T2 YU")
	if strings.Compare(strings.TrimSpace(nio.out.first().buf.String()), "YU=\"22\"") != 0 {
		t.Error("not expected", nio.out.first().buf.String())
	}
	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls YU@bla[T2]")
	if strings.Compare(strings.TrimSpace(nio.out.get("bla").buf.String()), "YU=\"22\"") != 0 {
		t.Error("not expected", nio.out.first().buf.String())
	}
	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls T2@bla YU@blu")
	if len(nio.out.first().buf.String()) != 0 {
		t.Error("not expected", nio.out.first().buf.String())
	}
	if nio.out.get("bla") != nil {
		t.Error("not expected", nio.out.first().buf.String())
	}
	if strings.Compare(strings.TrimSpace(nio.out.get("blu").buf.String()), "YU=\"22\"") != 0 {
		t.Error("not expected", nio.out.first().buf.String())
	}
	nio.out.first().buf.String() // make the compiler happy
}

func TestLsGhAddMask(t *testing.T) {
	createTestData(t)
	// --ghAddMask      set Value as github mask
	nio, _ := cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --shKeyValue --ghAddMask M@bla N@bla N M")
	if strings.Compare(strings.TrimSpace(nio.out.first().buf.String()), "M=\"3\";\nexport M;\necho ::add-mask::3;\nN=\"4711\";\nexport N;\necho ::add-mask::4711;") != 0 {
		t.Error(nio.out.first().buf.String())
	}
	if strings.Compare(strings.TrimSpace(nio.out.get("bla").buf.String()), "M=\"3\";\nexport M;\necho ::add-mask::3;\nN=\"4711\";\nexport N;\necho ::add-mask::4711;") != 0 {
		t.Error(nio.out.get("bla").buf.String())
	}

	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --ghAddMask M@bla N@bla N M")
	if strings.Compare(strings.TrimSpace(nio.out.first().buf.String()), "M=\"3\"\necho ::add-mask::3\nN=\"4711\"\necho ::add-mask::4711") != 0 {
		t.Error(nio.out.first().buf.String())
	}
	if strings.Compare(strings.TrimSpace(nio.out.get("bla").buf.String()), "M=\"3\"\necho ::add-mask::3\nN=\"4711\"\necho ::add-mask::4711") != 0 {
		t.Error(nio.out.get("bla").buf.String())
	}

}

func genRef(param string) string {
	// unresolved := kvpearl.ParseFuncsAndParams(param)
	ref := []flatKeyValue{
		{
			Key:   "M",
			Value: "3",
			Tags:  []string{},
		},
		{
			Key:   "N",
			Value: "4711",
			Tags:  []string{},
		},
	}
	jsb, _ := json.MarshalIndent(ref, "", "  ")
	ret := string(jsb)
	return ret
}
func TestLsJson(t *testing.T) {
	createTestData(t)
	// --json           select device keys
	nio, _ := cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --json M@bla N@bla M N")
	if strings.Compare(strings.TrimSpace(nio.out.first().buf.String()), genRef("")) != 0 {
		t.Error(nio.out.first().buf.String())
	}
	if strings.Compare(strings.TrimSpace(nio.out.get("bla").buf.String()), genRef("bla")) != 0 {
		t.Error(nio.out.get("bla").buf.String())
	}
}

func TestLsKeyValue(t *testing.T) {
	createTestData(t)
	// --keyValue           select device keys
	nio, _ := cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --keyValue M@bla N@bla N M")
	if strings.Compare(strings.TrimSpace(nio.out.first().buf.String()), "M=\"3\"\nN=\"4711\"") != 0 {
		t.Error(nio.out.first().buf.String())
	}
	if strings.Compare(strings.TrimSpace(nio.out.get("bla").buf.String()), "M=\"3\"\nN=\"4711\"") != 0 {
		t.Error(nio.out.get("bla").buf.String())
	}
}

func TestLsOnlyValue(t *testing.T) {
	createTestData(t)
	// --onlyValue           select device keys
	nio, _ := cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --onlyValue M@bla N@bla N M")
	if strings.Compare(strings.TrimSpace(nio.out.first().buf.String()), "3\n4711") != 0 {
		t.Error(nio.out.first().buf.String())
	}
	if strings.Compare(strings.TrimSpace(nio.out.get("bla").buf.String()), "3\n4711") != 0 {
		t.Error(nio.out.get("bla").buf.String())
	}
}
func TestLsShKeyValue(t *testing.T) {
	createTestData(t)
	// --shKeyValue           select device keys
	nio, _ := cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --shKeyValue M@bla N@bla N M")
	if strings.Compare(strings.TrimSpace(nio.out.first().buf.String()), "M=\"3\";\nexport M;\nN=\"4711\";\nexport N;") != 0 {
		t.Error(nio.out.first().buf.String())
	}
	if strings.Compare(strings.TrimSpace(nio.out.get("bla").buf.String()), "M=\"3\";\nexport M;\nN=\"4711\";\nexport N;") != 0 {
		t.Error(nio.out.get("bla").buf.String())
	}
}

func TestActions(t *testing.T) {
	createTestData(t)
	// --shKeyValue           select device keys
	nio, err := cmdNeckless(nil, "kv --casketFile casket.User1.json --file neckless.shared.json ls M@Action(Hund())")
	if err == nil {
		t.Error("Should be an error")
	}
	if err.Error() != "unknown action:Hund" {
		t.Error(err.Error())
	}
	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls TOTP@Len(Noop())")
	if strings.Compare(strings.TrimSpace(nio.out.first().buf.String()), "TOTP=\"16\"") != 0 {
		t.Error(nio.out.first().buf.String())
	}
	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --onlyValue TOTP@Totp()")
	if len(strings.TrimSpace(nio.out.first().buf.String())) != 6 {
		t.Error(nio.out.first().buf.String())
	}
	ref := []flatKeyValue{
		flatKeyValue{
			Key:   "TOTP",
			Value: "setit",
			Tags:  []string{},
		},
	}
	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --json TOTP@Totp()")
	var myRef []flatKeyValue
	err = json.Unmarshal(nio.out.first().buf.Bytes(), &myRef)
	if err != nil {
		t.Error("should not happend")
	}
	ref[0].Value = myRef[0].Value
	jsRef, _ := json.MarshalIndent(ref, "", "  ")
	if string(jsRef) == nio.out.first().buf.String() {
		t.Error(nio.out.first().buf.String())
	}
}
