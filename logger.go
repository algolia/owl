package owl

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	raven "github.com/getsentry/raven-go"
	"github.com/mgutz/ansi"
)

var (
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger

	logFile *os.File
)

func initLogger() error {
	if config.Logger == nil {
		useSentry = false
		return nil
	}

	ansi.DisableColors(!config.Logger.UseColors)

	if config.Logger.DisplayLogs {
		formatFlags := log.Ldate | log.Lmicroseconds
		infoLogger = log.New(os.Stdout, ansi.Color("[INF] ", "blue"), formatFlags)
		warningLogger = log.New(os.Stdout, ansi.Color("[WAR] ", "yellow"), formatFlags)
		errorLogger = log.New(os.Stderr, ansi.Color("[ERR] ", "red"), formatFlags)
	}

	if config.Logger.LogFilePath != "" {
		var err error
		logFile, err = os.OpenFile(config.Logger.LogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
		if err != nil {
			logFile = nil
			return Error("owl: cannot open or create log file: %s", err)
		}
	}

	if useSentry {
		if config.Logger.SentryDsn == "" {
			useSentry = false
			return Error("owl: Cannot use Sentry as the SentryDsn configuration field is empty")
		}
		raven.SetDSN(config.Logger.SentryDsn)
	}

	return nil
}

func stopLogger() {
	if logFile != nil {
		logFile.Close()
		logFile = nil
	}

	if infoLogger != nil {
		infoLogger = nil
		warningLogger = nil
		errorLogger = nil
	}
}

func Info(format string, args ...interface{}) {
	logEverywhere(infoLogger, "inf", format, args...)
}

func Warning(format string, args ...interface{}) {
	logEverywhere(warningLogger, "war", format, args...)
}

func Error(format string, args ...interface{}) error {
	message := logEverywhere(errorLogger, "err", format, args...)
	err := errors.New(message)
	if useSentry {
		raven.CaptureError(err, nil)
	}
	return err
}

func logEverywhere(l *log.Logger, lvl, format string, args ...interface{}) (msg string) {
	msg = fmt.Sprintf(format, args...)

	writeToLogFile(lvl, msg)
	writeToLogger(l, msg)

	return
}

func writeToLogger(l *log.Logger, msg string) {
	if l == nil {
		return
	}

	l.Println(msg)
}

func writeToLogFile(lvl, msg string) {
	if logFile == nil {
		return
	}

	jsonMsg := fmt.Sprintf("{\"ts\":\"%s\",\"lvl\":\"%s\",\"msg\":\"%s\",\"app\":\"%s\",\"logger\":\"%s\"}\n",
		time.Now().Format("2006-01-02T15:04:05.000"),
		lvl,
		msg,
		config.AppName,
		config.Logger.Logger,
	)

	if _, err := logFile.WriteString(jsonMsg); err != nil {
		writeToLogger(errorLogger, fmt.Sprintf("owl: cannot write log file to disk: %s", err))
	}
}
