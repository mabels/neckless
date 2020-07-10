package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/peterbourgon/ff/v2/ffcli"
	"neckless.adviser.com/casket"
	"neckless.adviser.com/key"
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

type GetPkmsArgs struct {
	casketFname string
	privEnvName string
	privKeyVal  string
	privIds     []string
	person      bool
	device      bool
}

func GetPkms(a GetPkmsArgs) ([]*member.PrivateMember, error) {
	casket, err := casket.Ls(a.casketFname)
	if err != nil {
		return nil, err
	}
	mtyps := []member.MemberType{}
	if a.person {
		mtyps = append(mtyps, member.Person)
	}
	if a.device {
		mtyps = append(mtyps, member.Device)
	}
	pkms := []*member.PrivateMember{}
	if len(os.Getenv(a.privEnvName)) > 0 || len(a.privKeyVal) > 0 {
		strPks := []string{}
		if len(a.privKeyVal) > 0 {
			strPks = append(strPks, a.privKeyVal)
		}
		envPkStr := os.Getenv(a.privEnvName)
		if len(envPkStr) > 0 {
			strPks = append(strPks, envPkStr)
		}
		for i := range strPks {
			pk, _, err := key.FromText(strPks[i], "from-cmd-line")
			if err != nil {
				return nil, err
			}
			if pk == nil {
				return nil, errors.New(fmt.Sprintf("we need a private key passed:%s", strPks[i]))
			}
			pkm, err := member.MakePrivateMember(&member.PrivateMemberArg{
				Member: member.MemberArg{
					Id:   "from-cmd-line",
					Name: "from-cmd-line",
				},
				PrivateKey: pk,
			})
			if err != nil {
				return nil, err
			}
			pkms = append(pkms, pkm)
		}
	}
	pkms = append(pkms, member.FilterByType(member.FilterById(casket.AsPrivateMembers(), a.privIds...), mtyps...)...)
	if len(pkms) == 0 {
		return nil, errors.New("you need a private key")
	}
	return pkms, nil
}

func casketCreateArgs(arg *NecklessArgs) *ffcli.Command {
	flags := flag.NewFlagSet("create", flag.ExitOnError)
	defaultName := uuid.New().String()
	flags.StringVar(&arg.Casket.create.Name, "name", defaultName, "name of the key")
	flags.StringVar(&arg.Casket.create.DeviceName, "deviceName", "", "name of the key")

	arg.Casket.create.DryRun = flags.Bool("dryRun", false, "set the dryrun flag")

	arg.Casket.create.DeviceType = flags.Bool("device", false, "is device")
	arg.Casket.create.PersonType = flags.Bool("person", false, "is person")

	flags.StringVar(&arg.Casket.create.ValidUntil, "valid", "", "not impl yet")

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
			if *arg.Casket.create.DeviceType {
				typ = member.Device
			}
			_, pk, err := casket.Create(casket.CreateArg{
				DryRun: *arg.Casket.create.DryRun,
				Fname:  &arg.Casket.Fname,

				MemberArg: member.MemberArg{
					Id:     uuid.New().String(),
					Type:   typ,
					Name:   arg.Casket.create.Name,
					Device: arg.Casket.create.DeviceName,
				},
			})
			if err != nil {
				log.Fatal(err)
				return nil
			}
			js, err := json.MarshalIndent(pk.AsJson(), "", "  ")
			fmt.Fprintln(arg.Nio.out, string(js))
			return nil
		},
	}
}

func casketGetArgs(arg *NecklessArgs) *ffcli.Command {
	flags := flag.NewFlagSet("get", flag.ExitOnError)
	flags.StringVar(&arg.Casket.get.PubFile, "outFile", "", "filename to write")
	arg.Casket.get.PrivateKey = flags.Bool("privateKey", false, "set export to private key")
	arg.Casket.get.Person = flags.Bool("person", false, "select person keys")
	arg.Casket.get.Device = flags.Bool("device", false, "select device keys")
	arg.Casket.get.KeyValue = flags.Bool("keyValue", false, "output as keyvalue")
	return &ffcli.Command{
		Name:       "get",
		ShortUsage: "get casket secrets",
		ShortHelp:  "get casket secrets",

		LongHelp: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
		FlagSet: flags,
		Exec: func(_ context.Context, args []string) error {
			pkms, err := GetPkms(GetPkmsArgs{
				casketFname: arg.Casket.Fname,
				privIds:     args,
				privEnvName: arg.Kvs.PrivKeyEnv,
				privKeyVal:  arg.Kvs.PrivKeyVal,
				person:      *arg.Casket.get.Person,
				device:      *arg.Casket.get.Device,
			})
			if err != nil {
				return err
			}

			// c, err := casket.Ls(arg.Casket.Fname)
			// pkms := c.AsPrivateMembers()
			// mtyps := []member.MemberType{}
			// if *arg.Casket.get.Person {
			// 	mtyps = append(mtyps, member.Person)
			// }
			// if *arg.Casket.get.Device {
			// 	mtyps = append(mtyps, member.Device)
			// }
			// // if len(mtyps) == 0 {
			// // mtyps = append(mtyps, member.Person)
			// // }
			// // fmt.Fprintln(arg.Nio.err, "xxxx", mtyps, args, len(pkms))
			// // fmt.Fprintln(arg.Nio.err, len(member.FilterByType(pkms, mtyps...)))
			// // fmt.Fprintln(arg.Nio.err, args[0], len(member.FilterById(pkms, args...)))
			// pkms = member.FilterById(member.FilterByType(pkms, mtyps...), args...)
			// // fmt.Fprintln(arg.Nio.err, "yyy", mtyps, args, len(pkms))
			var js []byte
			var err1 error
			if *arg.Casket.get.PrivateKey {
				jspkms := member.ToJsonPrivateMember(pkms...)
				js, err1 = json.MarshalIndent(jspkms, "", "  ")

			} else {
				jspkms := member.ToJsonPublicMember(pkms...)
				js, err1 = json.MarshalIndent(jspkms, "", "  ")
			}
			if err1 != nil {
				log.Fatal(err)
			}
			if *arg.Casket.get.KeyValue {
				out := make([]string, len(pkms))
				for i := range pkms {
					out[i] = pkms[i].PrivateKey.Marshal()
				}
				if len(arg.Casket.get.PubFile) > 0 {
					ioutil.WriteFile(arg.Casket.get.PubFile, []byte(strings.Join(out, "\n")), 0644)
				} else {
					fmt.Fprintln(arg.Nio.out, strings.Join(out, "\n"))
				}
			} else {
				if len(arg.Casket.get.PubFile) > 0 {
					ioutil.WriteFile(arg.Casket.get.PubFile, js, 0644)
				} else {
					fmt.Fprintln(arg.Nio.out, string(js))
				}
			}
			return nil
		},
	}
}

// func casketLsArgs(arg *NecklessArgs) *ffcli.Command {
// 	flags := flag.NewFlagSet("ls", flag.ExitOnError)
// 	// flags.Args()
// 	return &ffcli.Command{
// 		Name:       "ls",
// 		ShortUsage: "list casket secrets",
// 		ShortHelp:  "list casket secrets",

// 		LongHelp: strings.TrimSpace(`
// 	    This command is used to create and add user to the pipeline secret
// 	    `),
// 		FlagSet: flags,
// 		Exec: func(_ context.Context, _ []string) error {

// 			c, err := casket.Ls(arg.Casket.Fname)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			js, err := json.MarshalIndent(c.AsJson(), "", "  ")
// 			fmt.Fprintln(arg.Nio.out, string(js))
// 			return nil
// 		},
// 	}
// }

func casketRmArgs(arg *NecklessArgs) *ffcli.Command {
	flags := flag.NewFlagSet("rm", flag.ExitOnError)
	arg.Casket.rm.DryRun = flags.Bool("dryRun", false, "set the dryrun flag")
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
				DryRun: *arg.Casket.rm.DryRun,
				Fname:  &arg.Casket.Fname,
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
			fmt.Fprintln(arg.Nio.out, string(js))
			return nil
		},
	}
}
func casketArgs(arg *NecklessArgs) *ffcli.Command {
	flags := flag.NewFlagSet("casket", flag.ExitOnError)
	homeDir := os.Getenv("HOME")
	flags.StringVar(&arg.Casket.Fname, "file",
		fmt.Sprintf("%s/.neckless/casket.json", homeDir), "filename of the casket")

	// fmt.Fprintln("xXxxx", args.casket.create.MemberArg.Type)
	// var typ string
	// flags.StringVar(&arg.create.Type, "type", "xxx", "type of the key")
	// fmt.Fprintln(typ)
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
			casketGetArgs(arg),
			casketRmArgs(arg),
		},
		Exec: func(context.Context, []string) error { return flag.ErrHelp },
		/*
			Exec: func(context.Context, []string) error {
				fmt.Fprintln("Casket-Hello")
				return nil
			},
		*/
	}
}
