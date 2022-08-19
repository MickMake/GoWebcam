package mmSelfUpdate

import (
	"GoWebcam/defaults"
	"errors"
	"fmt"
	"github.com/logrusorgru/aurora"
	"log"
	"os"
)


type State struct {
	error
}


func (s *State) IsOk() bool {
	if s == nil {
		return true
	}
	return false
}
func (s *State) IsNotOk() bool {
	return !s.IsOk()
}

func (s *State) IsError() bool {
	if s == nil {
		return false
	}
	return true
}
func (s *State) IsNotError() bool {
	return !s.IsError()
}

func (s *State) SetError(format string, args ...interface{}) {
	s.error = errors.New(fmt.Sprintf(format, args...))
}

func (s *State) SetWarning(format string, args ...interface{}) {
	s.error = errors.New(fmt.Sprintf(format, args...))
}

func (s *State) SetOk(format string, args ...interface{}) {
	s.error = nil
}


// ******************************************************************************** //
type ux struct {
}
type typeColours struct {
	Ref           aurora.Aurora
	Defined       bool
	Name          string
	EnableColours bool
	// TemplateRef   *template.Template
	// TemplateFuncs template.FuncMap
	Prefix        string
}

var colours typeColours

func Open(name string, enable bool) (*typeColours, error) {
	var err error

	for range onlyOnce {
		if name == "" {
			name = defaults.BinaryVersion
		}
		name += ": "

		colours.Ref = aurora.NewAurora(enable)
		colours.Name = name
		colours.EnableColours = enable
		colours.Defined = true
		colours.Prefix = fmt.Sprintf("%s", aurora.BrightCyan(colours.Name).Bold())

		//err = termui.Init();
		//if err != nil {
		//      fmt.Printf("failed to initialize termui: %v", err)
		//      break
		//}

		err = CreateTemplate()
	}

	return &colours, err
}


func Close() {
	if colours.Defined {
		//termui.Close()
	}
}

func (u *ux) PrintflnBlue(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (u *ux) PrintfWhite(format string, args ...interface{}) {
	inline := fmt.Sprintf(format, args...)
	if colours.EnableColours {
		inline = fmt.Sprintf("%s", aurora.BrightWhite(inline))
	}
	_, _ = fmt.Fprint(os.Stdout, inline)
}

func (u *ux) PrintfCyan(format string, args ...interface{}) {
	inline := fmt.Sprintf(format, args...)
	if colours.EnableColours {
		inline = fmt.Sprintf("%s", aurora.BrightCyan(inline))
	}
	_, _ = fmt.Fprint(os.Stdout, inline)
}

func (u *ux) PrintfYellow(format string, args ...interface{}) {
	inline := fmt.Sprintf(format, args...)
	if colours.EnableColours {
		inline = fmt.Sprintf("%s", aurora.BrightYellow(inline))
	}
	_, _ = fmt.Fprint(os.Stdout, inline)
}

func (u *ux) PrintfRed(format string, args ...interface{}) {
	inline := fmt.Sprintf(format, args...)
	if colours.EnableColours {
		inline = fmt.Sprintf("%s", aurora.BrightRed(inline))
	}
	_, _ = fmt.Fprint(os.Stdout, inline)
}

func (u *ux) PrintfGreen(format string, args ...interface{}) {
	inline := fmt.Sprintf(format, args...)
	if colours.EnableColours {
		inline = fmt.Sprintf("%s", aurora.BrightGreen(inline))
	}
	_, _ = fmt.Fprint(os.Stdout, inline)
}

func (u *ux) PrintfBlue(format string, args ...interface{}) {
	inline := fmt.Sprintf(format, args...)
	if colours.EnableColours {
		inline = fmt.Sprintf("%s", aurora.BrightBlue(inline))
	}
	_, _ = fmt.Fprint(os.Stdout, inline)
}

func (u *ux) PrintfMagenta(format string, args ...interface{}) {
	inline := fmt.Sprintf(format, args...)
	if colours.EnableColours {
		inline = fmt.Sprintf("%s", aurora.BrightMagenta(inline))
	}
	_, _ = fmt.Fprint(os.Stdout, inline)
}

func (u *ux) PrintflnWhite(format string, args ...interface{}) {
	inline := fmt.Sprintf(format, args...)
	if colours.EnableColours {
		inline = fmt.Sprintf("%s", aurora.BrightWhite(inline))
	}
	_, _ = fmt.Fprint(os.Stdout, inline + "\n")
}

func (u *ux) PrintflnCyan(format string, args ...interface{}) {
	inline := fmt.Sprintf(format, args...)
	if colours.EnableColours {
		inline = fmt.Sprintf("%s", aurora.BrightCyan(inline))
	}
	_, _ = fmt.Fprint(os.Stdout, inline + "\n")
}

func (u *ux) PrintflnYellow(format string, args ...interface{}) {
	inline := fmt.Sprintf(format, args...)
	if colours.EnableColours {
		inline = fmt.Sprintf("%s", aurora.BrightYellow(inline))
	}
	_, _ = fmt.Fprint(os.Stdout, inline + "\n")
}

func (u *ux) PrintflnRed(format string, args ...interface{}) {
	inline := fmt.Sprintf(format, args...)
	if colours.EnableColours {
		inline = fmt.Sprintf("%s", aurora.BrightRed(inline))
	}
	_, _ = fmt.Fprint(os.Stdout, inline + "\n")
}

func (u *ux) PrintflnGreen(format string, args ...interface{}) {
	inline := fmt.Sprintf(format, args...)
	if colours.EnableColours {
		inline = fmt.Sprintf("%s", aurora.BrightGreen(inline))
	}
	_, _ = fmt.Fprint(os.Stdout, inline + "\n")
}

func (u *ux) PrintflnBlue(format string, args ...interface{}) {
	inline := fmt.Sprintf(format, args...)
	if colours.EnableColours {
		inline = fmt.Sprintf("%s", aurora.BrightBlue(inline))
	}
	_, _ = fmt.Fprint(os.Stdout, inline + "\n")
}

func (u *ux) PrintflnMagenta(format string, args ...interface{}) {
	inline := fmt.Sprintf(format, args...)
	if colours.EnableColours {
		inline = fmt.Sprintf("%s", aurora.BrightMagenta(inline))
	}
	_, _ = fmt.Fprint(os.Stdout, inline + "\n")
}


func SprintfWhite(format string, args ...interface{}) string {
	inline := fmt.Sprintf(format, args...)
	if colours.EnableColours {
		inline = fmt.Sprintf("%s", aurora.BrightWhite(inline))
	}
	return inline
}

func SprintfCyan(format string, args ...interface{}) string {
	inline := fmt.Sprintf(format, args...)
	if colours.EnableColours {
		inline = fmt.Sprintf("%s", aurora.BrightCyan(inline))
	}
	return inline
}

func SprintfYellow(format string, args ...interface{}) string {
	inline := fmt.Sprintf(format, args...)
	if colours.EnableColours {
		inline = fmt.Sprintf("%s", aurora.BrightYellow(inline))
	}
	return inline
}

func SprintfRed(format string, args ...interface{}) string {
	inline := fmt.Sprintf(format, args...)
	if colours.EnableColours {
		inline = fmt.Sprintf("%s", aurora.BrightRed(inline))
	}
	return inline
}

func SprintfGreen(format string, args ...interface{}) string {
	inline := fmt.Sprintf(format, args...)
	if colours.EnableColours {
		inline = fmt.Sprintf("%s", aurora.BrightGreen(inline))
	}
	return inline
}

func SprintfBlue(format string, args ...interface{}) string {
	inline := fmt.Sprintf(format, args...)
	if colours.EnableColours {
		inline = fmt.Sprintf("%s", aurora.BrightBlue(inline))
	}
	return inline
}

func SprintfMagenta(format string, args ...interface{}) string {
	inline := fmt.Sprintf(format, args...)
	if colours.EnableColours {
		inline = fmt.Sprintf("%s", aurora.BrightMagenta(inline))
	}
	return inline
}


func Sprintf(format string, args ...interface{}) string {
	inline := fmt.Sprintf(format, args...)
	if colours.EnableColours {
		inline = fmt.Sprintf("%s%s", colours.Prefix, inline)
	}
	return inline
}
func (u *ux) Printf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, Sprintf(format, args...))
}


func SprintfNormal(format string, args ...interface{}) string {
	inline := fmt.Sprintf(format, args...)
	if colours.EnableColours {
		inline = Sprintf("%s", aurora.BrightBlue(inline))
	}
	return inline
}

func (u *ux) PrintfNormal(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, SprintfNormal(format, args...))
}

func (u *ux) PrintflnNormal(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, SprintfNormal(format + "\n", args...))
}


func SprintfInfo(format string, args ...interface{}) string {
	inline := fmt.Sprintf(format, args...)
	return Sprintf("%s", aurora.BrightBlue(inline))
}

func (u *ux) PrintfInfo(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, SprintfInfo(format, args...))
}

func (u *ux) PrintflnInfo(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, SprintfInfo(format + "\n", args...))
}


func SprintfOk(format string, args ...interface{}) string {
	inline := fmt.Sprintf(format, args...)
	return Sprintf("%s", aurora.BrightGreen(inline))
}

func (u *ux) PrintfOk(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, SprintfOk(format, args...))
}

func (u *ux) PrintflnOk(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, SprintfOk(format + "\n", args...))
}

func SprintfDebug(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

func (u *ux) PrintfDebug(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, fmt.Sprintf(format + "\n", args...))
}


func SprintfWarning(format string, args ...interface{}) string {
	inline := fmt.Sprintf(format, args...)
	return Sprintf("%s", aurora.BrightYellow(inline))
}

func (u *ux) PrintfWarning(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, SprintfWarning(format, args...))
}

func (u *ux) PrintflnWarning(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, SprintfWarning(format + "\n", args...))
}


func SprintfError(format string, args ...interface{}) string {
	inline := fmt.Sprintf(format, args...)
	return Sprintf("%s", aurora.BrightRed(inline))
}

func (u *ux) PrintfError(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, SprintfError(format, args...))
}

func (u *ux) PrintflnError(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, SprintfError(format + "\n", args...))
}

func SprintError(err error) string {
	var s string

	for range onlyOnce {
		if err == nil {
			break
		}

		s = Sprintf("%s%s\n", aurora.BrightRed("ERROR: ").Framed(), aurora.BrightRed(err).Framed().SlowBlink().BgBrightWhite())
	}

	return s
}

func (u *ux) PrintError(err error) {
	_, _ = fmt.Fprintf(os.Stderr, SprintError(err))
}
