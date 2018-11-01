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
	l.printWithTag(tag_info, []int{36, 1}, mess...)
}

func (l *Logger) Warn(mess ...interface{}) {
	l.printWithTag(tag_warn, []int{33, 1}, mess...)
}

func (l *Logger) Err(mess ...interface{}) {
	l.printWithTag(tag_error, []int{31, 1}, mess...)
}

func (l *Logger) Log(mess ...interface{}) {
	l.log(mess...)
}

func (l *Logger) log(mess ...interface{}) {
	logTime := time.Now().Format("2006/01/02 15:04:05")

	if runtime.GOOS != "windows" {
		logTime = l.ansi(36) + logTime + ansi_clear
	}

	fmt.Print(logTime, "  " + l.ansi(2), l.Name, " |  " + ansi_clear)
	fmt.Println(mess...)
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

func (l *Logger) printWithTag(tag string, tagColor []int, mess ...interface{}) {
	tagLog := "[" + tag + "]"

	if runtime.GOOS != "windows" {
		tagLog = l.ansi(tagColor...) + tagLog + ansi_clear
	}

	fmt.Print(tagLog, " ")
	l.log(mess...)
}
