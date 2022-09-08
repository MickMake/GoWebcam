package main

import (
	"GoWebcam/defaults"
	"GoWebcam/mmWebcam"
	"fmt"
	"github.com/MickMake/GoUnify/Unify"
	"github.com/MickMake/GoUnify/Only"
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
				Description:   defaults.Description,
				BinaryName:    defaults.BinaryName,
				BinaryVersion: defaults.BinaryVersion,
				SourceRepo:    defaults.SourceRepo,
				BinaryRepo:    defaults.BinaryRepo,
				EnvPrefix:     defaults.EnvPrefix,
				HelpTemplate:  defaults.HelpTemplate,
				ReadMe:        defaults.Readme,
				Examples:      defaults.Examples,
			},
		Unify.Flags {},
		)

		wc := mmWebcam.New()
		wc.AttachCommands(unify.GetCmd())

		err = unify.Execute()
		if err != nil {
			break
		}
	}

	return err
}
