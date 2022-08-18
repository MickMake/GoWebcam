package cmd

import (
	"GoWebcam/Only"
	"GoWebcam/mmWebcam"
	"errors"
	"github.com/spf13/cobra"
	"time"
)


type CommandArgs struct {
	ConfigDir   string
	CacheDir    string
	ConfigFile  string
	WriteConfig bool
	Quiet       bool
	Debug       bool
	OutputType  string
	OutputFile  string

	// Web fetching
	WebHost     string
	WebPort     string
	WebUsername string
	WebPassword string
	WebTimeout  time.Duration
	WebPrefix   string

	Args []string

	Valid bool
	Error error
}

var Cmd CommandArgs
var Webcams     *mmWebcam.Config

func (ca *CommandArgs) IsValid() error {
	for range Only.Once {
		if !ca.Valid {
			ca.Error = errors.New("args are not valid")
			break
		}
	}

	return ca.Error
}

func (ca *CommandArgs) ProcessArgs(_ *cobra.Command, args []string) error {
	for range Only.Once {
		ca.Args = args

		ca.Valid = true
	}

	return ca.Error
}
