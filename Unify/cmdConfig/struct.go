package cmdConfig

import (
	"GoWebcam/Only"
	"GoWebcam/defaults"
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
	Error     error

	// Flags     map[string]interface{}

	viper     *viper.Viper
	cmd       *cobra.Command
	SelfCmd   *cobra.Command
}


func New() *Config {
	var ret *Config

	for range Only.Once {
		ret = &Config {
			Dir: ".",
			File: defaultConfigFile,
			Error: nil,
			// Flags: make(map[string]interface{}),
			viper: viper.New(),
		}

		ret.Dir, ret.Error = os.UserHomeDir()
		if ret.Error != nil {
			break
		}
		ret.SetDir(filepath.Join(ret.Dir, "." + defaults.BinaryName))

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
func (c *Config) Init(cmd *cobra.Command) error {
	var err error

	for range Only.Once {
		// If a config file is found, read it in.
		err = c.Open()
		if err != nil {
			break
		}

		c.viper.SetEnvPrefix(defaults.EnvPrefix)
		c.viper.AutomaticEnv() // read in environment variables that match
		err = bindFlags(c.cmd, c.viper)
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
			// for key, value := range c.Flags {
			// 	c.viper.SetDefault(key, value)
			// }
			c.cmd.Flags().VisitAll(func(f *pflag.Flag) {
				c.viper.SetDefault(f.Name, f.Value)
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

		// c.cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// 	c.Flags[f.Name] = f.Value
		// })

		// err = viper.Unmarshal(Cmd)
	}

	return c.Error
}

func (c *Config) Write() error {
	for range Only.Once {
		c.Error = c.viper.MergeInConfig()
		if c.Error != nil {
			break
		}

		// for key, value := range c.Flags {
		// 	c.viper.Set(key, value)
		// }
		c.cmd.Flags().VisitAll(func(f *pflag.Flag) {
			c.viper.Set(f.Name, f.Value)
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
			_, _ = fmt.Fprintf(os.Stderr, "%s:			%v\n", strings.ToTitle(f.Name), f.Value)
		})
		// for key := range c.Flags {
		// 	_, _ = fmt.Fprintf(os.Stderr, "%s:			%v\n", strings.ToTitle(key), c.viper.Get(key))
		// }
	}

	return c.Error
}

func (c *Config) SetDefault(key string, value interface{}) {
	for range Only.Once {
		c.viper.SetDefault(key, value)
		// c.Flags[key] = value
	}
}
