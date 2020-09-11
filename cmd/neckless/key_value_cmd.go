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

// KeyValueLsArgs defines the arguments to the kv ls command
type KeyValueLsArgs struct {
	json       *bool
	keyValue   *bool
	shKeyValue *bool
	ghAddMask  *bool
	onlyValue  *bool
	tags       *[]string
	writeFname string
}

// KeyValueArgs defines the global arguments
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
		kvp, err := kvpearl.Parse(args[i])
		if err != nil {
			errs = append(errs, err)
			continue
		}
		kvp, err = kvp.Resolv(func(key string, fparam kvpearl.FuncsAndParam) (*string, error) {
			c, err := ioutil.ReadFile(fparam.Param)
			if err != nil {
				return nil, err
			}
			ret := string(c)
			return &ret, nil
		})
		if err != nil {
			errs = append(errs, err)
		} else {
			sa, err := kvp.ToSetArgs()
			if err != nil {
				errs = append(errs, err)
			}
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

func parseArgs2KVpearl(args []string, write string, tags []string) (kvpearl.MapByToResolve, []error) {
	// resolv := func(key string, fname string) (*string, error) {
	// 	out := fmt.Sprintf("IGNORED:%s:%s", fname, key)
	// 	return &out, nil
	// }
	// kvp := kvpearl.Create()
	ret := kvpearl.MapByToResolve{}
	errs := []error{}
	if len(args) == 0 && len(tags) > 0 {
		args = []string{".*"}
	}
	for i := range args {
		m := matchNoAtOrEqual.FindStringSubmatch(args[i])
		if len(m) == 2 {
			arg := fmt.Sprintf("%s@%s[%s]", m[1], write, strings.Join(tags, ","))
			sa, err := kvpearl.Parse(arg)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			ret.Add(sa)
		} else {
			sa, err := kvpearl.Parse(args[i])
			if err != nil {
				errs = append(errs, err)
				continue
			}
			// if sa.ToResolve == nil || len(sa.ToResolve.Param) == 0 {
			// 	sa.ToResolve = kvpearl.ParseFuncsAndParams(write)
			// }
			ret.Add(sa)
		}
	}
	return ret, errs
}

func runActions(kv *kvpearl.JSONValue) (string, error) {
	val := kv.Value
	if kv.Unresolved != nil {
		var err error
		val, err = kv.Unresolved.RunFuncs(val)
		if err != nil {
			// fmt.Fprintln(narg.Nio.err.first().buf, err)
			return val, err
		}
		// val = fmt.Sprintf("(%s:%d)", val, len)
	}
	return val, nil
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
			// myOut, _ := json.MarshalIndent(keys, "", "  ")
			// fmt.Printf("%s:%s\n", args, string(myOut))
			// tags := narg.Kvs.Ls.tags
			// fmt.Fprintf(arg.Nio.out, "# %s\n", strings.Join(tags, ","))
			outputs := kvps.Match(keys)
			// out := kvps.AsJSON()
			err = nil
			// var err error
			// outputs := toKVPearl2Outputs(out)
			for fname := range outputs {
				kvs := outputs[fname]
				outValues := kvs.ToJSON()
				if *narg.Kvs.Ls.json {
					for i := range outValues {
						kv := outValues[i]
						for j := range kv.Vals {
							val, err := runActions(kv.Vals[j])
							if err != nil {
								return err
							}
							kv.Vals[j].Value = val
						}
					}
					jsStr, err := json.MarshalIndent(outValues, "", "  ")
					if err != nil {
						fmt.Fprintf(narg.Nio.err.first().buf, "%s", err)
					}
					fmt.Fprintln(narg.Nio.out.add(&fname).buf, string(jsStr))
				} else if *narg.Kvs.Ls.onlyValue {
					for i := range outValues {
						val, err := runActions(outValues[i].Vals[0])
						if err != nil {
							return err
						}
						out := narg.Nio.out.add(&fname)
						fmt.Fprintf(out.buf, "%s\n", val)
					}
					err = nil
				} else {
					eol := "\n"
					if *narg.Kvs.Ls.shKeyValue {
						eol = ";\n"
					}
					for i := range outValues {
						kv := outValues[i]
						val, err := runActions(kv.Vals[0])
						if err != nil {
							return err
						}
						var v []byte
						v, err = json.Marshal(val)
						fmt.Fprintf(narg.Nio.out.add(&fname).buf, "%s=%s%s", kv.Key, string(v), eol)
						if *narg.Kvs.Ls.shKeyValue {
							fmt.Fprintf(narg.Nio.out.add(&fname).buf, "export %s%s", kv.Key, eol)
						}
						if *narg.Kvs.Ls.ghAddMask {
							fmt.Fprintf(narg.Nio.out.add(&fname).buf, "echo ::add-mask::%s%s", kv.Vals[0].Value, eol)
						}
					}
				}
			}
			// fmt.Fprintf(arg.Nio.err, "XXXX:%d", len(out.Keys))
			// fmt.Fprintf(arg.Nio.out, "XXXX:%d", len(out.Keys))

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
