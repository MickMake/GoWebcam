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
	"strconv"
	"strings"
)


func AttachCmdWeb(cmd *cobra.Command) *cobra.Command {
	// ******************************************************************************** //
	var cmdWeb = &cobra.Command{
		Use:                   "web",
		Aliases:               []string{""},
		Short:                 fmt.Sprintf("Connect to a HASSIO broker."),
		Long:                  fmt.Sprintf("Connect to a HASSIO broker."),
		DisableFlagParsing:    false,
		DisableFlagsInUseLine: false,
		PreRunE:               Cmd.WebArgs,
		RunE:                  cmdWebFunc,
		Args:                  cobra.MinimumNArgs(1),
	}
	cmd.AddCommand(cmdWeb)
	cmdWeb.Example = cmdHelp.PrintExamples(cmdWeb, "run", "cron")

	// ******************************************************************************** //
	var cmdWebGet = &cobra.Command{
		Use:                   "get",
		Aliases:               []string{""},
		Short:                 fmt.Sprintf("One-off webcam fetch."),
		Long:                  fmt.Sprintf("One-off webcam fetch."),
		DisableFlagParsing:    false,
		DisableFlagsInUseLine: false,
		PreRunE:               Cmd.WebArgs,
		RunE:                  cmdWebGetFunc,
		Args:                  cobra.ExactArgs(2),
	}
	cmdWeb.AddCommand(cmdWebGet)
	cmdWebGet.Example = cmdHelp.PrintExamples(cmdWebGet, "Basin https://charlottepass.com.au/charlottepass/webcam/lucylodge/current.jpg")

	// ******************************************************************************** //
	var cmdWebRun = &cobra.Command{
		Use:                   "run",
		Aliases:               []string{""},
		Short:                 fmt.Sprintf("One-off webcam fetch from config."),
		Long:                  fmt.Sprintf("One-off webcam fetch from config."),
		DisableFlagParsing:    false,
		DisableFlagsInUseLine: false,
		PreRunE:               Cmd.WebArgs,
		RunE:                  cmdWebRunFunc,
		Args:                  cobra.RangeArgs(0, 1),
	}
	cmdWeb.AddCommand(cmdWebRun)
	cmdWebRun.Example = cmdHelp.PrintExamples(cmdWebRun, "")

	// ******************************************************************************** //
	var cmdWebCron = &cobra.Command{
		Use:                   "cron",
		Aliases:               []string{""},
		Short:                 fmt.Sprintf("Cron based webcam fetch from config."),
		Long:                  fmt.Sprintf("Cron based webcam fetch from config."),
		DisableFlagParsing:    false,
		DisableFlagsInUseLine: false,
		PreRunE:               Cmd.WebArgs,
		RunE:                  cmdWebCronFunc,
		Args:                  cobra.RangeArgs(0, 1),
	}
	cmdWeb.AddCommand(cmdWebCron)
	cmdWebCron.Example = cmdHelp.PrintExamples(cmdWebCron, "", "all")

	return cmdWeb
}

func (ca *CommandArgs) WebArgs(cmd *cobra.Command, args []string) error {
	for range Only.Once {
		ca.Error = ca.ProcessArgs(cmd, args)
		if ca.Error != nil {
			break
		}
	}

	return Cmd.Error
}

func cmdWebFunc(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}

func cmdWebGetFunc(_ *cobra.Command, args []string) error {
	for range Only.Once {
		prefix := Cmd.WebPrefix
		if args[0] != "" {
			prefix = args[0]
		}

		cmdLog.Printf("One-off fetch of webcam...\n")
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

func cmdWebRunFunc(_ *cobra.Command, _ []string) error {
	for range Only.Once {
		Webcams, Cmd.Error = mmWebcam.ReadConfig("config.json")
		if Cmd.Error != nil {
			break
		}

		cmdLog.Printf("One-off fetch of webcams from config...\n")
		Cmd.Error = Webcams.RunAll()
		if Cmd.Error != nil {
			break
		}
	}

	return Cmd.Error
}

func cmdWebCronFunc(cmd *cobra.Command, args []string) error {
	for range Only.Once {
		Webcams, Cmd.Error = mmWebcam.ReadConfig("config.json")
		if Cmd.Error != nil {
			break
		}

		cmdLog.Printf("Cron based webcam fetch from config...\n")

		crontab := make(map[string]*gocron.Job)
		for index, webcam := range Webcams.Images {
			var job *gocron.Job
			// job, Cmd.Error = CmdCron.Scheduler.CronWithSeconds(webcam.Cron).StartImmediately().Tag(webcam.Prefix).Do(Webcams.Images[index].GetImage)
			// job, Cmd.Error = CmdCron.AddJob(webcam.Cron, "Webcam:" + webcam.Prefix, Webcams.Images[index].GetImage)
			job, Cmd.Error = CmdCron.AddJob(webcam.Cron, webcam.Prefix, Webcams.Images[index].GetImage)
			if Cmd.Error != nil {
				cmdLog.Printf("crontab error: %s\n", Cmd.Error)
				break
			}
			crontab[webcam.Prefix] = job
		}

		if Webcams.Report.Cron != "" {
			// _, Cmd.Error = CmdCron.Scheduler.CronWithSeconds(Webcams.Report.Cron).StartImmediately().Tag("report").Do(CmdCron.PrintJobs)
			_, Cmd.Error = CmdCron.AddJob(Webcams.Report.Cron, "Report", CmdCron.PrintJobs)
			if Cmd.Error != nil {
				break
			}
		}

		for index := range Webcams.Scripts {
			name := fmt.Sprintf("script-%d", index)
			// _, Cmd.Error = CmdCron.Scheduler.CronWithSeconds(Webcams.Scripts[index].Cron).Tag(name).Do(RunScript, name)
			// _, Cmd.Error = CmdCron.AddJob(Webcams.Scripts[index].Cron, "Script:" + name, RunScript, name)
			_, Cmd.Error = CmdCron.AddJob(Webcams.Scripts[index].Cron, name, RunScript, name)
			if Cmd.Error != nil {
				break
			}
		}

		Cmd.Error = CmdCron.StartBlocking()
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

func RunScript(name string) error {
	var err error

	for range Only.Once {
		id := -1
		for _, key := range CmdCron.Jobs() {
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

		if id >= len(Webcams.Scripts) {
			break
		}

		job := Webcams.Scripts[id]

		err = cmdCron.Exec(job.Cmd, job.Args...)
	}

	return err
}
