package cmd

import (
	"GoWebcam/Only"
	"fmt"
	"github.com/spf13/cobra"
)


func AttachRootCmd(cmd *cobra.Command) *cobra.Command {
	// ******************************************************************************** //
	var rootCmd = &cobra.Command{
		Use:              DefaultBinaryName,
		Short:            fmt.Sprintf("%s - Webcam fetcher", DefaultBinaryName),
		Long:             fmt.Sprintf("%s - Webcam fetcher", DefaultBinaryName),
		Run:              gbRootFunc,
		TraverseChildren: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initConfig(cmd)
		},
	}
	if cmd != nil {
		cmd.AddCommand(rootCmd)
	}
	rootCmd.Example = PrintExamples(rootCmd, "")

	rootCmd.SetHelpTemplate(DefaultHelpTemplate)
	rootCmd.SetUsageTemplate(DefaultUsageTemplate)
	rootCmd.SetVersionTemplate(DefaultVersionTemplate)

	rootCmd.PersistentFlags().StringVarP(&Cmd.WebHost, flagWebHost, "", defaultHost, fmt.Sprintf("Web Host."))
	rootViper.SetDefault(flagWebHost, defaultHost)
	rootCmd.PersistentFlags().StringVarP(&Cmd.WebPort, flagWebPort, "", defaultPort, fmt.Sprintf("Web Port."))
	rootViper.SetDefault(flagWebPort, defaultPort)
	rootCmd.PersistentFlags().StringVarP(&Cmd.WebUsername, flagWebUsername, "u", defaultUsername, fmt.Sprintf("Web username."))
	rootViper.SetDefault(flagWebUsername, defaultUsername)
	rootCmd.PersistentFlags().StringVarP(&Cmd.WebPassword, flagWebPassword, "p", defaultPassword, fmt.Sprintf("Web password."))
	rootViper.SetDefault(flagWebPassword, defaultPassword)
	rootCmd.PersistentFlags().DurationVarP(&Cmd.WebTimeout, flagWebTimeout, "t", defaultTimeout, fmt.Sprintf("Web timeout."))
	rootViper.SetDefault(flagWebPassword, defaultTimeout)
	rootCmd.PersistentFlags().StringVarP(&Cmd.WebPrefix, flagWebPrefix, "", defaultPrefix, fmt.Sprintf("Web password."))
	rootViper.SetDefault(flagWebPrefix, defaultPrefix)

	rootCmd.PersistentFlags().StringVar(&Cmd.ConfigFile, flagConfigFile, Cmd.ConfigFile, fmt.Sprintf("%s: config file.", DefaultBinaryName))
	// _ = rootCmd.PersistentFlags().MarkHidden(flagConfigFile)
	rootCmd.PersistentFlags().BoolVarP(&Cmd.Debug, flagDebug, "", false, fmt.Sprintf("%s: Debug mode.", DefaultBinaryName))
	rootViper.SetDefault(flagDebug, false)
	rootCmd.PersistentFlags().BoolVarP(&Cmd.Quiet, flagQuiet, "q", false, fmt.Sprintf("%s: Silence all messages.", DefaultBinaryName))
	rootViper.SetDefault(flagQuiet, false)

	rootCmd.PersistentFlags().SortFlags = false
	rootCmd.Flags().SortFlags = false

	return rootCmd
}

func gbRootFunc(cmd *cobra.Command, args []string) {
	for range Only.Once {
		if len(args) == 0 {
			_ = cmd.Help()
			break
		}
	}
}
