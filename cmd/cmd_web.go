package cmd

import (
	"GoWebcam/Only"
	"GoWebcam/mmWebcam"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"sync"
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
		Webcams, Cmd.Error = mmWebcam.ReadConfig("config.json")
		if Cmd.Error != nil {
			break
		}

		LogPrintDate("One-off fetch of webcams from config...\n")
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

		LogPrintDate("Cron based webcam fetch from config...\n")
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
		LogPrintf("\n%s", buf.String())
	}
}

func RunScript(name string) error {
	var err error

	for range Only.Once {
		id := -1
		for _, key := range Cron.Scheduler.Jobs() {
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

		err = Exec(job.Cmd, job.Args...)
	}

	return err
}


func Exec(command string, args ...string) error {
	var err error

	for range Only.Once {
		LogPrintDate("Exec START: %s %v\n", command, args)

		cmd := exec.Command(command, args...)
		// out, err := cmd.CombinedOutput()
		// if err != nil {
		// 	break
		// }
		// LogPrintf("\n%s\n", string(out))

		var stdout io.ReadCloser
		stdout, err = cmd.StdoutPipe()
		if err != nil {
			break
		}

		// var stderr io.ReadCloser
		// stderr, err = cmd.StderrPipe()
		// if err != nil {
		// 	break
		// }

		// start the command after having set up the pipe
		err = cmd.Start()
		if err != nil {
			break
		}

		// read command's stdout line by line
		in := bufio.NewScanner(stdout)
		// inerr := bufio.NewScanner(stderr)

		// go func(){
		// 	for in.Scan() {
		// 		LogPrintf(inerr.Text()) // write each line to your log, or anything you need
		// 	}
		// }()

		for in.Scan() {
			LogPrintf(in.Text()) // write each line to your log, or anything you need
		}

		err = in.Err()
		if err != nil {
			LogPrintf("error: %s", err)
		}

		LogPrintDate("Exec STOP: %s %v\n", command, args)
	}

	return err
}


func Exec1(command string, args ...string) error {
	var err error

	for range Only.Once {
		LogPrintDate("Exec1 START: %s %v\n", command, args)

		cmd := exec.Command(command, args...)

		var stdoutBuf, stderrBuf bytes.Buffer
		cmd.Stdout = io.MultiWriter(&stdoutBuf)	// os.Stdout, &stdoutBuf)
		cmd.Stderr = io.MultiWriter(&stderrBuf)	// os.Stderr, &stderrBuf)

		err = cmd.Run()
		if err != nil {
			LogPrintf("cmd.Run() failed with %s\n", err)
		}
		outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
		fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)

		// out, err := cmd.CombinedOutput()
		// if err != nil {
		// 	break
		// }
		// LogPrintf("\n%s\n", string(out))

		LogPrintDate("Exec STOP: %s %v\n", command, args)
	}

	return err
}


func Exec2(command string, args ...string) error {
	var err error

	for range Only.Once {
		LogPrintDate("Exec2 START: %s %v\n", command, args)

		cmd := exec.Command(command, args...)

		var stdout, stderr []byte
		var errStdout, errStderr error
		stdoutIn, _ := cmd.StdoutPipe()
		stderrIn, _ := cmd.StderrPipe()
		err = cmd.Start()
		if err != nil {
			LogPrintf("cmd.Start() failed with '%s'\n", err)
			break
		}

		// cmd.Wait() should be called only after we finish reading
		// from stdoutIn and stderrIn.
		// wg ensures that we finish
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			stdout, errStdout = copyAndCapture(os.Stdout, stdoutIn)
			wg.Done()
		}()
		stderr, errStderr = copyAndCapture(os.Stderr, stderrIn)
		wg.Wait()

		err = cmd.Wait()
		if err != nil {
			LogPrintf("cmd.Run() failed with %s\n", err)
			break
		}
		if errStdout != nil || errStderr != nil {
			LogPrintf("failed to capture stdout or stderr\n")
			break
		}
		outStr, errStr := string(stdout), string(stderr)
		fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)

		LogPrintDate("Exec STOP: %s %v\n", command, args)
	}

	return err
}

func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			_, err := w.Write(d)
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
}


func Exec3(command string, args ...string) error {
	var err error

	for range Only.Once {
		LogPrintDate("Exec3 START: %s %v\n", command, args)

		cmd := exec.Command(command, args...)

		var errStdout, errStderr error
		stdoutIn, _ := cmd.StdoutPipe()
		stderrIn, _ := cmd.StderrPipe()
		stdout := NewCapturingPassThroughWriter(os.Stdout)
		stderr := NewCapturingPassThroughWriter(os.Stderr)
		err = cmd.Start()
		if err != nil {
			LogPrintf("cmd.Start() failed with '%s'\n", err)
			break
		}

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			_, errStdout = io.Copy(stdout, stdoutIn)
			wg.Done()
		}()

		_, errStderr = io.Copy(stderr, stderrIn)
		wg.Wait()

		err = cmd.Wait()
		if err != nil {
			LogPrintf("cmd.Run() failed with %s\n", err)
			break
		}
		if errStdout != nil || errStderr != nil {
			LogPrintf("failed to capture stdout or stderr\n")
			break
		}

		outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
		fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)

		LogPrintDate("Exec STOP: %s %v\n", command, args)
	}

	return err
}

// CapturingPassThroughWriter is a writer that remembers
// data written to it and passes it to w
type CapturingPassThroughWriter struct {
	buf bytes.Buffer
	w io.Writer
}

// NewCapturingPassThroughWriter creates new CapturingPassThroughWriter
func NewCapturingPassThroughWriter(w io.Writer) *CapturingPassThroughWriter {
	return &CapturingPassThroughWriter{
		w: w,
	}
}

func (w *CapturingPassThroughWriter) Write(d []byte) (int, error) {
	w.buf.Write(d)
	return w.w.Write(d)
}

// Bytes returns bytes written to the writer
func (w *CapturingPassThroughWriter) Bytes() []byte {
	return w.buf.Bytes()
}
