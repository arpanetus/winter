package core

import (
	"fmt"
	"runtime"
	"strconv"
	"time"
)

func NewLogger(name string) *Logger {
	return &Logger{
		Name: name,
	}
}

func (l *Logger) Info(mess ...interface{}) {
	l.printWithTag(false, "", tag_info, []int{36, 1}, mess...)
}

func (l *Logger) Infof(format string, mess ...interface{}) {
	l.printWithTag(true, format, tag_info, []int{36, 1}, mess...)
}

func (l *Logger) Warn(mess ...interface{}) {
	l.printWithTag(false, "", tag_warn, []int{33, 1}, mess...)
}

func (l *Logger) Warnf(format string, mess ...interface{}) {
	l.printWithTag(true, format, tag_warn, []int{33, 1}, mess...)
}

func (l *Logger) Err(mess ...interface{}) {
	l.printWithTag(false, "", tag_error, []int{31, 1}, mess...)
}

func (l *Logger) Errf(format string, mess ...interface{}) {
	l.printWithTag(true, format, tag_error, []int{31, 1}, mess...)
}

func (l *Logger) Note(mess ...interface{}) {
	l.printWithTag(false, "", tag_note, []int{34, 1}, mess...)
}

func (l *Logger) Notef(format string, mess ...interface{}) {
	l.printWithTag(true, format, tag_note, []int{34, 1}, mess...)
}

func (l *Logger) Log(mess ...interface{}) {
	l.log(false, "", mess...)
}

func (l *Logger) Logf(format string, mess ...interface{}) {
	l.log(true, format, mess...)
}

func (l *Logger) log(format bool, formatString string, mess ...interface{}) {
	logTime := time.Now().Format("2006/01/02 15:04:05")

	if runtime.GOOS != "windows" {
		logTime = l.ansi(36) + logTime + ansi_clear
	}

	fmt.Print(logTime, "  " + l.ansi(2), l.Name, " |  " + ansi_clear)
	if format {
		fmt.Printf(formatString, mess...)
		fmt.Println()
	} else {
		fmt.Println(mess...)
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

func (l *Logger) printWithTag(format bool, formatStirng string, tag string, tagColor []int, mess ...interface{}) {
	tagLog := "[" + tag + "]"

	if runtime.GOOS != "windows" {
		tagLog = l.ansi(tagColor...) + tagLog + ansi_clear
	}

	fmt.Print(tagLog, " ")
	l.log(format, formatStirng, mess...)
}
