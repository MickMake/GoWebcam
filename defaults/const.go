package defaults

import _ "embed"

// Need to execute `go generate -v -x defaults/const.go` OR `go generate -v -x ./...`
//go:generate cp ../README.md README.md
//go:generate cp ../EXAMPLES.md EXAMPLES.md

//go:embed README.md
var Readme string

//go:embed EXAMPLES.md
var Examples string

const (
	Description     = "Golang Webcam fetcher"
	BinaryName      = "GoWebcam"
	BinaryVersion   = "1.0.7"
	SourceRepo      = "github.com/MickMake/" + BinaryName
	BinaryRepo      = "github.com/MickMake/" + BinaryName

	EnvPrefix         = "WEBCAM"

	Debug           = false

	HelpTemplate = `
DefaultBinaryName - A simple automated webcam fetcher written in GoLang.

This is a simple webcam fetcher. It was an itch I needed to scratch.

What it does:
1. Regularly pull webcam images from any URL, (supports authentication).
2. Only creates a new file if the image has changed.
3. Allows for renaming of files based on time rounding, or Tesseract OCR.
4. Run scripts periodically. EG: To create MP4 videos from captured images.
5. Simple JSON config file.

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
)
