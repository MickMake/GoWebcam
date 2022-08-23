package cmdDaemon

import (
	"GoWebcam/Only"
	"GoWebcam/Unify/cmdHelp"
	"GoWebcam/defaults"
	"errors"
	"fmt"
	"github.com/sevlyar/go-daemon"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
	"syscall"
)


func (d *Daemon) AttachCommands(cmd *cobra.Command) *cobra.Command {
	for range Only.Once {
		if cmd == nil {
			break
		}
		d.cmd = cmd

		// ******************************************************************************** //
		d.SelfCmd = &cobra.Command{
			Use:                   CmdDaemon,
			Aliases:               []string{""},
			Short:                 fmt.Sprintf("Daemon - Daemonize commands."),
			Long:                  fmt.Sprintf("Daemon - Daemonize commands."),
			DisableFlagParsing:    false,
			DisableFlagsInUseLine: false,
			PreRunE:               d.InitArgs,
			Run: func(cmd *cobra.Command, args []string) {
				d.Error = d.CmdDaemon(cmd, args)
			},
			Args: cobra.MinimumNArgs(1),
		}
		cmd.AddCommand(d.SelfCmd)
		d.SelfCmd.Example = cmdHelp.PrintExamples(d.SelfCmd, "exec web run", "kill")

		// ******************************************************************************** //
		var cmdDaemonExec = &cobra.Command{
			Use:                   CmdDaemonExec,
			Aliases:               AliasesDaemonExec,
			Short:                 fmt.Sprintf("Daemon - Execute commands as a daemon."),
			Long:                  fmt.Sprintf("Daemon - Execute commands as a daemon."),
			DisableFlagParsing:    false,
			DisableFlagsInUseLine: false,
			PreRunE:               d.InitArgs,
			Run: func(cmd *cobra.Command, args []string) {
				d.Error = d.CmdDaemonExec(DummyFunc, nil, args)
			},
			Args: cobra.MinimumNArgs(1),
		}
		d.SelfCmd.AddCommand(cmdDaemonExec)
		cmdDaemonExec.Example = cmdHelp.PrintExamples(cmdDaemonExec, "")

		// ******************************************************************************** //
		var cmdDaemonKill = &cobra.Command{
			Use:                   CmdDaemonStop,
			Aliases:               AliasesDaemonStop,
			Short:                 fmt.Sprintf("Daemon - Terminate daemon."),
			Long:                  fmt.Sprintf("Daemon - Terminate daemon."),
			DisableFlagParsing:    false,
			DisableFlagsInUseLine: false,
			PreRunE:               d.InitArgs,
			Run: func(cmd *cobra.Command, args []string) {
				d.Error = d.CmdDaemonKill()
			},
			// Args:                  cobra.MinimumNArgs(1),
		}
		d.SelfCmd.AddCommand(cmdDaemonKill)
		cmdDaemonKill.Example = cmdHelp.PrintExamples(cmdDaemonKill, "")

		// ******************************************************************************** //
		var cmdDaemonReload = &cobra.Command{
			Use:                   CmdDaemonReload,
			Aliases:               AliasesDaemonReload,
			Short:                 fmt.Sprintf("Daemon - Reload daemon config."),
			Long:                  fmt.Sprintf("Daemon - Reload daemon config."),
			DisableFlagParsing:    false,
			DisableFlagsInUseLine: false,
			PreRunE:               d.InitArgs,
			Run: func(cmd *cobra.Command, args []string) {
				d.Error = d.CmdDaemonReload()
			},
			// Args:                  cobra.MinimumNArgs(1),
		}
		d.SelfCmd.AddCommand(cmdDaemonReload)
		cmdDaemonReload.Example = cmdHelp.PrintExamples(cmdDaemonReload, "")

		// ******************************************************************************** //
		var cmdDaemonList = &cobra.Command{
			Use:                   CmdDaemonList,
			Aliases:               AliasesDaemonList,
			Short:                 fmt.Sprintf("Daemon - List running daemon."),
			Long:                  fmt.Sprintf("Daemon - List running daemon."),
			DisableFlagParsing:    false,
			DisableFlagsInUseLine: false,
			PreRunE:               d.InitArgs,
			Run: func(cmd *cobra.Command, args []string) {
				d.Error = d.CmdDaemonList()
			},
			// Args:                  cobra.MinimumNArgs(1),
		}
		d.SelfCmd.AddCommand(cmdDaemonList)
		cmdDaemonList.Example = cmdHelp.PrintExamples(cmdDaemonList, "")
	}

	return d.SelfCmd
}


func (d *Daemon) InitArgs(cmd *cobra.Command, args []string) error {
	var err error
	for range Only.Once {
		//
	}
	return err
}

func (d *Daemon) CmdDaemon(cmd *cobra.Command, _ []string) error {
	for range Only.Once {
		d.Error = cmd.Help()
	}

	return d.Error
}

func (d *Daemon) CmdDaemonExec(fn DaemonFunc, _ *cobra.Command, args []string) error {
	for range Only.Once {
		nargs := []string{ fmt.Sprintf("[%s]", defaults.BinaryName) }
		nargs = append(nargs, args...)

		// @TODO - Sort out this Daemon mess.
		cntxt := &daemon.Context {
			PidFileName: pidFile,
			PidFilePerm: 0644,
			LogFileName: defaults.BinaryName + "Daemon.log",
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
		child, d.Error = cntxt.Reborn()
		if d.Error != nil {
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
			d.Error = WritePid(child.Pid)
			if d.Error != nil {
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

	return d.Error
}

func (d *Daemon) CmdDaemonKill() error {
	for range Only.Once {
		pid := ReadPid()
		if pid == -1 {
			d.Error = errors.New("PID file empty or no PID file")
			break
		}

		fmt.Printf("Killing daemon. PID: %d\n", pid)
		// Cmd.Error = syscall.Kill(pid, syscall.SIGTERM)
		var p *os.Process
		p, d.Error = os.FindProcess(pid)
		if d.Error != nil {
			break
		}

		d.Error = p.Signal(syscall.SIGTERM)
		if d.Error != nil {
			break
		}

		d.Error = os.Remove(pidFile)
		if d.Error != nil {
			break
		}
	}

	return d.Error
}

func (d *Daemon) CmdDaemonReload() error {
	for range Only.Once {
		pid := ReadPid()
		if pid == -1 {
			d.Error = errors.New("PID file empty or no PID file")
			break
		}

		fmt.Printf("Reloading daemon. PID: %d\n", pid)
		// Cmd.Error = syscall.Kill(pid, syscall.SIGHUP)
		var p *os.Process
		p, d.Error = os.FindProcess(pid)
		if d.Error != nil {
			break
		}

		d.Error = p.Signal(syscall.SIGTERM)
		if d.Error != nil {
			break
		}
	}

	return d.Error
}

func (d *Daemon) CmdDaemonList() error {
	for range Only.Once {
		// @TODO - Sort out this Daemon mess.
		cntxt := &daemon.Context {
			PidFileName: pidFile,
			PidFilePerm: 0644,
			LogFileName: defaults.BinaryName + "Daemon.log",
			LogFilePerm: 0640,
			WorkDir:     "./",
			Umask:       027,
			Args:        []string{ fmt.Sprintf("[%s]", defaults.BinaryName) },
		}

		var child *os.Process
		child, d.Error = cntxt.Search()
		if d.Error != nil {
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
			d.Error = os.Remove(pidFile)
			if d.Error != nil {
				break
			}

		// If discovered PID and no PID file.
		case (child != nil) && (pid == -1):
			fmt.Printf("Daemon running. PID: %d\n", child.Pid)
			fmt.Println("Creating PID file.")
			d.Error = WritePid(child.Pid)
			if d.Error != nil {
				break
			}

		// If discovered PID and a PID file.
		case (child != nil) && (pid != -1):
			fmt.Printf("Daemon running. PID: %d\n", child.Pid)
			if child.Pid == pid {
				break
			}
			fmt.Printf("Creating PID file. (Mismatch: %d != %d)\n", child.Pid, pid)
			d.Error = WritePid(child.Pid)
			if d.Error != nil {
				break
			}
		}
	}

	return d.Error
}


type DaemonFunc func(cmd *cobra.Command, args []string) error

func DummyFunc(cmd *cobra.Command, args []string) error {
	var err error
	for range Only.Once {
		//
	}
	return err
}
