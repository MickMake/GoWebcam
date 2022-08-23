package cmdCron

import (
	"GoWebcam/Only"
	"GoWebcam/Unify/cmdHelp"
	"GoWebcam/Unify/cmdLog"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"time"
)


func (c *Cron) AttachCommands(cmd *cobra.Command) *cobra.Command {
	for range Only.Once {
		if cmd == nil {
			break
		}
		c.cmd = cmd

		// ******************************************************************************** //
		c.SelfCmd = &cobra.Command{
			Use:                   "cron",
			Aliases:               []string{""},
			Short:                 fmt.Sprintf("Cron - Run a command via schedule."),
			Long:                  fmt.Sprintf("Cron - Run a command via schedule."),
			DisableFlagParsing:    false,
			DisableFlagsInUseLine: false,
			PreRunE:               c.InitArgs,
			Run:                   c.CmdCron,
			Args:                  cobra.MinimumNArgs(1),
		}
		cmd.AddCommand(c.SelfCmd)
		c.SelfCmd.Example = cmdHelp.PrintExamples(c.SelfCmd, "run ./5 . . . . web get Basin https://charlottepass.com.au/charlottepass/webcam/lucylodge/current.jpg", "run 00 12 . . . web get Basin https://charlottepass.com.au/charlottepass/webcam/lucylodge/current.jpg")

		// ******************************************************************************** //
		var cmdCronRun = &cobra.Command{
			Use:                   "run <minute> <hour> <month day> <month> <week day>  <command ...>",
			Aliases:               []string{""},
			Short:                 fmt.Sprintf("Cron - Run scheduled a command."),
			Long:                  fmt.Sprintf("Cron - Run scheduled a command."),
			DisableFlagParsing:    false,
			DisableFlagsInUseLine: false,
			PreRunE:               c.InitArgs,
			Run:                   c.CmdCronRun,
			Args:                  cobra.MinimumNArgs(6),
		}
		c.SelfCmd.AddCommand(cmdCronRun)
		cmdCronRun.Example = cmdHelp.PrintExamples(cmdCronRun, "./5 . . . . web get Basin https://charlottepass.com.au/charlottepass/webcam/lucylodge/current.jpg", "00 12 . . . web get Basin https://charlottepass.com.au/charlottepass/webcam/lucylodge/current.jpg")

		// ******************************************************************************** //
		var cmdConfigRead = &cobra.Command{}
		c.SelfCmd.AddCommand(cmdConfigRead)
		cmdConfigRead.Example = cmdHelp.PrintExamples(cmdConfigRead, "")

		// ******************************************************************************** //
		var cmdCronAdd = &cobra.Command{
			Use:                   "add",
			Aliases:               []string{""},
			Short:                 fmt.Sprintf("Cron - Add scheduled a command."),
			Long:                  fmt.Sprintf("Cron - Add scheduled a command."),
			DisableFlagParsing:    false,
			DisableFlagsInUseLine: false,
			PreRunE:               c.InitArgs,
			Run:                   c.CmdCronAdd,
			Args:                  cobra.MinimumNArgs(1),
		}
		c.SelfCmd.AddCommand(cmdCronAdd)
		cmdCronAdd.Example = cmdHelp.PrintExamples(cmdCronAdd, "add")

		// ******************************************************************************** //
		var cmdCronRemove = &cobra.Command{
			Use:                   "del",
			Aliases:               []string{"remove"},
			Short:                 fmt.Sprintf("Cron - Remove a scheduled command."),
			Long:                  fmt.Sprintf("Cron - Remove a scheduled command."),
			DisableFlagParsing:    false,
			DisableFlagsInUseLine: false,
			PreRunE:               c.InitArgs,
			Run:                   c.CmdCronRemove,
			Args:                  cobra.MinimumNArgs(1),
		}
		c.SelfCmd.AddCommand(cmdCronRemove)
		cmdCronRemove.Example = cmdHelp.PrintExamples(cmdCronRemove, "del")

		// ******************************************************************************** //
		var cmdCronList = &cobra.Command{
			Use:                   "list",
			Aliases:               []string{""},
			Short:                 fmt.Sprintf("Cron - List scheduled commands."),
			Long:                  fmt.Sprintf("Cron - List scheduled commands."),
			DisableFlagParsing:    false,
			DisableFlagsInUseLine: false,
			PreRunE:               c.InitArgs,
			Run:                   c.CmdCronList,
			Args:                  cobra.MinimumNArgs(1),
		}
		c.SelfCmd.AddCommand(cmdCronList)
		cmdCronList.Example = cmdHelp.PrintExamples(cmdCronList, "list")
	}

	return c.SelfCmd
}

func (c *Cron) InitArgs(cmd *cobra.Command, args []string) error {
	var err error
	for range Only.Once {
		//
	}
	return err
}

func (c *Cron) CmdCron(cmd *cobra.Command, args []string) {
	for range Only.Once {
		if len(args) == 0 {
			c.Error = cmd.Help()
			break
		}
	}
}

func (c *Cron) CmdCronRun(_ *cobra.Command, args []string) {
	for range Only.Once {
		// */1 * * * * /dir/command args args
		cronString := strings.Join(args[0:5], " ")
		cronString = strings.ReplaceAll(cronString, ".", "*")
		ResetArgs(args[5:]...)

		c.Scheduler = c.Scheduler.Cron(cronString)
		c.Scheduler = c.Scheduler.SingletonMode()

		c.Job, c.Error = c.Scheduler.Do(c.ReExecute)
		if c.Error != nil {
			break
		}

		cmdLog.Printf("Created job schedule using '%s'\n", cronString)
		cmdLog.Printf("Job command '%s'\n", strings.Join(os.Args, " "))

		c.Scheduler.StartBlocking()
		if c.Error != nil {
			break
		}
	}
}

func (c *Cron) CmdCronAdd(_ *cobra.Command, _ []string) {
	for range Only.Once {
		fmt.Println("Not yet implemented.")

		// var msg string
		// switch {
		// 	case args[0] == "":
		// 		fallthrough
		// 	case args[0] == "default":
		// 		//u, _ := user.Current()
		// 		//msg = fmt.Sprintf("Regular sync by %s", u.ApiUsername)
		// 	default:
		// 		msg = args[0]
		// }
		//
		// args = args[1:]
		//
		// //Cmd.Error = Cmd.CronAdd(msg, args...)
		// if Cmd.Error != nil {
		// 	break
		// }
	}
}

func (c *Cron) CmdCronRemove(_ *cobra.Command, _ []string) {
	for range Only.Once {
		fmt.Println("Not yet implemented.")

		// var msg string
		// switch {
		// 	case args[0] == "":
		// 		fallthrough
		// 	case args[0] == "default":
		// 		//u, _ := user.Current()
		// 		//msg = fmt.Sprintf("Regular sync by %s", u.ApiUsername)
		// 	default:
		// 		msg = args[0]
		// }
		//
		// args = args[1:]
		//
		// //Cmd.Error = Cmd.CronAdd(msg, args...)
		// if Cmd.Error != nil {
		// 	break
		// }
	}
}

func (c *Cron) CmdCronList(_ *cobra.Command, _ []string) {
	for range Only.Once {
		fmt.Println("Not yet implemented.")

		// var msg string
		// 	switch {
		// 		case args[0] == "":
		// 			fallthrough
		// 		case args[0] == "default":
		// 			//u, _ := user.Current()
		// 			//msg = fmt.Sprintf("Regular sync by %s", u.ApiUsername)
		// 		default:
		// 			msg = args[0]
		// }
		//
		// args = args[1:]
		//
		// Cmd.Error = Cmd.CronList(msg, args...)
		// if Cmd.Error != nil {
		// 	break
		// }
	}
}

func (c *Cron) ReExecute() error {
	for range Only.Once {
		cmdLog.Printf("Running scheduled command '%s'\n", strings.Join(os.Args, " "))
		// LogPrint("Last run '%s'\n", Cron.Job.LastRun().Format(time.UnixDate))
		cmdLog.Printf("Next run '%s'\n", c.Job.ScheduledTime().Format(time.UnixDate))
		cmdLog.Printf("Run count '%d'\n", c.Job.RunCount())

		rootCmd := FindRoot(c.SelfCmd)
		c.Error = rootCmd.Execute()
		if c.Error != nil {
			cmdLog.Printf("ERROR: %s\n", c.Error)
			break
		}
	}

	return c.Error
}

func FindRoot(cmd *cobra.Command) *cobra.Command {
	var ret *cobra.Command
	for range Only.Once {
		if !cmd.HasParent() {
			ret = cmd
			break
		}

		ret = FindRoot(cmd.Parent())
	}

	return ret
}

func ResetArgs(args ...string) {
	for range Only.Once {
		newArgs := []string{os.Args[0]}
		newArgs = append(newArgs, args...)
		os.Args = newArgs
	}
}
