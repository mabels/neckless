package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/peterbourgon/ff/v2/ffcli"
)

// CrazyBeeArgs Toplevel Command Args
type NecklessArgs struct {
	casket CasketArgs
	kvs    KeyValueArgs
	gems   GemArgs
}

func buildArgs(osArgs []string, args *NecklessArgs) *ffcli.Command {
	// fmt.Println(osArgs)
	rootFlags := flag.NewFlagSet("neckless", flag.ExitOnError)
	rootCmd := &ffcli.Command{
		Name:       "neckless",
		ShortUsage: "neckless subcommand [flags]",
		ShortHelp:  "neckless short help",
		LongHelp:   strings.TrimSpace("neckless long help"),
		Subcommands: []*ffcli.Command{
			casketArgs(&args.casket),
			gemArgs(&args.gems),
			keyValueArgs(&args.kvs),
		},
		FlagSet: rootFlags,
		Exec:    func(context.Context, []string) error { return flag.ErrHelp },
	}

	if err := rootCmd.ParseAndRun(context.Background(), osArgs); err != nil && err != flag.ErrHelp {
		log.Fatal(err)
	}
	return rootCmd
}

// func add(a, b int) int {
// 	return a + b
// }

func main() {
	args := NecklessArgs{}
	buildArgs(os.Args[1:], &args)
	// fmt.Println(">>>", args, cmd.FlagSet.Args())
	// fmt.Println("DryRun", args.casket.create.DryRun)
	// fmt.Println("File", *&args.casket.create.Fname)
	// fmt.Println("Name", args.casket.create.MemberArg.Name)
	// fmt.Println("Device", args.casket.create.MemberArg.Device)
	// fmt.Println("Type", args.casket.create.MemberArg.Type)
}
