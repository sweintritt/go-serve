package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type loglevel int

const (
	loglevel_debug   = 0
	loglevel_info    = 1
	loglevel_warning = 2
	loglevel_error   = 3
	loglevel_fatal   = 4
)

func toString(l loglevel) string {
	switch l {
	case loglevel_debug:
		return "DEBUG"
	case loglevel_info:
		return "INFO "
	case loglevel_warning:
		return "WARN "
	case loglevel_error:
		return "ERROR"
	case loglevel_fatal:
		return "FATAL"
	default:
		return "DEBUG"
	}
}

func getPrefix(l loglevel, t time.Time) string {
	prefix := toString(l)
	prefix += " "
	prefix += fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.%3d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), t.Nanosecond()/1000000)

	// getPrefix -> Print -> Debug -> Caller => 3
	_, file, line, ok := runtime.Caller(3)
	if ok {
		folders := strings.Split(file, "/")
		file = folders[len(folders)-1]
		prefix += " [" + file + ":" + strconv.Itoa(line) + "]"
	}

	prefix += " "
	return prefix
}

func Debug(v ...interface{}) {
	Print(os.Stdout, loglevel_debug, time.Now(), v...)
}

func Debugf(format string, v ...interface{}) {
	Printf(os.Stdout, loglevel_debug, time.Now(), format, v...)
}

func Info(v ...interface{}) {
	Print(os.Stdout, loglevel_info, time.Now(), v...)
}

func Infof(format string, v ...interface{}) {
	Printf(os.Stdout, loglevel_info, time.Now(), format, v...)
}

func Warning(v ...interface{}) {
	Print(os.Stdout, loglevel_warning, time.Now(), v...)
}

func Warningf(format string, v ...interface{}) {
	Printf(os.Stdout, loglevel_warning, time.Now(), format, v...)
}

func Error(v ...interface{}) {
	Print(os.Stderr, loglevel_error, time.Now(), v...)
}

func Errorf(format string, v ...interface{}) {
	Printf(os.Stderr, loglevel_error, time.Now(), format, v...)
}

func Fatal(v ...interface{}) {
	Print(os.Stderr, loglevel_fatal, time.Now(), v...)
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	Printf(os.Stderr, loglevel_fatal, time.Now(), format, v...)
	os.Exit(1)
}

func Print(out *os.File, l loglevel, t time.Time, v ...interface{}) {
	fmt.Fprintf(out, getPrefix(l, t)+fmt.Sprintln(v...))
}

func Printf(out *os.File, l loglevel, t time.Time, format string, v ...interface{}) {
	fmt.Fprintf(out, getPrefix(l, t)+format+"\n", v...)
}
