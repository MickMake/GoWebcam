package mmWebcam

import (
	"GoWebcam/Only"
	"GoWebcam/Unify/cmdLog"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	// mqtt "github.com/eclipse/paho.mqtt.golang"
	"net/url"
	"time"
)


type Webcam struct {
	Prefix   string        `json:"prefix"`
	Url      string        `json:"url"`
	Username string        `json:"username,omitempty"`
	Password string        `json:"password,omitempty"`
	Host     string        `json:"host,omitempty"`
	Port     string        `json:"port,omitempty"`

	// Also Global options.
	Timeout  time.Duration `json:"timeout,omitempty"`
	Dir      string        `json:"dir,omitempty"`
	Cron     string        `json:"cron,omitempty"`
	Delay    string        `json:"delay,omitempty"`
	Rename   Rename        `json:"rename,omitempty"`

	lastRefresh time.Time
	url         *url.URL
	firstRun    bool
	err         error
	logfile     cmdLog.Log
}


func NewWebcam(req Webcam) Webcam {
	var ret Webcam

	for range Only.Once {
		ret = req

		ret.err = ret.setUrl(req)
		if ret.err != nil {
			break
		}

		ret.firstRun = true
		ret.lastRefresh = time.Time{}

		if ret.Dir == "" {
			ret.Dir = "images"
		}
		ret.err = Mkdir(ret.Dir)
		if ret.err != nil {
			break
		}
	}

	return ret
}


func (m *Webcam) IsFirstRun() bool {
	return m.firstRun
}

func (m *Webcam) IsNotFirstRun() bool {
	return !m.firstRun
}

func (m *Webcam) UnsetFirstRun() {
	m.firstRun = false
}

func (m *Webcam) GetError() error {
	return m.err
}

func (m *Webcam) IsError() bool {
	if m.err != nil {
		return true
	}
	return false
}

func (m *Webcam) IsNewDay() bool {
	var yes bool

	for range Only.Once {
		last := m.lastRefresh.Format("20060102")
		now := time.Now().Format("20060102")

		if last != now {
			yes = true
			break
		}
	}

	return yes
}

func (m *Webcam) setUrl(req Webcam) error {

	for range Only.Once {
		if req.Host != "" {
			m.Username = req.Username
			m.Password = req.Password
			m.Host = req.Host

			if req.Port == "" {
				req.Port = "443"
			}
			if req.Port == "80" {
				req.Port = "http"
			} else if req.Port == "443" {
				req.Port = "https"
			}
			m.Port = req.Port

			req.Url = fmt.Sprintf("%s://%s:%s@%s:%s",
				req.Port,
				m.Username,
				m.Password,
				m.Host,
				m.Port,
			)
		}

		m.url, m.err = url.Parse(req.Url)
	}

	return m.err
}

func (m *Webcam) SetAuth(username string, password string) error {

	for range Only.Once {
		if username == "" {
			m.err = errors.New("username empty")
			break
		}
		m.Username = username

		if password == "" {
			m.err = errors.New("password empty")
			break
		}
		m.Password = password
	}

	return m.err
}

func (m *Webcam) GetImage() error {
	for range Only.Once {
		m.err = m.LogOpen()
		if m.err != nil {
			break
		}

		m.LogPrintfDate("Fetch webcam image: %s", m.Prefix)
		var response *http.Response
		response, m.err = http.Get(m.url.String())
		if m.err != nil {
			m.LogPrintf("	- ERROR: %s\n", m.err)
			break
		}
		//goland:noinspection ALL
		defer response.Body.Close()
		m.LogPrintf("	- GET")

		if response.StatusCode != http.StatusOK {
			m.LogPrintf("	- ERROR STATUS: %v\n", response.StatusCode)
			break
		}

		var body []byte
		body, m.err = io.ReadAll(response.Body)
		if m.err != nil {
			m.LogPrintf("	- ERROR: %s\n", m.err)
			break
		}

		lf := m.GetLastFile()
		lfmd5 := Md5HashFile(lf)
		tfmd5 := MD5HashString(string(body))
		if lfmd5 == tfmd5 {
			// m.err = os.Remove(fp)
			m.LogPrintf("	- DUPE (%s)\n", filepath.Base(lf))
			break
		}


		filename := m.GetFilename()
		dir := m.GetBaseDir()
		fp := filepath.Join(dir, filename)

		m.err = Mkdir(dir)
		if m.err != nil {
			break
		}

		// Open a file for writing
		var file *os.File
		file, m.err = os.Create(fp)
		if m.err != nil {
			m.LogPrintf("	- ERROR: %s\n", m.err)
			break
		}
		//goland:noinspection ALL
		defer file.Close()

		_, m.err = file.Write(body)
		if m.err != nil {
			m.LogPrintf("	- ERROR: %s\n", m.err)
			break
		}


		// // Use io.Copy to just dump the response body to the file. This supports huge files
		// _, m.err = io.Copy(file, body)
		// if m.err != nil {
		// 	m.LogPrintf("	- ERROR: %s\n", m.err)
		// 	break
		// }
		//
		//
		// tfmd5 := Md5HashFile(fp)
		// if lfmd5 == tfmd5 {
		// 	m.err = os.Remove(fp)
		// 	m.LogPrintf("	- DUPE (%s)\n", lf)
		// 	break
		// }

		m.LogPrintf("	- SAVED\n")
		m.lastRefresh = time.Now()
		if m.firstRun == true {
			m.firstRun = false
		}

	}

	m.LogClose()

	return m.err
}

func (m *Webcam) GetLastFile() string {
	var ret string

	for range Only.Once {
		dir := m.GetBaseDir()

		var files []string
		re := regexp.MustCompile(`\w+-\d{8}_\d{6}(-\d)*\.[a-zA-Z]+$`)
		m.err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if re.MatchString(path) {
				files = append(files, path)
			}
			return nil
		})

		if files == nil {
			break
		}

		// sort.Strings(files)
		// ret = files[len(files)-1]

		var when time.Time
		for _, file := range files {
			f, err := os.Stat(file)
			if err != nil {
				continue
			}
			fm := f.ModTime()
			if fm.After(when) {
				when = f.ModTime()
				ret = file
			}
		}
	}

	return ret
}

func (m *Webcam) GetBaseDir() string {
	var ret string

	for range Only.Once {
		d := time.Now().Format("20060102")
		ret = filepath.Join(m.Dir, m.Prefix, d)
	}

	return ret
}

func (m *Webcam) GetFilename() string {
	var ret string

	for range Only.Once {
		now := time.Now()

		for range Only.Once {
			if m.Rename.OCR != "" {
				break
			}

			if m.Rename.duration != 0 {
				now = now.Truncate(m.Rename.duration)
				break
			}
		}

		nowDate := now.Format("20060102")
		nowTime := now.Format("150405")
		ext := filepath.Ext(m.url.Path)
		ret = fmt.Sprintf("%s-%s_%s%s", m.Prefix, nowDate, nowTime, ext)

		for index := 1; index < 32; index++ {
			if FileExists(filepath.Join(m.GetBaseDir(), ret)) {
				ret = fmt.Sprintf("%s-%s_%s-%d%s", m.Prefix, nowDate, nowTime, index, ext)
				continue
			}
			break
		}
	}

	return ret
}

func (m *Webcam) GetDir() string {
	var ret string

	for range Only.Once {
		d := time.Now().Format("20060102")
		ret = filepath.Join(m.Dir, m.Prefix, d)

		if !DirExists(ret) {
			m.err = Mkdir(ret)
			if m.err != nil {
				break
			}
		}
	}

	return ret
}

func (m *Webcam) LogOpen() error {
	return m.logfile.Open(m.GetDir(), m.Prefix + ".log")
}

func (m *Webcam) LogPrintfDate(format string, args ...interface{}) {
	m.logfile.LogPrintfDate(format, args...)
}

func (m *Webcam) LogPrintf(format string, args ...interface{}) {
	m.logfile.LogPrintf(format, args...)
}

func (m *Webcam) LogClose() error {
	return m.logfile.Close()
}
