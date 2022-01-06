package neckless

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/mabels/neckless/gem"
	"github.com/mabels/neckless/member"
	"github.com/mabels/neckless/necklace"
	"github.com/mabels/neckless/pearl"
)

// GemAddArgs defines arguments for the gem add command
type GemAddArgs struct {
	PubFile  string
	Device   *bool
	Person   *bool
	KeyValue *bool
	ToKeyIds []string
}

// GemLsArgs defines arguments for the gem ls command
type GemLsArgs struct {
	Device *bool
	Person *bool
}

// GemArgs defines the global arguments for the gem command
type GemArgs struct {
	Fname       string
	CasketFname string
	PrivKeyIds  []string
	Add         GemAddArgs
	Ls          GemLsArgs
}

// GetGems retrievs the gems from the neckless
func GetGems(pkms []*member.PrivateMember, nl *necklace.Necklace) ([]*gem.Gem, []error) {
	closedGems := nl.FilterByType(gem.Type)
	out := []*gem.Gem{}
	errs := []error{}
	for i := range closedGems {
		tmp := closedGems[i]
		openGem, err := gem.OpenPearl(member.ToPrivateKeys(pkms), tmp)
		// jstmp, _ := json.Marshal(tmp)
		// jsmem, _ := json.Marshal(member.ToPrivateKeys(pkms))
		if err != nil {
			// fmt.Printf("GetGem:ERR:%d:%s\n%s\n%s\n", i, err, jstmp, jsmem)
			errs = append(errs, err)
		} else {
			// jsopen, _ := json.Marshal(openGem)
			// fmt.Printf("GetGem:Open:%d:%s,%s:%s\n", i, jsopen, jstmp, jsmem)
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
	// jspms, _ := json.Marshal(pms)
	// fmt.Printf("updateGem:%s:%s\n", jspms, toIds)
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

func gemAddCmd(arg *NecklessArgs) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "add",
		Short: "manage a gem stone in neckless",

		Long: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
		RunE: func(*cobra.Command, []string) error {
			pkms, err := GetPkms(GetPkmsArgs{
				casketFname: arg.Gems.CasketFname,
				filter:      member.Matcher(arg.Gems.PrivKeyIds...),
				person:      *arg.Gems.Add.Person,
				device:      *arg.Gems.Add.Device})
			if err != nil {
				return err
			}
			var jsStr []byte
			if strings.Compare(arg.Gems.Add.PubFile, "stdin") == 0 {
				jsStr, err = ioutil.ReadAll(arg.Nio.in)
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
				// fmt.Printf("X---->%s\n", gems)
				prl, err := updateGem(gems[i], pkms, pubMembers, arg.Gems.Add.ToKeyIds...)
				if err != nil {
					return err
				}
				nl.Reset(prl, gems[i].Pearl.Closed.FingerPrint)
			}
			if len(gems) == 0 {
				myGem := gem.Create()
				prl, err := updateGem(myGem, pkms, pubMembers, arg.Gems.Add.ToKeyIds...)
				// jsprl, _ := json.Marshal(prl)
				// fmt.Printf("C---->%s:%s:%s\n", gems, err, jsprl)
				if err != nil {
					return err
				}
				nl.Reset(prl)
			}
			nl.Save(arg.Gems.Fname)
			jsStr, err = json.MarshalIndent(pubMembers, "", "  ")
			fmt.Fprintln(arg.Nio.out.first().buf, string(jsStr))
			return err
		},
	}
	flags := cmd.PersistentFlags()
	// homeDir := os.Getenv("HOME")
	flags.StringVar(&arg.Gems.Add.PubFile, "pubFile", "stdin", "the pubMemberFile to add")
	arg.Gems.Add.Person = flags.Bool("person", false, "select person keys")
	arg.Gems.Add.Device = flags.Bool("device", false, "select device keys")
	flags.StringSliceVar(&arg.Gems.Add.ToKeyIds, "toKeyId", []string{}, "the neckless file")

	return cmd
}
func gemRmCmd(arg *NecklessArgs) *cobra.Command {
	return &cobra.Command{
		Use:   "rm",
		Short: "manage a gem stone in neckless",

		Long: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
		RunE: func(_ *cobra.Command, args []string) error {
			pkms, err := GetPkms(GetPkmsArgs{
				casketFname: arg.Gems.CasketFname,
				filter:      member.Matcher(arg.Gems.PrivKeyIds...),
			})
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
			fmt.Fprintln(arg.Nio.out.first().buf, string(jsStr))
			return err
		},
	}
}
func gemLsCmd(arg *NecklessArgs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "manage a gem stone in neckless",

		Long: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
		RunE: func(*cobra.Command, []string) error {
			pkms, err := GetPkms(GetPkmsArgs{
				casketFname: arg.Gems.CasketFname,
				filter:      member.Matcher(arg.Gems.PrivKeyIds...),
				person:      *arg.Gems.Ls.Person,
				device:      *arg.Gems.Ls.Device})
			if err != nil {
				return err
			}
			nl, _ := necklace.Read(arg.Gems.Fname)
			// fmt.Fprintln(arg.Nio.err, pkms[0].Id)
			gems, _ := GetGems(pkms, &nl)
			jsStr, err := json.MarshalIndent(gem.ToJsonGems(gems...), "", "  ")
			if err != nil {
				return err
			}
			fmt.Fprintln(arg.Nio.out.first().buf, string(jsStr))
			return nil
		},
	}
	flags := cmd.PersistentFlags()
	arg.Gems.Ls.Person = flags.Bool("person", false, "select person keys")
	arg.Gems.Ls.Device = flags.Bool("device", false, "select device keys")
	return cmd
}

func gemCmd(arg *NecklessArgs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gem",
		Short: "manage a gem stone in neckless",

		Long: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
	}
	cmd.AddCommand(gemAddCmd(arg), gemRmCmd(arg), gemLsCmd(arg))

	flags := cmd.PersistentFlags()
	// homeDir := os.Getenv("HOME")
	necklessFile := findFile(".neckless")
	flags.StringVar(&arg.Gems.Fname, "file", necklessFile, "the neckless file")
	homeDir := os.Getenv("HOME")
	flags.StringVar(&arg.Gems.CasketFname, "casketFile",
		fmt.Sprintf("%s/.neckless/casket.json", homeDir), "filename of the casket")
	// arg.Gems.PrivKeyIds = arrayFlags{}
	flags.StringSliceVar(&arg.Gems.PrivKeyIds, "privkeyid", []string{}, "the neckless file")
	return cmd
}
