package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"

	"crazybee.adviser.com/pipeline"
	"github.com/google/uuid"
	"github.com/peterbourgon/ff/v2/ffcli"
)

// CrazyBeeArgs Toplevel Command Args
type CrazyBeeArgs struct {
	pipeline pipeline.PipelineArgs
}

func runTest(ctx context.Context, args []string) error {
	log.Printf("runTest:\n")
	return nil
}
func pipeLineArgs(arg *pipeline.PipelineArgs) *ffcli.Command {
	flags := flag.NewFlagSet("pipeline", flag.ExitOnError)
	defaultName := uuid.New().String()
	flags.StringVar(&arg.name, "name", defaultName, "name of the key")
	// defaultDur, _ := time.ParseDuration("5y")
	// flags.Duration("valid", defaultDur, "defined the validilty of the key")
	return &ffcli.Command{
		Name:       "Pipeline",
		ShortUsage: "manage a pipeline secrets",
		ShortHelp:  "manage pipeline secrets",

		LongHelp: strings.TrimSpace(`
    This command is used to create and add user to the pipeline secret
    `),
		FlagSet: flags,
	}
}

func buildArgs(osArgs []string, args *CrazyBeeArgs) *ffcli.Command {
	rootFlags := flag.NewFlagSet("crazybee", flag.ExitOnError)
	rootCmd := &ffcli.Command{
		Name:       "crazybee",
		ShortUsage: "crazybee subcommand [flags]",
		ShortHelp:  "crazybee short help",
		LongHelp:   strings.TrimSpace("crazybee long help"),
		Subcommands: []*ffcli.Command{
			pipeLineArgs(&args.pipeline),
		},
		FlagSet: rootFlags,
		Exec:    func(context.Context, []string) error { return flag.ErrHelp },
	}

	if err := rootCmd.Parse(osArgs); err != nil && err != flag.ErrHelp {
		log.Fatal(err)
	}
	return rootCmd
}

func main() {
	var args CrazyBeeArgs
	buildArgs(os.Args[1:], &args)
}
