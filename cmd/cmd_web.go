package cmd

import (
	"GoWebcam/Only"
	"GoWebcam/mmWebcam"
	"bytes"
	"errors"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
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
	cmdWeb.Example = PrintExamples(cmdWeb, "run", "cron")

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
	cmdWebGet.Example = PrintExamples(cmdWebGet, "Basin https://charlottepass.com.au/charlottepass/webcam/lucylodge/current.jpg")

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
	cmdWebRun.Example = PrintExamples(cmdWebRun, "")

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
	cmdWebCron.Example = PrintExamples(cmdWebCron, "", "all")

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

		LogPrintDate("One-off fetch of webcam...\n")
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
		LogPrintDate("One-off fetch of webcams from config...\n")
		Webcams, Cmd.Error = mmWebcam.ReadConfig("config.json")
		if Cmd.Error != nil {
			break
		}

		Cmd.Error = Webcams.RunAll()
		if Cmd.Error != nil {
			break
		}


		// updateCounter := 0
		// timer := time.NewTicker(60 * time.Second)
		// for t := range timer.C {
		// 	if updateCounter < 5 {
		// 		updateCounter++
		// 		LogPrintDate("Sleeping: %d\n", updateCounter)
		// 		continue
		// 	}
		//
		// 	updateCounter = 0
		// 	LogPrintDate("Update: %s\n", t.String())
		// 	Cmd.Error = WebCron()
		// 	if Cmd.Error != nil {
		// 		break
		// 	}
		// }
	}

	return Cmd.Error
}

func cmdWebCronFunc(_ *cobra.Command, args []string) error {

	for range Only.Once {
		LogPrintDate("Cron based webcam fetch from config...\n")
		Webcams, Cmd.Error = mmWebcam.ReadConfig("config.json")
		if Cmd.Error != nil {
			break
		}

		Cron.Scheduler = gocron.NewScheduler(time.Local)

		crontab := make(map[string]*gocron.Job)
		for index, webcam := range Webcams.Images {
			var job *gocron.Job
			job, Cmd.Error = Cron.Scheduler.CronWithSeconds(webcam.Cron).StartImmediately().Tag(webcam.Prefix).Do(Webcams.Images[index].GetImage)
			if Cmd.Error != nil {
				LogPrintDate("crontab error: %s\n", Cmd.Error)
				break
			}
			crontab[webcam.Prefix] = job
		}

		if Webcams.Report.Cron != "" {
			_, Cmd.Error = Cron.Scheduler.CronWithSeconds(Webcams.Report.Cron).StartImmediately().Tag("report").Do(PrintJobs)
		}

		for index := range Webcams.Scripts {
			name := fmt.Sprintf("script-%d", index)
			// name := Webcams.Scripts[index].Cmd + " " + strings.Join(Webcams.Scripts[index].Args, " ")
			// _, Cmd.Error = Cron.Scheduler.CronWithSeconds(Webcams.Scripts[index].Cron).Tag(name).Do(func() {
			// 	RunCommand(Webcams.Scripts[index].Cmd, Webcams.Scripts[index].Args...)
			// })
			_, Cmd.Error = Cron.Scheduler.CronWithSeconds(Webcams.Scripts[index].Cron).Tag(name).Do(RunScript, name)
		}

		// fmt.Println(Cron.Scheduler.Location())
		// fmt.Println(Cron.Scheduler.Jobs())
		// fmt.Println(Cron.Scheduler.NextRun())

		Cron.Scheduler.RunAll()
		// PrintJobs()
		Cron.Scheduler.StartBlocking()

		if !Cron.Scheduler.IsRunning() {
			Cmd.Error = errors.New("cron scheduler has not started")
			break
		}

		// updateCounter := 0
		// timer := time.NewTicker(60 * time.Second)
		// for range timer.C {
		// 	if updateCounter < 5 {
		// 		updateCounter++
		// 		// LogPrintDate("Sleeping: %d\n", updateCounter)
		// 		continue
		// 	}
		// 	updateCounter = 0
		//
		// 	PrintJobs()
		// }

		time.Sleep(time.Hour * 5)
		Cron.Scheduler.Stop()
		fmt.Println(Cron.Scheduler.IsRunning())

		// Cmd.Error = config.Write("config.json")
		// if Cmd.Error != nil {
		// 	break
		// }
	}

	return Cmd.Error
}

func PrintJobs() {
	for range Only.Once {
		LogPrintDate("PrintJobs: %s\n", time.Now().Format("2006/01/02 15:04:05"))

		crontab := make(map[string]*gocron.Job)
		var jobs []string
		for _, key := range Cron.Scheduler.Jobs() {
			name := strings.Join(key.Tags(), " ")
			crontab[name] = key
			jobs = append(jobs, name)
		}
		sort.Strings(jobs)

		buf := new(bytes.Buffer)
		table := tablewriter.NewWriter(buf)
		table.SetHeader([]string{"Webcam", "Last Run", "Next Run", "Run Count", "Running", "Error"})
		for _, key := range jobs {
			job := crontab[key]
			table.Append([]string {
				strings.Join(job.Tags(), " "),
				job.LastRun().Format("2006/01/02 15:04:05"),
				job.NextRun().Format("2006/01/02 15:04:05"),
				fmt.Sprintf("%d", job.RunCount()),
				fmt.Sprintf("%v", job.IsRunning()),
				fmt.Sprintf("%v", job.Error()),
				// job.ScheduledAtTime(),
			})
		}
		table.Render()
		fmt.Println(buf.String())
	}
}

func RunScript(name string) {

	for range Only.Once {
		id := -1
		for _, key := range Cron.Scheduler.Jobs() {
			if strings.Join(key.Tags(), " ") != name {
				continue
			}
			ids := strings.Join(key.Tags(), "")
			ids = strings.TrimPrefix(ids, "script-")
			i, err := strconv.Atoi(ids)
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
		LogPrintDate("RunCommand: %s\n", time.Now().Format("2006/01/02 15:04:05"))
		LogPrintDate("Exec START: %s %v\n", job.Cmd, job.Args)

		c := exec.Command(job.Cmd, job.Args...)
		out, err := c.CombinedOutput()
		if err != nil {
			break
		}
		LogPrint("\n%s\n", string(out))
		LogPrintDate("Exec STOP: %s %v\n", job.Cmd, job.Args)
	}
}

func WebCron() error {
	for range Only.Once {
		// if Cmd.Web == nil {
		// 	Cmd.Error = errors.New("mqtt not available")
		// 	break
		// }
		//
		// if Cmd.Web.IsFirstRun() {
		// 	Cmd.Web.UnsetFirstRun()
		// } else {
		// 	time.Sleep(time.Second * 40) // Takes up to 40 seconds for data to come in.
		// }

		// web.Init(Cmd.Web, "config.json")
		// if ep.IsError() {
		// 	Cmd.Error = ep.GetError()
		// 	break
		// }
		//
		// data := ep.GetData()
		//
		// if Cmd.Mqtt.IsNewDay() {
		// 	LogPrintDate("New day: Configuring %d entries in HASSIO.\n", len(data.Entries))
		// 	for _, r := range data.Entries {
		// 		fmt.Printf(".")
		// 		// Cmd.Error = Cmd.Mqtt.SensorPublishConfig(r.PointId, r.PointName, r.Unit, i)
		// 		Cmd.Error = Cmd.Mqtt.SensorPublishConfig(r)
		// 		if Cmd.Error != nil {
		// 			break
		// 		}
		// 	}
		// 	fmt.Println()
		// }
		//
		// LogPrintDate("Updating %d entries to HASSIO.\n", len(data.Entries))
		// for _, r := range data.Entries {
		// 	fmt.Printf(".")
		// 	// Cmd.Error = Cmd.Mqtt.SensorPublishState(r.PointId, r.Value)
		// 	Cmd.Error = Cmd.Mqtt.SensorPublishValue(r)
		// 	if Cmd.Error != nil {
		// 		break
		// 	}
		// }
		// fmt.Println()
		// Cmd.Web.LastRefresh = time.Now()

		if Cmd.Error != nil {
			break
		}
	}

	if Cmd.Error != nil {
		LogPrintDate("Error: %s\n", Cmd.Error)
	}
	return Cmd.Error
}
