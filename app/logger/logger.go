package logger

import (
    "log"
    "os"
    "strings"
    "time"
)

var CustomLogger *log.Logger

func init() {
    CustomLogger = log.New(os.Stdout, "", 0)
}

type LogLevel int
const (
    DebugLevel LogLevel = 3
    InfoLevel  LogLevel = 2
    WarnLevel  LogLevel = 1
    ErrorLevel LogLevel = 0
)

var logLevel = InfoLevel

func Log(severity string, message string, a ...any) {
    if len(a) == 0 {
        CustomLogger.Printf("[%s] %s %s\n", strings.ToUpper(severity), time.Now().Format("2006-01-02T15:04:05"), message)
        return
    } else {
        CustomLogger.Printf("[%s] %s %s %s\n", strings.ToUpper(severity), time.Now().Format("2006-01-02T15:04:05"), message, a)
    }
}

func Info(message string, a ...any) {
    if logLevel >= InfoLevel {
        Log("info", message, a...)
    }
}

func Error(message string, a ...any) {
    Log("error", message, a...)
}

func Debug(message string, a ...any) {
    if logLevel >= DebugLevel {
        Log("debug", message, a...)
    }
}

func SetLevel(level string) {
    switch strings.ToLower(level) {
    case "debug":
        logLevel = DebugLevel
    case "info":
        logLevel = InfoLevel
    case "warn":
        logLevel = WarnLevel
    case "error":
        logLevel = ErrorLevel
    }
}
