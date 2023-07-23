package internal

type ExportType uint8

const (
	ExportTypeRelease ExportType = iota
	ExportTypeDebug
	ExportTypePack
)

type GodotArgBuilder interface {
	AddHeadlessFlag()
	AddDebugFlag()
	AddVerboseFlag()
	AddQuietFlag()

	AddDumpGDExtensionInterfaceFlag()
	AddDumpExtensionApiFlag()
	AddCheckOnlyFlag()

	AddExportFlag(exportType ExportType)

	GenerateArgs() string
}
