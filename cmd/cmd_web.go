package cmd

import (
	"GoWebcam/Only"
	"GoWebcam/Unify/cmdCron"
	"GoWebcam/Unify/cmdHelp"
	"GoWebcam/Unify/cmdLog"
	"GoWebcam/mmWebcam"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"strings"
)


type Webcams struct {
	Config    *mmWebcam.Config
	Error     error

	cmd       *cobra.Command
	SelfCmd   *cobra.Command
}


func NewWeb() *Webcams {
	var ret *Webcams

	for range Only.Once {
		ret = &Webcams{
			Config: &mmWebcam.Config{},
			Error: nil,

			cmd: nil,
			SelfCmd: nil,
		}
	}

	return ret
}


func (w *Webcams) AttachCmdWeb(cmd *cobra.Command) *cobra.Command {
	for range Only.Once {
		if cmd == nil {
			break
		}
		w.cmd = cmd

		// ******************************************************************************** //
		w.SelfCmd = &cobra.Command{
			Use:                   "web",
			Aliases:               []string{""},
			Short:                 fmt.Sprintf("Webcam fetcher."),
			Long:                  fmt.Sprintf("Webcam fetcher."),
			DisableFlagParsing:    false,
			DisableFlagsInUseLine: false,
			PreRunE:               w.InitArgs,
			RunE:                  w.CmdWeb,
			Args:                  cobra.MinimumNArgs(1),
		}
		cmd.AddCommand(w.SelfCmd)
		w.SelfCmd.Example = cmdHelp.PrintExamples(w.SelfCmd, "run", "cron")

		// ******************************************************************************** //
		var cmdWebGet = &cobra.Command{
			Use:                   "get",
			Aliases:               []string{""},
			Short:                 fmt.Sprintf("One-off webcam fetch."),
			Long:                  fmt.Sprintf("One-off webcam fetch."),
			DisableFlagParsing:    false,
			DisableFlagsInUseLine: false,
			PreRunE:               w.InitArgs,
			RunE:                  w.CmdWebGet,
			Args:                  cobra.ExactArgs(2),
		}
		w.SelfCmd.AddCommand(cmdWebGet)
		cmdWebGet.Example = cmdHelp.PrintExamples(cmdWebGet, "Basin https://charlottepass.com.au/charlottepass/webcam/lucylodge/current.jpg")

		// ******************************************************************************** //
		var cmdWebRun = &cobra.Command{
			Use:                   "run",
			Aliases:               []string{""},
			Short:                 fmt.Sprintf("One-off webcam fetch from config."),
			Long:                  fmt.Sprintf("One-off webcam fetch from config."),
			DisableFlagParsing:    false,
			DisableFlagsInUseLine: false,
			PreRunE:               w.InitArgs,
			RunE:                  w.CmdWebRun,
			Args:                  cobra.RangeArgs(0, 1),
		}
		w.SelfCmd.AddCommand(cmdWebRun)
		cmdWebRun.Example = cmdHelp.PrintExamples(cmdWebRun, "")

		// ******************************************************************************** //
		var cmdWebCron = &cobra.Command{
			Use:                   "cron",
			Aliases:               []string{""},
			Short:                 fmt.Sprintf("Cron based webcam fetch from config."),
			Long:                  fmt.Sprintf("Cron based webcam fetch from config."),
			DisableFlagParsing:    false,
			DisableFlagsInUseLine: false,
			PreRunE:               w.InitArgs,
			RunE:                  w.CmdWebCron,
			Args:                  cobra.RangeArgs(0, 1),
		}
		w.SelfCmd.AddCommand(cmdWebCron)
		cmdWebCron.Example = cmdHelp.PrintExamples(cmdWebCron, "", "all")
	}

	return w.SelfCmd
}

func (w *Webcams) InitArgs(_ *cobra.Command, _ []string) error {
	var err error
	for range Only.Once {
		//
	}
	return err
}

func (w *Webcams) CmdWeb(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}

func (w *Webcams) CmdWebGet(_ *cobra.Command, args []string) error {
	for range Only.Once {
		prefix := Cmd.WebPrefix
		if args[0] != "" {
			prefix = args[0]
		}

		Cmd.Error = cmdLog.LogFileSet("")
		log.Printf("One-off fetch of webcam...\n")
		webcam := mmWebcam.New(mmWebcam.Webcam{
			Url:   args[1],

			Prefix:   prefix,
			Username: Cmd.WebUsername,
			Password: Cmd.WebPassword,
			Host:     Cmd.WebHost,
			Port:     Cmd.WebPort,
		})
		Cmd.Error = webcam.GetError()
		if Cmd.Error != nil {
			break
		}

		Cmd.Error = webcam.GetImage()
		if Cmd.Error != nil {
			break
		}
	}

	return Cmd.Error
}

func (w *Webcams) CmdWebRun(_ *cobra.Command, _ []string) error {
	for range Only.Once {
		w.Config, Cmd.Error = mmWebcam.ReadConfig("config.json")
		if Cmd.Error != nil {
			break
		}

		Cmd.Error = cmdLog.LogFileSet("")
		log.Printf("One-off fetch of webcams from config...\n")
		Cmd.Error = w.Config.RunAll()
		if Cmd.Error != nil {
			break
		}
	}

	return Cmd.Error
}

func (w *Webcams) CmdWebCron(_ *cobra.Command, _ []string) error {
	for range Only.Once {
		w.Config, Cmd.Error = mmWebcam.ReadConfig("config.json")
		if Cmd.Error != nil {
			break
		}

		Cmd.Error = cmdLog.LogFileSet("")
		cmdLog.Printf("Cron based webcam fetch from config...\n")

		crontab := make(map[string]*gocron.Job)
		for index, webcam := range w.Config.Images {
			var job *gocron.Job
			// job, Cmd.Error = CmdCron.AddJob(webcam.Cron, "Webcam:" + webcam.Prefix, Webcams.Images[index].GetImage)
			job, Cmd.Error = Cmd.CmdCron.AddJob(webcam.Cron, webcam.Prefix, w.Config.Images[index].GetImage)
			if Cmd.Error != nil {
				cmdLog.Printf("crontab error: %s\n", Cmd.Error)
				break
			}
			crontab[webcam.Prefix] = job
		}

		if w.Config.Report.Cron != "" {
			_, Cmd.Error = Cmd.CmdCron.AddJob(w.Config.Report.Cron, "Report", Cmd.CmdCron.PrintJobs)
			if Cmd.Error != nil {
				break
			}
		}

		for index := range w.Config.Scripts {
			name := fmt.Sprintf("script-%d", index)
			// _, Cmd.Error = CmdCron.AddJob(Webcams.Scripts[index].Cron, "Script:" + name, RunScript, name)
			_, Cmd.Error = Cmd.CmdCron.AddJob(w.Config.Scripts[index].Cron, name, w.RunScript, name)
			if Cmd.Error != nil {
				break
			}
		}

		Cmd.Error = Cmd.CmdCron.StartBlocking()
		if Cmd.Error != nil {
			break
		}

		// updateCounter := 0
		// timer := time.NewTicker(60 * time.Second)
		// for range timer.C {
		// 	if updateCounter < 5 {
		// 		updateCounter++
		// 		// cmdLog.Printf("Sleeping: %d\n", updateCounter)
		// 		continue
		// 	}
		// 	updateCounter = 0
		//
		// 	PrintJobs()
		// }
		//
		// time.Sleep(time.Hour * 5)
		// Cron.Scheduler.Stop()
		// fmt.Println(Cron.Scheduler.IsRunning())
		//
		// Cmd.Error = config.Write("config.json")
		// if Cmd.Error != nil {
		// 	break
		// }
	}

	return Cmd.Error
}

func (w *Webcams) RunScript(name string) error {
	var err error

	for range Only.Once {
		id := -1
		for _, key := range Cmd.CmdCron.Jobs() {
			if strings.Join(key.Tags(), " ") != name {
				continue
			}
			ids := strings.Join(key.Tags(), "")
			ids = strings.TrimPrefix(ids, "script-")
			var i int
			i, err = strconv.Atoi(ids)
			if err != nil {
				break
			}
			id = i
			break
		}
		if id == -1 {
			break
		}

		if id >= len(w.Config.Scripts) {
			break
		}

		job := w.Config.Scripts[id]

		err = cmdCron.Exec(job.Cmd, job.Args...)
	}

	return err
}
