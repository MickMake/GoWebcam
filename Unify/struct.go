package Unify

import (
	"GoWebcam/Only"
	"GoWebcam/Unify/cmdConfig"
	"GoWebcam/Unify/cmdCron"
	"GoWebcam/Unify/cmdDaemon"
	"GoWebcam/Unify/cmdHelp"
	"GoWebcam/Unify/cmdVersion"
	"GoWebcam/defaults"
	"fmt"
	"github.com/spf13/cobra"
	"time"
)


func New(options Options, flags Flags, ) *Unify {
	var unify Unify

	for range Only.Once {
		unify.Options = options
		unify.Flags = flags


		unify.Error = unify.InitCmds()
		if unify.Error != nil {
			break
		}

		unify.Error = unify.InitFlags()
		if unify.Error != nil {
			break
		}
	}

	return &unify
}

func (u *Unify) InitCmds() error {

	for range Only.Once {
		// ******************************************************************************** //
		u.Commands.CmdRoot = &cobra.Command {
			Use:              defaults.BinaryName,
			Short:            fmt.Sprintf("%s - Webcam fetcher", defaults.BinaryName),
			Long:             fmt.Sprintf("%s - Webcam fetcher", defaults.BinaryName),
			Run:              CmdRoot,
			TraverseChildren: true,
		}
		u.Commands.CmdRoot.Example = cmdHelp.PrintExamples(u.Commands.CmdRoot, "")

		u.Commands.CmdVersion = cmdVersion.New(u.Options.BinaryName, u.Options.BinaryVersion, false)
		u.Commands.CmdVersion.SetBinaryRepo(u.Options.BinaryRepo)
		u.Commands.CmdVersion.SetSourceRepo(u.Options.SourceRepo)

		u.Commands.CmdDaemon = cmdDaemon.New()

		u.Commands.CmdCron = cmdCron.New()

		u.Commands.CmdConfig = cmdConfig.New(u.Options.BinaryName)

		u.Commands.CmdHelp = cmdHelp.New()
		u.Commands.CmdHelp.SetCommand(u.Options.BinaryName)
		u.Commands.CmdHelp.SetExtendedHelpTemplate(u.Options.HelpTemplate)
	}

	return u.Error
}

func (u *Unify) InitFlags() error {

	for range Only.Once {
		// SelfCmd.PersistentFlags().StringVarP(&Cmd.WebHost, flagWebHost, "", defaultHost, fmt.Sprintf("Web Host."))
		// Cmd.CmdConfig.SetDefault(flagWebHost, defaultHost)
		// SelfCmd.PersistentFlags().StringVarP(&Cmd.WebPort, flagWebPort, "", defaultPort, fmt.Sprintf("Web Port."))
		// Cmd.CmdConfig.SetDefault(flagWebPort, defaultPort)
		// SelfCmd.PersistentFlags().StringVarP(&Cmd.WebUsername, flagWebUsername, "u", defaultUsername, fmt.Sprintf("Web username."))
		// Cmd.CmdConfig.SetDefault(flagWebUsername, defaultUsername)
		// SelfCmd.PersistentFlags().StringVarP(&Cmd.WebPassword, flagWebPassword, "p", defaultPassword, fmt.Sprintf("Web password."))
		// Cmd.CmdConfig.SetDefault(flagWebPassword, defaultPassword)
		// SelfCmd.PersistentFlags().StringVarP(&Cmd.WebPrefix, flagWebPrefix, "", defaultPrefix, fmt.Sprintf("Web password."))
		// Cmd.CmdConfig.SetDefault(flagWebPrefix, defaultPrefix)

		u.Commands.CmdRoot.PersistentFlags().StringVar(&u.Flags.ConfigFile, flagConfigFile, defaultConfig, fmt.Sprintf("%s: config file.", defaults.BinaryName))
		// _ = rootCmd.PersistentFlags().MarkHidden(flagConfigFile)
		u.Commands.CmdRoot.PersistentFlags().BoolVarP(&u.Flags.Debug, flagDebug, "", defaultDebug, fmt.Sprintf("%s: Debug mode.", defaults.BinaryName))
		u.Commands.CmdConfig.SetDefault(flagDebug, false)
		u.Commands.CmdRoot.PersistentFlags().BoolVarP(&u.Flags.Quiet, flagQuiet, "", defaultQuiet, fmt.Sprintf("%s: Silence all messages.", defaults.BinaryName))
		u.Commands.CmdConfig.SetDefault(flagQuiet, false)
		u.Commands.CmdRoot.PersistentFlags().DurationVarP(&u.Flags.Timeout, flagTimeout, "", defaultTimeout, fmt.Sprintf("Web timeout."))
		u.Commands.CmdConfig.SetDefault(flagTimeout, defaultTimeout)

		u.Commands.CmdRoot.PersistentFlags().SortFlags = false
		u.Commands.CmdRoot.Flags().SortFlags = false

		// cobra.OnInitialize(initConfig)	// Bound to rootCmd now.
		cobra.EnableCommandSorting = false
	}

	return u.Error
}

func (u *Unify) Execute() error {
	var err error
	for range Only.Once {
		u.Commands.CmdVersion.AttachCommands(u.Commands.CmdRoot, false)
		u.Commands.CmdDaemon.AttachCommands(u.Commands.CmdRoot)
		u.Commands.CmdCron.AttachCommands(u.Commands.CmdRoot)
		u.Commands.CmdConfig.AttachCommands(u.Commands.CmdRoot)
		u.Commands.CmdHelp.AttachCommands(u.Commands.CmdRoot)
		u.Commands.CmdConfig.SetDir(u.Flags.ConfigDir)
		u.Commands.CmdConfig.SetFile(u.Flags.ConfigFile)
		u.Commands.CmdRoot.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return u.Commands.CmdConfig.Init(cmd)
		}

		err = u.Commands.Execute()
		if err != nil {
			break
		}
	}
	return err
}

func (u *Unify) GetRootCmd() *cobra.Command {
	var ret *cobra.Command
	for range Only.Once {
		ret = u.Commands.CmdRoot
	}
	return ret
}

func (c *Commands) Execute() error {
	return c.CmdRoot.Execute()
}

func CmdRoot(cmd *cobra.Command, args []string) {
	for range Only.Once {
		if len(args) == 0 {
			_ = cmd.Help()
			break
		}
	}
}


type Unify struct {
	Options  Options  `json:"options"`
	Flags    Flags    `json:"flags"`
	Commands Commands `json:"commands"`

	Error    error    `json:"-"`
}

type Options struct {
	BinaryName    string `json:"binary_name"`
	BinaryVersion string `json:"binary_version"`
	SourceRepo    string `json:"source_repo"`
	BinaryRepo    string `json:"binary_repo"`
	EnvPrefix     string `json:"env_prefix"`
	HelpTemplate  string `json:"help_template"`
}

type Flags struct {
	ConfigFile string        `json:"config_file"`
	ConfigDir  string        `json:"config_dir"`
	CacheDir   string        `json:"cache_dir"`
	Quiet      bool          `json:"quiet"`
	Debug      bool          `json:"debug"`
	Timeout    time.Duration `json:"timeout"`
}

type Commands struct {
	CmdRoot    *cobra.Command
	CmdVersion *cmdVersion.Version
	CmdDaemon  *cmdDaemon.Daemon
	CmdCron    *cmdCron.Cron
	CmdConfig  *cmdConfig.Config
	CmdHelp    *cmdHelp.Help
}

// func (c *Commands) IsValid() error {
// 	for range Only.Once {
// 		if !c.Valid {
// 			c.Error = errors.New("args are not valid")
// 			break
// 		}
// 	}
//
// 	return c.Error
// }
//
// func (c *Commands) ProcessArgs(_ *cobra.Command, _ []string) error {
// 	for range Only.Once {
// 		// ca.Args = args
//
// 		c.Valid = true
// 	}
//
// 	return c.Error
// }
