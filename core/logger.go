package core

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"
)

func NewLogger(name string) *Logger {
	return &Logger{
		Name: name,
		writer: fmt.Print,
		writerf: fmt.Printf,
		writerln: fmt.Println,
	}
}

func (l *Logger) LogIntoFile(filePath string, ) *Logger {
	logFile, err := os.Create(filePath + "-" + time.Now().Format("2006-01-02") + ".log")
	if err != nil {
		l.Err("Could not create new log file, logging into terminal")
		return l
	}

	l.fileWriter = bufio.NewWriter(logFile)

	l.writer = func(a ...interface{}) (n int, err error) {
		num, err := fmt.Fprint(l.fileWriter, a...)
		return num, err
	}
	l.writerf = func(format string, a ...interface{}) (n int, err error) {
		num, err := fmt.Fprintf(l.fileWriter, format, a...)
		return num, err
	}
	l.writerln = func(a ...interface{}) (n int, err error) {
		num, err := fmt.Fprintln(l.fileWriter, a...)
		return num, err
	}

	l.logFile = logFile

	l.logIntoFile = true
	return l
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
	loggerName := l.Name + " |  "

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
