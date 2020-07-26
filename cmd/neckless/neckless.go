package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
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

func versionStr(args *NecklessArgs) string {
	return fmt.Sprintf("Version: %s:%s\n", args.Version, args.GitCommit)
}

func versionCmd(arg *NecklessArgs) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "version",
		Long:  strings.TrimSpace(`print version`),
		Args:  cobra.MinimumNArgs(0),
		RunE: func(*cobra.Command, []string) error {
			fmt.Fprintf(arg.Nio.out, "Version: %s:%s\n", arg.Version, arg.GitCommit)
			return nil
		},
	}
}

func buildArgs(osArgs []string, args *NecklessArgs) (*cobra.Command, error) {
	// fmt.Println(osArgs)
	// rootFlags := flag.NewFlagSet("neckless", flag.ExitOnError)
	// rootFlags.SetOutput(args.Nio.err)
	// // fmt.Fprintf(args.Nio.err, "kfkdkfkfd\n")
	// // args.Nio.err.Write([]byte("menox"))
	// // rootFlags.Output().Write([]byte("meno"))
	// rootCmd := &ffcli.Command{
	// 	Name:       "neckless",
	// 	ShortUsage: "neckless subcommand [flags]",
	// 	ShortHelp:  "neckless short help",
	// 	LongHelp:   strings.TrimSpace("neckless long help"),
	// 	Subcommands: []*ffcli.Command{
	// 		versionCmd(args),
	// 		casketCmd(args),
	// 		gemCmd(args),
	// 		keyValueCmd(args),
	// 		gpgCmd(args),
	// 	},
	// 	FlagSet: rootFlags,
	// 	Exec:    func(context.Context, []string) error { return flag.ErrHelp },
	// }

	f, err := os.OpenFile("logfile", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.SetOutput(io.Writer(f))
	log.Println("Called with=>", osArgs)
	// err = rootCmd.ParseAndRun(context.Background(), osArgs)
	// // // fmt.Printf(">>>>>", osArgs)
	// // if  err != nil && err != flag.ErrHelp {
	// // 	fmt.Fprintln(args.Nio.err, err)
	// // }

	rootCmd := &cobra.Command{
		Use: path.Base(osArgs[0]),
		// 	Name:       "neckless",
		// 	ShortUsage: "neckless subcommand [flags]",
		Short:   "neckless short help",
		Long:    strings.TrimSpace("neckless long help"),
		Version: versionStr(args),
		Args:    cobra.MinimumNArgs(0),
		RunE:    gpgRunE(args),
	}
	rootCmd.SetArgs(osArgs[1:])
	// rootCmd.PersistentFlags().BoolVarP(&args.Gpg.Armor, "armor", "a", false, "Author name for copyright attribution")
	// rootCmd.PersistentFlags().BoolVarP(&args.Gpg.DetachSign, "detach-sign", "b", false, "Author name for copyright attribution")
	// rootCmd.PersistentFlags().BoolVarP(&args.Gpg.Sign, "sign", "s", false, "Author name for copyright attribution")
	// rootCmd.PersistentFlags().StringVarP(&args.Gpg.UserID, "user-id", "u", "", "Author name for copyright attribution")
	gpgFlags(rootCmd.PersistentFlags(), args)
	rootCmd.AddCommand(gpgCmd(args))
	rootCmd.AddCommand(versionCmd(args))
	rootCmd.AddCommand(casketCmd(args))
	rootCmd.AddCommand(gemCmd(args))
	rootCmd.AddCommand(keyValueCmd(args))

	// cmdEcho.AddCommand(cmdTimes)
	err = rootCmd.Execute()
	// fmt.Println(args.Gpg)
	f.Close()

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
	_, err := buildArgs(os.Args, &args)
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
