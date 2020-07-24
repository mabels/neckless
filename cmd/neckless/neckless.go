package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/peterbourgon/ff/v2/ffcli"
)

// CrazyBeeArgs Toplevel Command Args
type NecklessIO struct {
	// in  *bufio.Reader
	// out *bufio.Writer
	// err *bufio.Writer
	in  *bufio.Reader
	out *bytes.Buffer
	err *bytes.Buffer
}
type NecklessArgs struct {
	GitCommit string
	Version   string
	Nio       NecklessIO
	Casket    CasketArgs
	Kvs       KeyValueArgs
	Gems      GemArgs
	Gpg       GpgArgs
}

func versionCmd(arg *NecklessArgs) *ffcli.Command {
	return &ffcli.Command{
		Name:       "version",
		ShortUsage: "version",
		ShortHelp:  "print version help",
		LongHelp:   strings.TrimSpace("print version help"),
		Exec: func(context.Context, []string) error {
			fmt.Fprintf(arg.Nio.out, "Version: %s:%s\n", arg.Version, arg.GitCommit)
			return nil
		},
	}
}

func buildArgs(osArgs []string, args *NecklessArgs) (*ffcli.Command, error) {
	// fmt.Println(osArgs)
	rootFlags := flag.NewFlagSet("neckless", flag.ExitOnError)
	rootFlags.SetOutput(args.Nio.err)
	// fmt.Fprintf(args.Nio.err, "kfkdkfkfd\n")
	// args.Nio.err.Write([]byte("menox"))
	// rootFlags.Output().Write([]byte("meno"))
	rootCmd := &ffcli.Command{
		Name:       "neckless",
		ShortUsage: "neckless subcommand [flags]",
		ShortHelp:  "neckless short help",
		LongHelp:   strings.TrimSpace("neckless long help"),
		Subcommands: []*ffcli.Command{
			versionCmd(args),
			casketCmd(args),
			gemCmd(args),
			keyValueCmd(args),
			gpgCmd(args),
		},
		FlagSet: rootFlags,
		Exec:    func(context.Context, []string) error { return flag.ErrHelp },
	}

	err := rootCmd.ParseAndRun(context.Background(), osArgs)
	// // fmt.Printf(">>>>>", osArgs)
	// if  err != nil && err != flag.ErrHelp {
	// 	fmt.Fprintln(args.Nio.err, err)
	// }
	return rootCmd, err
}

// func add(a, b int) int {
// 	return a + b
// }

var GitCommit string
var Version string

func main() {
	nio := NecklessIO{
		in:  bufio.NewReader(os.Stdin),
		out: new(bytes.Buffer),
		err: new(bytes.Buffer),
	}
	args := NecklessArgs{
		GitCommit: GitCommit,
		Version:   Version,
		Nio:       nio,
	}
	_, err := buildArgs(os.Args[1:], &args)
	// fmt.Println("xxxx", nio.out.String())
	os.Stdout.WriteString(nio.out.String())
	// os.Stderr.WriteString("Hallo")
	os.Stderr.WriteString(nio.err.String())
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintln(err.Error()))
		os.Exit(1)
	}
	// nio.out.Flush()
	// nio.err.Flush()
	// fmt.Println(">>>", args, cmd.FlagSet.Args())
	// fmt.Println("DryRun", args.casket.create.DryRun)
	// fmt.Println("File", *&args.casket.create.Fname)
	// fmt.Println("Name", args.casket.create.MemberArg.Name)
	// fmt.Println("Device", args.casket.create.MemberArg.Device)
	// fmt.Println("Type", args.casket.create.MemberArg.Type)
}
