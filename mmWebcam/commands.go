package mmWebcam

import (
	"fmt"
	"github.com/MickMake/GoUnify/Only"
	"github.com/MickMake/GoUnify/cmdCron"
	"github.com/MickMake/GoUnify/cmdExec"
	"github.com/MickMake/GoUnify/cmdHelp"
	"github.com/MickMake/GoUnify/cmdLog"
	"github.com/go-co-op/gocron"
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"strings"
)


type Webcams struct {
	Config    *Config
	Cron      *cmdCron.Cron
	Error     error

	cmd       *cobra.Command
	SelfCmd   *cobra.Command
}


func New() *Webcams {
	var ret *Webcams

	for range Only.Once {
		ret = &Webcams{
			Config: &Config{},
			Cron: cmdCron.New(),
			Error: nil,

			cmd: nil,
			SelfCmd: nil,
		}
	}

	return ret
}

const Group = "Webcam"
func (w *Webcams) AttachCommands(cmd *cobra.Command) *cobra.Command {
	for range Only.Once {
		if cmd == nil {
			break
		}
		w.cmd = cmd

		// ******************************************************************************** //
		w.SelfCmd = &cobra.Command{
			Use:                   "webcam",
			Aliases:               []string{"web"},
			Short:                 fmt.Sprintf("Webcam fetcher."),
			Long:                  fmt.Sprintf("Webcam fetcher."),
			DisableFlagParsing:    false,
			DisableFlagsInUseLine: false,
			PreRunE:               w.InitArgs,
			RunE:                  w.CmdWeb,
			Args:                  cobra.MinimumNArgs(1),
		}
		cmd.AddCommand(w.SelfCmd)
		w.SelfCmd.Example = cmdHelp.PrintExamples(w.SelfCmd, "get Basin https://charlottepass.com.au/charlottepass/webcam/lucylodge/current.jpg", "run", "cron")
		w.SelfCmd.Annotations = map[string]string{"group": Group}

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
		cmdWebGet.Annotations = map[string]string{"group": Group}

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
		cmdWebRun.Annotations = map[string]string{"group": Group}

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
		cmdWebCron.Example = cmdHelp.PrintExamples(cmdWebCron, "")
		cmdWebCron.Annotations = map[string]string{"group": Group}
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
		// prefix := Cmd.WebPrefix
		var prefix string
		if args[0] != "" {
			prefix = args[0]
		}

		w.Error = cmdLog.LogFileSet("")
		log.Printf("One-off fetch of webcam...\n")
		webcam := NewWebcam(Webcam {
			Url:      args[1],
			Prefix:   prefix,
			// Username: Cmd.WebUsername,
			// Password: Cmd.WebPassword,
			// Host:     Cmd.WebHost,
			// Port:     Cmd.WebPort,
		})
		w.Error = webcam.GetError()
		if w.Error != nil {
			break
		}

		w.Error = webcam.GetImage()
		if w.Error != nil {
			break
		}
	}

	return w.Error
}

func (w *Webcams) CmdWebRun(_ *cobra.Command, _ []string) error {
	for range Only.Once {
		w.Config, w.Error = ReadConfig("config.json")
		if w.Error != nil {
			break
		}

		w.Error = cmdLog.LogFileSet("")
		log.Printf("One-off fetch of webcams from config...\n")
		w.Error = w.Config.RunAll()
		if w.Error != nil {
			break
		}
	}

	return w.Error
}

func (w *Webcams) CmdWebCron(_ *cobra.Command, _ []string) error {
	for range Only.Once {
		w.Config, w.Error = ReadConfig("config.json")
		if w.Error != nil {
			break
		}

		w.Error = cmdLog.LogFileSet("")
		cmdLog.Printf("Cron based webcam fetch from config...\n")

		// crontab := make(map[string]*gocron.Job)
		for index, webcam := range w.Config.Images {
			var job *gocron.Job
			// job, w.Error = CmdCron.AddJob(webcam.Cron, "Webcam:" + webcam.Prefix, Webcams.Images[index].GetImage)
			job, w.Error = w.Cron.AddJob(webcam.Cron, webcam.Prefix, w.Config.Images[index].GetImage)
			if w.Error != nil {
				cmdLog.Printf("crontab error: %s - %s\n", w.Error, job.Error())
				break
			}
			// crontab[webcam.Prefix] = job
		}

		if w.Config.Report.Cron != "" {
			_, w.Error = w.Cron.AddJob(w.Config.Report.Cron, "Report", w.Cron.PrintJobs)
			if w.Error != nil {
				break
			}
		}

		for index := range w.Config.Scripts {
			name := fmt.Sprintf("script-%d", index)
			// _, w.Error = CmdCron.AddJob(Webcams.Scripts[index].Cron, "Script:" + name, RunScript, name)
			_, w.Error = w.Cron.AddJob(w.Config.Scripts[index].Cron, name, w.RunScript, name)
			if w.Error != nil {
				break
			}
		}

		w.Error = w.Cron.StartBlocking()
		if w.Error != nil {
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
		// w.Error = config.Write("config.json")
		// if w.Error != nil {
		// 	break
		// }
	}

	return w.Error
}

func (w *Webcams) RunScript(name string) error {
	var err error

	for range Only.Once {
		id := -1
		for _, key := range w.Cron.Jobs() {
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

		err = cmdExec.Exec(job.Cmd, job.Args...)
	}

	return err
}
