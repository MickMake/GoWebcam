package mmVersion

import (
	"GoWebcam/Only"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
)


func printVersionSummary(release *selfupdate.Release) string {
	var ret string

	for range Only.Once {
		ret += fmt.Sprintf("\nExecutable: ")
		ret += fmt.Sprintf("%s ", release.RepoName)
		ret += fmt.Sprintf("%s\n", release.Version.String())

		ret += fmt.Sprintf("Url: ")
		ret += fmt.Sprintf("%s\n", release.URL)

		ret += fmt.Sprintf("Binary Size: ")
		ret += fmt.Sprintf("%d\n", release.AssetByteSize)

		ret += fmt.Sprintf("Published Date: ")
		ret += fmt.Sprintf("%s\n", release.PublishedAt.String())
	}

	return ret
}


func printVersion(release *selfupdate.Release) string {
	var ret string

	for range Only.Once {
		ret += fmt.Sprintf("Repository release information:\n")
		ret += fmt.Sprintf("Executable: %s v%s\n",
			fmt.Sprintf(release.RepoName),
			fmt.Sprintf(release.Version.String()),
		)

		ret += fmt.Sprintf("Url: %s\n", fmt.Sprintf(release.URL))

		//ret += fmt.Sprintf("TypeRepo Owner: %s\n", ux.SprintfBlue(release.RepoOwner))
		//ret += fmt.Sprintf("TypeRepo Name: %s\n", ux.SprintfBlue(release.RepoName))

		ret += fmt.Sprintf("Binary Size: %s\n", fmt.Sprintf("%d", release.AssetByteSize))

		ret += fmt.Sprintf("Published Date: %s\n", fmt.Sprintf(release.PublishedAt.String()))

		if release.ReleaseNotes != "" {
			ret += fmt.Sprintf("Release Notes: %s\n", fmt.Sprintf(release.ReleaseNotes))
		}
	}

	return ret
}


func stripUrlPrefix(url ...string) string {
	u := strings.Join(url, "/")
	u = strings.ReplaceAll(u, "//", "/")

	u = strings.TrimPrefix(u, "https://")
	u = strings.TrimPrefix(u, DefaultRepoServer)
	u = strings.TrimPrefix(u, "/")
	u = strings.TrimSuffix(u, "/")
	u = strings.TrimSpace(u)
	return u
}


func addUrlPrefix(url ...string) string {
	u := strings.Join(url, "/")

	switch {
		case strings.HasPrefix(u, "/"):
			u = "https://" + DefaultRepoServer + u

		case strings.HasPrefix(u, "github.com"):
			u = "https://" + u

		case strings.HasPrefix(u, "http"):
			// Leave url as is.

		default:
			u = "https://" + DefaultRepoServer + "/" + u
	}
	return u
}


func dropVprefix(v string) string {
	return strings.TrimPrefix(v, "v")
}

// Try and force the version array to conform to three values.
func fixVersion(v string) string {
	v = dropVprefix(v)
	sa := [3]string{"0", "0", "0"}
	for i, sav := range strings.Split(v, ".") {
		sa[i] = sav
	}
	return fmt.Sprintf("%s.%s.%s", sa[0], sa[1], sa[2])
}


func addVprefix(v string) string {
	return "v" + strings.TrimPrefix(v, "v")
}


func (v *Version) IsBootstrapBinary() bool {
	var ok bool
	for range Only.Once {
		if v.CmdName != v.CmdFile {
			break
		}
		if v.CmdName != BootstrapBinaryName {
			break
		}
		ok = true
	}
	return ok
}


func CopyFile(runtimeBin string, targetBin string) error {
	var err error

	for range Only.Once {
		var input []byte
		input, err = ioutil.ReadFile(runtimeBin)
		if err != nil {
			break
		}

		err = ioutil.WriteFile(targetBin, input, 0755)
		if err != nil {
			fmt.Println("Error creating", targetBin)
			break
		}
	}

	return err
}


func CompareBinary(runtimeBin string, newBin string) error {
	var err error

	for range Only.Once {
		var srcBin []byte
		srcBin, err = ioutil.ReadFile(runtimeBin)
		if err != nil {
			break
		}
		if srcBin == nil {
			break
		}

		var targetBin []byte
		targetBin, err = ioutil.ReadFile(newBin)
		if err != nil {
			break
		}
		if targetBin == nil {
			break
		}

		if len(srcBin) != len(targetBin) {
			break
		}

		for i := range srcBin {
			if srcBin[i] != targetBin[i] {
				err = errors.New("binary files differ")
				break
			}
		}
	}

	return err
}


func (v *Version) AutoRun() State {
	for range Only.Once {
		if !v.AutoExec {
			break
		}

		if v.IsBootstrapBinary() {
			// Let's avoid an endless loop.
			break
		}

		if len(v.FullArgs) > 0 {
			if v.FullArgs[0] == CmdVersion {
				// Let's avoid another endless loop.
				break
			}
		}

		fmt.Printf("Executing the real binary: '%s'\n", v.RuntimeBinary)
		c := exec.Command(v.TargetBinary, v.FullArgs...)

		var stdoutBuf, stderrBuf bytes.Buffer
		c.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
		c.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
		err := c.Run()
		waitStatus := c.ProcessState.Sys().(syscall.WaitStatus)
		waitStatus.ExitStatus()

		if err != nil {
			v.State.SetError(err.Error())
			break
		}

		v.State.SetOk("")
	}

	return v.State
}


func newHTTPClient(ctx context.Context, token string) *http.Client {
	if token == "" {
		return http.DefaultClient
	}
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	return oauth2.NewClient(ctx, src)
}


// 		updater := selfupdate.DefaultUpdater()
//		updater, err := selfupdate.NewUpdater()
//		selfupdate.UncompressCommand()
//		release, err := selfupdate.UpdateCommand()
//		release, err := selfupdate.UpdateSelf(semver.MustParse(su.version.ToString()), su.useRepo)
//		err := selfupdate.UpdateTo()