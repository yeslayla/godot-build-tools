package internal

import "strings"

type DefaultGodotArgBuilder struct {
	args []string
}

func NewGodotArgBuilder(projectDir string) GodotArgBuilder {
	return &DefaultGodotArgBuilder{
		args: []string{"--path", projectDir},
	}
}

func (b *DefaultGodotArgBuilder) AddHeadlessFlag() {
	b.args = append(b.args, "--headless")
}

func (b *DefaultGodotArgBuilder) AddDebugFlag() {
	b.args = append(b.args, "--debug")
}

func (b *DefaultGodotArgBuilder) AddVerboseFlag() {
	b.args = append(b.args, "--verbose")
}

func (b *DefaultGodotArgBuilder) AddQuietFlag() {
	b.args = append(b.args, "--quiet")
}

func (b *DefaultGodotArgBuilder) AddDumpGDExtensionInterfaceFlag() {
	b.args = append(b.args, "--dump-gdextension-interface")
}

func (b *DefaultGodotArgBuilder) AddDumpExtensionApiFlag() {
	b.args = append(b.args, "--dump-extension-api")
}

func (b *DefaultGodotArgBuilder) AddCheckOnlyFlag() {
	b.args = append(b.args, "--check-only")
}

func (b *DefaultGodotArgBuilder) AddExportFlag(exportType ExportType) {
	switch exportType {
	case ExportTypeRelease:
		b.args = append(b.args, "--export")
	case ExportTypeDebug:
		b.args = append(b.args, "--export-debug")
	case ExportTypePack:
		b.args = append(b.args, "--export-pack")
	}
}

func (b *DefaultGodotArgBuilder) GenerateArgs() string {
	return strings.Join(b.args, " ")
}
