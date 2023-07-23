package logging

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Logger is an interface for logging.
type Logger interface {
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debugf(format string, args ...interface{})

	Mask(value string)

	StartGroup(name string)
	EndGroup()

	NoticeMessage(message string, input NoticeMessageInput)
	SetOutput(name string, value string)
	SetSummary(summary string)
}

// DefaultLogger is a logger that logs to the console.
type DefaultLogger struct {
	info  *log.Logger
	warn  *log.Logger
	err   *log.Logger
	debug *log.Logger

	outputsFile string
	summaryFile string

	groups []string
	masks  []string
}

// NoticeMessageInput holds optional parameters for a notice message.
type NoticeMessageInput struct {
	Title    *string
	Filename *string
	Line     *int
	EndLine  *int
	Col      *int
	EndCol   *int
}

// LoggerOptions holds options for creating a new logger.
type LoggerOptions struct {
	OutputsFile string
	SummaryFile string
	Debug       bool
}

// NewLogger creates a new default logger.
func NewLogger(options *LoggerOptions) Logger {

	var debugLogger *log.Logger
	if options.Debug {
		debugLogger = log.New(os.Stdout, "DEBUG ", log.LstdFlags)
	}

	return &DefaultLogger{
		info:   log.New(os.Stdout, "INFO ", log.LstdFlags),
		warn:   log.New(os.Stdout, "WARNING ", log.LstdFlags),
		err:    log.New(os.Stderr, "ERROR ", log.LstdFlags),
		debug:  debugLogger,
		groups: []string{},
		masks:  []string{},

		outputsFile: options.OutputsFile,
		summaryFile: options.SummaryFile,
	}
}

// formatMessage wraps the message with the current group and removes any masked values.
func (l *DefaultLogger) formatMessage(message string) string {
	message = l.removeMasks(message)
	message = l.addGroups(message)
	return message
}

// removeMasks removes any masked values from the message.
func (l *DefaultLogger) removeMasks(message string) string {
	for _, mask := range l.masks {
		message = strings.ReplaceAll(message, mask, "********")
	}
	return message
}

// addGroups adds the current group to the message.
func (l *DefaultLogger) addGroups(message string) string {
	for _, group := range l.groups {
		message = fmt.Sprint(group, " - ", message)
	}
	return message
}

// Infof logs an info message.
func (l *DefaultLogger) Infof(format string, args ...interface{}) {
	l.info.Printf(l.formatMessage(format), args...)
}

// Warnf logs a warning message.
func (l *DefaultLogger) Warnf(format string, args ...interface{}) {
	l.warn.Printf(l.formatMessage(format), args...)
}

// Errorf logs an error message.
func (l *DefaultLogger) Errorf(format string, args ...interface{}) {
	l.err.Printf(l.formatMessage(format), args...)
}

// Debugf logs a debug message if debug logging is enabled.
func (l *DefaultLogger) Debugf(format string, args ...interface{}) {
	if l.debug != nil {
		l.debug.Printf(l.formatMessage(format), args...)
	}
}

// NoticeMessage sends a notice about a line.
func (l *DefaultLogger) NoticeMessage(message string, input NoticeMessageInput) {
	var prefix string = ""
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

	l.info.Printf(l.formatMessage(fmt.Sprintf("%s %s", prefix, message)))
}

// Mask hides a value in the log output.
func (l *DefaultLogger) Mask(value string) {
	l.masks = append(l.masks, value)
}

// StartGroup groups together log messages.
func (l *DefaultLogger) StartGroup(name string) {
	l.groups = append(l.groups, name)
}

// EndGroup ends a group.
func (l *DefaultLogger) EndGroup() {
	if len(l.groups) > 0 {
		l.groups = l.groups[:len(l.groups)-1]
	}
}

// SetOutput outputs a key-value pair to the outputs file.
func (l *DefaultLogger) SetOutput(name string, value string) {
	if l.outputsFile == "" {
		return
	}

	f, err := os.OpenFile(l.outputsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		l.Errorf("failed to open outputs file: %s", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(fmt.Sprintf("%s=%s", name, value)); err != nil {
		l.Errorf("failed to write to outputs file: %s", err)
		return
	}
}

// SetSummary outputs a markdown-formatted summary to the summary file.
func (l *DefaultLogger) SetSummary(summary string) {
	if l.summaryFile == "" {
		return
	}

	f, err := os.OpenFile(l.summaryFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		l.Errorf("failed to open summary file: %s", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(summary); err != nil {
		l.Errorf("failed to write to summary file: %s", err)
		return
	}
}
