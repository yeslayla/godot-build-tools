package internal

import "strings"

type Godot4ArgBuilder struct {
	args []string
}

func NewGodot4ArgBuilder(projectDir string) GodotArgBuilder {
	return &Godot4ArgBuilder{
		args: []string{"--path", projectDir},
	}
}

func (b *Godot4ArgBuilder) AddHeadlessFlag() {
	b.args = append(b.args, "--headless")
}

func (b *Godot4ArgBuilder) AddDebugFlag() {
	b.args = append(b.args, "--debug")
}

func (b *Godot4ArgBuilder) AddVerboseFlag() {
	b.args = append(b.args, "--verbose")
}

func (b *Godot4ArgBuilder) AddQuietFlag() {
	b.args = append(b.args, "--quiet")
}

func (b *Godot4ArgBuilder) AddDumpGDExtensionInterfaceFlag() {
	b.args = append(b.args, "--dump-gdextension-interface")
}

func (b *Godot4ArgBuilder) AddDumpExtensionApiFlag() {
	b.args = append(b.args, "--dump-extension-api")
}

func (b *Godot4ArgBuilder) AddCheckOnlyFlag() {
	b.args = append(b.args, "--check-only")
}

func (b *Godot4ArgBuilder) AddExportFlag(exportType ExportType) {
	switch exportType {
	case ExportTypeRelease:
		b.args = append(b.args, "--export")
	case ExportTypeDebug:
		b.args = append(b.args, "--export-debug")
	case ExportTypePack:
		b.args = append(b.args, "--export-pack")
	}
}

func (b *Godot4ArgBuilder) GenerateArgs() string {
	return strings.Join(b.args, " ")
}
