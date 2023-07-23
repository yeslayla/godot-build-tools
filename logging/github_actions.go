package logging

import (
	"fmt"
	"log"
	"os"
)

// GitHubActionsLogger is a logger that logs to GitHub Actions.
type GitHubActionsLogger struct {
	info  *log.Logger
	warn  *log.Logger
	err   *log.Logger
	debug *log.Logger
}

// NewGitHubActionsLogger creates a new GitHubActionsLogger.
func NewGitHubActionsLogger(debug bool) Logger {
	var debugLogger *log.Logger
	if debug {
		debugLogger = log.New(os.Stdout, "::debug::", 0)
	}

	return &GitHubActionsLogger{
		info:  log.New(os.Stdout, "", 0),
		warn:  log.New(os.Stdout, "::warning::", 0),
		err:   log.New(os.Stderr, "::error::", 0),
		debug: debugLogger,
	}
}

// Infof logs an info message.
func (l *GitHubActionsLogger) Infof(format string, args ...interface{}) {
	l.info.Printf(format, args...)
}

// Warnf logs a warning message.
func (l *GitHubActionsLogger) Warnf(format string, args ...interface{}) {
	l.warn.Printf(format, args...)
}

// Errorf logs an error message.
func (l *GitHubActionsLogger) Errorf(format string, args ...interface{}) {
	l.err.Printf(format, args...)
}

// Debugf logs a debug message if debug logging is enabled.
func (l *GitHubActionsLogger) Debugf(format string, args ...interface{}) {
	if l.debug != nil {
		l.debug.Printf(format, args...)
	}
}

// NoticeMessage sends a notice message to GitHub Actions.
func (l *GitHubActionsLogger) NoticeMessage(message string, input NoticeMessageInput) {
	var prefix string = "::notice"
	if input.Title != nil {
		prefix += " title=" + *input.Title
	}
	if input.Filename != nil {
		prefix += " file=" + *input.Filename
	}
	if input.Line != nil {
		prefix += " line=" + fmt.Sprint(*input.Line)
	}
	if input.EndLine != nil {
		prefix += " endLine=" + fmt.Sprint(*input.EndLine)
	}
	if input.Col != nil {
		prefix += " col=" + fmt.Sprint(*input.Col)
	}
	if input.EndCol != nil {
		prefix += " endColumn=" + fmt.Sprint(*input.EndCol)
	}

	l.info.Printf("%s::%s", prefix, message)
}

// StartGroup groups together log messages.
func (l *GitHubActionsLogger) StartGroup(name string) {
	l.info.Printf("::group::%s", name)
}

// EndGroup ends a group.
func (l *GitHubActionsLogger) EndGroup() {
	l.info.Println("::endgroup::")
}

// Mask masks a value in log output.
func (l *GitHubActionsLogger) Mask(value string) {
	l.info.Printf("::add-mask::%s", value)
}

// SetOutput sets an output parameter.
func (l *GitHubActionsLogger) SetOutput(name string, value string) {
	outputFile := os.Getenv("GITHUB_OUTPUT")
	if outputFile == "" {
		l.Errorf("GITHUB_OUTPUT is not set")
		return
	}

	f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		l.Errorf("failed to open output file: %v", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(fmt.Sprint(name, "=", value)); err != nil {
		l.Errorf("failed to write output file: %v", err)
		return
	}

}

// SetSummary sets a job's summary in markdown format.
func (l *GitHubActionsLogger) SetSummary(summary string) {
	summaryFile := os.Getenv("GITHUB_STEP_SUMMARY")
	if summaryFile == "" {
		l.Errorf("GITHUB_STEP_SUMMARY is not set")
		return
	}

	f, err := os.OpenFile(summaryFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		l.Errorf("failed to open summary file: %v", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(summary); err != nil {
		l.Errorf("failed to write summary file: %v", err)
		return
	}

}
