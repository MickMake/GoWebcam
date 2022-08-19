package mmSelfUpdate

import (
	"fmt"
	"github.com/google/go-github/v30/github"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"path/filepath"
	"runtime"
	"strings"
)


func (su *TypeSelfUpdate) Update() error {
	for range onlyOnce {
		if su.IsNotValid() {
			break
		}

		ux.PrintflnBlue("Checking '%s' for version greater than v%s", su.useRepo.GetUrl(), su.OldVersion.String())
		previous := su.OldVersion.ToSemVer()
		latest, err := su.ref.UpdateSelf(previous, su.useRepo.GetShortUrl())
		if err != nil {
			su.State.SetError(err)
			break
		}

		if previous.Equals(latest.Version) {
			ux.PrintflnOk("%s is up to date: v%s", su.useRepo.Name, latest.Version.String())
		} else {
			ux.PrintflnOk("%s updated to v%s", su.useRepo.Name, latest.Version)
			if latest.ReleaseNotes != "" {
				ux.PrintflnOk("%s %s Release Notes:%s", su.useRepo.Name, latest.Version, latest.ReleaseNotes)
			}
		}
	}

	return su.State
}


func (su *TypeSelfUpdate) PrintVersion(version *VersionValue) error {
	for range onlyOnce {
		if su.IsNotValid() {
			break
		}

		release := su.FetchVersion(version)
		if su.State.IsNotOk() {
			break
		}
		if release == nil {
			break
		}

		fmt.Printf(printVersion(release))
	}

	return su.State
}


func (su *TypeSelfUpdate) IsUpdated(print bool) error {
	for range onlyOnce {
		if su.IsNotValid() {
			break
		}

		latest := su.FetchVersion(su.GetVersionValue())
		if su.State.IsNotOk() {
			break
		}
		su.SetVersion(latest.Version.String())

		current := su.FetchVersion(su.OldVersion)
		if current == nil {
			if su.OldVersion.GT(latest.Version) {
				su.State.SetError("%s has version, (v%s), greater than repository, (v%s).",
					su.useRepo.Name,
					su.OldVersion.String(),
					latest.Version.String(),
				)
				if print {
					ux.PrintflnBlue("Current version (v%s)", su.OldVersion.String())
					ux.PrintflnYellow("Current version info unknown.")

					ux.PrintflnBlue("Updated version (v%s)", latest.Version.String())
					fmt.Printf(printVersion(latest))
					}
				break
			}

			su.State.SetWarning("%s can be updated to v%s.",
				su.useRepo.Name,
				latest.Version.String())
			if print {
				ux.PrintflnBlue("Current version (v%s)", su.OldVersion.String())
				ux.PrintflnYellow("Current version info unknown.")

				ux.PrintflnBlue("Updated version (v%s)", latest.Version.String())
				fmt.Printf(printVersion(latest))
			}
			su.State.SetOk()
			break
		}

		if current.Version.Equals(latest.Version) {
			su.State.SetOk("%s is up to date at v%s.",
				su.useRepo.Name,
				su.OldVersion.String())
			if print {
				ux.PrintflnOk("%s", su.State.GetOk())
				fmt.Printf(printVersion(current))
			}
			break
		}

		if current.Version.LE(latest.Version) {
			su.State.SetWarning("%s can be updated to v%s.",
				su.useRepo.Name,
				latest.Version.String())
			if print {
				ux.PrintflnWarning("%s", su.State.GetWarning())
				ux.PrintflnBlue("Current version (v%s)", current.Version.String())
				fmt.Printf(printVersion(current))

				ux.PrintflnBlue("Updated version (v%s)", latest.Version.String())
				fmt.Printf(printVersion(latest))
			}
			break
		}

		if current.Version.GT(latest.Version) {
			su.State.SetWarning("%s is more recent at v%s, (latest is %s).",
				su.useRepo.Name,
				su.OldVersion.String(),
				latest.Version.String())
			if print {
				ux.PrintflnWarning("%s", su.State.GetWarning())
				ux.PrintflnBlue("Current version (v%s)", current.Version.String())
				fmt.Printf(printVersion(current))

				ux.PrintflnBlue("Updated version (v%s)", latest.Version.String())
				fmt.Printf(printVersion(latest))
			}
			break
		}
	}

	return su.State
}


func (su *TypeSelfUpdate) Set(s SelfUpdateArgs) error {
	if s.owner != nil {
		su.useRepo.Owner = *s.owner
	}

	if s.name != nil {
		su.useRepo.Name = *s.name
	}

	if s.version != nil {
		su.SetVersion(*s.version)
	}

	if s.binaryRepo != nil {
		_ = su.binaryRepo.Set(*s.binaryRepo)
	}

	if s.sourceRepo != nil {
		_ = su.sourceRepo.Set(*s.sourceRepo)
	}

	if s.logging != nil {
		su.logging = (*FlagValue)(s.logging)
	} else {
		su.logging = &defaultFalse
	}

	if su.IsNotValid() {
		su.State.SetError("Invalid repo")
	}

	return su.State
}


func (su *TypeSelfUpdate) SetDebug(value bool) error {
	su.logging = (*FlagValue)(&value)
	if su.IsNotValid() {
		su.State.SetError("Invalid value")
	}
	return su.State
}


func (su *TypeSelfUpdate) SetName(value string) error {
	su.useRepo.Name = value
	if su.IsNotValid() {
		su.State.SetError("Invalid value")
	}
	return su.State
}

func (su *TypeSelfUpdate) GetName() string {
	return su.useRepo.Name
}


func (su *TypeSelfUpdate) SetOldVersion(value string) error {
	for range onlyOnce {
		if su.OldVersion != nil {
			su.State.SetOk()
			break
		}
		su.OldVersion = toVersionValue(value)
		if su.IsNotValid() {
			su.State.SetError("Invalid value")
		}
	}
	return su.State
}


func (su *TypeSelfUpdate) SetVersion(value string) error {
	if su.OldVersion == nil {
		su.OldVersion = su.useRepo.Version
	}

	su.useRepo.Version = toVersionValue(value)
	if su.IsNotValid() {
		su.State.SetError("Invalid value")
	}
	return su.State
}

func (su *TypeSelfUpdate) GetVersion() string {
	return su.useRepo.Version.String()
}

func (su *TypeSelfUpdate) FetchVersion(version *VersionValue) *selfupdate.Release {
	var release *selfupdate.Release

	for range onlyOnce {
		if su.IsNotValid() {
			break
		}

		var ok bool
		var err error

		repo := su.useRepo.GetShortUrl()

		switch {
			case version.IsNotValid():
				fallthrough
			case version.IsLatest():
				release, ok, err = su.ref.DetectLatest(repo)

			default:
				v := version.String()
				if !strings.HasPrefix(v, "v") {
					v = "v" + v
				}
				release, ok, err = su.ref.DetectVersion(repo, v)
		}

		if !ok {
			su.State.SetWarning("Version '%s' not found within '%s'", version.String(), su.useRepo.GetUrl())
			break
		}
		if err != nil {
			su.State.SetWarning("Version '%s' not found within '%s' - ERROR:%s", version.String(), su.useRepo.GetUrl(), err)
			break
		}

		su.State.SetOutput(release)
	}

	return release
}


func (su *TypeSelfUpdate) SetRepo(value ...string) error {
	for range onlyOnce {
		if su.OldVersion == nil {
			su.OldVersion = su.useRepo.Version
		}

		err := su.useRepo.Set(value...)
		if err != nil {
			su.State.SetError("Invalid value")
			break
		}

		su.config.Filters = addFilters(su.useRepo.Name, runtime.GOOS, runtime.GOARCH)
		su.ref, _ = selfupdate.NewUpdater(*su.config)

		su.TargetBinary = filepath.Join(su.Runtime.CmdDir, su.useRepo.Name)

		su.State.SetOk()
	}

	return su.State
}

func (su *TypeSelfUpdate) GetRepo() string {
	return su.useRepo.GetUrl()
}


func (su *TypeSelfUpdate) SetSourceRepo(value ...string) error {
	for range onlyOnce {
		if su.OldVersion == nil {
			su.OldVersion = su.sourceRepo.Version
		}

		_ = su.sourceRepo.Set(value...)
		if su.IsNotValid() {
			su.State.SetError("Invalid value")
		}
	}
	return su.State
}

func (su *TypeSelfUpdate) GetSourceRepo() string {
	return su.sourceRepo.GetUrl()
}


func (su *TypeSelfUpdate) SetBinaryRepo(value ...string) error {
	for range onlyOnce {
		if su.OldVersion == nil {
			su.OldVersion = su.binaryRepo.Version
		}

		_ = su.binaryRepo.Set(value...)
		if su.IsNotValid() {
			su.State.SetError("Invalid value")
		}
	}
	return su.State
}

func (su *TypeSelfUpdate) GetBinaryRepo() string {
	return su.binaryRepo.GetUrl()
}


func (su *TypeSelfUpdate) UpdateTo(newVersion *VersionValue) error {
	for range onlyOnce {
		if su.IsNotValid() {
			break
		}

		newRelease := su.FetchVersion(newVersion)
		if newRelease == nil {
			ux.PrintflnError("Version v%s doesn't exist for '%s'", su.useRepo.Version.String(), su.useRepo.Name)
			break
		}
		ux.PrintflnBlue("Updated version v%s => v%s", su.OldVersion.String(), newRelease.Version.String())

		err := su.ref.UpdateTo(newRelease, su.TargetBinary)
		if err != nil {
			su.State.SetError(err)
			break
		}

		ux.PrintflnOk("%s updated from v%s to v%s", su.useRepo.Name, su.OldVersion.String(), newRelease.Version.String())
	}

	return su.State
}


func (su *TypeSelfUpdate) GetVersionValue() *VersionValue {
	var ret *VersionValue

	for range onlyOnce {
		ret = su.useRepo.Version
	}

	return ret
}


func (su *TypeSelfUpdate) PrintVersionSummary(version string) error {
	for range onlyOnce {
		if su.IsNotValid() {
			break
		}

		var release *selfupdate.Release
		if version == LatestVersion {
			release = su.FetchVersion(nil)
		} else {
			release = su.FetchVersion(toVersionValue(version))
		}

		if release == nil {
			su.State.SetWarning("Version '%s' (%s) is empty.", version, toVersionValue(version))
			break
		}

		if version != release.Name {
			su.State.SetWarning("Version '%s' (%s) differs to repo '%s'.", version, toVersionValue(version), release.Name)
			// Don't show mismatched versions.
			break
		}

		fmt.Print(printVersionSummary(release))
	}

	return su.State
}


func (su *TypeSelfUpdate) PrintSummary(release *github.RepositoryRelease) error {
	for range onlyOnce {
		if su.IsNotValid() {
			break
		}

		if release == nil {
			break
		}

		ux.PrintfBlue("\nExecutable: ")
		ux.PrintfCyan("%s ", su.useRepo.Name)
		ux.PrintflnWhite("%s", release.GetName())

		ux.PrintfBlue("Url: ")
		ux.PrintflnWhite("%s", release.GetHTMLURL())

		ux.PrintfBlue("Binary Size: ")
		ux.PrintflnWhite("unknown")

		ux.PrintfBlue("Published Date: ")
		ux.PrintflnWhite("%s", release.GetPublishedAt().String())
	}

	return su.State
}
