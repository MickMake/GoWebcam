package mmVersion

import (
	"GoWebcam/Only"
	"bufio"
	"errors"
	"github.com/kardianos/osext"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)


//goland:noinspection ALL
type Version struct {
	CmdName        string		`json:"cmd_name" mapstructure:"cmd_name"`
	CmdVersion     string		`json:"cmd_version" mapstructure:"cmd_version"`
	CmdSourceRepo  UrlValue		`json:"cmd_source_repo" mapstructure:"cmd_source_repo"`
	CmdBinaryRepo  UrlValue		`json:"cmd_binary_repo" mapstructure:"cmd_binary_repo"`

	Cmd            string		`json:"cmd" mapstructure:"cmd"`
	CmdDir         string	    `json:"cmd_dir" mapstructure:"cmd_dir"`
	CmdFile        string	    `json:"cmd_file" mapstructure:"cmd_file"`

	WorkingDir     Path			`json:"working_dir" mapstructure:"working_dir"`
	BaseDir        Path			`json:"base_dir" mapstructure:"base_dir"`
	BinDir         Path			`json:"bin_dir" mapstructure:"bin_dir"`
	ConfigDir      Path			`json:"config_dir" mapstructure:"config_dir"`
	CacheDir       Path			`json:"cache_dir" mapstructure:"cache_dir"`
	TempDir        Path			`json:"temp_dir" mapstructure:"temp_dir"`

	FullArgs       ExecArgs		`json:"full_args" mapstructure:"full_args"`
	Args           ExecArgs		`json:"args" mapstructure:"args"`
	ArgFiles       ExecArgs		`json:"arg_files" mapstructure:"arg_files"`

	Env            ExecEnv		`json:"env" mapstructure:"env"`
	EnvMap         Environment	`json:"env_map" mapstructure:"env_map"`

	TimeStamp      time.Time	`json:"timestamp" mapstructure:"timestamp"`
	Timeout        time.Duration `json:"timeout" mapstructure:"timeout"`

	GoRuntime      GoRuntime	`json:"go_runtime" mapstructure:"go_runtime"`

	User           User			`json:"user" mapstructure:"user"`

	Debug          bool			`json:"debug" mapstructure:"debug"`
	Verbose        bool			`json:"verbose" mapstructure:"verbose"`
	State          State	    `json:"state" mapstructure:"state"`


	useRepo       *UrlValue
	OldVersion    *VersionValue
	TargetBinary  string
	RuntimeBinary string
	AutoExec      bool

	logging    *FlagValue
	config     *selfupdate.Config
	ref        *selfupdate.Updater

	// Runtime    *Version
	cmd        *cobra.Command
	SelfCmd    *cobra.Command
}


type Path string
func (p *Path) DirExists() bool {
	var ok bool

	for range Only.Once {
		stat, err := os.Stat(string(*p))
		if os.IsNotExist(err) {
			break
		}

		if !stat.IsDir() {
			break
		}

		ok = true
	}

	return ok
}

func (p *Path) FileExists() bool {
	var ok bool

	for range Only.Once {
		stat, err := os.Stat(string(*p))
		if os.IsNotExist(err) {
			break
		}

		if stat.IsDir() {
			break
		}

		ok = true
	}

	return ok
}

func (p *Path) Chmod(mode os.FileMode) bool {
	var ok bool

	for range Only.Once {
		err := os.Chmod(string(*p), mode)
		if err != nil {
			break
		}

		ok = true
	}

	return ok
}

func (p *Path) Set(elem ...string) {

	for range Only.Once {
		dir := filepath.Join(elem...)
		if strings.HasPrefix(dir, "~/") {
			u, err := user.Current()
			if err != nil {
				break
			}
			dir = strings.TrimPrefix(dir, "~/")
			dir = filepath.Join(u.HomeDir, dir)
		}

		*p = Path(dir)
	}
}

func (p *Path) String() string {
	return (string)(*p)
}

//func (p *Path) Set(elem ...string) Path {
//	return (Path)(filepath.Join(elem...))
//}

func (p *Path) Join(elem ...string) Path {
	var pa []string
	//if p == nil {
	//	*p = "/"
	//}
	pa = append(pa, (string)(*p))
	pa = append(pa, elem...)
	return (Path)(filepath.Join(pa...))
}

func (p *Path) MkdirAll() error {
	var err error

	for range Only.Once {
		if p.DirExists() {
			break
		}

		dir := string(*p)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			break
		}

		if !p.DirExists() {
			err = errors.New("no dir")
			break
		}
	}

	return err
}

func (p *Path) Copy(fp string) error {
	var err error

	for range Only.Once {
		var stat os.FileInfo
		stat, err = os.Stat(fp)
		if os.IsNotExist(err) {
			break
		}
		if stat.IsDir() {
			err = errors.New("file is a dir")
			break
		}

		var input []byte
		input, err = ioutil.ReadFile(fp)
		if err != nil {
			break
		}

		dfp := filepath.Join(string(*p), filepath.Base(fp))
		err = ioutil.WriteFile(dfp, input, stat.Mode())
		if err != nil {
			break
		}
	}

	return err
}

func (p *Path) Move(fp string) error {
	var err error

	for range Only.Once {
		err = p.Copy(fp)
		if err != nil {
			break
		}

		err = os.Remove(fp)
		if err != nil {
			break
		}
	}

	return err
}

//goland:noinspection SpellCheckingInspection
var RcFiles = []Path {
	// BASH
	"/etc/profile",
	"/etc/bashrc",
	"~/.profile",
	"~/.bash_profile",
	"~/.bashrc",
	"~/.bash_login",
	"~/.bash_logout",

	// ZSH
	"/etc/zlogin",
	"/etc/zlogout",
	"/etc/zprofile",
	"/etc/zshenv",
	"/etc/zshrc",
	"~/.zlogin",
	"~/.zlogout",
	"~/.zprofile",
	"~/.zshenv",
	"~/.zshrc",

	// CSH
	"/etc/csh.cshrc",
	"/etc/csh.login",
	"/etc/csh.logout",
	"~/.cshrc",
	"~/.login",
	"~/.logout",
}

//goland:noinspection GoUnusedExportedFunction
func GrepFiles(search string, fps ...Path) ([]string, error) {
	var files []string
	var err error

	if fps == nil {
		fps = RcFiles
	}
	if len(fps) == 0 {
		fps = RcFiles
	}

	for _, p := range fps {
		var line int
		line, err = p.GrepFile(search)
		if line > 0 {
			files = append(files, p.String() + " line:" + strconv.Itoa(line))
		}
	}

	return files, err
}

func (p *Path) GrepFile(search string) (int, error) {
	var line int
	var err error

	for range Only.Once {
		p.Set(string(*p))

		var f *os.File
		f, err = os.Open(string(*p))
		if err != nil {
			// Silently ignore missing files.
			err = nil
			break
		}
		//goland:noinspection ALL
		defer f.Close()

		// Splits on newlines by default.
		scanner := bufio.NewScanner(f)
		line = 1
		// https://golang.org/pkg/bufio/#Scanner.Scan
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), search) {
				break
			}

			line++
		}

		err = scanner.Err()
		if err != nil {
			break
		}
	}

	return line, err
}


type ExecEnv []string
type Environment map[string]string
type GoRuntime struct {
	Os string
	Arch string
	Root string
	Version string
	Compiler string
	NumCpus int
}

type User struct {
	*user.User
}


// Instead of creating every time, let's cache the initial result in a global variable.
var globalRuntime *Version

func New(binary string, version string, debugFlag bool) *Version {
	var ret *Version

	for range Only.Once {
		if globalRuntime != nil {
			// Instead of creating every time, let's cache the initial result in a global variable.
			//globalRuntime.TimeStamp = time.Now()
			ret = globalRuntime
			break
		}

		ret = &Version{
			CmdName:    binary,
			CmdVersion: version,

			Cmd:        "",
			CmdDir:     "",
			CmdFile:    "",

			WorkingDir: ".",
			BaseDir:    ".",
			BinDir:     ".",
			ConfigDir:  ".",
			CacheDir:   ".",
			TempDir:    ".",

			FullArgs:   os.Args,
			Args:       os.Args[1:],
			ArgFiles:   []string{},

			Env:        os.Environ(),
			EnvMap:     make(Environment),

			TimeStamp:  time.Now(),

			GoRuntime: GoRuntime{
				Os:       runtime.GOOS,
				Arch:     runtime.GOARCH,
				Root:     runtime.GOROOT(),
				Version:  runtime.Version(),
				Compiler: runtime.Compiler,
				NumCpus:  runtime.NumCPU(),
			},

			Debug:      debugFlag,
			Verbose:    false,
			State:      State{},
		}

		for _, item := range os.Environ() {
			s := strings.SplitN(item, "=", 2)
			ret.EnvMap[s[0]] = s[1]
		}

		var err error
		var exe string
		var p string
		//ret.Cmd, err = os.Executable()
		//if err != nil {
		//	ret.State.SetError(err)
		//	break
		//}
		//ret.Cmd, err = filepath.Abs(ret.Cmd)
		//if err != nil {
		//	ret.State.SetError(err)
		//	break
		//}
		exe, err = osext.Executable()
		if err != nil {
			ret.State.SetError(err.Error())
			break
		}
		//if ret.GoRuntime.Os == "windows" {
		//	exe = strings.TrimSuffix(exe,".exe")
		//}
		ret.Cmd = exe
		ret.CmdDir = filepath.Dir(exe)
		ret.CmdFile = filepath.Base(exe)

		ret.User.User, err = user.Current()
		if err != nil {
			ret.State.SetError(err.Error())
			break
		}

		p, err = os.Getwd()
		if err != nil {
			ret.State.SetError(err.Error())
			break
		}
		ret.WorkingDir.Set(p)

		//if runtime.GOOS == "windows" {
		//	ret.BaseDir = ""
		//} else {
			ret.BaseDir.Set(ret.User.HomeDir, "." + ret.CmdName)
		//}

		//if runtime.GOOS == "windows" {
		//	ret.BinDir = ""
		//} else {
			ret.BinDir = ret.BaseDir.Join("bin")
		//}

		p, err = os.UserConfigDir()
		if err != nil {
			if runtime.GOOS == "windows" {
				ret.ConfigDir = ""
			} else {
				ret.ConfigDir = "."
			}
		} else {
			ret.ConfigDir = ret.BaseDir.Join("etc")
		}
		//ret.ConfigDir = filepath.Join(ret.ConfigDir, ret.CmdName)

		p, err = os.UserCacheDir()
		if err != nil {
			if runtime.GOOS == "windows" {
				ret.CacheDir = ""
			} else {
				ret.CacheDir = "."
			}
		} else {
			ret.CacheDir = ret.BaseDir.Join("cache")
		}
		//ret.CacheDir = filepath.Join(ret.CacheDir, ret.CmdName)

		p = os.TempDir()
		if ret.TempDir == "" {
			if runtime.GOOS == "windows" {
				ret.TempDir = "C:\\tmp"
			} else {
				ret.TempDir = "/tmp"
			}
		} else {
			ret.TempDir = ret.BaseDir.Join("tmp")
		}


		// ******************************************************************************** //
		ret.TargetBinary = ret.Cmd
		ret.RuntimeBinary = ResolveFile(ret.Cmd)
		ret.AutoExec =     false

		ret.logging =    toBoolValue(ret.Debug)
		ret.config =     &selfupdate.Config {
			APIToken:            "",
			EnterpriseBaseURL:   "",
			EnterpriseUploadURL: "",
			Validator:           nil, 	// &MyValidator{},
			Filters:             []string{},
		}

		ret.useRepo = &ret.CmdBinaryRepo

		// Workaround for selfupdate not being flexible enough to support variable asset names
		// Should enable a template similar to GoReleaser.
		// EG: {{ .ProjectName }}-{{ .Os }}_{{ .Arch }}
		//var asset string
		//asset, te.State = toolGhr.GetAsset(rt.CmdBinaryRepo, "latest")
		//te.config.Filters = append(te.config.Filters, asset)

		// Ignore the above and just make sure all filenames are lowercase.
		ret.config.Filters = append(ret.config.Filters, addFilters(ret.CmdFile, runtime.GOOS, runtime.GOARCH)...)
		ret.ref, _ = selfupdate.NewUpdater(*ret.config)
		if *ret.logging {
			selfupdate.EnableLog()
		}
		// ******************************************************************************** //

		// Instead of creating every time, let's cache the initial result in a global variable.
		globalRuntime = ret
	}

	return ret
}

func (v *Version) SetRepos(source string, binary string) State {
	// if state := ux.IfNilReturnError(r); state.IsError() {
	// 	return state
	// }
	v.CmdSourceRepo = toUrlValue(source)
	v.CmdBinaryRepo = toUrlValue(binary)

	return v.State
}

func (v *Version) EnsureNotNil() *Version {
	if v == nil {
		return New("binary", "version", false)
	}
	return v
}

func (v *Version) IsWindows() bool {
	var ok bool
	if v.GoRuntime.Os == "windows" {
		ok = true
	}
	return ok
}
func (v *Version) IsMac() bool {
	var ok bool
	if v.GoRuntime.Os == "darwin" {
		ok = true
	}
	return ok
}
func (v *Version) IsOsx() bool {
	var ok bool
	if v.GoRuntime.Os == "darwin" {
		ok = true
	}
	return ok
}
