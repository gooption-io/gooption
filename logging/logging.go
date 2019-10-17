// Package logging handles logging to StdOut and Writer vscan.log
package logging

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
)

// Prettify StdOut with colors
var (
	// private color properties
	warningMessage = color.New(color.FgHiYellow).SprintFunc()
	errorMessage   = color.New(color.FgHiRed).SprintFunc()
	fatalMessage   = color.New(color.BgRed, color.FgHiWhite).SprintFunc()

	// public color properties
	InfoMessage   = color.New(color.FgHiGreen).SprintFunc()
	UnderlineText = color.New(color.Underline).SprintFunc()
)

func logToStdOut(level string, fields ...interface{}) {

	var log = logrus.New()

	log.SetReportCaller(true)

	localTimezone, _ := time.Now().In(time.Local).Zone()

	formatter := &prefixed.TextFormatter{
		FullTimestamp:   true,
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05 " + localTimezone,
		ForceFormatting: true,
		SpacePadding:    0,
	}
	formatter.SetColorScheme(&prefixed.ColorScheme{
		TimestampStyle:  "white+u",
		InfoLevelStyle:  "white:28",
		WarnLevelStyle:  "white:208",
		ErrorLevelStyle: "white:red",
		FatalLevelStyle: "white:red",
	})

	log.Formatter = formatter
	switch level {
	case "warning":
		log.Warningln(warningMessage(fields...))
	case "info":
		log.Infoln(InfoMessage(fields...))
	case "error":
		log.Errorln(errorMessage(fields...))
	case "fatal":
		log.Fatalln(fatalMessage(fields...))
	default:
		log.Errorln(errorMessage(fields...))
	}

}
func Log(level string, fields string, args ...interface{}) {

	// Use String Builder for more efficient strings concat
	s := strings.Builder{}

	_, f, l, ok := runtime.Caller(1)

	if ok {
		_, file := filepath.Split(f)

		line := l

		// Don't write runtime caller for HTTP requests log
		if file != "requestslog.go" {

			_, _ = fmt.Fprintf(&s, "caller=%v line=%d msg=", file, line)
		}

	}

	_, err := fmt.Fprintf(&s, fields, args...)

	if err != nil {
		logToStdOut("error", fmt.Sprintf("Failed to write logging string in strings.Builder buffer %v", err))
	}

	logToStdOut(level, s.String())
}
