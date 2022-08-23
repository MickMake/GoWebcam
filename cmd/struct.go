package cmd

import (
	"GoWebcam/Only"
	"GoWebcam/Unify/cmdConfig"
	"GoWebcam/Unify/cmdCron"
	"GoWebcam/Unify/cmdDaemon"
	"GoWebcam/Unify/cmdHelp"
	"GoWebcam/Unify/cmdVersion"
	"GoWebcam/mmWebcam"
	"errors"
	"github.com/spf13/cobra"
	"time"
)


type CommandArgs struct {
	ConfigFile  string
	ConfigDir   string
	CacheDir    string

	Quiet       bool
	Debug       bool
	Daemonize	bool

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

var Cmd        CommandArgs
var Webcams    *mmWebcam.Config
var CmdVersion *cmdVersion.Version
var CmdDaemon  *cmdDaemon.Daemon
var CmdCron    *cmdCron.Cron
var CmdConfig  *cmdConfig.Config
var CmdHelp    *cmdHelp.Help


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
