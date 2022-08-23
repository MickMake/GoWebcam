package cmdHelp

import (
	"GoWebcam/Only"
	"GoWebcam/defaults"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
	"strings"
)


func PrintConfig(cmd *cobra.Command) {
	for range Only.Once {
		// fmt.Printf("Config file '%s':\n", Cmd.ConfigFile)	// @TODO - fixup.

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Flag", "Short flag", "Environment", "Description", "Value"})
		table.SetBorder(true)

		cmd.PersistentFlags().SortFlags = false
		cmd.Flags().SortFlags = false
		cmd.Flags().VisitAll(func(flag *pflag.Flag) {
			if flag.Hidden {
				return
			}

			sh := ""
			if flag.Shorthand != "" {
				sh = "-" + flag.Shorthand
			}

			// fmt.Printf("key: %s => %v (%s)\n", flag.Name, flag.Value, flag.Value.String())
			table.Append([]string{
				"--" + flag.Name,
				sh,
				PrintFlagEnv(flag.Name),
				flag.Usage,
				flag.Value.String(),
				// flag.Value.String(),
			})
		})

		table.Render()
	}
}

func PrintFlagEnv(flag string) string {
	fenv := strings.ReplaceAll(flag, "-", "_")
	fenv = strings.ToUpper(fenv)

	// ret := fmt.Sprintf("--%s\t%s_%s\n", flag, EnvPrefix, fenv)
	ret := fmt.Sprintf("%s_%s", defaults.EnvPrefix, fenv)
	return ret
}

func PrintFlags(cmd *cobra.Command) {
	for range Only.Once {
		fmt.Printf("\nUsing environment variables instad of flags.\n")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Flag", "Short flag", "Environment", "Description", "Default"})
		table.SetBorder(true)

		cmd.PersistentFlags().SortFlags = false
		cmd.Flags().SortFlags = false
		cmd.Flags().VisitAll(func(flag *pflag.Flag) {
			if flag.Hidden {
				return
			}

			sh := ""
			if flag.Shorthand != "" {
				sh = "-" + flag.Shorthand
			}

			table.Append([]string{
				"--" + flag.Name,
				sh,
				PrintFlagEnv(flag.Name),
				flag.Usage,
				flag.DefValue,
				// flag.Value.String(),
			})
		})

		table.Render()
	}
}


func PrintExamples(cmd *cobra.Command, examples ...string) string {
	var ret string

	c := BuildCmd(cmd)
	for _, example := range examples {
		ret += fmt.Sprintf("\t%s %s\n", c, example)
	}

	return ret
}

func BuildCmd(cmd *cobra.Command) string {
	var ret string
	if cmd.HasParent() {
		ret += BuildCmd(cmd.Parent())
	}
	ret += cmd.Name() + " "
	return ret
}
