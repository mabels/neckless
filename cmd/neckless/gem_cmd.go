package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/peterbourgon/ff/v2/ffcli"
	"neckless.adviser.com/gem"
	"neckless.adviser.com/member"
	"neckless.adviser.com/necklace"
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
	ToKeyIds arrayFlags
}
type GemLsArgs struct {
	Device *bool
	Person *bool
}

type GemArgs struct {
	Fname       string
	CasketFname string
	PrivKeyIds  arrayFlags
	Add         GemAddArgs
	Ls          GemLsArgs
}

func GetGems(pkms []*member.PrivateMember, nl *necklace.Necklace) ([]*gem.Gem, []error) {
	closedGems := nl.FilterByType(gem.Type)
	out := []*gem.Gem{}
	errs := []error{}
	for i := range closedGems {
		tmp := closedGems[i]
		openGem, err := gem.OpenPearl(member.ToPrivateKeys(pkms), tmp)
		if err != nil {
			errs = append(errs, err)
		} else {
			out = append(out, openGem)
		}
	}
	return out, errs
}

func updateGem(myGem *gem.Gem, pkms []*member.PrivateMember, jpms []member.JsonPublicMember, toIds ...string) (*pearl.Pearl, error) {
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
	var pms []*member.PublicMember
	if len(toIds) > 0 {
		pms = myGem.Ls(toIds...)
	} else {
		pms = myGem.LsByType(member.Person)
	}
	mo := pearl.PearlOwner{
		Signer: &pkms[0].PrivateKey,
		Owners: member.ToPublicKeys(pms),
	}
	p, err := myGem.ClosePearl(&mo)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func gemAddCmd(arg *NecklessArgs) *ffcli.Command {
	flags := flag.NewFlagSet("gem.add", flag.ExitOnError)
	// homeDir := os.Getenv("HOME")
	flags.StringVar(&arg.Gems.Add.PubFile, "pubFile", "stdin", "the pubMemberFile to add")
	arg.Gems.Add.Person = flags.Bool("person", false, "select person keys")
	arg.Gems.Add.Device = flags.Bool("device", false, "select device keys")
	flags.Var(&arg.Gems.Add.ToKeyIds, "toKeyId", "the neckless file")

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
			pkms, err := GetPkms(GetPkmsArgs{
				casketFname: arg.Gems.CasketFname,
				privIds:     arg.Gems.PrivKeyIds,
				person:      *arg.Gems.Add.Person,
				device:      *arg.Gems.Add.Device})
			if err != nil {
				return err
			}
			var jsStr []byte
			if strings.Compare(arg.Gems.Add.PubFile, "stdin") == 0 {
				jsStr, err = ioutil.ReadAll(os.Stdin)
			} else {
				jsStr, err = ioutil.ReadFile(arg.Gems.Add.PubFile)
			}
			if err != nil {
				return err
			}
			// fmt.Fprintln(arg.Nio.err, "-1:", string(jsStr))
			pubMembers := []member.JsonPublicMember{}
			if len(jsStr) != 0 {
				err = json.Unmarshal(jsStr, &pubMembers)
				if err != nil {
					return err
				}
			}
			// fmt.Fprintln(arg.Nio.err, "-2")
			nl, _ := necklace.Read(arg.Gems.Fname)
			// fmt.Fprintln(arg.Nio.err, "-3")
			gems, _ := GetGems(pkms, &nl)
			for i := range gems {
				prl, err := updateGem(gems[i], pkms, pubMembers, arg.Gems.Add.ToKeyIds...)
				if err != nil {
					return err
				}
				nl.Reset(prl, gems[i].Pearl.Closed.FingerPrint)
			}
			if len(gems) == 0 {
				myGem := gem.Create()
				prl, err := updateGem(myGem, pkms, pubMembers)
				if err != nil {
					return err
				}
				nl.Reset(prl)
			}
			nl.Save(arg.Gems.Fname)
			jsStr, err = json.MarshalIndent(pubMembers, "", "  ")
			fmt.Fprintln(arg.Nio.out, string(jsStr))
			return err
		},
	}
}
func gemRmCmd(arg *NecklessArgs) *ffcli.Command {
	flags := flag.NewFlagSet("gem.rm", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "rm",
		ShortUsage: "manage a gem stone in neckless",
		ShortHelp:  "manage gem stone",

		LongHelp: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
		FlagSet:     flags,
		Subcommands: []*ffcli.Command{},
		Exec: func(_ context.Context, args []string) error {
			pkms, err := GetPkms(GetPkmsArgs{
				casketFname: arg.Gems.CasketFname,
				privIds:     arg.Gems.PrivKeyIds})
			nl, _ := necklace.Read(arg.Gems.Fname)
			gems, _ := GetGems(pkms, &nl)
			// fmt.Fprintln(arg.Nio.err, pkms[0].Id)
			myGems := []*gem.JsonGem{}
			for i := range gems {
				myGem := gems[i]
				myGem.Rm(args...)
				mo := pearl.PearlOwner{
					Signer: &pkms[0].PrivateKey,
					Owners: member.ToPublicKeys(myGem.LsByType(member.Person)),
				}
				p, err := myGem.ClosePearl(&mo)
				if err != nil {
					return err
				}
				nl.Reset(p, myGem.Pearl.Closed.FingerPrint)
				myGems = append(myGems, myGem.AsJSON())
			}
			nl.Save(arg.Gems.Fname)
			jsStr, err := json.MarshalIndent(myGems, "", "  ")
			fmt.Fprintln(arg.Nio.out, string(jsStr))
			return err
		},
	}
}
func gemLsCmd(arg *NecklessArgs) *ffcli.Command {
	flags := flag.NewFlagSet("gem.ls", flag.ExitOnError)
	arg.Gems.Ls.Person = flags.Bool("person", false, "select person keys")
	arg.Gems.Ls.Device = flags.Bool("device", false, "select device keys")

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
			pkms, err := GetPkms(GetPkmsArgs{
				casketFname: arg.Gems.CasketFname,
				privIds:     arg.Gems.PrivKeyIds,
				person:      *arg.Gems.Ls.Person,
				device:      *arg.Gems.Ls.Device})
			if err != nil {
				return err
			}
			nl, _ := necklace.Read(arg.Gems.Fname)
			fmt.Fprintln(arg.Nio.err, pkms[0].Id)
			gems, _ := GetGems(pkms, &nl)
			jsStr, err := json.MarshalIndent(gem.ToJsonGems(gems...), "", "  ")
			if err != nil {
				return err
			}
			fmt.Fprintln(arg.Nio.out, string(jsStr))
			return nil
		},
	}
}

func gemCmd(arg *NecklessArgs) *ffcli.Command {
	flags := flag.NewFlagSet("gem", flag.ExitOnError)
	// homeDir := os.Getenv("HOME")
	necklessFile := findFile(".neckless")
	flags.StringVar(&arg.Gems.Fname, "file", necklessFile, "the neckless file")
	homeDir := os.Getenv("HOME")
	flags.StringVar(&arg.Gems.CasketFname, "casketFile",
		fmt.Sprintf("%s/.neckless/casket.json", homeDir), "filename of the casket")
	arg.Gems.PrivKeyIds = arrayFlags{}
	flags.Var(&arg.Gems.PrivKeyIds, "privkeyid", "the neckless file")

	return &ffcli.Command{
		Name:       "gem",
		ShortUsage: "manage a gem stone in neckless",
		ShortHelp:  "manage gem stone",

		LongHelp: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
		FlagSet: flags,
		Subcommands: []*ffcli.Command{
			gemAddCmd(arg),
			gemRmCmd(arg),
			gemLsCmd(arg),
		},
		Exec: func(context.Context, []string) error { return flag.ErrHelp },
	}
}
