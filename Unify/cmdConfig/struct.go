package cmdConfig

import (
	"GoWebcam/Only"
	"GoWebcam/Unify/cmdVersion"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)


type Config struct {
	Dir       string
	File      string
	EnvPrefix string
	Error     error

	// Flags     map[string]interface{}

	viper     *viper.Viper
	cmd       *cobra.Command
	SelfCmd   *cobra.Command
}


func New(name string) *Config {
	var ret *Config

	for range Only.Once {
		ret = &Config {
			Dir: ".",
			File: defaultConfigFile,
			EnvPrefix: cmdVersion.GetEnvPrefix(),
			Error: nil,

			// Flags: make(map[string]interface{}),

			viper: viper.New(),
			cmd: nil,
			SelfCmd: nil,
		}

		ret.Dir, ret.Error = os.UserHomeDir()
		if ret.Error != nil {
			break
		}
		ret.SetDir(filepath.Join(ret.Dir, "." + name))

		// Cmd.CacheDir, ret.Error = os.UserHomeDir()
		// if ret.Error != nil {
		// 	break
		// }
		// Cmd.CacheDir = filepath.Join(Cmd.CacheDir, "." + defaults.BinaryName, "cache")
		// _, ret.Error = os.Stat(Cmd.CacheDir)
		// if os.IsExist(ret.Error) {
		// 	break
		// }
		// ret.Error = os.MkdirAll(Cmd.CacheDir, 0700)
		// if ret.Error != nil {
		// 	break
		// }

		ret.SetFile(filepath.Join(ret.Dir, defaultConfigFile))
	}

	return ret
}

func (c *Config) SetDir(path string) {
	for range Only.Once {
		if path == "" {
			break
		}
		c.Dir = path
		c.viper.AddConfigPath(c.Dir)
		_, c.Error = os.Stat(c.Dir)
		if os.IsExist(c.Error) {
			break
		}
		c.Error = os.MkdirAll(c.Dir, 0700)
		if c.Error != nil {
			break
		}
	}
}

func (c *Config) SetFile(fn string) {
	for range Only.Once {
		if fn == "" {
			break
		}
		c.File = fn
		c.viper.SetConfigFile(c.File)
		// c.viper.SetConfigName("config")
	}
}

// Init reads in config file and ENV variables if set.
func (c *Config) Init(_ *cobra.Command) error {
	var err error

	for range Only.Once {
		// If a config file is found, read it in.
		err = c.Open()
		if err != nil {
			break
		}

		c.viper.SetEnvPrefix(c.EnvPrefix)
		c.viper.AutomaticEnv() // read in environment variables that match
		c.cmd.Flags().VisitAll(func(f *pflag.Flag) {
			// Environment variables can't have dashes in them, so bind them to their equivalent
			// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
			if strings.Contains(f.Name, "-") {
				envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
				err = c.viper.BindEnv(f.Name, fmt.Sprintf("%s_%s", c.EnvPrefix, envVarSuffix))
			}

			// Apply the viper config value to the flag when the flag is not set and viper has a value
			if !f.Changed && c.viper.IsSet(f.Name) {
				// val := c.viper.Get(f.Name)	// Doesn't handle time.Duration well.
				// val := c.cmd.Flag(f.Name).Value.String()
				val := f.Value.String()
				err = c.cmd.Flags().Set(f.Name, val)
			}
		})

		if err != nil {
			break
		}
	}

	return err
}

func (c *Config) Open() error {
	for range Only.Once {
		c.Error = c.viper.ReadInConfig()
		if _, ok := c.Error.(viper.UnsupportedConfigError); ok {
			break
		}

		if _, ok := c.Error.(viper.ConfigParseError); ok {
			break
		}

		if _, ok := c.Error.(viper.ConfigMarshalError); ok {
			break
		}

		if os.IsNotExist(c.Error) {
			c.cmd.Flags().VisitAll(func(f *pflag.Flag) {
				switch f.Value.Type() {
					case "duration":
						c.viper.SetDefault(f.Name, f.Value.String())
					default:
						c.viper.SetDefault(f.Name, f.Value)
					}
			})

			c.Error = c.viper.WriteConfig()
			if c.Error != nil {
				break
			}

			c.Error = c.viper.ReadInConfig()
		}
		if c.Error != nil {
			break
		}

		c.Error = c.viper.MergeInConfig()
		if c.Error != nil {
			break
		}
	}

	return c.Error
}

func (c *Config) Write() error {
	for range Only.Once {
		c.Error = c.viper.MergeInConfig()
		if c.Error != nil {
			break
		}

		c.cmd.Flags().VisitAll(func(f *pflag.Flag) {
			switch f.Value.Type() {
				case "duration":
					c.viper.Set(f.Name, f.Value.String())
				default:
					c.viper.Set(f.Name, f.Value)
			}
		})

		c.Error = c.viper.WriteConfig()
		if c.Error != nil {
			break
		}
	}

	return c.Error
}

func (c *Config) Read() error {
	for range Only.Once {
		c.Error = c.viper.ReadInConfig()
		if c.Error != nil {
			break
		}

		_, _ = fmt.Fprintln(os.Stderr, "Config file settings:")
		c.cmd.Flags().VisitAll(func(f *pflag.Flag) {
			_, _ = fmt.Fprintf(os.Stderr, "%s:			%s\n", strings.ToTitle(f.Name), f.Value.String())
		})
	}

	return c.Error
}

func (c *Config) SetDefault(key string, value interface{}) {
	for range Only.Once {
		c.viper.SetDefault(key, value)
		// c.Flags[key] = value
	}
}