package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

// NecklessOutput defines the structure of an output to a file or stdout/stderr
// this is needed for testing or multiple outputs of a kv ls command
type NecklessOutput struct {
	buf   *bytes.Buffer
	Size  int
	Name  string
	Perm  os.FileMode
	Error error
}

func (no *NecklessOutput) writer() io.Writer {
	return nil
}

// NecklessOutputs defines the list of requested outputs
type NecklessOutputs struct {
	nos []NecklessOutput
}

func (no *NecklessOutputs) first() *NecklessOutput {
	return &no.nos[0]
}

func (no *NecklessOutputs) get(fname string) *NecklessOutput {
	for i := range no.nos {
		if no.nos[i].Name == fname {
			return &no.nos[i]
		}
	}
	return nil
}

func (no *NecklessOutputs) add(fname *string) *NecklessOutput {
	if fname != nil && len(*fname) != 0 {
		nos := no.get(*fname)
		if nos != nil {
			return nos
		}
		no.nos = append(no.nos, NecklessOutput{
			buf:  &bytes.Buffer{},
			Size: 0,
			Name: *fname,
			Perm: 0644,
		})
		// fmt.Printf("add:append:%p:%s:%d\n", no, *fname, len(no.nos))
		return &no.nos[len(no.nos)-1]
	}
	return &no.nos[0]
}

// NecklessIO defines the io for the application
type NecklessIO struct {
	quite bool
	in    *bufio.Reader
	out   NecklessOutputs
	err   NecklessOutputs
}

func (no *NecklessOutputs) write(quite ...bool) {
	status := [](*NecklessOutput){}
	for i := range no.nos {
		out := no.nos[i]
		// fmt.Printf(">>>>:%p:%s:%d:%d\n", outs, out, i, len(outs.nos))
		if out.Name == "/dev/stderr" {
			os.Stderr.WriteString(out.buf.String())
		} else if out.Name == "/dev/stdout" {
			os.Stdout.WriteString(out.buf.String())
		} else {
			if out.Perm == 0000 {
				out.Perm = 0644
			}
			out.Error = ioutil.WriteFile(out.Name, out.buf.Bytes(), out.Perm)
			status = append(status, &out)
		}
		out.Size = len(out.buf.Bytes())
	}
	if (len(quite) == 0 || !quite[0]) && len(status) > 0 {
		hasError := false
		for i := range status {
			if status[i].Error != nil {
				hasError = true
			}
		}
		if hasError {
			data, err := json.MarshalIndent(status, "", "  ")
			if err != nil {
				os.Stderr.WriteString(err.Error())
				return
			}
			os.Stderr.Write(data)
		}
	}
}

// NecklessArgs defines the global args of the neckless command
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
			fmt.Fprintf(arg.Nio.out.first().buf, "Version: %s:%s\n", arg.Version, arg.GitCommit)
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
		Short:        "neckless short help",
		Long:         strings.TrimSpace("neckless long help"),
		Version:      versionStr(args),
		Args:         cobra.MinimumNArgs(0),
		RunE:         gpgRunE(args),
		SilenceUsage: true,
	}
	// rootCmd.SetOut(args.Nio.out.first().writer())
	// rootCmd.SetErr(args.Nio.err.first().writer())
	rootCmd.SetArgs(osArgs[1:])
	// rootCmd.PersistentFlags().BoolVarP(&args.Gpg.Armor, "armor", "a", false, "Author name for copyright attribution")
	// rootCmd.PersistentFlags().BoolVarP(&args.Gpg.DetachSign, "detach-sign", "b", false, "Author name for copyright attribution")
	// rootCmd.PersistentFlags().BoolVarP(&args.Gpg.Sign, "sign", "s", false, "Author name for copyright attribution")
	// rootCmd.PersistentFlags().StringVarP(&args.Gpg.UserID, "user-id", "u", "", "Author name for copyright attribution")
	rootCmd.PersistentFlags().BoolVarP(&args.Nio.quite, "quite", "q", false, "no output")
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

// GitCommit is injected during compile time
var GitCommit string

// Version is injected during compile time
var Version string

func main() {
	args := NecklessArgs{
		GitCommit: GitCommit,
		Version:   Version,
		Nio: NecklessIO{
			in: bufio.NewReader(os.Stdin),
			out: NecklessOutputs{
				nos: []NecklessOutput{{
					buf:  new(bytes.Buffer),
					Name: "/dev/stdout",
				}}},
			err: NecklessOutputs{
				nos: []NecklessOutput{{
					buf:  new(bytes.Buffer),
					Name: "/dev/stderr",
				}}},
		},
	}
	_, err := buildArgs(os.Args, &args)
	args.Nio.out.write(args.Nio.quite)
	args.Nio.err.write(args.Nio.quite)

	if err != nil {
		os.Stderr.WriteString(fmt.Sprintln(err.Error()))
		os.Exit(1)
	}
}
