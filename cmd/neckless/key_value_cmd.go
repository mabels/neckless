package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
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
	tags       *[]string
	writeFname string
}

type KeyValueArgs struct {
	Fname       string
	CasketFname string
	PrivKeyIds  *[]string
	PrivKeyEnv  string
	PrivKeyVal  string
	Ls          KeyValueLsArgs
}

func toKVreadFile(args []string) ([]*kvpearl.SetArg, []error) {
	// ret := []kvpearl.Key{}
	errs := []error{}
	sas := []*kvpearl.SetArg{}
	for i := range args {
		sa, err := kvpearl.Parse(args[i], func(key string, fname string) (*string, error) {
			c, err := ioutil.ReadFile(fname)
			if err != nil {
				return nil, err
			}
			ret := string(c)
			return &ret, nil
		})
		if err != nil {
			errs = append(errs, err)
		} else {
			sas = append(sas, sa)
		}
	}
	return sas, errs
}
func kvAddCmd(narg *NecklessArgs) *cobra.Command {
	// Args VAL=Wert[TAGS,TAGS]
	return &cobra.Command{
		SilenceErrors: true,
		Use:           "add",
		Short:         "manage a key value secrets args <KEY=VAL[TAGS,]>",

		Long: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				return flag.ErrHelp
			}
			sas, errs := toKVreadFile(args)
			for i := range errs {
				fmt.Fprintln(narg.Nio.err.first().buf, errs[i])
			}
			kvp := kvpearl.Create()
			for i := range sas {
				kvp.Set(*sas[i])
			}
			jsStr, err := json.MarshalIndent(kvp.AsJSON(), "", "  ")
			if err != nil {
				return err
			}
			fmt.Fprintln(narg.Nio.out.first().buf, string(jsStr))
			nl, errnl := necklace.Read(narg.Kvs.Fname)
			if len(errnl) > 0 {
				out := make([]string, len(errnl))
				for i := range errnl {
					out[i] = errnl[i].Error()
				}
				return errors.New(strings.Join(out, "|"))
			}
			pkms, err := GetPkms(GetPkmsArgs{
				casketFname: narg.Kvs.CasketFname,
				filter:      member.Matcher(*narg.Kvs.PrivKeyIds...),
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
				// fmt.Println("signer:", pkms[0].Id)
				// fmt.Fprintf(narg.Nio.err.first().buf, "%s:%s\n", pkms[0].Id, gems[i].LsByType(member.Device)[0].Id)
				p, err := kvp.ClosePearl(&mo)
				if err != nil {
					return err
				}
				nl.Reset(p)
			}
			_, err = nl.Save(narg.Kvs.Fname)
			return err
		},
	}
}

var matchNoAtOrEqual = regexp.MustCompile("^([^@=]+)$")

func parseArgs2KVpearl(args []string, write string, tags []string) (map[string]*kvpearl.KVParsed, []error) {
	resolv := func(key string, fname string) (*string, error) {
		out := fmt.Sprintf("IGNORED:%s:%s", fname, key)
		return &out, nil
	}
	// kvp := kvpearl.Create()
	ret := map[string]*kvpearl.KVParsed{}
	errs := []error{}
	for i := range args {
		m := matchNoAtOrEqual.FindStringSubmatch(args[i])
		if len(m) == 2 {
			arg := fmt.Sprintf("%s@%s[%s]", m[1], write, strings.Join(tags, ","))
			// fmt.Println("parsed:", arg)
			sa, err := kvpearl.Parse(arg, resolv)
			if err != nil {
				errs = append(errs, err)
			} else {
				ret[sa.Key] = sa.ToKVParsed()
			}

		} else {
			sa, err := kvpearl.Parse(args[i], resolv)
			// myOut, _ := json.MarshalIndent(kvp, "", "  ")
			// fmt.Println("parsed:", args[i], string(myOut))
			if err != nil {
				errs = append(errs, err)
			} else {
				ret[sa.Key] = sa.ToKVParsed()
			}
		}
	}
	return ret, errs
}

func kvLsCmd(narg *NecklessArgs) *cobra.Command {
	cmd := &cobra.Command{
		SilenceErrors: true,
		Use:           "ls",
		Short:         "manage a key value secrets",

		Long: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
		RunE: func(_ *cobra.Command, args []string) error {
			nl, errnl := necklace.Read(narg.Kvs.Fname)
			if len(errnl) > 0 {
				out := make([]string, len(errnl))
				for i := range errnl {
					out[i] = errnl[i].Error()
				}
				return errors.New(strings.Join(out, "|"))
			}
			closedKvps := nl.FilterByType(kvpearl.Type)
			pkms, err := GetPkms(GetPkmsArgs{
				casketFname: narg.Kvs.CasketFname,
				filter:      member.Matcher(*narg.Kvs.PrivKeyIds...),
				privEnvName: narg.Kvs.PrivKeyEnv,
				privKeyVal:  narg.Kvs.PrivKeyVal,
				person:      false,
				device:      false,
			})
			if err != nil {
				return err
			}
			// fmt.Fprintf(arg.Nio.err, "%d:%s\n", len(closedKvps), pkms[0].Id)
			kvps := kvpearl.KVPearls{}
			for i := range closedKvps {
				closedKvp := closedKvps[i]
				kvp, err := kvpearl.OpenPearl(member.ToPrivateKeys(pkms), closedKvp)
				if err == nil {
					kvps = append(kvps, kvp)
				} else {
					fmt.Fprintln(narg.Nio.err.first().buf, err)
				}
			}
			keys, errs := parseArgs2KVpearl(args, narg.Kvs.Ls.writeFname, *narg.Kvs.Ls.tags)
			for i := range errs {
				fmt.Fprintln(narg.Nio.err.first().buf, errs[i])
			}
			myOut, _ := json.MarshalIndent(keys, "", "  ")
			fmt.Printf("%s:%s\n", args, myOut)
			// tags := narg.Kvs.Ls.tags
			// fmt.Fprintf(arg.Nio.out, "# %s\n", strings.Join(tags, ","))
			out := kvps.Merge(keys).AsJSON()
			// out := kvps.AsJSON()
			err = nil
			// var err error
			if *narg.Kvs.Ls.json {
				var jsStr []byte
				jsStr, err = json.MarshalIndent(out, "", "  ")
				fmt.Fprintln(narg.Nio.out.first().buf, string(jsStr))
			} else if *narg.Kvs.Ls.onlyValue {
				for i := range out.Keys {
					key := out.Keys[i]
					// fmt.Println("writeTo:", *key.Values[0].Unresolved)
					out := narg.Nio.out.add(key.Values[0].Unresolved)
					fmt.Fprintf(out.buf, "%s\n", key.Values[0].Value)
				}
				err = nil
			} else {
				eol := "\n"
				if *narg.Kvs.Ls.shKeyValue {
					eol = ";\n"
				}
				for i := range out.Keys {
					key := out.Keys[i]
					var v []byte
					v, err = json.Marshal(key.Values[0].Value)

					fmt.Fprintf(narg.Nio.out.first().buf, "%s=%s%s", key.Key, string(v), eol)
					if *narg.Kvs.Ls.shKeyValue {
						fmt.Fprintf(narg.Nio.out.first().buf, "export %s%s", key.Key, eol)
					}
					if *narg.Kvs.Ls.ghAddMask {
						fmt.Fprintf(narg.Nio.out.first().buf, "echo ::add-mask::%s%s", out.Keys[i].Values[0].Value, eol)
					}
				}
			}
			if len(narg.Kvs.Ls.writeFname) > 0 {
				narg.Nio.out.first().Name = narg.Kvs.Ls.writeFname
			}
			// fmt.Fprintf(arg.Nio.err, "XXXX:%d", len(out.Keys))
			// fmt.Fprintf(arg.Nio.out, "XXXX:%d", len(out.Keys))
			if len(out.Keys) == 0 {
				return errors.New(fmt.Sprintf("There was nothing found:[%s]", strings.Join(args, "],[")))
			}

			return err
		},
	}
	flags := cmd.PersistentFlags()
	narg.Kvs.Ls.json = flags.Bool("json", false, "select device keys")
	narg.Kvs.Ls.keyValue = flags.Bool("keyValue", true, "select device keys")
	narg.Kvs.Ls.onlyValue = flags.Bool("onlyValue", false, "select device keys")
	narg.Kvs.Ls.shKeyValue = flags.Bool("shKeyValue", false, "select device keys")
	narg.Kvs.Ls.ghAddMask = flags.Bool("ghAddMask", false, "set Value as github mask")
	narg.Kvs.Ls.tags = flags.StringSlice("tag", []string{}, "list of tags to filter")
	flags.StringVar(&narg.Kvs.Ls.writeFname, "write", "", "name of the file to write to")
	return cmd
}

func kvRmCmd(arg *NecklessArgs) *cobra.Command {
	return &cobra.Command{
		SilenceErrors: true,
		Use:           "rm",
		Short:         "manage a key value secrets",

		Long: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
	}
}

func keyValueCmd(arg *NecklessArgs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kv",
		Short: "manage a key value secrets",
		Long: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
	}
	cmd.AddCommand(kvAddCmd(arg), kvLsCmd(arg), kvRmCmd(arg))
	flags := cmd.PersistentFlags()
	necklessFile := findFile(".neckless")
	flags.StringVar(&arg.Kvs.Fname, "file", necklessFile, "the neckless file")
	homeDir := os.Getenv("HOME")
	flags.StringVar(&arg.Kvs.CasketFname, "casketFile",
		fmt.Sprintf("%s/.neckless/casket.json", homeDir), "filename of the casket")
	arg.Kvs.PrivKeyIds = flags.StringSlice("privkeyid", []string{}, "the neckless file")
	flags.StringVar(&arg.Kvs.PrivKeyEnv, "privkeyenv", "NECKLESS_PRIVKEY", "the neckless file")
	flags.StringVar(&arg.Kvs.PrivKeyVal, "privkeyval", "", "the neckless file")
	return cmd
}
