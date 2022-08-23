package cmdHelp

import (
	"GoWebcam/Only"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

type Help struct {
	Error     error

	Command      string
	HelpTemplate string
	UsageTemplate string
	FlagHelpTemplate string
	ExtendedHelpTemplate string

	cmd       *cobra.Command
	SelfCmd   *cobra.Command
}

func New() *Help {
	var ret *Help

	for range Only.Once {
		ret = &Help {
			Error: nil,

			Command:              "DefaultBinaryName",
			HelpTemplate:         DefaultHelpTemplate,
			UsageTemplate:        DefaultUsageTemplate,
			FlagHelpTemplate:     DefaultFlagHelpTemplate,
			ExtendedHelpTemplate: ExtendedHelpTemplate,

			cmd: nil,
			SelfCmd: nil,
		}
	}

	return ret
}

func (h *Help) SetCommand(text string) {
	for range Only.Once {
		if text == "" {
			break
		}

		h.Command = text
	}
}

func (h *Help) SetHelpTemplate(text string) {
	for range Only.Once {
		if text == "" {
			break
		}

		h.HelpTemplate = strings.ReplaceAll(text, "DefaultBinaryName", h.Command)
	}
}

func (h *Help) SetUsageTemplate(text string) {
	for range Only.Once {
		if text == "" {
			break
		}

		h.UsageTemplate = strings.ReplaceAll(text, "DefaultBinaryName", h.Command)
	}
}

func (h *Help) SetFlagHelpTemplate(text string) {
	for range Only.Once {
		if text == "" {
			break
		}

		h.FlagHelpTemplate = strings.ReplaceAll(text, "DefaultBinaryName", h.Command)
	}
}

func (h *Help) SetExtendedHelpTemplate(text string) {
	for range Only.Once {
		if text == "" {
			break
		}

		h.ExtendedHelpTemplate = strings.ReplaceAll(text, "DefaultBinaryName", h.Command)
	}
}

func (h *Help) ExtendedHelp() {
	fmt.Println(h.ExtendedHelpTemplate)
}