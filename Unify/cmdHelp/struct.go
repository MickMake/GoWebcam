package cmdHelp

import (
	"GoWebcam/Only"
	"github.com/spf13/cobra"
)

type Help struct {
	Error     error

	cmd       *cobra.Command
	SelfCmd   *cobra.Command
}

func New() *Help {
	var ret *Help

	for range Only.Once {
		ret = &Help {
			Error: nil,
		}
	}

	return ret
}
