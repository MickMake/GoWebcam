package mmWebcam

import (
	"GoWebcam/Only"
	"encoding/json"
	"os"
	"time"
)


type Config struct {
	Images   []Webcam      `json:"images"`

	Timeout  time.Duration `json:"timeout,omitempty"`
	Dir      string        `json:"dir,omitempty"`
	Cron     string        `json:"cron,omitempty"`
	Rename   Rename        `json:"rename,omitempty"`

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

// func (c *Rename) Filename(fn string) string {
// 	for range Only.Once {
// 		if c.ByTime != "" {
// 			d, err := time.ParseDuration(c.ByTime)
// 			if err != nil {
// 				break
// 			}
// 			fmt.Printf("%v\n", d)
//
// 			// match := strings.FieldsFunc(fn, func(r rune) bool {
// 			// 	if r == '_' {
// 			// 		return true
// 			// 	}
// 			// 	if r == '-' {
// 			// 		return true
// 			// 	}
// 			// 	return false
// 			// })
//
// 			// Gotta be a better way.
// 			match := strings.Split(fn,"-")
//
//
// 			// re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
// 			re := regexp.MustCompile(`[\d{6,8}]+`)
// 			if !re.MatchString(fn) {
// 				// break
// 			}
//
// 			match = re.FindAllString(fn, -1)
// 			if len(match) != 2 {
// 				break
// 			}
//
// 			var ft time.Time
// 			ft, err = time.Parse("20060102 150405", fmt.Sprintf("%s %s", match[0], match[1]))
// 			if err != nil {
// 				break
// 			}
//
// 			fmt.Printf("Before: %v\n", ft)
// 			ft = ft.Round(d)
// 			fmt.Printf("After: %v\n", ft)
// 			fn = ft.Format("20060102 150405")
//
// 			break
// 		}
//
// 		if c.ByOCR != "" {
// 			break
// 		}
//
// 		d := time.Now().Format("20060102")
// 		t := time.Now().Format("150405")
// 		s := filepath.Ext(m.url.Path)
// 		ret = fmt.Sprintf("%s-%s_%s%s", m.Prefix, d, t, s)
// 	}
//
// 	return fn
// }


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
