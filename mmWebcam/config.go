package mmWebcam

import (
	"GoWebcam/Only"
	"GoWebcam/Unify/cmdLog"
	"encoding/json"
	"os"
	"time"
)


type Config struct {
	Images   []Webcam      `json:"images"`

	Timeout  time.Duration `json:"timeout,omitempty"`
	Dir      string        `json:"dir,omitempty"`
	Cron     string        `json:"cron,omitempty"`
	Logfile  string        `json:"logfile,omitempty"`

	Rename   Rename        `json:"rename,omitempty"`

	Report struct {
		Cron  string `json:"cron,omitempty"`
		Level string `json:"level,omitempty"`
	} `json:"report,omitempty"`

	Scripts []struct {
		Cron  string `json:"cron,omitempty"`
		Cmd string `json:"cmd,omitempty"`
		Args []string `json:"args,omitempty"`
	} `json:"scripts,omitempty"`

	Error    error         `json:"-"`
}

func (c *Config) Read(fp string) (*Config, error) {
	var err error

	for range Only.Once {
		if c == nil {
			c = &Config{}
		}

		// Open a file for reading
		var data []byte
		data, err = os.ReadFile(fp)
		if err != nil {
			break
		}

		err = json.Unmarshal(data, c)
		if err != nil {
			break
		}

		err = cmdLog.LogFileSet(c.Logfile)
		if err != nil {
			break
		}

		if c.Dir == "" {
			c.Dir = "."
		}
		if c.Timeout == 0 {
			c.Timeout = time.Second * 30
		}
		if c.Cron == "" {
			c.Cron = "00 */5 * * * *"
		}
		if c.Rename.RoundTime != "" {
			c.Rename.duration, err = time.ParseDuration(c.Rename.RoundTime)
			if err != nil {
				break
			}
		}

		for i := range c.Images {
			c.Images[i] = New(c.Images[i])

			if c.Images[i].Dir == "" {
				c.Images[i].Dir = c.Dir
			}
			if c.Images[i].Timeout == 0 {
				c.Images[i].Timeout = c.Timeout
			}
			if c.Images[i].Cron == "" {
				c.Images[i].Cron = c.Cron
			}
			if c.Images[i].Rename.RoundTime == "" {
				c.Images[i].Rename.RoundTime = c.Rename.RoundTime
				c.Images[i].Rename.duration = c.Rename.duration
			}


			if c.Images[i].Rename.RoundTime != "" {
				c.Images[i].Rename.duration, err = time.ParseDuration(c.Images[i].Rename.RoundTime)
				if err != nil {
					break
				}
			}

			if c.Images[i].Rename.OCR == "" {
				// @TODO - Do stuff.
				if err != nil {
					break
				}
			}
		}
	}

	return c, err
}

func (c *Config) Write(fp string) error {
	for range Only.Once {
		if fp == "" {
			fp = "config.json"
		}

		var data []byte
		data, c.Error = json.MarshalIndent(c, "", "\t")
		if c.Error != nil {
			break
		}

		// Open a file for writing
		var file *os.File
		file, c.Error = os.Create(fp)
		if c.Error != nil {
			break
		}
		//goland:noinspection ALL
		defer file.Close()

		_, c.Error = file.Write(data)
		if c.Error != nil {
			break
		}
	}

	return c.Error
}

func (c *Config) Run(index int) error {
	for range Only.Once {
		if index >= len(c.Images) {
			break
		}

		if c.Images[index].Dir == "" {
			c.Images[index].Dir = c.Dir
		}
		if c.Images[index].Timeout == 0 {
			c.Images[index].Timeout = c.Timeout
		}
		if c.Images[index].Cron == "" {
			c.Images[index].Cron = c.Cron
		}

		c.Error = c.Images[index].GetImage()
		if c.Error != nil {
			break
		}
	}

	return c.Error
}

func (c *Config) RunAll() error {
	for index := range c.Images {
		c.Error = c.Run(index)
		if c.Error != nil {
			break
		}
	}

	return c.Error
}


type Rename struct {
	RoundTime string `json:"by_time,omitempty"`
	duration time.Duration
	OCR      string `json:"by_ocr,omitempty"`
}


func ReadConfig(fp string) (*Config, error) {
	c := &Config{}
	return c.Read(fp)
}

func NewWebcam(req Webcam) *Webcam {
	var ret Webcam

	for range Only.Once {
	}

	return &ret
}
