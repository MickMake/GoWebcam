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
	BinaryVersion   = "1.0.8"
	SourceRepo      = "github.com/MickMake/" + BinaryName
	BinaryRepo      = "github.com/MickMake/" + BinaryName

	EnvPrefix         = "WEBCAM"

	Debug           = false

	HelpSummary = `
DefaultBinaryName - A simple automated webcam fetcher written in GoLang.

This is a simple webcam fetcher. It was an itch I needed to scratch.

What it does:
1. Regularly pull webcam images from any URL, (supports authentication).
2. Only creates a new file if the image has changed.
3. Allows for renaming of files based on time rounding, or Tesseract OCR.
4. Run scripts periodically. EG: To create MP4 videos from captured images.
5. Simple JSON config file.
`
)
