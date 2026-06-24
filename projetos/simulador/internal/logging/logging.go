package logging

import (
    "fmt"
    "os"
)

type Level int

const (
    ErrorLevel Level = iota
    WarnLevel
    InfoLevel
    DebugLevel
)

var currentLevel = InfoLevel

func Init(verbose, quiet bool) {
    switch {
    case quiet:
        currentLevel = ErrorLevel
    case verbose:
        currentLevel = DebugLevel
    default:
        currentLevel = InfoLevel
    }
}

func Debugf(format string, args ...interface{}) {
    if currentLevel >= DebugLevel {
        fmt.Fprintf(os.Stderr, "DEBUG: "+format+"\n", args...)
    }
}

func Infof(format string, args ...interface{}) {
    if currentLevel >= InfoLevel {
        fmt.Fprintf(os.Stderr, "INFO: "+format+"\n", args...)
    }
}

func Warnf(format string, args ...interface{}) {
    if currentLevel >= WarnLevel {
        fmt.Fprintf(os.Stderr, "WARN: "+format+"\n", args...)
    }
}

func Errorf(format string, args ...interface{}) {
    fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", args...)
}
