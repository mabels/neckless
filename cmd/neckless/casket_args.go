package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/peterbourgon/ff/v2/ffcli"
	"neckless.adviser.com/casket"
	"neckless.adviser.com/member"
)

type CasketRmArgs struct {
	DryRun *bool
}
type CasketCreateArgs struct {
	Name       string
	DryRun     *bool
	DeviceName string
	DeviceType *bool
	PersonType *bool
	ValidUntil string
}
type CasketGetArgs struct {
	PubFile    string
	PrivateKey *bool
	Device     *bool
	Person     *bool
	KeyValue   *bool
}
type CasketArgs struct {
	Fname  string
	create CasketCreateArgs
	rm     CasketRmArgs
	get    CasketGetArgs
}

func casketCreateArgs(arg *CasketArgs) *ffcli.Command {
	flags := flag.NewFlagSet("create", flag.ExitOnError)
	defaultName := uuid.New().String()
	flags.StringVar(&arg.create.Name, "name", defaultName, "name of the key")
	flags.StringVar(&arg.create.DeviceName, "deviceName", "", "name of the key")

	arg.create.DryRun = flags.Bool("dryRun", false, "set the dryrun flag")

	arg.create.DeviceType = flags.Bool("device", false, "is device")
	arg.create.PersonType = flags.Bool("person", false, "is person")

	flags.StringVar(&arg.create.ValidUntil, "valid", "", "not impl yet")

	return &ffcli.Command{
		Name:       "create",
		ShortUsage: "manage a casket secrets",
		ShortHelp:  "manage casket secrets",

		LongHelp: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
		FlagSet: flags,
		Exec: func(ctx context.Context, args []string) error {
			typ := member.Person
			if *arg.create.DeviceType {
				typ = member.Device
			}
			_, pk, err := casket.Create(casket.CreateArg{
				DryRun: *arg.create.DryRun,
				Fname:  &arg.Fname,

				MemberArg: member.MemberArg{
					Id:     uuid.New().String(),
					Type:   typ,
					Name:   arg.create.Name,
					Device: arg.create.DeviceName,
				},
			})
			if err != nil {
				log.Fatal(err)
				return nil
			}
			js, err := json.MarshalIndent(pk.AsJson(), "", "  ")
			fmt.Println(string(js))
			return nil
		},
	}
}

func casketGetArgs(arg *CasketArgs) *ffcli.Command {
	flags := flag.NewFlagSet("get", flag.ExitOnError)
	flags.StringVar(&arg.get.PubFile, "outFile", "", "filename to write")
	arg.get.PrivateKey = flags.Bool("privateKey", false, "set export to private key")
	arg.get.Person = flags.Bool("person", false, "select person keys")
	arg.get.Device = flags.Bool("device", false, "select device keys")
	arg.get.KeyValue = flags.Bool("keyValue", false, "output as keyvalue")
	return &ffcli.Command{
		Name:       "get",
		ShortUsage: "get casket secrets",
		ShortHelp:  "get casket secrets",

		LongHelp: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
		FlagSet: flags,
		Exec: func(_ context.Context, args []string) error {
			c, err := casket.Ls(arg.Fname)
			if err != nil {
				log.Fatal(err)
			}
			pkms := c.AsPrivateMembers()
			mtyps := []member.MemberType{}
			if *arg.get.Person {
				mtyps = append(mtyps, member.Person)
			}
			if *arg.get.Device {
				mtyps = append(mtyps, member.Device)
			}
			if len(mtyps) == 0 {
				mtyps = append(mtyps, member.Person)
			}
			// fmt.Println("xxxx", mtyps)
			pkms = member.FilterById(member.FilterByType(pkms, mtyps...), args...)
			var js []byte
			var err1 error
			if *arg.get.PrivateKey {
				jspkms := member.ToJsonPrivateMember(pkms...)
				js, err1 = json.MarshalIndent(jspkms, "", "  ")

			} else {
				jspkms := member.ToJsonPublicMember(pkms...)
				js, err1 = json.MarshalIndent(jspkms, "", "  ")
			}
			if err1 != nil {
				log.Fatal(err)
			}
			if *arg.get.KeyValue {
				for i := range pkms {
					fmt.Println(string(pkms[i].PrivateKey.Marshal()))
				}
			} else {
				fmt.Println(string(js))
			}
			return nil
		},
	}
}

func casketLsArgs(arg *CasketArgs) *ffcli.Command {
	flags := flag.NewFlagSet("ls", flag.ExitOnError)
	// flags.Args()
	return &ffcli.Command{
		Name:       "ls",
		ShortUsage: "list casket secrets",
		ShortHelp:  "list casket secrets",

		LongHelp: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
		FlagSet: flags,
		Exec: func(_ context.Context, _ []string) error {
			c, err := casket.Ls(arg.Fname)
			if err != nil {
				log.Fatal(err)
			}
			js, err := json.MarshalIndent(c.AsJson(), "", "  ")
			fmt.Println(string(js))
			return nil
		},
	}
}

func casketRmArgs(arg *CasketArgs) *ffcli.Command {
	flags := flag.NewFlagSet("rm", flag.ExitOnError)
	arg.rm.DryRun = flags.Bool("dryRun", false, "set the dryrun flag")
	return &ffcli.Command{
		Name:       "rm",
		ShortUsage: "manage a casket secrets",
		ShortHelp:  "manage casket secrets",

		LongHelp: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
		FlagSet: flags,
		Exec: func(_ context.Context, args []string) error {
			_, pks, err := casket.Rm(casket.RmArg{
				Ids:    args,
				DryRun: *arg.rm.DryRun,
				Fname:  &arg.Fname,
			})
			if err != nil {
				log.Fatal(err)
			}
			out := make([]*member.JsonPrivateMember, len(pks))
			for i := range pks {
				out[i] = pks[i].AsJson()
			}
			js, err := json.MarshalIndent(out, "", "  ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(js))
			return nil
		},
	}
}
func casketArgs(arg *CasketArgs) *ffcli.Command {
	flags := flag.NewFlagSet("casket", flag.ExitOnError)
	homeDir := os.Getenv("HOME")
	flags.StringVar(&arg.Fname, "file",
		fmt.Sprintf("%s/.neckless/casket.json", homeDir), "filename of the casket")

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
		Name:       "casket",
		ShortUsage: "manage a casket secrets",
		ShortHelp:  "manage casket secrets",

		LongHelp: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
		FlagSet: flags,
		Subcommands: []*ffcli.Command{
			casketCreateArgs(arg),
			casketLsArgs(arg),
			casketGetArgs(arg),
			casketRmArgs(arg),
		},
		Exec: func(context.Context, []string) error {
			fmt.Println("Casket-Hello")
			return nil
		},
	}
}
