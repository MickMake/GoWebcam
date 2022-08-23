package cmd

import "time"

//goland:noinspection SpellCheckingInspection
const (
	flagConfigFile = "config"
	flagDebug      = "debug"
	flagQuiet      = "quiet"

	flagWebHost     = "host"
	flagWebPort     = "port"
	flagWebUsername = "user"
	flagWebPassword = "password"
	flagWebTimeout  = "timeout"
	flagWebPrefix   = "prefix"

	defaultHost     = ""
	defaultPort     = ""
	defaultUsername = ""
	defaultPassword = ""
	defaultTimeout  = time.Second * 30
	defaultPrefix   = ""
)


const ExtendedHelpTemplate = `
DefaultBinaryName - A simple automated webcam fetcher written in GoLang.

This tool is a basic webcam image puller. It can:
1. Pull any arbitrary image.
2. Can handle username/passwords.
3. Rename image files based on rounding, or tesseract OCR.
4. Create movies of images.

Use case example:
# Simple cron pulling an image every 5 minutes.
	% DefaultBinaryName cron run  . ./5 . . . .  web get Basin https://charlottepass.com.au/charlottepass/webcam/lucylodge/current.jpg

# Once-off run of all webcams defined in config.json file.
	% DefaultBinaryName web run

# Pull webcam images as defined in config.json file, via cron.
	% DefaultBinaryName web cron

# Same as above, but run as a daemon.
	% DefaultBinaryName daemon exec  web cron

# List currently scheduled jobs.
	% DefaultBinaryName daemon list

# Config file.
	Show current config.
	% DefaultBinaryName config read
`
