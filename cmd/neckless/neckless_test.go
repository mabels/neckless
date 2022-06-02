package neckless

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/mabels/neckless/gem"
	"github.com/mabels/neckless/member"
)

/*
NECKLESS="go run github.com/mabels/neckless/cmd/neckless"
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
	// fmt.Printf(">>>%v\n", cargs[len(cargs)-1])
	_, err := buildArgs(cargs, &args)
	if t != nil && err != nil {
		// pwd, _ := os.Getwd()
		t.Error(fmt.Errorf("%s=>%v", strargs, err))
	}
	return &args.Nio, err
}

func TestAddUserToGem(t *testing.T) {
	for i := 0; i < 43; i++ {
		// test for unsort maps quirks
		os.Remove("casket.User1.json")
		_, err := cmdNeckless(t, "casket --file casket.User1.json create --person --name Person.User1 --email test@test.com")
		if err != nil {
			t.Error(err)
		}
		_, err = cmdNeckless(t, "casket --file casket.User1.json create --device --name Device.User1")
		if err != nil {
			t.Error(err)
		}
		nio, err := cmdNeckless(t, "casket --file casket.User1.json get test@test")
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
	_, _ = cmdNeckless(t, "casket --file casket.User1.json create --person --name Person.User1")
	_, _ = cmdNeckless(t, "casket --file casket.User1.json create --device --name Device.User1")
	_, _ = cmdNeckless(t, "casket --file casket.User1.json get --outFile device1.pub.json --device")
	os.Remove("neckless.shared.json")
	dev1, _ := ioutil.ReadFile("device1.pub.json")
	nio, _ := cmdNeckless(t, "gem --casketFile casket.User1.json --file neckless.shared.json add", string(dev1))
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
	_, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json add M=1 M=2")
	nio, err = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json add N=4711 M=3 QUOTE=$<DO\r\n\"OF>")
	if err != nil {
		t.Error("there should be an error")
	}
	// t.Error("CMDOUT:", nio.out.first().buf.String())
	// t.Error("CMDERR:", nio.err.first().buf.String())
	_, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json add TOTP=4S62BZNFXXSZLCRO")
}

func TestKvs(t *testing.T) {
	createTestData(t)
	nio, _ := cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --json")
	var mya []FlatKeyValue
	// inputJS := string(nio.out.first().buf.Bytes())
	// t.Error(inputJS)
	json.Unmarshal(nio.out.first().buf.Bytes(), &mya)
	if len(mya) != 4 {
		t.Error("should be not ", len(mya))
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

	_, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json add YU=1[T1]")
	_, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json add YU=2[T2]")
	nio, err := cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --tag T1 YU")
	if err != nil {
		t.Error(err)
	}
	if strings.Compare(strings.TrimSpace(nio.out.first().buf.String()), "YU=1") != 0 {
		t.Error("not expected", nio.out.first().buf.String())
	}
	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --tag T2 YU")
	if strings.Compare(strings.TrimSpace(nio.out.first().buf.String()), "YU=2") != 0 {
		t.Error("not expected", nio.out.first().buf.String())
	}
	_, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json add YU=22[T2]")
	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --tag T2 YU")
	if strings.Compare(strings.TrimSpace(nio.out.first().buf.String()), "YU=22") != 0 {
		t.Error("not expected", nio.out.first().buf.String())
	}
	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls YU@bla[T2]")
	if strings.Compare(strings.TrimSpace(nio.out.get("bla").buf.String()), "YU=22") != 0 {
		t.Error("not expected", nio.out.first().buf.String())
	}
	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --tag T2 T2@bla YU@blu")
	if len(nio.out.first().buf.String()) != 0 {
		t.Error("not expected", nio.out.first().buf.String())
	}
	if nio.out.get("bla") != nil {
		t.Error("not expected", nio.out.first().buf.String())
	}
	if strings.Compare(strings.TrimSpace(nio.out.get("blu").buf.String()), "YU=22") != 0 {
		t.Error("not expected", nio.out.first().buf.String())
	}
	// nio.out.first().buf.String() // make the compiler happy
}

func TestLsGhAddMask(t *testing.T) {
	createTestData(t)
	// --ghAddMask      set Value as github mask
	nio, _ := cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --shKeyValue --ghAddMask M@bla N@bla N M")
	if strings.Compare(strings.TrimSpace(nio.out.first().buf.String()), "M=3;\nexport M;\necho ::add-mask::3;\nN=4711;\nexport N;\necho ::add-mask::4711;") != 0 {
		t.Error(nio.out.first().buf.String())
	}
	if strings.Compare(strings.TrimSpace(nio.out.get("bla").buf.String()), "M=3;\nexport M;\necho ::add-mask::3;\nN=4711;\nexport N;\necho ::add-mask::4711;") != 0 {
		t.Error(nio.out.get("bla").buf.String())
	}

	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --ghAddMask M@bla N@bla N M")
	if strings.Compare(strings.TrimSpace(nio.out.first().buf.String()), "M=3\necho ::add-mask::3\nN=4711\necho ::add-mask::4711") != 0 {
		t.Error(nio.out.first().buf.String())
	}
	if strings.Compare(strings.TrimSpace(nio.out.get("bla").buf.String()), "M=3\necho ::add-mask::3\nN=4711\necho ::add-mask::4711") != 0 {
		t.Error(nio.out.get("bla").buf.String())
	}

	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --ghAddMask QUOTE")
	if strings.Compare(strings.TrimSpace(nio.out.first().buf.String()), "QUOTE='$<DO\r\n\"OF>'\necho ::add-mask::'$<DO\r\n\"OF>'") != 0 {
		t.Error(nio.out.first().buf.String())
	}

}

func TestLsGhEvalAddMask(t *testing.T) {
	createTestData(t)
	// --ghAddMask      set Value as github mask
	nio, _ := cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --shKeyValue --ghEvalAddMask M@bla N@bla N M")
	if strings.Compare(strings.TrimSpace(nio.out.first().buf.String()), "M=\"3\";\nexport M;\necho ::add-mask::\"3\";\nN=\"4711\";\nexport N;\necho ::add-mask::\"4711\";") != 0 {
		t.Error(nio.out.first().buf.String())
	}

	if strings.Compare(strings.TrimSpace(nio.out.get("bla").buf.String()), "M=\"3\";\nexport M;\necho ::add-mask::\"3\";\nN=\"4711\";\nexport N;\necho ::add-mask::\"4711\";") != 0 {
		t.Error(nio.out.get("bla").buf.String())
	}

	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --ghEvalAddMask M@bla N@bla N M")
	if strings.Compare(strings.TrimSpace(nio.out.first().buf.String()), "M=\"3\"\necho ::add-mask::\"3\"\nN=\"4711\"\necho ::add-mask::\"4711\"") != 0 {
		t.Error(nio.out.first().buf.String())
	}
	if strings.Compare(strings.TrimSpace(nio.out.get("bla").buf.String()), "M=\"3\"\necho ::add-mask::\"3\"\nN=\"4711\"\necho ::add-mask::\"4711\"") != 0 {
		t.Error(nio.out.get("bla").buf.String())
	}

	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --ghEvalAddMask QUOTE")
	if strings.Compare(strings.TrimSpace(nio.out.first().buf.String()), "QUOTE=\"$<DO\\r\\n\\\"OF>\"\necho ::add-mask::\"$<DO\\r\\n\\\"OF>\"") != 0 {
		t.Error(nio.out.first().buf.String())
	}

}

func genRef(param string) string {
	// unresolved := kvpearl.ParseFuncsAndParams(param)
	ref := []FlatKeyValue{
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
		{
			Key:   "QUOTE",
			Value: "$<DO\r\n\"OF>",
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
	nio, _ := cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --json M@bla N@bla M N QUOTE QUOTE@bla")
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
	if strings.Compare(strings.TrimSpace(nio.out.first().buf.String()), "M=3\nN=4711") != 0 {
		t.Error(nio.out.first().buf.String())
	}
	if strings.Compare(strings.TrimSpace(nio.out.get("bla").buf.String()), "M=3\nN=4711") != 0 {
		t.Error(nio.out.get("bla").buf.String())
	}
}

func TestLsOnlyValue(t *testing.T) {
	createTestData(t)
	// --onlyValue           select device keys
	nio, _ := cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --onlyValue M@bla N@bla N M QUOTE")
	if strings.Compare(strings.TrimSpace(nio.out.first().buf.String()), "3\n4711\n'$<DO\r\n\"OF>'") != 0 {
		t.Error(nio.out.first().buf.String())
	}
	if strings.Compare(strings.TrimSpace(nio.out.get("bla").buf.String()), "3\n4711") != 0 {
		t.Error(nio.out.get("bla").buf.String())
	}
}

func TestLsRawValue(t *testing.T) {
	createTestData(t)
	nio, _ := cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --rawValue M@bla N@bla N M QUOTE")
	if strings.Compare(strings.TrimSpace(nio.out.first().buf.String()), "3\n4711\n$<DO\r\n\"OF>") != 0 {
		t.Error(nio.out.first().buf.String())
	}
	if strings.Compare(strings.TrimSpace(nio.out.get("bla").buf.String()), "3\n4711") != 0 {
		t.Error(nio.out.get("bla").buf.String())
	}
}

func TestLsEvalValue(t *testing.T) {
	createTestData(t)
	nio, _ := cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --evalValue M@bla N@bla N M QUOTE")
	if strings.Compare(strings.TrimSpace(nio.out.first().buf.String()), "\"3\"\n\"4711\"\n\"$<DO\\r\\n\\\"OF>\"") != 0 {
		t.Error(nio.out.first().buf.String())
	}
	if strings.Compare(strings.TrimSpace(nio.out.get("bla").buf.String()), "\"3\"\n\"4711\"") != 0 {
		t.Error(nio.out.get("bla").buf.String())
	}
}

func TestLsShKeyValue(t *testing.T) {
	createTestData(t)
	// --shKeyValue           select device keys
	nio, _ := cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --shKeyValue M@bla N@bla N M QUOTE QUOTE@bla")
	ref := "M=3;\nexport M;\nN=4711;\nexport N;\nQUOTE='$<DO\r\n\"OF>';\nexport QUOTE;"

	if strings.Compare(strings.TrimSpace(nio.out.first().buf.String()), ref) != 0 {
		t.Errorf("%s!=%s", nio.out.first().buf.String(), ref)
	}
	if strings.Compare(strings.TrimSpace(nio.out.get("bla").buf.String()), ref) != 0 {
		t.Errorf("%s!=%s", nio.out.get("bla").buf.String(), ref)
	}
}

func TestLsShEvalKeyValue(t *testing.T) {
	createTestData(t)
	// --shKeyValue           select device keys
	nio, _ := cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --shEvalKeyValue M@bla N@bla N M QUOTE QUOTE@bla")
	ref := "M=\"3\";\nexport M;\nN=\"4711\";\nexport N;\nQUOTE=\"$<DO\\r\\n\\\"OF>\";\nexport QUOTE;"
	if strings.Compare(strings.TrimSpace(nio.out.first().buf.String()), ref) != 0 {
		t.Errorf("%s!=%s", nio.out.first().buf.String(), ref)
	}
	if strings.Compare(strings.TrimSpace(nio.out.get("bla").buf.String()), ref) != 0 {
		t.Errorf("%s!=%s", nio.out.get("bla").buf.String(), ref)
	}
}

func TestLsEmptyTagDefectOutput(t *testing.T) {
	createTestData(t)
	// --shKeyValue           select device keys
	nio, _ := cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --emptyTag")
	if strings.HasPrefix(nio.out.first().buf.String(), "[]:{") {
		t.Error("i forgot something in the source")
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
	if strings.Compare(strings.TrimSpace(nio.out.first().buf.String()), "TOTP=16") != 0 {
		t.Error(nio.out.first().buf.String())
	}
	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --onlyValue TOTP@Totp()")
	if len(strings.TrimSpace(nio.out.first().buf.String())) != 6 {
		t.Error(nio.out.first().buf.String())
	}
	ref := []FlatKeyValue{
		{
			Key:   "TOTP",
			Value: "setit",
			Tags:  []string{},
		},
	}
	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --json TOTP@Totp()")
	var myRef []FlatKeyValue
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

func TestSingleUserPerson(t *testing.T) {
	os.Remove("casket.SingleUser.json")
	os.Remove("gem.SingleUser.json")
	_, err := cmdNeckless(t, "casket --file casket.SingleUser.json create --person --name Person.User1 --email test@test.com")
	if err != nil {
		t.Error(err)
	}
	nio, err := cmdNeckless(t, "casket --file casket.SingleUser.json get")
	if err != nil {
		t.Error(err)
	}
	nio, err = cmdNeckless(t, "gem --casketFile casket.SingleUser.json --file gem.SingleUser.json add", nio.out.first().buf.String())
	if err != nil {
		t.Error(err)
	}
	nio, err = cmdNeckless(nil, "kv --casketFile casket.SingleUser.json --file gem.SingleUser.json add MENO=DOOF")
	if !strings.Contains(err.Error(), "Owners must be found") {
		t.Error("Unknown Error:", nio.err.first().buf.String())
	}
	nio, err = cmdNeckless(nil, "kv --casketFile casket.SingleUser.json --file gem.SingleUser.json ls")
	if err != nil {
		t.Error(err)
	}
	// t.Error(nio.err.first().buf.String())
	// t.Error(nio.out.first().buf.String())
}

func TestSingleUserDevice(t *testing.T) {
	os.Remove("casket.SingleUser.json")
	os.Remove("gem.SingleUser.json")
	_, err := cmdNeckless(t, "casket --file casket.SingleUser.json create --device --name Person.User1 --email test@test.com")
	if err != nil {
		t.Error(err)
	}
	nio, err := cmdNeckless(t, "casket --file casket.SingleUser.json get")
	if err != nil {
		t.Error(err)
	}
	nio, err = cmdNeckless(t, "gem --casketFile casket.SingleUser.json --file gem.SingleUser.json add", nio.out.first().buf.String())
	if err != nil {
		t.Error(err)
	}
	nio, err = cmdNeckless(nil, "kv --casketFile casket.SingleUser.json --file gem.SingleUser.json add MENO=DOOF", nio.out.first().buf.String())
	if !strings.Contains(err.Error(), "you need a private key") {
		t.Error(err)
	}
	nio, err = cmdNeckless(t, "kv --casketFile casket.SingleUser.json --file gem.SingleUser.json ls")
	if err != nil {
		t.Error(err)
	}
	if strings.Compare(nio.out.first().buf.String(), "") != 0 {
		t.Errorf("Compare:%s", nio.out.first().buf.String())
	}
}

func TestSingleUserID(t *testing.T) {
	os.Remove("casket.SingleUser.json")
	os.Remove("gem.SingleUser.json")
	nio, err := cmdNeckless(t, "casket --file casket.SingleUser.json create --device --name Person.User1 --email test@test.com")
	if err != nil {
		t.Error(err)
	}
	nio, err = cmdNeckless(t, "casket --file casket.SingleUser.json get")
	if err != nil {
		t.Error(err)
	}
	var jscreates []member.JsonPublicMember
	err = json.Unmarshal(nio.out.first().buf.Bytes(), &jscreates)
	if err != nil {
		t.Error(err)
	}
	jscreate := jscreates[0]
	nio, err = cmdNeckless(t, fmt.Sprintf(
		"gem --casketFile casket.SingleUser.json --file gem.SingleUser.json --privkeyid %s add --toKeyId %s",
		jscreate.Id, jscreate.Id), nio.out.first().buf.String())
	if err != nil {
		t.Error(err)
	}
	nio, err = cmdNeckless(t, fmt.Sprintf(
		"kv --casketFile casket.SingleUser.json --file gem.SingleUser.json --privkeyid %s add MENO=$D>O<OF",
		jscreate.Id), nio.out.first().buf.String())
	if err != nil {
		t.Error(err)
	}
	nio, err = cmdNeckless(t, fmt.Sprintf(
		"kv --casketFile casket.SingleUser.json --file gem.SingleUser.json --privkeyid %s ls --json",
		jscreate.Id))
	if err != nil {
		t.Error(err)
	}
	var jsadds []FlatKeyValue
	err = json.Unmarshal(nio.out.first().buf.Bytes(), &jsadds)
	if err != nil {
		t.Error(err)
	}
	jsadd := jsadds[0]
	if strings.Compare(jsadd.Key, "MENO") != 0 {
		t.Errorf("Expect MENO got %s", jsadd.Key)
	}
	if strings.Compare(jsadd.Value, "$D>O<OF") != 0 {
		t.Errorf("Expect $DOOF got %s", jsadd.Value)
	}
	if len(jsadd.Tags) != 0 {
		t.Errorf("Expect Tags length %d", len(jsadd.Tags))
	}
	// t.Error(nio.err.first().buf.String())
	// t.Error(nio.out.first().buf.String())

	nio, err = cmdNeckless(t, fmt.Sprintf(
		"kv --casketFile casket.SingleUser.json --file gem.SingleUser.json --privkeyid %s ls --keyValue MENO",
		jscreate.Id))
	if err != nil {
		t.Error(err)
	}
	inStr := string(nio.out.first().buf.Bytes())
	kv := strings.Split(strings.TrimSuffix(inStr, "\n"), "=")
	if len(kv) != 2 {
		t.Errorf("Split should be 2 %s", inStr)
	}
	if strings.Compare(kv[0], "MENO") != 0 {
		t.Errorf("Key should MENO %s", kv[0])
	}
	if strings.Compare(kv[1], "'$D>O<OF'") != 0 {
		t.Errorf("Key should '$D>O<OF' %s", kv[1])
	}
}

func equalFlatKeyValues(t *testing.T, buf *bytes.Buffer, values []FlatKeyValue, caze string) {
	var fromJSON []FlatKeyValue
	if err := json.Unmarshal(buf.Bytes(), &fromJSON); err != nil {
		t.Error(caze, err)
	}
	if len(values) != len(fromJSON) {
		t.Error(caze, "Return should equal length:", len(values), len(fromJSON))
	}
	for i := range values {
		v1 := values[i]
		v2 := fromJSON[i]
		if strings.Compare(v1.Key, v2.Key) != 0 {
			t.Error(caze, "Key should be the same:", i, v1.Key, v2.Key)
		}
		if strings.Compare(v1.Value, v2.Value) != 0 {
			t.Error(caze, "Value should be the same:", i, v1.Value, v2.Value)
		}
	}
}

func TestLsEmptyTag(t *testing.T) {
	createTestData(t)
	// --onlyValue           select device keys
	nio, _ := cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json add M1=34, M2=35[,] M3=36[YU] M4=37[XU] M5=38[,ZU] M6=39[,AU]")

	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --json M[123456]")

	equalFlatKeyValues(t, nio.out.first().buf, []FlatKeyValue{
		{Key: "M1", Value: "34"},
		{Key: "M2", Value: "35"},
		{Key: "M5", Value: "38"},
		{Key: "M6", Value: "39"},
	}, "1")
	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --json --emptyTag M[123456]")
	equalFlatKeyValues(t, nio.out.first().buf, []FlatKeyValue{
		{Key: "M1", Value: "34"},
		{Key: "M2", Value: "35"},
		{Key: "M5", Value: "38"},
		{Key: "M6", Value: "39"},
	}, "2")

	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --json --tag YU --tag XU --tag ZU --tag AU M[123456]")
	equalFlatKeyValues(t, nio.out.first().buf, []FlatKeyValue{
		{Key: "M3", Value: "36"},
		{Key: "M4", Value: "37"},
		{Key: "M5", Value: "38"},
		{Key: "M6", Value: "39"},
	}, "3")

	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls --json --emptyTag --tag YU --tag XU M[123456]")
	equalFlatKeyValues(t, nio.out.first().buf, []FlatKeyValue{
		{Key: "M1", Value: "34"},
		{Key: "M2", Value: "35"},
		{Key: "M3", Value: "36"},
		{Key: "M4", Value: "37"},
		{Key: "M5", Value: "38"},
		{Key: "M6", Value: "39"},
	}, "4")
}
