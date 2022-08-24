package cmdHelp

import (
	"GoWebcam/Only"
	"fmt"
	"github.com/spf13/cobra"
)


const Command = "help"
func (h *Help) AttachCommands(cmd *cobra.Command) *cobra.Command {
	for range Only.Once {
		if cmd == nil {
			break
		}
		h.cmd = cmd

		// ******************************************************************************** //
		h.SelfCmd = &cobra.Command {
			Use:                   "help-all",
			// Aliases:               []string{"flags"},
			Short:                 fmt.Sprintf("Help - Extended help"),
			Long:                  fmt.Sprintf("Help - Extended help"),
			DisableFlagParsing:    false,
			DisableFlagsInUseLine: false,
			PreRunE:               h.InitArgs,
			RunE:                  h.CmHelpAll,
			Args:                  cobra.RangeArgs(0, 0),
		}
		cmd.AddCommand(h.SelfCmd)
		h.SelfCmd.Example = PrintExamples(h.SelfCmd, "")
		h.SelfCmd.Annotations = map[string]string{"command":Command}

		h.cmd.SetHelpTemplate(DefaultHelpTemplate)
		h.cmd.SetUsageTemplate(DefaultUsageTemplate)
	}

	return h.SelfCmd
}

func (h *Help) InitArgs(_ *cobra.Command, _ []string) error {
	var err error
	for range Only.Once {
		//
	}
	return err
}

func (h *Help) CmHelpAll(cmd *cobra.Command, args []string) error {
	for range Only.Once {
		if len(args) > 0 {
			fmt.Println("Unknown sub-command.")
		}

		h.ExtendedHelp()

		// cmd.SetUsageTemplate(DefaultFlagHelpTemplate)
		cmd.SetUsageTemplate("")
		h.Error = cmd.Help()

		h.PrintConfig(h.cmd)
	}

	return h.Error
}
