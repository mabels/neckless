package main

import (
	"flag"
	"strings"

	"github.com/peterbourgon/ff/v2/ffcli"
)

type KeyValueArgs struct {
}

func keyValueArgs(arg *KeyValueArgs) *ffcli.Command {
	flags := flag.NewFlagSet("kv", flag.ExitOnError)
	// homeDir := os.Getenv("HOME")
	// flags.StringVar(&arg.create.Fname, "file",
	// 	fmt.Sprintf("%s/.neckless/casket.json", homeDir), "filename of the casket")

	// fmt.Println("xXxxx", args.casket.create.MemberArg.Type)
	// var typ string
	// flags.StringVar(&arg.create.Type, "type", "xxx", "type of the key")
	// fmt.Println(typ)
	// if strings.Compare(typ, string(member.Person)) == 0 {
	// 	arg.create.Type = member.Person
	// }
	// if strings.Compare(typ, string(member.Device)) == 0 {
	// 	arg.create.Type = member.Device
	// }
	// 	// defaultDur, _ := time.ParseDuration("5y")
	// 	// flags.Duration("valid", defaultDur, "defined the validilty of the key")
	return &ffcli.Command{
		Name:       "kv",
		ShortUsage: "manage a key value secrets",
		ShortHelp:  "manage key value secrets",

		LongHelp: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
		FlagSet:     flags,
		Subcommands: []*ffcli.Command{},
	}
}
