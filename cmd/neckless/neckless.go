package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/peterbourgon/ff/v2/ffcli"
	"neckless.adviser.com/casket"
	"neckless.adviser.com/member"
)

type CasketArgs struct {
	create casket.CreateArg
}
type NecklessCmdArgs struct {
}

// CrazyBeeArgs Toplevel Command Args
type NecklessArgs struct {
	casket   CasketArgs
	neckless NecklessCmdArgs
}

// func runTest(ctx context.Context, args []string) error {
// 	log.Printf("runTest:\n")
// 	return nil
// }

type CreateArg struct {
	member.MemberArg
	CasketDryRun bool    // if dryrun don't write
	CasketFname  *string //
}

func casketCreateArgs(arg *CasketArgs) *ffcli.Command {
	flags := flag.NewFlagSet("create", flag.ExitOnError)
	defaultName := uuid.New().String()
	flags.StringVar(&arg.create.Name, "name", defaultName, "name of the key")

	// fmt.Println("xXxxx", args.casket.create.MemberArg.Device)
	defaultDevice := uuid.New().String()
	flags.StringVar(&arg.create.Device, "device", defaultDevice, "device of the key")
	return &ffcli.Command{
		Name:       "create",
		ShortUsage: "manage a casket secrets",
		ShortHelp:  "manage casket secrets",

		LongHelp: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
		FlagSet: flags,
	}
}

func casketLsArgs(arg *CasketArgs) *ffcli.Command {
	flags := flag.NewFlagSet("ls", flag.ExitOnError)
	return &ffcli.Command{
		Name:       "ls",
		ShortUsage: "manage a casket secrets",
		ShortHelp:  "manage casket secrets",

		LongHelp: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
		FlagSet: flags,
	}
}

func casketRmArgs(arg *CasketArgs) *ffcli.Command {
	flags := flag.NewFlagSet("rm", flag.ExitOnError)
	return &ffcli.Command{
		Name:       "rm",
		ShortUsage: "manage a casket secrets",
		ShortHelp:  "manage casket secrets",

		LongHelp: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
		FlagSet: flags,
	}
}
func casketArgs(arg *CasketArgs) *ffcli.Command {
	flags := flag.NewFlagSet("casket", flag.ExitOnError)
	homeDir := os.Getenv("HOME")
	flags.StringVar(&arg.create.CasketFname, "file",
		fmt.Sprintf("%s/.neckless/casket.json", homeDir), "filename of the casket")
	arg.create.CasketDryRun = flags.Bool("dryRun", false, "set the dryrun flag")

	// fmt.Println("xXxxx", args.casket.create.MemberArg.Type)
	// var typ string
	// flags.StringVar(&arg.create.Type, "type", "xxx", "type of the key")
	// fmt.Println(typ)
	// if strings.Compare(typ, string(member.Person)) == 0 {
	// 	arg.create.Type = member.Person
	// }
	// if strings.Compare(typ, string(member.Device)) == 0 {
	// 	arg.create.Type = member.Device
	// }
	// 	// defaultDur, _ := time.ParseDuration("5y")
	// 	// flags.Duration("valid", defaultDur, "defined the validilty of the key")
	return &ffcli.Command{
		Name:       "casket",
		ShortUsage: "manage a casket secrets",
		ShortHelp:  "manage casket secrets",

		LongHelp: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
		FlagSet: flags,
		Subcommands: []*ffcli.Command{
			casketCreateArgs(arg),
			casketLsArgs(arg),
			casketRmArgs(arg),
		},
	}
}

func necklessArgs(arg *NecklessCmdArgs) *ffcli.Command {
	flags := flag.NewFlagSet("casket", flag.ExitOnError)
	// ret := flags.Bool("dryRun", arg.create.CasketDryRun, "set the dryrun flag")
	// fmt.Println("xxxxx", ret, *ret)
	// if ret != nil {
	// 	arg.create.CasketDryRun = *ret
	// }
	// defaultName := uuid.New().String()
	// flags.StringVar(&arg.create.Name, "name", defaultName, "name of the key")

	// // fmt.Println("xXxxx", args.casket.create.MemberArg.Device)
	// defaultDevice := uuid.New().String()
	// flags.StringVar(&arg.create.Device, "device", defaultDevice, "device of the key")

	// fmt.Println("xXxxx", args.casket.create.MemberArg.Type)
	// var typ string
	// flags.StringVar(&arg.create.Type, "type", "xxx", "type of the key")
	// fmt.Println(typ)
	// if strings.Compare(typ, string(member.Person)) == 0 {
	// 	arg.create.Type = member.Person
	// }
	// if strings.Compare(typ, string(member.Device)) == 0 {
	// 	arg.create.Type = member.Device
	// }
	// 	// defaultDur, _ := time.ParseDuration("5y")
	// 	// flags.Duration("valid", defaultDur, "defined the validilty of the key")
	return &ffcli.Command{
		Name:       "neckless",
		ShortUsage: "manage a neckless",
		ShortHelp:  "manage a neckless",

		LongHelp: strings.TrimSpace(`
	    This command is used to create and add user to the pipeline secret
	    `),
		FlagSet: flags,
	}
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
			// missing GemStone
			// missing KV
			necklessArgs(&args.neckless),
		},
		FlagSet: rootFlags,
		Exec:    func(context.Context, []string) error { return flag.ErrHelp },
	}

	if err := rootCmd.Parse(osArgs); err != nil && err != flag.ErrHelp {
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
	fmt.Println("DryRun", *args.casket.create.CasketDryRun)
	fmt.Println("File", *&args.casket.create.CasketFname)
	fmt.Println("Name", args.casket.create.MemberArg.Name)
	fmt.Println("Device", args.casket.create.MemberArg.Device)
	fmt.Println("Type", args.casket.create.MemberArg.Type)
}
