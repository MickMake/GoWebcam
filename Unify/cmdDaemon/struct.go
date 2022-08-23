package cmdDaemon

import (
	"GoWebcam/Only"
	"github.com/spf13/cobra"
)


type Daemon struct {
	Error error

	cmd        *cobra.Command
	SelfCmd    *cobra.Command
}

func New() *Daemon {
	var ret *Daemon

	for range Only.Once {
		ret = &Daemon{
			Error: nil,

			cmd: nil,
			SelfCmd: nil,
		}
	}

	return ret
}
