package mmSelfUpdate

import (
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

	for range onlyOnce {
		ret += ux.SprintfBlue("\nExecutable: ")
		ret += ux.SprintfCyan("%s ", release.RepoName)
		ret += ux.SprintfWhite("%s\n", release.Version.String())

		ret += ux.SprintfBlue("Url: ")
		ret += ux.SprintfWhite("%s\n", release.URL)

		ret += ux.SprintfBlue("Binary Size: ")
		ret += ux.SprintfWhite("%d\n", release.AssetByteSize)

		ret += ux.SprintfBlue("Published Date: ")
		ret += ux.SprintfWhite("%s\n", release.PublishedAt.String())
	}

	return ret
}


func printVersion(release *selfupdate.Release) string {
	var ret string

	for range onlyOnce {
		ret += ux.SprintfBlue("Repository release information:\n")
		ret += fmt.Sprintf("Executable: %s v%s\n",
			ux.SprintfBlue(release.RepoName),
			ux.SprintfWhite(release.Version.String()),
		)

		ret += fmt.Sprintf("Url: %s\n", ux.SprintfBlue(release.URL))

		//ret += fmt.Sprintf("TypeRepo Owner: %s\n", ux.SprintfBlue(release.RepoOwner))
		//ret += fmt.Sprintf("TypeRepo Name: %s\n", ux.SprintfBlue(release.RepoName))

		ret += fmt.Sprintf("Binary Size: %s\n", ux.SprintfBlue("%d", release.AssetByteSize))

		ret += fmt.Sprintf("Published Date: %s\n", ux.SprintfBlue(release.PublishedAt.String()))

		if release.ReleaseNotes != "" {
			ret += fmt.Sprintf("Release Notes: %s\n", ux.SprintfBlue(release.ReleaseNotes))
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


func (su *TypeSelfUpdate) IsBootstrapBinary() bool {
	var ok bool
	for range onlyOnce {
		if su.Runtime.CmdName != su.Runtime.CmdFile {
			break
		}
		if su.Runtime.CmdName != BootstrapBinaryName {
			break
		}
		ok = true
	}
	return ok
}


func CopyFile(runtimeBin string, targetBin string) error {
	var err error

	for range onlyOnce {
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

	for range onlyOnce {
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


func (su *TypeSelfUpdate) AutoRun() error {
	for range onlyOnce {
		if !su.AutoExec {
			break
		}

		if su.IsBootstrapBinary() {
			// Let's avoid an endless loop.
			break
		}

		if len(su.Runtime.FullArgs) > 0 {
			if su.Runtime.FullArgs[0] == CmdVersion {
				// Let's avoid another endless loop.
				break
			}
		}

		ux.PrintflnNormal("Executing the real binary: '%s'", su.RuntimeBinary)
		c := exec.Command(su.TargetBinary, su.Runtime.FullArgs...)

		var stdoutBuf, stderrBuf bytes.Buffer
		c.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
		c.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
		err := c.Run()
		waitStatus := c.ProcessState.Sys().(syscall.WaitStatus)
		waitStatus.ExitStatus()

		if err != nil {
			su.State.SetError(err)
			break
		}

		su.State.SetOk()
	}

	return su.State
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
