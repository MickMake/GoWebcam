package cmdDaemon

import (
	"GoWebcam/Only"
	"GoWebcam/mmWebcam"
	"fmt"
	"github.com/sevlyar/go-daemon"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)


func ReadPid() int {
	var ret int

	for range Only.Once {
		if !mmWebcam.FileExists(pidFile) {
			break
		}

		os.Stat(pidFile)
		// Open PID file
		pid, err := os.ReadFile(pidFile)
		if err != nil {
			break
		}

		ps := strings.TrimSpace(string(pid))
		ret, err = strconv.Atoi(ps)
		if err != nil {
			ret = -1
			break
		}
	}

	return ret
}

func WritePid(pid int) error {
	var err error

	for range Only.Once {
		// Open a file for writing
		var file *os.File
		file, err = os.Create(pidFile)
		if err != nil {
			break
		}
		//goland:noinspection GoUnhandledErrorResult,GoDeferInLoop
		defer file.Close()

		_, err = file.Write([]byte(fmt.Sprintf("%d", pid)))
		if err != nil {
			break
		}
	}

	return err
}

//goland:noinspection GoUnusedExportedFunction
func DaemonizeClose(cntxt *daemon.Context) error {
	return cntxt.Release()
}


func worker() {
	fmt.Println("DEBUG")
LOOP:
	for {
		time.Sleep(time.Second) // this is work to be done by worker.
		select {
		case <-stop:
			break LOOP
		default:
		}
	}
	done <- struct{}{}
}

var (
	stop = make(chan struct{})
	done = make(chan struct{})
)

func termHandler(sig os.Signal) error {
	log.Println("Daemon terminating...")
	stop <- struct{}{}
	if sig == syscall.SIGQUIT {
		<-done
	}
	return daemon.ErrStop
}

func reloadHandler(_ os.Signal) error {
	log.Println("Daemon configuration reloaded")
	return nil
}
