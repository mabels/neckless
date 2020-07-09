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

	"github.com/peterbourgon/ff/v2/ffcli"
	"neckless.adviser.com/casket"
	"neckless.adviser.com/gem"
	"neckless.adviser.com/member"
	"neckless.adviser.com/neckless"
	"neckless.adviser.com/pearl"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	// change this, this is just can example to satisfy the interface
	return strings.Join(*i, ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, strings.TrimSpace(value))
	return nil
}

type GemAddArgs struct {
	PubFile  string
	Device   *bool
	Person   *bool
	KeyValue *bool
}

type GemArgs struct {
	Fname       string
	CasketFname string
	PrivKeyIds  arrayFlags
	Add         GemAddArgs
}

func updateGem(myGem *gem.Gem, pkms []*member.PrivateMember, jpms []member.JsonPublicMember) (*pearl.Pearl, error) {
	// pms := make([]*member.PublicMember, len(jpms))
	for i := range jpms {
		pm, err := member.JsToPublicMember(&jpms[i])
		if err != nil {
			return nil, err
		}
		myGem.Add(pm)
		// pms[i] = pm
	}
	for j := range pkms {
		myGem.Add(pkms[j].Public())
	}
	mo := pearl.PearlOwner{
		Signer: &pkms[0].PrivateKey,
		Owners: member.ToPublicKeys(myGem.LsByType(member.Person)),
	}
	p, err := myGem.ClosePearl(&mo)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func gemAddArgs(arg *GemArgs) *ffcli.Command {
	flags := flag.NewFlagSet("gem.add", flag.ExitOnError)
	// homeDir := os.Getenv("HOME")
	flags.StringVar(&arg.Add.PubFile, "pubFile", "stdin", "the pubMemberFile to add")
	arg.Add.Person = flags.Bool("person", false, "select person keys")
	arg.Add.Device = flags.Bool("device", false, "select device keys")

	return &ffcli.Command{
		Name:       "add",
		ShortUsage: "manage a gem stone in neckless",
		ShortHelp:  "manage gem stone",

		LongHelp: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
		FlagSet:     flags,
		Subcommands: []*ffcli.Command{},
		Exec: func(context.Context, []string) error {
			casket, err := casket.Ls(arg.CasketFname)
			if err != nil {
				log.Fatal(err)
				return nil
			}
			mtyps := []member.MemberType{}
			if *arg.Add.Person {
				mtyps = append(mtyps, member.Person)
			}
			if *arg.Add.Device {
				mtyps = append(mtyps, member.Device)
			}
			if len(mtyps) == 0 {
				mtyps = append(mtyps, member.Person)
			}
			pkms := member.FilterByType(member.FilterById(casket.AsPrivateMembers(), arg.PrivKeyIds...), mtyps...)
			if len(pkms) == 0 {
				log.Fatal(errors.New("you need a private key"))
			}
			var jsStr []byte
			if strings.Compare(arg.Add.PubFile, "stdin") == 0 {
				jsStr, err = ioutil.ReadAll(os.Stdin)
			} else {
				jsStr, err = ioutil.ReadFile(arg.Add.PubFile)
			}
			if err != nil {
				log.Fatal(err)
				return nil
			}
			pubMembers := []member.JsonPublicMember{}
			err = json.Unmarshal(jsStr, &pubMembers)
			if err != nil {
				log.Fatal(err)
				return nil
			}
			nl := neckless.GetAndOpen(arg.Fname)
			gems := nl.FilterByType(gem.Type)
			for i := range gems {
				tmp := gems[i]
				myGem, err := gem.OpenPearl(member.ToPrivateKeys(pkms), tmp)
				if err != nil {
					log.Fatal(err)
					return nil
				}
				prl, err := updateGem(myGem, pkms, pubMembers)
				if err != nil {
					log.Fatal(err)
					return nil
				}
				nl.Reset(prl, gems[i].FingerPrint)
			}
			if len(gems) == 0 {
				myGem := gem.Create()
				prl, err := updateGem(myGem, pkms, pubMembers)
				if err != nil {
					log.Fatal(err)
					return nil
				}
				nl.Reset(prl)
			}
			nl.Save(arg.Fname)
			jsStr, err = json.MarshalIndent(pubMembers, "", "  ")
			fmt.Println(string(jsStr))
			return nil
		},
	}
}
func gemRmArgs(arg *GemArgs) *ffcli.Command {
	flags := flag.NewFlagSet("gem.rm", flag.ExitOnError)
	// homeDir := os.Getenv("HOME")
	flags.StringVar(&arg.Fname, "file", ".neckless", "the neckless file")

	return &ffcli.Command{
		Name:       "rm",
		ShortUsage: "manage a gem stone in neckless",
		ShortHelp:  "manage gem stone",

		LongHelp: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
		FlagSet:     flags,
		Subcommands: []*ffcli.Command{},
	}
}
func gemLsArgs(arg *GemArgs) *ffcli.Command {
	flags := flag.NewFlagSet("gem.ls", flag.ExitOnError)
	// homeDir := os.Getenv("HOME")

	return &ffcli.Command{
		Name:       "ls",
		ShortUsage: "manage a gem stone in neckless",
		ShortHelp:  "manage gem stone",

		LongHelp: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
		FlagSet:     flags,
		Subcommands: []*ffcli.Command{},
		Exec: func(context.Context, []string) error {
			casket, err := casket.Ls(arg.CasketFname)
			if err != nil {
				log.Fatal(err)
				return nil
			}
			mtyps := []member.MemberType{}
			if *arg.Add.Person {
				mtyps = append(mtyps, member.Person)
			}
			if *arg.Add.Device {
				mtyps = append(mtyps, member.Device)
			}
			if len(mtyps) == 0 {
				mtyps = append(mtyps, member.Person)
			}
			pkms := member.FilterByType(member.FilterById(casket.AsPrivateMembers(), arg.PrivKeyIds...), mtyps...)
			if len(pkms) == 0 {
				log.Fatal(errors.New("you need a private key"))
			}
			nl := neckless.GetAndOpen(arg.Fname)
			gems := nl.FilterByType(gem.Type)
			out := []*gem.Gem{}
			for i := range gems {
				tmp := gems[i]
				openGem, err := gem.OpenPearl(member.ToPrivateKeys(pkms), tmp)
				if err != nil {
					log.Fatal(err)
					return nil
				}
				out = append(out, openGem)
			}
			jsStr, err := json.MarshalIndent(gem.ToJsonGems(out...), "", "  ")
			if err != nil {
				log.Fatal(err)
				return nil
			}
			fmt.Println(string(jsStr))
			return nil
		},
	}
}

func gemArgs(arg *GemArgs) *ffcli.Command {
	flags := flag.NewFlagSet("gem", flag.ExitOnError)
	// homeDir := os.Getenv("HOME")
	flags.StringVar(&arg.Fname, "file", ".neckless", "the neckless file")
	homeDir := os.Getenv("HOME")
	flags.StringVar(&arg.CasketFname, "casketFile",
		fmt.Sprintf("%s/.neckless/casket.json", homeDir), "filename of the casket")
	arg.PrivKeyIds = arrayFlags{}
	flags.Var(&arg.PrivKeyIds, "privkeyid", "the neckless file")

	return &ffcli.Command{
		Name:       "gem",
		ShortUsage: "manage a gem stone in neckless",
		ShortHelp:  "manage gem stone",

		LongHelp: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
		FlagSet: flags,
		Subcommands: []*ffcli.Command{
			gemAddArgs(arg),
			gemRmArgs(arg),
			gemLsArgs(arg),
		},
	}
}
