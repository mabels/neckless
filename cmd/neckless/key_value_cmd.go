package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/peterbourgon/ff/v2/ffcli"
	"neckless.adviser.com/kvpearl"
	"neckless.adviser.com/member"
	"neckless.adviser.com/necklace"
	"neckless.adviser.com/pearl"
)

type KeyValueLsArgs struct {
	json       *bool
	keyValue   *bool
	shKeyValue *bool
	ghAddMask  *bool
	onlyValue  *bool
	tags       arrayFlags
}

type KeyValueArgs struct {
	Fname       string
	CasketFname string
	PrivKeyIds  arrayFlags
	PrivKeyEnv  string
	PrivKeyVal  string
	Ls          KeyValueLsArgs
}

func toKV(kvp *kvpearl.KVPearl, args []string) (*kvpearl.KVPearl, []error) {
	// ret := []kvpearl.Key{}
	errs := []error{}
	for i := range args {
		_, err := kvp.Parse(args[i])
		if err != nil {
			errs = append(errs, err)
		}
	}
	return kvp, errs
}
func kvAddCmd(arg *NecklessArgs) *ffcli.Command {
	flags := flag.NewFlagSet("kv.Add", flag.ExitOnError)
	// Args VAL=Wert[TAGS,TAGS]
	return &ffcli.Command{
		Name:       "add",
		ShortUsage: "manage a key value secrets args <KEY=VAL[TAGS,]>",
		ShortHelp:  "manage key value secrets",

		LongHelp: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
		FlagSet:     flags,
		Subcommands: []*ffcli.Command{},
		Exec: func(_ context.Context, args []string) error {
			if len(args) == 0 {
				return flag.ErrHelp
			}
			kvp := kvpearl.Create()
			_, errs := toKV(kvp, args)
			for i := range errs {
				fmt.Fprintln(arg.Nio.err, errs[i])
			}
			jsStr, err := json.MarshalIndent(kvp.AsJSON(), "", "  ")
			if err != nil {
				return err
			}
			fmt.Fprintln(arg.Nio.out, string(jsStr))
			nl, _ := necklace.Read(arg.Kvs.Fname)
			pkms, err := GetPkms(GetPkmsArgs{
				casketFname: arg.Kvs.CasketFname,
				privIds:     arg.Kvs.PrivKeyIds,
				person:      true,
				device:      false,
			})
			if err != nil {
				return nil
			}
			gems, _ := GetGems(pkms, &nl)
			// fmt.Fprintln(arg.Nio.err, len(gems), pkms[0].Id)
			for i := range gems {
				mo := pearl.PearlOwner{
					Signer: &pkms[0].PrivateKey,
					Owners: member.ToPublicKeys(gems[i].LsByType(member.Device)),
				}
				fmt.Fprintf(arg.Nio.err, "%s:%s\n", pkms[0].Id, gems[i].LsByType(member.Device)[0].Id)
				p, err := kvp.ClosePearl(&mo)
				if err != nil {
					return err
				}
				nl.Reset(p)
			}
			_, err = nl.Save(arg.Kvs.Fname)
			return err
		},
	}
}

func kvLsCmd(arg *NecklessArgs) *ffcli.Command {
	flags := flag.NewFlagSet("kv.Ls", flag.ExitOnError)
	arg.Kvs.Ls.json = flags.Bool("json", false, "select device keys")
	arg.Kvs.Ls.keyValue = flags.Bool("keyValue", true, "select device keys")
	arg.Kvs.Ls.onlyValue = flags.Bool("onlyValue", false, "select device keys")
	arg.Kvs.Ls.shKeyValue = flags.Bool("shKeyValue", false, "select device keys")
	arg.Kvs.Ls.ghAddMask = flags.Bool("ghAddMask", false, "set Value as github mask")
	flags.Var(&arg.Kvs.Ls.tags, "tag", "list of tags to filter")
	return &ffcli.Command{
		Name:       "ls",
		ShortUsage: "manage a key value secrets",
		ShortHelp:  "manage key value secrets",

		LongHelp: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
		FlagSet:     flags,
		Subcommands: []*ffcli.Command{},
		Exec: func(_ context.Context, args []string) error {
			nl, _ := necklace.Read(arg.Kvs.Fname)
			closedKvps := nl.FilterByType(kvpearl.Type)
			pkms, err := GetPkms(GetPkmsArgs{
				casketFname: arg.Kvs.CasketFname,
				privIds:     arg.Kvs.PrivKeyIds,
				privEnvName: arg.Kvs.PrivKeyEnv,
				privKeyVal:  arg.Kvs.PrivKeyVal,
				person:      false,
				device:      false,
			})
			if err != nil {
				return err
			}
			// fmt.Fprintf(arg.Nio.err, "%d:%s\n", len(closedKvps), pkms[0].Id)
			kvps := []*kvpearl.KVPearl{}
			for i := range closedKvps {
				closedKvp := closedKvps[i]
				kvp, err := kvpearl.OpenPearl(member.ToPrivateKeys(pkms), closedKvp)
				if err == nil {
					kvps = append(kvps, kvp)
				} else {
					fmt.Fprintln(arg.Nio.err, err)
				}
			}
			keys := args
			tags := arg.Kvs.Ls.tags
			// fmt.Fprintf(arg.Nio.out, "# %s\n", strings.Join(tags, ","))
			out := kvpearl.Merge(kvps, keys, tags).AsJSON()
			err = nil
			// var err error
			if *arg.Kvs.Ls.json {
				var jsStr []byte
				jsStr, err = json.MarshalIndent(out, "", "  ")
				fmt.Fprintln(arg.Nio.out, string(jsStr))
			} else if *arg.Kvs.Ls.onlyValue {
				for i := range out.Keys {
					key := out.Keys[i]
					fmt.Fprintf(arg.Nio.out, "%s\n", key.Values[0].Value)
				}
				err = nil
			} else {
				for i := range out.Keys {
					key := out.Keys[i]
					var v []byte
					v, err = json.Marshal(key.Values[0].Value)
					fmt.Fprintf(arg.Nio.out, "%s=%s\n", key.Key, string(v))
					if *arg.Kvs.Ls.shKeyValue {
						fmt.Fprintf(arg.Nio.out, "export %s\n", key.Key)
					}
					if *arg.Kvs.Ls.ghAddMask {
						fmt.Fprintf(arg.Nio.out, "echo ::add-mask::%s\n", out.Keys[i].Values[0].Value)
					}
				}
			}
			// fmt.Fprintf(arg.Nio.err, "XXXX:%d", len(out.Keys))
			// fmt.Fprintf(arg.Nio.out, "XXXX:%d", len(out.Keys))
			if len(out.Keys) == 0 {
				return errors.New(fmt.Sprintf("There was nothing found:[%s]", strings.Join(args, "],[")))
			}

			return err
		},
	}
}

func kvRmCmd(arg *NecklessArgs) *ffcli.Command {
	flags := flag.NewFlagSet("kv.Rm", flag.ExitOnError)
	return &ffcli.Command{
		Name:       "rm",
		ShortUsage: "manage a key value secrets",
		ShortHelp:  "manage key value secrets",

		LongHelp: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
		FlagSet:     flags,
		Subcommands: []*ffcli.Command{},
		Exec:        func(context.Context, []string) error { return flag.ErrHelp },
	}
}

func keyValueCmd(arg *NecklessArgs) *ffcli.Command {
	flags := flag.NewFlagSet("kv", flag.ExitOnError)
	necklessFile := findFile(".neckless")
	flags.StringVar(&arg.Kvs.Fname, "file", necklessFile, "the neckless file")
	homeDir := os.Getenv("HOME")
	flags.StringVar(&arg.Kvs.CasketFname, "casketFile",
		fmt.Sprintf("%s/.neckless/casket.json", homeDir), "filename of the casket")
	arg.Gems.PrivKeyIds = arrayFlags{}
	flags.Var(&arg.Kvs.PrivKeyIds, "privkeyid", "the neckless file")
	flags.StringVar(&arg.Kvs.PrivKeyEnv, "privkeyenv", "NECKLESS_PRIVKEY", "the neckless file")
	flags.StringVar(&arg.Kvs.PrivKeyVal, "privkeyval", "", "the neckless file")
	return &ffcli.Command{
		Name:       "kv",
		ShortUsage: "manage a key value secrets",
		ShortHelp:  "manage key value secrets",

		LongHelp: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
		FlagSet: flags,
		Subcommands: []*ffcli.Command{
			kvAddCmd(arg),
			kvLsCmd(arg),
			kvRmCmd(arg),
		},
		Exec: func(context.Context, []string) error { return flag.ErrHelp },
	}
}
