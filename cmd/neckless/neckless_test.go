package main

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
