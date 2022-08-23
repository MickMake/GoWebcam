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


func AttachRootCmd() *cobra.Command {
	var SelfCmd *cobra.Command
	for range Only.Once {
		if Cmd.CmdVersion == nil {
			Cmd.CmdWebcams = NewWeb()
		}

		if Cmd.CmdVersion == nil {
			Cmd.CmdVersion = cmdVersion.New(defaults.BinaryName, defaults.BinaryVersion, defaults.Debug)
			Cmd.CmdVersion.SetBinaryRepo(defaults.BinaryRepo)
			Cmd.CmdVersion.SetSourceRepo(defaults.SourceRepo)
		}

		if Cmd.CmdDaemon == nil {
			Cmd.CmdDaemon = cmdDaemon.New()
		}

		if Cmd.CmdCron == nil {
			Cmd.CmdCron = cmdCron.New()
		}

		if Cmd.CmdConfig == nil {
			Cmd.CmdConfig = cmdConfig.New(defaults.BinaryName)
		}

		if Cmd.CmdHelp == nil {
			Cmd.CmdHelp = cmdHelp.New()
			Cmd.CmdHelp.SetCommand(defaults.BinaryName)
			Cmd.CmdHelp.SetExtendedHelpTemplate(ExtendedHelpTemplate)
		}

		// ******************************************************************************** //
		SelfCmd = &cobra.Command{
			Use:              defaults.BinaryName,
			Short:            fmt.Sprintf("%s - Webcam fetcher", defaults.BinaryName),
			Long:             fmt.Sprintf("%s - Webcam fetcher", defaults.BinaryName),
			Run:              CmdRoot,
			TraverseChildren: true,
		}
		SelfCmd.Example = cmdHelp.PrintExamples(SelfCmd, "")

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

		SelfCmd.PersistentFlags().StringVar(&Cmd.ConfigFile, flagConfigFile, Cmd.ConfigFile, fmt.Sprintf("%s: config file.", defaults.BinaryName))
		// _ = rootCmd.PersistentFlags().MarkHidden(flagConfigFile)
		SelfCmd.PersistentFlags().BoolVarP(&Cmd.Debug, flagDebug, "", false, fmt.Sprintf("%s: Debug mode.", defaults.BinaryName))
		Cmd.CmdConfig.SetDefault(flagDebug, false)
		SelfCmd.PersistentFlags().BoolVarP(&Cmd.Quiet, flagQuiet, "q", false, fmt.Sprintf("%s: Silence all messages.", defaults.BinaryName))
		Cmd.CmdConfig.SetDefault(flagQuiet, false)
		SelfCmd.PersistentFlags().DurationVarP(&Cmd.WebTimeout, flagWebTimeout, "t", defaultTimeout, fmt.Sprintf("Web timeout."))
		Cmd.CmdConfig.SetDefault(flagWebTimeout, defaultTimeout)

		SelfCmd.PersistentFlags().SortFlags = false
		SelfCmd.Flags().SortFlags = false

		Cmd.CmdWebcams.AttachCmdWeb(SelfCmd)
		Cmd.CmdVersion.AttachCommands(SelfCmd, false)
		Cmd.CmdDaemon.AttachCommands(SelfCmd)
		Cmd.CmdCron.AttachCommands(SelfCmd)
		Cmd.CmdConfig.AttachCommands(SelfCmd)
		Cmd.CmdHelp.AttachCommands(SelfCmd)
		Cmd.CmdConfig.SetDir(Cmd.ConfigDir)
		Cmd.CmdConfig.SetFile(Cmd.ConfigFile)
		SelfCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return Cmd.CmdConfig.Init(cmd)
		}
	}

	return SelfCmd
}

func CmdRoot(cmd *cobra.Command, args []string) {
	for range Only.Once {
		if len(args) == 0 {
			_ = cmd.Help()
			break
		}
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	for range Only.Once {
		rootCmd := AttachRootCmd()

		// cobra.OnInitialize(initConfig)	// Bound to rootCmd now.
		cobra.EnableCommandSorting = false

		err := rootCmd.Execute()
		if err != nil {
			break
		}
		if Cmd.Error != nil {
			break
		}
	}

	return Cmd.Error
}
