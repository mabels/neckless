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
	"context"
	"flag"
	"io/ioutil"
	"strings"

	"github.com/peterbourgon/ff/v2/ffcli"
)

type GpgArgs struct {
	KeyIdFormat string
	StatusFd    int
	Verify      string
}

func gpgCmd(arg *NecklessArgs) *ffcli.Command {
	flags := flag.NewFlagSet("gpg.cmd", flag.ExitOnError)
	flags.StringVar(&arg.Gpg.KeyIdFormat, "keyid-format", "long", "set keyid-format")
	flags.IntVar(&arg.Gpg.StatusFd, "status-fd", 1, "status file descriptor")
	flags.StringVar(&arg.Gpg.Verify, "verify", "", "verify filename")

	return &ffcli.Command{
		Name:       "git",
		ShortUsage: "set this command to plugin in that into git as gpg",
		ShortHelp:  "gpg substitude",

		LongHelp: strings.TrimSpace(`
			set this command to plugin in that into git as gpg
        `),
		FlagSet:     flags,
		Subcommands: []*ffcli.Command{},
		Exec: func(_ context.Context, args []string) error {
			if len(arg.Gpg.Verify) > 0 {
				content, err := ioutil.ReadFile(arg.Gpg.KeyIdFormat)
				if err != nil {
					return err
				}
				if strings.Contains(string(content), "BEGIN PGP SIGNATURE") {
					cmd := exec.Command("/usr/local/bin/gpg",
						"keyid-format", arg.Gpg.KeyIdFormat,
						"status-fd", strconv.Itoa(arg.Gpg.StatusFd),
						"verify", arg.Gpg.Verify,
						args...)
					err := cmd.Run()
						log.Printf("Command finished with error: %v", err)
					}
					if err != nil {
						return err
					}
				}
			}
			switch arg.Gpg.KeyIdFormat {
			case "long":
				break
			case "short":
				break
			}
			return nil
		},
	}
}
