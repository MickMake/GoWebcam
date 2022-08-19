package mmSelfUpdate

import (
	"context"
	"fmt"
	"github.com/google/go-github/v30/github"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
	"strings"
)


var onlyTwice = []string{"", ""}


func (su *TypeSelfUpdate) LoadCommands(cmd *cobra.Command, disableVflag bool) error {
	for range onlyOnce {
		if cmd == nil {
			break
		}
		su.cmd = cmd

		var versionCmd = &cobra.Command{
			Use:                   CmdVersion,
			Short:                 ux.SprintfMagenta(su.Runtime.CmdName) + ux.SprintfBlue(" - Self-manage this executable."),
			Long:                  ux.SprintfMagenta(su.Runtime.CmdName) + ux.SprintfBlue(" - Self-manage this executable."),
			DisableFlagParsing:    true,
			DisableFlagsInUseLine: true,
			Run: func(cmd *cobra.Command, args []string) {
				su.OldVersion = toVersionValue(su.Runtime.CmdVersion)
				su.State = su.Version(cmd, args...)
			},
		}
		su.SelfCmd = versionCmd

		var selfUpdateCmd = &cobra.Command{
			Use:                   CmdSelfUpdate,
			Short:                 ux.SprintfMagenta(su.Runtime.CmdName) + ux.SprintfBlue(" - Update version of executable."),
			Long:                  ux.SprintfMagenta(su.Runtime.CmdName) + ux.SprintfBlue(" - Check and update the latest version."),
			DisableFlagParsing:    true,
			DisableFlagsInUseLine: true,
			Run: func(cmd *cobra.Command, args []string) {
				su.OldVersion = toVersionValue(su.Runtime.CmdVersion)
				su.State = su.VersionUpdate()
			},
		}

		su.cmd.AddCommand(versionCmd)
		su.cmd.AddCommand(selfUpdateCmd)

		var versionCheckCmd = &cobra.Command{
			Use:                   CmdVersionCheck,
			Short:                 ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Check and show any version updates."),
			Long:                  ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Check and show any version updates."),
			DisableFlagParsing:    true,
			DisableFlagsInUseLine: true,
			Run: func(cmd *cobra.Command, args []string) {
				su.OldVersion = toVersionValue(su.Runtime.CmdVersion)
				su.State = su.VersionCheck()
			},
		}
		var versionListCmd = &cobra.Command{
			Use:                   CmdVersionList,
			Short:                 ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - List available versions."),
			Long:                  ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - List available versions."),
			DisableFlagParsing:    true,
			DisableFlagsInUseLine: true,
			Run: func(cmd *cobra.Command, args []string) {
				su.OldVersion = toVersionValue(su.Runtime.CmdVersion)
				su.State = su.VersionList(args...)
			},
		}
		var versionInfoCmd = &cobra.Command{
			Use:                   CmdVersionInfo,
			Short:                 ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Info on current version."),
			Long:                  ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Info on current version."),
			DisableFlagParsing:    true,
			DisableFlagsInUseLine: true,
			Run: func(cmd *cobra.Command, args []string) {
				su.OldVersion = toVersionValue(su.Runtime.CmdVersion)
				su.State = su.VersionInfo(args...)
			},
		}
		var versionLatestCmd = &cobra.Command{
			Use:                   CmdVersionLatest,
			Short:                 ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Info on latest version."),
			Long:                  ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Info on latest version."),
			DisableFlagParsing:    true,
			DisableFlagsInUseLine: true,
			Run: func(cmd *cobra.Command, args []string) {
				su.OldVersion = toVersionValue(su.Runtime.CmdVersion)
				su.State = su.VersionInfo(CmdVersionLatest)
			},
		}
		var versionUpdateCmd = &cobra.Command{
			Use:                   CmdVersionUpdate,
			Short:                 ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Update version of this executable."),
			Long:                  ux.SprintfMagenta(CmdVersion) + ux.SprintfBlue(" - Check and update the latest version of this executable."),
			DisableFlagParsing:    true,
			DisableFlagsInUseLine: true,
			Run: func(cmd *cobra.Command, args []string) {
				su.OldVersion = toVersionValue(su.Runtime.CmdVersion)
				su.State = su.VersionUpdate()
			},
		}

		versionCmd.AddCommand(versionCheckCmd)
		versionCmd.AddCommand(versionListCmd)
		versionCmd.AddCommand(versionInfoCmd)
		versionCmd.AddCommand(versionLatestCmd)
		versionCmd.AddCommand(versionUpdateCmd)

		if !disableVflag {
			su.cmd.Flags().BoolP(FlagVersion, "v", false, ux.SprintfBlue("Display version of %s", su.Runtime.CmdName))
		}
	}

	return su.State
}


func (su *TypeSelfUpdate) Version(cmd *cobra.Command, args ...string) error {
	for range onlyOnce {
		su.VersionShow()
		su.SetHelp(cmd)
		_ = cmd.Help()
		su.State.SetOk()
	}
	return su.State
}


func (su *TypeSelfUpdate) VersionShow() error {
	su.Runtime.PrintNameVersion()
	su.State.SetOk()
	return su.State
}


func (su *TypeSelfUpdate) VersionInfo(args ...string) error {
	for range onlyOnce {
		if len(args) == 0 {
			args = []string{CmdVersionLatest}
		}

		for _, v := range args {
			fv := GetSemVer(v)
			su.State = su.PrintVersion(fv)
			if su.State.IsNotOk() {
				break
			}
		}
	}

	return su.State
}


func (su *TypeSelfUpdate) VersionList(args ...string) error {
	for range onlyOnce {
		token := os.Getenv("GITHUB_TOKEN")
		if token == "" {
			token, _ = gitconfig.GithubToken()
		}

		var err error
		gh := github.NewClient(newHTTPClient(context.Background(), token))
		var rels []*github.RepositoryRelease
		rels, _, err = gh.Repositories.ListReleases(context.Background(), su.useRepo.Owner, su.useRepo.Name, nil)
		if err != nil {
			su.State.SetError(err)
			break
		}

		for _, rel := range rels {
			su.State = su.PrintVersionSummary(*rel.TagName)
			if su.State.IsOk() {
				continue
			}
			su.State = su.PrintVersionSummary(*rel.Name)
			if su.State.IsOk() {
				continue
			}
			// WORKAROUND: (selfupdate) - If selfupdate.Release returns nil, then print direct.
			su.State = su.PrintSummary(rel)
		}
	}

	return su.State
}


func (su *TypeSelfUpdate) VersionCheck() error {
	for range onlyOnce {
		su.State = su.IsUpdated(true)
		if su.State.IsError() {
			break
		}
	}

	return su.State
}


func (su *TypeSelfUpdate) VersionUpdate() error {
	for range onlyOnce {
		su.State = su.IsUpdated(true)
		if su.State.IsError() {
			break
		}

		su.State = su.CreateDummyBinary()
		if su.State.IsError() {
			// Only break on error, NOT warning.
			break
		}

		su.State = su.UpdateTo(su.GetVersionValue())
		if su.State.IsNotOk() {
			break
		}

		if !su.AutoExec {
			break
		}

		// AutoExec will execute the new binary with the same args as given.
		su.State = su.AutoRun()
		if su.State.IsNotOk() {
			break
		}
	}

	return su.State
}


func (su *TypeSelfUpdate) FlagCheckVersion(cmd *cobra.Command) bool {
	var ok bool
	for range onlyOnce {
		var fl *pflag.FlagSet
		if cmd == nil {
			fl = su.cmd.Flags()
		} else {
			fl = cmd.Flags()
		}

		// Show version.
		ok, _ = fl.GetBool(FlagVersion)
		if ok {
			su.VersionShow()
			break
		}
	}
	return ok
}


func (su *TypeSelfUpdate) GetCmd() *cobra.Command {
	var ret *cobra.Command
	if state := su.IsNil(); state.IsError() {
		return ret
	}
	return su.SelfCmd
}


func (su *TypeSelfUpdate) CmdHelp() error {
	if state := su.IsNil(); state.IsError() {
		return state
	}

	err := su.SelfCmd.Help()
	if err != nil {
		su.State.SetError(err)
	}
	return su.State
}


func _GetUsage(c *cobra.Command) string {
	var str string

	if c.Parent() == nil {
		str += ux.SprintfCyan("%s ", c.Name())
	} else {
		str += ux.SprintfCyan("%s ", c.Parent().Name())
		str += ux.SprintfGreen("%s ", c.Use)
	}

	if c.HasAvailableSubCommands() {
		str += ux.SprintfGreen("[command] ")
		str += ux.SprintfCyan("<args> ")
	}

	return str
}


func _GetVersion(c *cobra.Command) string {
	var str string

	if c.Parent() == nil {
		//str = ux.SprintfBlue("%s ", su.Runtime.CmdName)
		//str += ux.SprintfCyan("v%s", su.Runtime.CmdVersion)
	}

	return str
}


func (su *TypeSelfUpdate) SetHelp(c *cobra.Command) {
	var tmplHelp string
	var tmplUsage string

	cobra.AddTemplateFunc("GetUsage", _GetUsage)
	cobra.AddTemplateFunc("GetVersion", _GetVersion)

	cobra.AddTemplateFunc("SprintfBlue", ux.SprintfBlue)
	cobra.AddTemplateFunc("SprintfCyan", ux.SprintfCyan)
	cobra.AddTemplateFunc("SprintfGreen", ux.SprintfGreen)
	cobra.AddTemplateFunc("SprintfMagenta", ux.SprintfMagenta)
	cobra.AddTemplateFunc("SprintfRed", ux.SprintfRed)
	cobra.AddTemplateFunc("SprintfWhite", ux.SprintfWhite)
	cobra.AddTemplateFunc("SprintfYellow", ux.SprintfYellow)

	tmplUsage += `
{{ SprintfBlue "Usage: " }}
	{{ GetUsage . }}

{{- if gt (len .Aliases) 0 }}
{{ SprintfBlue "\nAliases:" }} {{ .NameAndAliases }}
{{- end }}

{{- if .HasExample }}
{{ SprintfBlue "\nExamples:" }}
	{{ .Example }}
{{- end }}

{{- if .HasAvailableSubCommands }}
{{ SprintfBlue "\nWhere " }}{{ SprintfGreen "[command]" }}{{ SprintfBlue " is one of:" }}
{{- range .Commands }}
{{- if (or .IsAvailableCommand (eq .Name "help")) }}
	{{ rpad (SprintfGreen .Name) .NamePadding}}	- {{ .Short }}{{ end }}
{{- end }}
{{- end }}

{{- if .HasHelpSubCommands }}
{{- SprintfBlue "\nAdditional help topics:" }}
{{- range .Commands }}
{{- if .IsAdditionalHelpTopicCommand }}
	{{ rpad (SprintfGreen .CommandPath) .CommandPathPadding }} {{ .Short }}
{{- end }}
{{- end }}
{{- end }}

{{- if .HasAvailableSubCommands }}
{{ SprintfBlue "\nUse" }} {{ SprintfCyan .CommandPath }} {{ SprintfCyan "help" }} {{ SprintfGreen "[command]" }} {{ SprintfBlue "for more information about a command." }}
{{- end }}
`

	tmplHelp = `{{ GetVersion . }}

{{ SprintfBlue "Commmand:" }} {{ SprintfCyan .Use }}

{{ SprintfBlue "Description:" }} 
	{{ with (or .Long .Short) }}
{{- . | trimTrailingWhitespaces }}
{{- end }}

{{- if or .Runnable .HasSubCommands }}
{{ .UsageString }}
{{- end }}
`

	//c.SetHelpCommand(c)
	//c.SetHelpFunc(PrintHelp)
	c.SetHelpTemplate(tmplHelp)
	c.SetUsageTemplate(tmplUsage)
}


type Example struct {
	Command string
	Args []string
	Info string
}
type Examples []Example


func (su *TypeSelfUpdate) HelpExamples() {
	su.State.SetOk()
	return

	for range onlyOnce {
		var examples Examples


		examples = append(examples, Example {
			Command: "load",
			Args:    []string{"-json", "config.json", "-template", "'{{ .Json.dir }}'"},
			Info:    "Print to STDOUT the .dir key from config.json.",
		})
		examples = append(examples, Example {
			Command: "load",
			Args:    []string{"-template", "'{{ .Json.dir }}'", "config.json"},
			Info:    "The same thing, but with less arguments.",
		})

		examples = append(examples, Example {
			Command: "load",
			Args:    []string{"-template", "'{{ .Json.hello }}'", "-json", `'{ "hello": "world" }'`},
			Info:    "Template and JSON arguments can be either string or file reference.",
		})
		examples = append(examples, Example {
			Command: "load",
			Args:    []string{"-template", "hello_world.tmpl", "-json", `'{ "hello": "world" }'`},
			Info:    "The same again...",
		})
		examples = append(examples, Example {
			Command: "load",
			Args:    []string{"-template", "'{{ .Json.hello }}'", "-json", "hello.json"},
			Info:    "The same again...",
		})


		examples = append(examples, Example {
			Command: "load",
			Args:    []string{"-json", "config.json", "-template", "DockerFile.tmpl", "-out", "Dockerfile"},
			Info:    "Process Dockerfile.tmpl file and output to Dockerfile.",
		})
		examples = append(examples, Example {
			Command: "load",
			Args:    []string{"-out", "Dockerfile", "config.json", "DockerFile.tmpl"},
			Info:    "And again with less arguments..",
		})
		examples = append(examples, Example {
			Command: "convert",
			Args:    []string{"config.json", "DockerFile.tmpl"},
			Info:    "'convert' does the same , but removes the template file afterwards...",
		})


		examples = append(examples, Example {
			Command: "load",
			Args:    []string{"-out", "MyScript.sh", "MyScript.sh.tmpl", "config.json"},
			Info:    "Process the MyScript.sh.tmpl template file and write the result to MyScript.sh.",
		})
		examples = append(examples, Example {
			Command: "convert",
			Args:    []string{"MyScript.sh.tmpl", "config.json"},
			Info:    "Same again using 'convert'. Template and json files can be in any order.",
		})
		examples = append(examples, Example {
			Command: "run",
			Args:    []string{"MyScript.sh.tmpl", "config.json"},
			Info:    "Same again using 'run'. This will execute the MyScript.sh output file afterwards.",
		})


		ux.PrintflnBlue("Examples:")
		for _, v := range examples {
			fmt.Printf("# %s\n\t%s %s\n\n",
				ux.SprintfBlue(v.Info),
				ux.SprintfCyan("%s %s", su.Runtime.CmdName, v.Command),
				ux.SprintfWhite(strings.Join(v.Args, " ")),
			)
		}
	}

	su.State.SetOk()
}
