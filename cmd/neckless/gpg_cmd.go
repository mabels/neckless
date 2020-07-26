//
// verify signature
/*
--keyid-format=long
--status-fd=1
--verify /var/folders/c9/2bhsxr6j7dd0ykz7cxnk0t2m0000gq/T//.git_vtag_tmp3lrQDc
 -
tree ee115078405d7f2ca8c8c5c5851788b5d8ac3f10
parent 8863222541bdd4ec2d0892c174098eb148bc2393
author Meno Abels <meno.abels@adviser.com> 1591382918 +0200
committer Meno Abels <meno.abels@adviser.com> 1591382918 +0200

cleanup for publish
*/

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type GpgArgs struct {
	KeyIdFormat string
	StatusFd    int
	Verify      string
	Sign        bool
	DetachSign  bool
	Armor       bool
	UserID      string
	GpgCli      string
}

func gpgFlags(flags *pflag.FlagSet, arg *NecklessArgs) *pflag.FlagSet {
	// flags := flag.NewFlagSet("gpg.cmd", flag.ExitOnError)
	flags.StringVar(&arg.Gpg.KeyIdFormat, "keyid-format", "long", "set keyid-format")
	flags.StringVar(&arg.Gpg.GpgCli, "gpg-cli", "/usr/local/bin/gpg", "gpg programm")
	flags.IntVar(&arg.Gpg.StatusFd, "status-fd", 1, "status file descriptor")
	flags.StringVar(&arg.Gpg.Verify, "verify", "", "verify filename")
	flags.BoolVarP(&arg.Gpg.Sign, "sign", "s", false, "make a signature")
	flags.BoolVarP(&arg.Gpg.DetachSign, "detach-sign", "b", false, "make a detached signature")
	flags.BoolVarP(&arg.Gpg.Armor, "armor", "a", false, "create ascii armored output")
	flags.StringVarP(&arg.Gpg.UserID, "user-id", "u", "", "use USER-ID to sign or decrypt")
	return flags
}

func runCmd(cmd *exec.Cmd, args *NecklessArgs) error {
	var sout bytes.Buffer
	cmd.Stdout = &sout
	var serr bytes.Buffer
	cmd.Stderr = &serr
	err := cmd.Run()
	switch args.Gpg.StatusFd {
	case 1:
		log.Println("Exec(1)=", cmd.Args, sout.String(), serr.String(), err)
		args.Nio.out.WriteString(sout.String())
		args.Nio.err.WriteString(serr.String())
		break
	case 2:
		log.Println("Exec(2)=", cmd.Args, sout.String(), serr.String(), err)
		args.Nio.out.WriteString(sout.String())
		args.Nio.err.WriteString(serr.String())
		break
	default:
		return errors.New("Unknown statusfd")
	}
	return err
}

func gpgRunE(args *NecklessArgs) func(_ *cobra.Command, arg []string) error {
	return func(_ *cobra.Command, arg []string) error {
		if args.Gpg.Sign {
			opt := []string{fmt.Sprintf("--status-fd=%d", args.Gpg.StatusFd), "-s"}
			if args.Gpg.DetachSign {
				opt = append(opt, "-b")
			}
			if args.Gpg.Armor {
				opt = append(opt, "-a")
			}
			if len(args.Gpg.UserID) > 0 {
				opt = append(opt, "-u", args.Gpg.UserID)
			}
			cmd := exec.Command(args.Gpg.GpgCli, opt...)
			if len(arg) == 0 {
				data, _ := ioutil.ReadAll(os.Stdin)
				cmd.Stdin = strings.NewReader(string(data))
			}
			return runCmd(cmd, args)
		}
		// func(_ *cobra.Command, arg []string) error {
		if len(args.Gpg.Verify) > 0 {
			content, err := ioutil.ReadFile(args.Gpg.Verify)
			if err != nil {
				return err
			}
			// log.Printf("---\n%s\n---\n%s", string(content), string(data))
			// log.Printf("---\n%s\n---", string(content))
			if strings.Contains(string(content), "BEGIN PGP SIGNATURE") {
				cmd := exec.Command(args.Gpg.GpgCli,
					fmt.Sprintf("--keyid-format=%s", args.Gpg.KeyIdFormat),
					fmt.Sprintf("--status-fd=%d", args.Gpg.StatusFd),
					"--verify", args.Gpg.Verify, arg[0])
				if arg[0] == "-" {
					data, _ := ioutil.ReadAll(os.Stdin)
					cmd.Stdin = strings.NewReader(string(data))
				}
				return runCmd(cmd, args)
			}
		}
		return nil
	}
}

func gpgCmd(arg *NecklessArgs) *cobra.Command {
	return &cobra.Command{
		Use:   "gpg",
		Short: "to plugin in that into git as gpg",
		Long: strings.TrimSpace(`
			set this command to plugin in that into git as gpg
        `),
		RunE: gpgRunE(arg),
		Args: cobra.MinimumNArgs(0),
	}
}
