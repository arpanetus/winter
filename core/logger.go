package core

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strconv"
	"time"
)

var (
	longestLogger = 0
)

func NewLogger(name string) *Logger {
	nameLen := len(name)
	if nameLen > longestLogger {
		longestLogger = nameLen
	}

	logger := &Logger{
		Name: name,
		writer: fmt.Print,
		writerf: fmt.Printf,
		writerln: fmt.Println,
	}

	return logger
}

func LogDefaultToFiles(filePath string, prefix ...string) {
	MainLogger.LogIntoFile(filePath, prefix...)
	RequestLogger.LogIntoFile(filePath, prefix...)
	RouterLogger.LogIntoFile(filePath, prefix...)
	WebSocketLogger.LogIntoFile(filePath, prefix...)
}

func (l *Logger) LogIntoFile(filePath string, prefix ...string) *Logger {
	defaultPrefix := "log"
	if len(prefix) > 0 {
		defaultPrefix = prefix[0]
	}

	logFilePath := defaultPrefix + "-" + l.Name + "-" + time.Now().Format("2006-01-02") + ".log"
	fullPath := path.Join(filePath, logFilePath)

	l.Info("Logging into file:", fullPath)

	logFile, err := os.Create(fullPath)
	if err != nil {
		l.Err("Could not create new log file, logging into terminal")
		return l
	}

	l.fileWriter = bufio.NewWriter(logFile)

	l.overrideWithWriter(l.fileWriter)

	l.logFile = logFile

	l.logIntoFile = true
	return l
}

func (l *Logger) Finfo(writer io.Writer, mess ...interface{}) {
	l.overrideWithWriter(writer)
	l.log(tag_info, []int{36, 1}, false, "", mess...)
}

func (l *Logger) Finfof(writer io.Writer, format string, mess ...interface{}) {
	l.overrideWithWriter(writer)
	l.log(tag_info, []int{36, 1}, true, format, mess...)
}

func (l *Logger) Fwarn(writer io.Writer, mess ...interface{}) {
	l.overrideWithWriter(writer)
	l.log(tag_warn, []int{33, 1}, false, "", mess...)
}

func (l *Logger) Fwarnf(writer io.Writer, format string, mess ...interface{}) {
	l.overrideWithWriter(writer)
	l.log(tag_warn, []int{33, 1}, true, format, mess...)
}

func (l *Logger) Ferr(writer io.Writer, mess ...interface{}) {
	l.overrideWithWriter(writer)
	l.log(tag_error, []int{31, 1}, false, "", mess...)
}

func (l *Logger) Ferrf(writer io.Writer, format string, mess ...interface{}) {
	l.overrideWithWriter(writer)
	l.log(tag_error, []int{31, 1}, true, format, mess...)
}

func (l *Logger) Fnote(writer io.Writer, mess ...interface{}) {
	l.overrideWithWriter(writer)
	l.log(tag_note, []int{34, 1}, false, "", mess...)
}

func (l *Logger) Fnotef(writer io.Writer, format string, mess ...interface{}) {
	l.overrideWithWriter(writer)
	l.log(tag_note, []int{34, 1}, true, format, mess...)
}

func (l *Logger) Flog(writer io.Writer, mess ...interface{}) {
	l.overrideWithWriter(writer)
	l.log("", []int{}, false, "", mess...)
}

func (l *Logger) Flogf(writer io.Writer, format string, mess ...interface{}) {
	l.overrideWithWriter(writer)
	l.log("", []int{}, true, format, mess...)
}

func (l *Logger) Info(mess ...interface{}) {
	l.log(tag_info, []int{36, 1}, false, "", mess...)
}

func (l *Logger) Infof(format string, mess ...interface{}) {
	l.log(tag_info, []int{36, 1}, true, format, mess...)
}

func (l *Logger) Warn(mess ...interface{}) {
	l.log(tag_warn, []int{33, 1}, false, "", mess...)
}

func (l *Logger) Warnf(format string, mess ...interface{}) {
	l.log(tag_warn, []int{33, 1}, true, format, mess...)
}

func (l *Logger) Err(mess ...interface{}) {
	l.log(tag_error, []int{31, 1}, false, "", mess...)
}

func (l *Logger) Errf(format string, mess ...interface{}) {
	l.log(tag_error, []int{31, 1}, true, format, mess...)
}

func (l *Logger) Note(mess ...interface{}) {
	l.log(tag_note, []int{34, 1}, false, "", mess...)
}

func (l *Logger) Notef(format string, mess ...interface{}) {
	l.log(tag_note, []int{34, 1}, true, format, mess...)
}

func (l *Logger) Log(mess ...interface{}) {
	l.log("", []int{}, false, "", mess...)
}

func (l *Logger) Logf(format string, mess ...interface{}) {
	l.log("", []int{}, true, format, mess...)
}

func (l *Logger) log(tag string, tagColor []int, format bool, formatString string, mess ...interface{}) {
	logTime := time.Now().Format("2006/01/02 15:04:05")
	ansiRequired := runtime.GOOS != bad_os && !l.logIntoFile
	loggerName := l.Name + l.getFillerSpaces() + " |  "

	if ansiRequired {
		logTime = l.ansi(36) + logTime + ansi_clear
		loggerName = l.ansi(2) + loggerName + ansi_clear
	}

	if len(tag) > 0 {
		tagLog := "[" + tag + "]"

		if ansiRequired {
			tagLog = l.ansi(tagColor...) + tagLog + ansi_clear
		}

		l.writer(tagLog, " ")
	}

	l.writer(logTime, " ", loggerName)

	if format {
		l.writerf(formatString, mess...)
		l.writerln()
	} else {
		l.writerln(mess...)
	}

	if l.logIntoFile {
		l.fileWriter.Flush()
	}
}

func (l *Logger) ansi(codes ...int) string {
	ansiSpaceCode := ansi_prefix
	for _, n := range codes {
		if ansiSpaceCode == ansi_prefix {
			ansiSpaceCode = ansiSpaceCode + strconv.Itoa(n)
		}

		ansiSpaceCode = ansiSpaceCode + ";" + strconv.Itoa(n)
	}
	ansiSpaceCode = ansiSpaceCode + ansi_suffix

	return ansiSpaceCode
}

func (l *Logger) getFillerSpaces() string {
	spaces := ""
	count := longestLogger - len(l.Name)
	for i := 0; i < count; i++ {
		spaces = spaces + " "
	}
	return spaces
}

func (l *Logger) overrideWithWriter(writer io.Writer) {
	l.writer = func(a ...interface{}) (n int, err error) {
		return fmt.Fprint(writer, a...)
	}
	l.writerln = func(a ...interface{}) (n int, err error) {
		return fmt.Fprintln(writer, a...)
	}
	l.writerf = func(format string, a ...interface{}) (n int, err error) {
		return fmt.Fprintf(writer, format, a...)
	}
}
