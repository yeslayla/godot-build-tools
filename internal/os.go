package internal

type TargetOS uint8

const (
	TargetOSLinux TargetOS = iota
	TargetOSWindows
	TargetOSMacOS
)
