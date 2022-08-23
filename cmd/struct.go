package cmd

import (
	"GoWebcam/Only"
	"GoWebcam/Unify/cmdConfig"
	"GoWebcam/Unify/cmdCron"
	"GoWebcam/Unify/cmdDaemon"
	"GoWebcam/Unify/cmdHelp"
	"GoWebcam/Unify/cmdVersion"
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

	// Webcams fetching
	WebHost     string
	WebPort     string
	WebUsername string
	WebPassword string
	WebTimeout  time.Duration
	WebPrefix   string

	CmdWebcams *Webcams
	CmdVersion *cmdVersion.Version
	CmdDaemon  *cmdDaemon.Daemon
	CmdCron    *cmdCron.Cron
	CmdConfig  *cmdConfig.Config
	CmdHelp    *cmdHelp.Help

	Valid bool
	Error error
}

var Cmd        CommandArgs
// var Webcams    *mmWebcam.Config
// var CmdWebcams *Webcams
// var CmdVersion *cmdVersion.Version
// var CmdDaemon  *cmdDaemon.Daemon
// var CmdCron    *cmdCron.Cron
// var CmdConfig  *cmdConfig.Config
// var CmdHelp    *cmdHelp.Help


func (ca *CommandArgs) IsValid() error {
	for range Only.Once {
		if !ca.Valid {
			ca.Error = errors.New("args are not valid")
			break
		}
	}

	return ca.Error
}

func (ca *CommandArgs) ProcessArgs(_ *cobra.Command, _ []string) error {
	for range Only.Once {
		// ca.Args = args

		ca.Valid = true
	}

	return ca.Error
}
