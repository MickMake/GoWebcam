package cmd

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
)


func AttachRootCmd(cmd *cobra.Command) *cobra.Command {

	if CmdVersion == nil {
		CmdVersion = cmdVersion.New(defaults.BinaryName, defaults.BinaryVersion, defaults.Debug)
		CmdVersion.SetBinaryRepo(defaults.BinaryRepo)
		CmdVersion.SetSourceRepo(defaults.SourceRepo)
	}

	if CmdDaemon == nil {
		CmdDaemon = cmdDaemon.New()
	}

	if CmdCron == nil {
		CmdCron = cmdCron.New()
	}

	if CmdConfig == nil {
		CmdConfig = cmdConfig.New()
	}

	if CmdHelp == nil {
		CmdHelp = cmdHelp.New()
	}

	// ******************************************************************************** //
	var rootCmd = &cobra.Command {
		Use:              defaults.BinaryName,
		Short:            fmt.Sprintf("%s - Webcam fetcher", defaults.BinaryName),
		Long:             fmt.Sprintf("%s - Webcam fetcher", defaults.BinaryName),
		Run:              CmdRoot,
		TraverseChildren: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return CmdConfig.Init(cmd)
		},
	}
	if cmd != nil {
		cmd.AddCommand(rootCmd)
	}
	rootCmd.Example = cmdHelp.PrintExamples(rootCmd, "")

	rootCmd.SetHelpTemplate(cmdHelp.DefaultHelpTemplate)
	rootCmd.SetUsageTemplate(cmdHelp.DefaultUsageTemplate)
	// rootCmd.SetVersionTemplate(DefaultVersionTemplate)

	rootCmd.PersistentFlags().StringVarP(&Cmd.WebHost, flagWebHost, "", defaultHost, fmt.Sprintf("Web Host."))
	CmdConfig.SetDefault(flagWebHost, defaultHost)
	rootCmd.PersistentFlags().StringVarP(&Cmd.WebPort, flagWebPort, "", defaultPort, fmt.Sprintf("Web Port."))
	CmdConfig.SetDefault(flagWebPort, defaultPort)
	rootCmd.PersistentFlags().StringVarP(&Cmd.WebUsername, flagWebUsername, "u", defaultUsername, fmt.Sprintf("Web username."))
	CmdConfig.SetDefault(flagWebUsername, defaultUsername)
	rootCmd.PersistentFlags().StringVarP(&Cmd.WebPassword, flagWebPassword, "p", defaultPassword, fmt.Sprintf("Web password."))
	CmdConfig.SetDefault(flagWebPassword, defaultPassword)
	rootCmd.PersistentFlags().DurationVarP(&Cmd.WebTimeout, flagWebTimeout, "t", defaultTimeout, fmt.Sprintf("Web timeout."))
	CmdConfig.SetDefault(flagWebTimeout, defaultTimeout)
	rootCmd.PersistentFlags().StringVarP(&Cmd.WebPrefix, flagWebPrefix, "", defaultPrefix, fmt.Sprintf("Web password."))
	CmdConfig.SetDefault(flagWebPrefix, defaultPrefix)

	// rootCmd.PersistentFlags().BoolVarP(&Cmd.Daemonize, flagDaemonize, "d", false, fmt.Sprintf("Daemonize program."))
	// rootViper.SetDefault(flagDaemonize, false)
	rootCmd.PersistentFlags().StringVar(&Cmd.ConfigFile, flagConfigFile, Cmd.ConfigFile, fmt.Sprintf("%s: config file.", defaults.BinaryName))
	// _ = rootCmd.PersistentFlags().MarkHidden(flagConfigFile)
	rootCmd.PersistentFlags().BoolVarP(&Cmd.Debug, flagDebug, "", false, fmt.Sprintf("%s: Debug mode.", defaults.BinaryName))
	CmdConfig.SetDefault(flagDebug, false)
	rootCmd.PersistentFlags().BoolVarP(&Cmd.Quiet, flagQuiet, "q", false, fmt.Sprintf("%s: Silence all messages.", defaults.BinaryName))
	CmdConfig.SetDefault(flagQuiet, false)

	rootCmd.PersistentFlags().SortFlags = false
	rootCmd.Flags().SortFlags = false

	_ = AttachCmdWeb(rootCmd)
	CmdVersion.AttachCommands(rootCmd, false)
	CmdDaemon.AttachCommands(rootCmd)
	CmdCron.AttachCommands(rootCmd)
	CmdConfig.AttachCommands(rootCmd)
	CmdHelp.AttachCommands(rootCmd)
	CmdConfig.SetDir(Cmd.ConfigDir)
	CmdConfig.SetFile(Cmd.ConfigFile)

	return rootCmd
}

func CmdRoot(cmd *cobra.Command, args []string) {
	for range Only.Once {
		if len(args) == 0 {
			_ = cmd.Help()
			break
		}
	}
}
