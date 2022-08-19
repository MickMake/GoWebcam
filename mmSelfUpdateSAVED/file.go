package mmSelfUpdate

import (
	"os"
	"path/filepath"
)


func (su *TypeSelfUpdate) CreateDummyBinary() error {
	for range onlyOnce {
		var err error

		result := FileStat(su.RuntimeBinary, su.TargetBinary)
		if result.CopyOfRuntime {
			su.AutoExec = true
			break
		}

		//if result.IsRuntimeBinary {
		//	// We are running as the bootstrap binary.
		//	su.State.SetOk()
		//	break
		//}

		if result.LinkToRuntime {
			err = os.Remove(su.TargetBinary)
			if err != nil {
				su.State.SetError(err)
				break
			}
			result.IsMissing = true
			su.AutoExec = true
		}

		if result.IsMissing {
			err = CopyFile(su.RuntimeBinary, su.TargetBinary)
			if err != nil {
				su.State.SetError(err)
				break
			}
			su.AutoExec = true
		}
	}

	return su.State
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

	for range onlyOnce {
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

	for range onlyOnce {
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
