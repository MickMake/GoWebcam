package main

import (
	"GoWebcam/Only"
	"GoWebcam/Unify"
	"GoWebcam/defaults"
	"GoWebcam/mmWebcam"
	"fmt"
	"os"
)


func main() {
	var err error

	for range Only.Once {
		err = Execute()
		if err != nil {
			break
		}

	}

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	var err error

	for range Only.Once {
		unify := Unify.New(
			Unify.Options {
				BinaryName:    defaults.BinaryName,
				BinaryVersion: defaults.BinaryVersion,
				SourceRepo:    defaults.SourceRepo,
				BinaryRepo:    defaults.BinaryRepo,
				EnvPrefix:     defaults.EnvPrefix,
				HelpTemplate:  defaults.HelpTemplate,
			},
			Unify.Flags {},
		)

		wc := mmWebcam.New()
		wc.AttachCommands(unify.GetRootCmd())

		err = unify.Execute()
		if err != nil {
			break
		}
		// if Cmd.Error != nil {
		// 	break
		// }
	}

	return err
}
