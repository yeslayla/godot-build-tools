package internal

import "runtime"

type TargetOS uint8

const (
	TargetOSLinux TargetOS = iota
	TargetOSWindows
	TargetOSMacOS
)

func (t TargetOS) String() string {
	switch t {
	case TargetOSLinux:
		return "linux"
	case TargetOSWindows:
		return "windows"
	case TargetOSMacOS:
		return "macos"
	}
	return ""
}

func NewTargetOSFromRuntime(GOOSRuntime string) TargetOS {
	switch GOOSRuntime {
	case "linux":
		return TargetOSLinux
	case "windows":
		return TargetOSWindows
	case "darwin":
		return TargetOSMacOS
	}
	return TargetOSLinux
}

func CurrentTargetOS() TargetOS {
	return NewTargetOSFromRuntime(runtime.GOOS)
}
