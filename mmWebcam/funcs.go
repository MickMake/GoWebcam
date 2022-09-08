package mmWebcam

import (
	"github.com/MickMake/GoUnify/Only"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)


func Mkdir(dir ...string) error {
	var err error

	for range Only.Once {
		d := filepath.Join(dir...)
		var fi os.FileInfo
		fi, err = os.Stat(d)
		if fi != nil {
			if !fi.IsDir() {
				err = errors.New("directory exists as a file")
				break
			}
		}

		if errors.Is(err, os.ErrNotExist) {
			err = os.MkdirAll(d, os.ModePerm)
			if err != nil {
				break
			}
		}

		break
	}

	return err
}

func MD5HashString(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func Md5HashFile(filename string) string {
	var ret string

	for range Only.Once {
		if filename == "" {
			break
		}

		f, err := os.Open(filename)
		if err != nil {
			break
		}
		defer f.Close()

		h := md5.New()
		if _, err := io.Copy(h, f); err != nil {
			break
		}
		ret = fmt.Sprintf("%x", h.Sum(nil))
		// ret = hex.EncodeToString()
	}

	return ret
}

// func LogSprintf(format string, args ...interface{}) string {
// 	// format = timeStamp() + format
// 	return fmt.Sprintf(format, args...)
// }
//
// func LogSprintfDate(format string, args ...interface{}) string {
// 	ret := fmt.Sprintf("%s ", TimeNow())
// 	ret += fmt.Sprintf(format, args...)
// 	return ret
// }
//
// func TimeNow() string {
// 	return time.Now().Format("2006-01-02 15:04:05")
// }

func DirExists(fn string) bool {
	var yes bool

	for range Only.Once {
		f, err := os.Stat(fn)

		if errors.Is(err, os.ErrNotExist) {
			yes = false
			break
		}

		if !f.IsDir() {
			yes = false
			break
		}

		if err == nil {
			yes = true
			break
		}

		if errors.Is(err, os.ErrNotExist) {
			yes = false
			break
		}

		if err != nil {
			// Schrodinger: file may or may not exist. See err for details.
			// Do *NOT* use !os.IsNotExist(err) to test for file existence
			yes = false
			break
		}
	}

	return yes
}

func FileExists(fn string) bool {
	var yes bool
	for range Only.Once {
		_, err := os.Stat(fn)

		if err == nil {
			yes = true
			break
		}

		if errors.Is(err, os.ErrNotExist) {
			yes = false
			break
		}

		if err != nil {
			// Schrodinger: file may or may not exist. See err for details.
			// Do *NOT* use !os.IsNotExist(err) to test for file existence
			yes = false
			break
		}
	}
	return yes
}

func GetSortedKeys(s []interface{}) []string {
	var sorted []string
	for range Only.Once {
	}
	return sorted
}
