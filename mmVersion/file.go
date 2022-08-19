package mmVersion

import (
	"GoWebcam/Only"
	"os"
	"path/filepath"
)


func (v *Version) CreateDummyBinary() State {
	for range Only.Once {
		var err error

		result := FileStat(v.RuntimeBinary, v.TargetBinary)
		if result.CopyOfRuntime {
			v.AutoExec = true
			break
		}

		//if result.IsRuntimeBinary {
		//	// We are running as the bootstrap binary.
		//	su.State.SetOk()
		//	break
		//}

		if result.LinkToRuntime {
			err = os.Remove(v.TargetBinary)
			if err != nil {
				v.State.SetError(err.Error())
				break
			}
			result.IsMissing = true
			v.AutoExec = true
		}

		if result.IsMissing {
			err = CopyFile(v.RuntimeBinary, v.TargetBinary)
			if err != nil {
				v.State.SetError(err.Error())
				break
			}
			v.AutoExec = true
		}
	}

	return v.State
}

type TargetFile struct {
	IsMissing bool
	IsRuntimeBinary bool
	FileMatches bool
	IsSymlink bool
	LinkTo string
	LinkEval string
	LinkToRuntime bool
	CopyOfRuntime bool

	Error error
	Info os.FileInfo
}


func FileStat(runtimeBinary string, targetBinary string) *TargetFile {
	var targetFile TargetFile

	for range Only.Once {
		targetFile.Info, targetFile.Error = os.Stat(targetBinary)
		if os.IsNotExist(targetFile.Error) {
			targetFile.IsMissing = true
		} else {
			targetFile.IsMissing = false

			if filepath.Base(runtimeBinary) == BootstrapBinaryName {
				targetFile.IsRuntimeBinary = true
			} else if runtimeBinary == targetBinary {
				targetFile.IsRuntimeBinary = true
				targetFile.CopyOfRuntime = true
			} else {
				targetFile.IsRuntimeBinary = false

				targetFile.Error = CompareBinary(runtimeBinary, targetBinary)
				if targetFile.Error == nil {
					targetFile.FileMatches = true
				} else {
					targetFile.FileMatches = false
				}
			}
		}

		targetFile.LinkTo, targetFile.Error = os.Readlink(targetBinary)
		if targetFile.LinkTo != "" {
			targetFile.IsSymlink = true

			targetFile.LinkEval, targetFile.Error = filepath.EvalSymlinks(targetBinary)
			if targetFile.LinkEval == "" {
				targetFile.LinkToRuntime = false
			} else {
				targetFile.LinkEval, targetFile.Error = filepath.Abs(targetFile.LinkEval)
				if targetFile.LinkEval == runtimeBinary {
					targetFile.LinkToRuntime = true
				} else if filepath.Base(targetFile.LinkEval) == BootstrapBinaryName {
					targetFile.LinkToRuntime = true
				} else {
					targetFile.LinkToRuntime = false
				}
			}
		} else {
			targetFile.IsSymlink = false
		}
	}

	return &targetFile
}


func ResolveFile(file string) string {
	var result string
	var err error

	for range Only.Once {
		_, err = os.Stat(file)
		if os.IsNotExist(err) {
			break
		}

		result, err = os.Readlink(file)
		if result == "" {
			result = file
			break
		}

		result, err = filepath.EvalSymlinks(file)
		if result == "" {
			result = file
			break
		}

		result, err = filepath.Abs(result)
		if result == "" {
			result = file
			break
		}
	}

	return result
}
