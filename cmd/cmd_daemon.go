package cmd

import (
	"GoWebcam/Only"
	"GoWebcam/mmWebcam"
	"errors"
	"fmt"
	"github.com/sevlyar/go-daemon"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)


func AttachCmdDaemon(cmd *cobra.Command) *cobra.Command {
	// ******************************************************************************** //
	var cmdDaemon = &cobra.Command{
		Use:                   "daemon",
		Aliases:               []string{""},
		Short:                 fmt.Sprintf("Daemonize commands."),
		Long:                  fmt.Sprintf("Daemonize commands."),
		DisableFlagParsing:    false,
		DisableFlagsInUseLine: false,
		PreRunE:               Cmd.DaemonArgs,
		RunE:                  cmdDaemonFunc,
		Args:                  cobra.MinimumNArgs(1),
	}
	cmd.AddCommand(cmdDaemon)
	cmdDaemon.Example = PrintExamples(cmdDaemon, "exec web run", "kill")

	// ******************************************************************************** //
	var cmdDaemonExec = &cobra.Command{
		Use:                   "exec",
		Aliases:               []string{"run"},
		Short:                 fmt.Sprintf("Execute commands as a daemon."),
		Long:                  fmt.Sprintf("Execute commands as a daemon."),
		DisableFlagParsing:    false,
		DisableFlagsInUseLine: false,
		PreRunE:               Cmd.DaemonArgs,
		RunE:                  cmdDaemonExecFunc,
		Args:                  cobra.MinimumNArgs(1),
	}
	cmdDaemon.AddCommand(cmdDaemonExec)
	cmdDaemonExec.Example = PrintExamples(cmdDaemonExec, "")

	// ******************************************************************************** //
	var cmdDaemonKill = &cobra.Command{
		Use:                   "kill",
		Aliases:               []string{"stop"},
		Short:                 fmt.Sprintf("Terminate daemon."),
		Long:                  fmt.Sprintf("Terminate daemon."),
		DisableFlagParsing:    false,
		DisableFlagsInUseLine: false,
		PreRunE:               Cmd.DaemonArgs,
		RunE:                  cmdDaemonKillFunc,
		// Args:                  cobra.MinimumNArgs(1),
	}
	cmdDaemon.AddCommand(cmdDaemonKill)
	cmdDaemonKill.Example = PrintExamples(cmdDaemonKill, "")

	// ******************************************************************************** //
	var cmdDaemonReload = &cobra.Command{
		Use:                   "reload",
		Aliases:               []string{"hup"},
		Short:                 fmt.Sprintf("Reload daemon config."),
		Long:                  fmt.Sprintf("Reload daemon config."),
		DisableFlagParsing:    false,
		DisableFlagsInUseLine: false,
		PreRunE:               Cmd.DaemonArgs,
		RunE:                  cmdDaemonReloadFunc,
		// Args:                  cobra.MinimumNArgs(1),
	}
	cmdDaemon.AddCommand(cmdDaemonReload)
	cmdDaemonReload.Example = PrintExamples(cmdDaemonReload, "")

	// ******************************************************************************** //
	var cmdDaemonList = &cobra.Command{
		Use:                   "list",
		Aliases:               []string{"ls"},
		Short:                 fmt.Sprintf("List running daemon."),
		Long:                  fmt.Sprintf("List running daemon."),
		DisableFlagParsing:    false,
		DisableFlagsInUseLine: false,
		PreRunE:               Cmd.DaemonArgs,
		RunE:                  cmdDaemonListFunc,
		// Args:                  cobra.MinimumNArgs(1),
	}
	cmdDaemon.AddCommand(cmdDaemonList)
	cmdDaemonList.Example = PrintExamples(cmdDaemonList, "")

	return cmdDaemon
}

func (ca *CommandArgs) DaemonArgs(cmd *cobra.Command, args []string) error {
	for range Only.Once {
		ca.Error = ca.ProcessArgs(cmd, args)
		if ca.Error != nil {
			break
		}
	}

	return Cmd.Error
}

func cmdDaemonFunc(cmd *cobra.Command, _ []string) error {
	for range Only.Once {
		_ = cmd.Help()
	}

	return Cmd.Error
}

func cmdDaemonExecFunc(cmd *cobra.Command, args []string) error {
	return Daemonize(cmdDaemonExecFunc, cmd, args)
}

func cmdDaemonKillFunc(_ *cobra.Command, _ []string) error {
	for range Only.Once {
		pid := ReadPid()
		if pid == -1 {
			Cmd.Error = errors.New("PID file empty or no PID file")
			break
		}

		fmt.Printf("Killing daemon. PID: %d\n", pid)
		Cmd.Error = syscall.Kill(pid, syscall.SIGTERM)
		if Cmd.Error != nil {
			break
		}

		Cmd.Error = os.Remove(pidFile)
		if Cmd.Error != nil {
			break
		}
	}

	return Cmd.Error
}

func cmdDaemonReloadFunc(_ *cobra.Command, _ []string) error {
	for range Only.Once {
		pid := ReadPid()
		if pid == -1 {
			Cmd.Error = errors.New("PID file empty or no PID file")
			break
		}

		fmt.Printf("Reloading daemon. PID: %d\n", pid)
		Cmd.Error = syscall.Kill(pid, syscall.SIGHUP)
	}

	return Cmd.Error
}

func cmdDaemonListFunc(_ *cobra.Command, _ []string) error {
	for range Only.Once {
		// @TODO - Sort out this Daemon mess.
		cntxt := &daemon.Context {
			PidFileName: pidFile,
			PidFilePerm: 0644,
			LogFileName: DefaultBinaryName + "Daemon.log",
			LogFilePerm: 0640,
			WorkDir:     "./",
			Umask:       027,
			Args:        []string{ fmt.Sprintf("[%s]", DefaultBinaryName) },
		}

		var child *os.Process
		child, Cmd.Error = cntxt.Search()
		if Cmd.Error != nil {
			break
		}

		pid := ReadPid()
		switch {
			// If no discovered PID and no PID file.
			case (child == nil) && (pid == -1):
				fmt.Println("No daemon running.")

			// If no discovered PID and a PID file.
			case (child == nil) && (pid != -1):
				fmt.Println("Removing stale PID file.")
				Cmd.Error = os.Remove(pidFile)
				if Cmd.Error != nil {
					break
				}

			// If discovered PID and no PID file.
			case (child != nil) && (pid == -1):
				fmt.Printf("Daemon running. PID: %d\n", child.Pid)
				fmt.Println("Creating PID file.")
				Cmd.Error = WritePid(child.Pid)

			// If discovered PID and a PID file.
			case (child != nil) && (pid != -1):
				fmt.Printf("Daemon running. PID: %d\n", child.Pid)
				if child.Pid != pid {
					fmt.Printf("Creating PID file. (Mismatch: %d != %d)\n", child.Pid, pid)
					Cmd.Error = WritePid(child.Pid)
				}
		}
	}

	return Cmd.Error
}


const pidFile = DefaultBinaryName + ".pid"

type DaemonFunc func(cmd *cobra.Command, args []string) error

func Daemonize(fn DaemonFunc, cmd *cobra.Command, args []string) error {
	var err error

	for range Only.Once {
		nargs := []string{ fmt.Sprintf("[%s]", DefaultBinaryName) }
		nargs = append(nargs, args...)

		// @TODO - Sort out this Daemon mess.
		cntxt := &daemon.Context {
			PidFileName: pidFile,
			PidFilePerm: 0644,
			LogFileName: DefaultBinaryName + "Daemon.log",
			LogFilePerm: 0640,
			WorkDir:     "./",
			Umask:       027,
			Args:        nargs,
		}
		daemon.SetSigHandler(termHandler, syscall.SIGQUIT)
		daemon.SetSigHandler(termHandler, syscall.SIGTERM)
		daemon.SetSigHandler(reloadHandler, syscall.SIGHUP)

		// go worker()

		fmt.Printf("Starting daemon: %s\n", strings.Join(nargs, " "))
		var child *os.Process
		child, err = cntxt.Reborn()
		if err != nil {
			// log.Printf("Error: %s\n", err)
			log.Println("Daemon already running.")

			pid := ReadPid()
			if pid != -1 {
				log.Printf("PID: %d\n", pid)
			}
			break
		}

		if child != nil {
			fmt.Printf("Daemon started. PID: %d\n", child.Pid)
			err = WritePid(child.Pid)
			if err != nil {
				break
			}
			break
		}
		//goland:noinspection GoUnhandledErrorResult,GoDeferInLoop
		defer cntxt.Release()

		// @TODO - Never seems to get to here!
		log.Println("Daemon started.")
		// Cmd.Error = fn(cmd, args)
	}

	return err
}

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
