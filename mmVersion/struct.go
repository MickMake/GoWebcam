package mmVersion

import (
	"GoWebcam/Only"
	"fmt"
)


type SelfUpdateGetter interface {
}

type SelfUpdateArgs struct {
	owner      *string
	name       *string
	version    *string
	sourceRepo *string
	binaryRepo *string

	logging    *bool
}

// type TypeSelfUpdate struct {
// 	sourceRepo *UrlValue
// 	binaryRepo *UrlValue
// 	useRepo    *UrlValue
//
// 	OldVersion    *VersionValue
// 	TargetBinary  string
// 	RuntimeBinary string
// 	AutoExec      bool
//
// 	logging    *FlagValue
// 	config     *selfupdate.Config
// 	ref        *selfupdate.Updater
//
// 	Runtime    *TypeRuntime
// 	State      State
// 	cmd        *cobra.Command
// 	SelfCmd    *cobra.Command
// }

// func (su *TypeSelfUpdate) IsNil() error {
// 	return Ux.IfNilReturnError(su)
// }


// func New(rt *TypeRuntime) *TypeRuntime {
// 	rt = rt.EnsureNotNil()
//
// 	te := TypeRuntime {
// 		sourceRepo: toUrlValue(rt.CmdSourceRepo),
// 		binaryRepo: toUrlValue(rt.CmdBinaryRepo),
// 		useRepo:    nil,
//
// 		OldVersion:   nil,
// 		TargetBinary: rt.Cmd,
// 		RuntimeBinary: ResolveFile(rt.Cmd),
// 		AutoExec:     false,
//
// 		logging:    toBoolValue(rt.Debug),
// 		config:     &selfupdate.Config{
// 			APIToken:            "",
// 			EnterpriseBaseURL:   "",
// 			EnterpriseUploadURL: "",
// 			Validator:           nil, 	// &MyValidator{},
// 			Filters:             []string{},
// 		},
//
// 		Runtime: rt,
// 		State:   State{},
// 		cmd:     nil,
// 	}
// 	// te.State.SetPackage("")
// 	// te.State.SetFunctionCaller()
// 	te.useRepo = te.binaryRepo
//
// 	// Workaround for selfupdate not being flexible enough to support variable asset names
// 	// Should enable a template similar to GoReleaser.
// 	// EG: {{ .ProjectName }}-{{ .Os }}_{{ .Arch }}
// 	//var asset string
// 	//asset, te.State = toolGhr.GetAsset(rt.CmdBinaryRepo, "latest")
// 	//te.config.Filters = append(te.config.Filters, asset)
//
// 	// Ignore the above and just make sure all filenames are lowercase.
// 	te.config.Filters = append(te.config.Filters, addFilters(rt.CmdFile, runtime.GOOS, runtime.GOARCH)...)
// 	te.ref, _ = selfupdate.NewUpdater(*te.config)
// 	if *te.logging {
// 		selfupdate.EnableLog()
// 	}
//
// 	return &te
// }


//type MyValidator struct {
//}
//func (v *MyValidator) Validate(release, asset []byte) error {
//	calculatedHash := fmt.Sprintf("%x", sha256.Sum256(release))
//	hash := fmt.Sprintf("%s", asset[:sha256.BlockSize])
//	if calculatedHash != hash {
//		return fmt.Errorf("sha2: validation failed: hash mismatch: expected=%q, got=%q", calculatedHash, hash)
//	}
//	return nil
//}
//func (v *MyValidator) Suffix() string {
//	return ".gz"
//}


func addFilters(Binary string, Os string, Arch string) []string {
	var ret []string
	ret = append(ret, fmt.Sprintf("(?i)%s_.*_%s_%s.*", Binary, Os, Arch))
	ret = append(ret, fmt.Sprintf("(?i)%s_%s_%s.*", Binary, Os, Arch))
	ret = append(ret, fmt.Sprintf("(?i)%s-.*_%s_%s.*", Binary, Os, Arch))
	ret = append(ret, fmt.Sprintf("(?i)%s-%s_%s.*", Binary, Os, Arch))
	if Arch == "amd64" {
		// This is recursive - so be careful what you place in the "Arch" argument.
		ret = append(ret, addFilters(Binary, Os, "x86_64.*")...)
		ret = append(ret, addFilters(Binary, Os, "64.*")...)
		ret = append(ret, addFilters(Binary, Os, "64bit.*")...)
	}
	return ret
}


func (v *Version) IsValid() bool {
	var ok bool
	for range Only.Once {
		if v.useRepo.Owner == "" {
			v.State.SetWarning("rep owner is not defined - selfupdate disabled")
			break
		}

		if v.useRepo.Name == "" {
			v.State.SetWarning("repo name is not defined - selfupdate disabled")
			break
		}

		// Refer to binary repo definition first.
		if v.CmdBinaryRepo.IsValid() {
			v.useRepo = &v.CmdBinaryRepo
			v.State.SetOk("")
			ok = true
			break
		}

		// If binary repo is not set, use source repo.
		if v.CmdSourceRepo.IsValid() {
			v.useRepo = &v.CmdSourceRepo
			v.State.SetOk("")
			ok = true
			break
		}

		v.State.SetWarning(errorNoRepo)
	}

	return ok
}
func (v *Version) IsNotValid() bool {
	return !v.IsValid()
}

func (v *Version) getRepo() string {
	var ret string

	for range Only.Once {
		if v.CmdBinaryRepo.IsValid() {
			ret = v.CmdBinaryRepo.String()
			break
		}
		if v.CmdSourceRepo.IsValid() {
			ret = v.CmdSourceRepo.String()
			break
		}
	}

	return ret
}
