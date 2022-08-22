package cmd


const DefaultHelpTemplate = `
{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`

const DefaultUsageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags: Use "{{.Root.CommandPath}} help-all" for more info.

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand }}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} help [command]" for more information about a command.{{end}}
`

// const DefaultVersionTemplate = `
// {{with .Name}}{{printf "%s " .}}{{end}}{{printf "version %s" .Version}}
// `

//goland:noinspection GoUnusedConst
const DefaultFlagHelpTemplate = `{{if .HasAvailableInheritedFlags}}Flags available for all commands:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}
`

const ExtendedHelpText = `
DefaultBinaryName - A simple automated webcam fetcher written in GoLang.

This tool is a basic webcam image puller. It can:
1. Pull any arbitrary image.
2. Can handle username/passwords.
3. Rename image files based on rounding, or tesseract OCR.
4. Create movies of images.

Use case example:
# Simple cron pulling an image every 5 minutes.
	% DefaultBinaryName cron run . ./5 . . . . web get Basin https://charlottepass.com.au/charlottepass/webcam/lucylodge/current.jpg

# Pull webcam images as defined in config.json file, via cron.
	% DefaultBinaryName web cron

# Once-off run of all webcams defined in config.json file.
	% DefaultBinaryName web run

# Config file.
	Show current config.
	% DefaultBinaryName config read
`