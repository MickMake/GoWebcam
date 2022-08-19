package mmVersion

import (
	"GoWebcam/Only"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)


// type VersionValue semver.Version

func (v *Version) GetSemVer() *VersionValue {
	// v := semver.MustParse(r.CmdVersion)
	// return semver.Version(v.String())
	return GetSemVer(v.CmdVersion)
}

func (v *Version) PrintNameVersion() {
	fmt.Printf("%s ", v.CmdName)
	fmt.Printf("v%s", v.CmdVersion)
}

func (v *Version) TimeStampString() string {
	return v.TimeStamp.Format("2006-01-02T15:04:05-0700")
}

func (v *Version) TimeStampEpoch() int64 {
	return v.TimeStamp.Unix()
}

func (v *Version) GetTimeout() string {
	//d := r.Timeout.Round(time.Second)
	d := v.Timeout
	h := d / time.Hour
	d -= h * time.Hour

	m := d / time.Minute

	s := m / time.Second

	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}


func (v *Version) GetEnvMap() *Environment {
	return &v.EnvMap
}


func (v *Version) SetArgs(a ...string) {
	v.Args.Set(a...)
	//var err error
	//
	//for range onlyOnce {
	//	r.Args = a
	//}
	//
	//return err
}


func (v *Version) AddArgs(a ...string) {
	v.Args.Append(a...)
	//var err error
	//
	//for range onlyOnce {
	//	r.Args = append(r.Args, a...)
	//}
	//
	//return err
}


func (v *Version) GetArgs() []string {
	return v.Args.GetAll()
}


func (v *Version) GetArg(index int) string {
	return v.Args.Get(index)
}

func (v *Version) GetArgRange(lower int, upper int) []string {
	return v.Args.Range(lower, upper)
}

func (v *Version) SprintfArgRange(lower int, upper int) string {
	return v.Args.SprintfRange(lower, upper)
}

func (v *Version) SprintfArgsFrom(lower int) string {
	return v.Args.SprintfFrom(lower)
}

//goland:noinspection SpellCheckingInspection
func (v *Version) GetNargs(begin int, size int) []string {
	return v.Args.GetFromSize(begin, size)
}

//goland:noinspection SpellCheckingInspection
func (v *Version) SprintfNargs(lower int, upper int) string {
	return v.Args.SprintfFromSize(lower, upper)
}


func (v *Version) SetFullArgs(a ...string) {
	v.FullArgs.Set(a...)
	//var err error
	//
	//for range onlyOnce {
	//	r.FullArgs = a
	//}
	//
	//return err
}


func (v *Version) AddFullArgs(a ...string) {
	v.FullArgs.Append(a...)
	//var err error
	//
	//for range onlyOnce {
	//	r.FullArgs = append(r.FullArgs, a...)
	//}
	//
	//return err
}


func (v *Version) GetFullArgs() []string {
	return v.FullArgs.GetAll()
}


func (v *Version) SetCmd(a ...string) error {
	var err error

	for range Only.Once {
		v.Cmd, err = filepath.Abs(filepath.Join(a...))
		if err != nil {
			break
		}

		v.CmdDir = filepath.Dir(v.Cmd)
		v.CmdFile = filepath.Base(v.Cmd)
	}

	return err
}


func (v *Version) IsRunningAs(run string) bool {
	var ok bool
	// If OK - running executable file matches the string 'run'.
	//ok, err := regexp.MatchString("^" + run, r.CmdFile)

	if v.IsWindows() {
		//fmt.Printf("DEBUG: WINDOWS!\n")
		ok = strings.HasPrefix(run, strings.TrimSuffix(v.CmdFile, ".exe"))
		//run = strings.TrimSuffix(run, ".exe")
	} else {
		ok = strings.HasPrefix(run, v.CmdFile)
	}
	//fmt.Printf("DEBUG: Cmd.Runtime.IsRunningAs?? %s\n", ok)
	//fmt.Printf("DEBUG: run: %s\n", run)
	//fmt.Printf("DEBUG: r.CmdName: %s\n", r.CmdName)
	//fmt.Printf("DEBUG: r.CmdFile: %s\n", r.CmdFile)
	return ok
}
func (v *Version) IsRunningAsFile() bool {
	// If OK - running executable file matches the application binary name.
	//ok, err := regexp.MatchString("^" + r.CmdName, r.CmdFile)
	ok := strings.HasPrefix(v.CmdName, v.CmdFile)
	return ok
}
func (v *Version) IsRunningAsLink() bool {
	return !v.IsRunningAsFile()
}
