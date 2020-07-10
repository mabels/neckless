package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"neckless.adviser.com/gem"
	"neckless.adviser.com/kvpearl"
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

func cmdNeckless(t *testing.T, strargs string) (*NecklessIO, error) {
	cargs := strings.Split(strargs, " ")
	nio := NecklessIO{
		// in:  bytes.Buffer,
		out: new(bytes.Buffer),
		err: new(bytes.Buffer),
	}
	args := NecklessArgs{
		Nio: nio,
	}
	// fmt.Println(">>>", cargs)
	_, err := buildArgs(cargs, &args)
	if t != nil && err != nil {
		t.Error(fmt.Sprintf("%s=>%s", strargs, err))
	}
	return &nio, err
}

func TestAddUserToGem(t *testing.T) {
	for i := 0; i < 43; i++ {
		// test for unsort maps quirks
		os.Remove("casket.User1.json")
		nio, _ := cmdNeckless(t, "casket --file casket.User1.json create --person --name Person.User1")
		nio, _ = cmdNeckless(t, "casket --file casket.User1.json create --device --name Device.User1")
		nio, _ = cmdNeckless(t, "casket --file casket.User1.json get")
		nio, _ = cmdNeckless(t, "casket --file casket.User1.json get --outFile device1.pub.json --device")
		os.Remove("casket.User2.json")
		nio, _ = cmdNeckless(t, "casket --file casket.User2.json create --name Person.User2")
		nio, _ = cmdNeckless(t, "casket --file casket.User2.json get --outFile user2.pub.json")
		os.Remove("neckless.shared.json")
		// t.Error("X0.out:", nio.out.String(), "\nX0.err:", nio.err.String())
		nio, _ = cmdNeckless(t, "gem --casketFile casket.User1.json --file neckless.shared.json add")
		nio, _ = cmdNeckless(t, "gem --casketFile casket.User1.json --file neckless.shared.json ls")
		// t.Error("X1.out:", nio.out.String(), "\nX1.err:", nio.err.String())
		nio, _ = cmdNeckless(t, "gem --casketFile casket.User1.json --file neckless.shared.json add --pubFile user2.pub.json")

		nio, _ = cmdNeckless(t, "gem --casketFile casket.User1.json --file neckless.shared.json add --pubFile device1.pub.json")
		// t.Error("X2", nio.out.String(), nio.err.String())
		// t.Error("X2", nio.out.String(), nio.err.String())
		u1nio, _ := cmdNeckless(t, "gem --casketFile casket.User1.json --file neckless.shared.json ls")
		u2nio, _ := cmdNeckless(t, "gem --casketFile casket.User2.json --file neckless.shared.json ls")
		tmp := []gem.JsonGem{}
		// t.Error(string(u1nio.out.Bytes()))
		err := json.Unmarshal(u1nio.out.Bytes(), &tmp)
		if err != nil {
			t.Error(err)
		}
		if len(tmp) != 1 {
			t.Error("not expected len", len(tmp))
		}
		if len(tmp[0].PubKeys) != 3 {
			t.Error("not expected len pubkeys", len(tmp[0].PubKeys))
		}
		if bytes.Compare(u1nio.out.Bytes(), u2nio.out.Bytes()) != 0 {
			t.Error("YYYY", u1nio.out.String(), u1nio.err.String())
			t.Error("XXXX", u2nio.out.String(), u2nio.err.String())
		}
		var toDelId string
		for i := range tmp[0].PubKeys {
			if strings.Compare(string(tmp[0].PubKeys[i].Type), string(member.Device)) == 0 {
				toDelId = tmp[0].PubKeys[i].Id
			}
		}
		nio, _ = cmdNeckless(t, fmt.Sprintf("gem --casketFile casket.User2.json --file neckless.shared.json rm %s", toDelId))
		// t.Error("X2", tmp[0].PubKeys[0].Id, nio.out.String(), nio.err.String())

		u1nio, _ = cmdNeckless(t, "gem --casketFile casket.User1.json --file neckless.shared.json ls")
		u2nio, _ = cmdNeckless(t, "gem --casketFile casket.User2.json --file neckless.shared.json ls")
		tmp = []gem.JsonGem{}
		// t.Error(string(u1nio.out.Bytes()))
		err = json.Unmarshal(u1nio.out.Bytes(), &tmp)
		if err != nil {
			t.Error(err)
		}
		if len(tmp) != 1 {
			t.Error("YYYY", u1nio.out.String(), u1nio.err.String())
			t.Error("XXXX", u2nio.out.String(), u2nio.err.String())
			t.Error("not expected len", len(tmp))
		}
		if len(tmp[0].PubKeys) != 2 {
			t.Error("YYYY", u1nio.out.String(), u1nio.err.String())
			t.Error("XXXX", u2nio.out.String(), u2nio.err.String())
			t.Error("not expected len pubkeys", len(tmp[0].PubKeys))
		}
		if bytes.Compare(u1nio.out.Bytes(), u2nio.out.Bytes()) != 0 {
			t.Error("YYYY", u1nio.out.String(), u1nio.err.String())
			t.Error("XXXX", u2nio.out.String(), u2nio.err.String())
		}
		nio.out.String() // make the compiler happy
	}
	// t.Error("xx", nio.out.String())

}

func TestKvs(t *testing.T) {
	os.Remove("casket.User1.json")
	nio, _ := cmdNeckless(t, "casket --file casket.User1.json create --person --name Person.User1")
	nio, _ = cmdNeckless(t, "casket --file casket.User1.json create --device --name Device.User1")
	nio, _ = cmdNeckless(t, "casket --file casket.User1.json get --outFile device1.pub.json --device")
	os.Remove("neckless.shared.json")
	nio, _ = cmdNeckless(t, "gem --casketFile casket.User1.json --file neckless.shared.json add")
	nio, _ = cmdNeckless(t, "gem --casketFile casket.User1.json --file neckless.shared.json add --pubFile device1.pub.json")

	nio, err := cmdNeckless(nil, "kv --casketFile casket.User1.json --file neckless.shared.json ls")
	if err == nil {
		t.Error("there should be an error")
	}
	if len(nio.out.Bytes()) != 0 {
		t.Error("should be empty")
	}
	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json add M=1 M=2")
	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json add N=4711 M=3")
	nio, _ = cmdNeckless(t, "kv --casketFile casket.User1.json --file neckless.shared.json ls -json")
	var my kvpearl.JsonKVPearl
	json.Unmarshal(nio.out.Bytes(), &my)
	if len(my.Keys) != 2 {
		t.Error("not ok")
	}
	if strings.Compare(my.Keys[0].Key, "M") != 0 {
		t.Error("not M")
	}
	if strings.Compare(my.Keys[0].Values[0].Value, "3") != 0 {
		t.Error("not M")
	}
	if strings.Compare(my.Keys[1].Key, "N") != 0 {
		t.Error("not N")
	}
	if strings.Compare(my.Keys[1].Values[0].Value, "4711") != 0 {
		t.Error("not N")
	}

	// t.Error(string(nio.out.Bytes()))

	nio.out.String() // make the compiler happy

}

// gem add -file User1
// casket get -file User2 | gem add -file User1
// gem ls -file User1
// gem ls -file User2

// import (
// 	"testing"
// )

// func TestCrazyBee(t *testing.T) {
// 	args := CrazyBeeArgs{}
// 	buildArgs([]string{}, &args)
// }

// func TestCrazyBeeCreatePipeLine(t *testing.T) {
// 	args := CrazyBeeArgs{}
// 	buildArgs([]string{"pipeline"}, &args)
// 	if len(args.pipeline.name) == 48 {
// 		t.Errorf("pipeline.name not ok: %s", args.pipeline.name)
// 	}
// }

// func TestCrazyBeeCreatePipeLineSet(t *testing.T) {
// 	args := CrazyBeeArgs{}
// 	buildArgs([]string{
// 		"pipeline",
// 		"--name", "Test.Name",
// 	}, &args)
// 	if args.pipeline.name != "Test.Name" {
// 		t.Errorf("pipeline.name not ok: %s", args.pipeline.name)
// 	}
// }
