package glog

import (
    "os"
)

type TLogger interface {
    V(level Level) Verbose
    Info(args ...interface{})
    Infoln(args ...interface{})
    Infof(format string, args ...interface{})
    Warning(args ...interface{})
    Warningln(args ...interface{})
    Warningf(format string, args ...interface{})
    Error(args ...interface{})
    Errorln(args ...interface{})
    Errorf(format string, args ...interface{})
    Fatal(args ...interface{})
    Fatalln(args ...interface{})
    Fatalf(format string, args ...interface{})
    Exit(args ...interface{})
    Exitln(args ...interface{})
    Exitf(format string, args ...interface{})
    Flush()
}

var logger TLogger

func SetLogger(l TLogger) {
    logger = l
}

// V reports whether verbosity level l is at least the requested verbose level.
func V(l int) Verbose {
    return logger.V(l)
}

 // Flush flushes all pending log I/O.
func  Flush() {
    logger.Flush()
    return
}

// Info logs to the INFO log.
func Info(args ...interface{}) {
    logger.Info(args...)
}

// Infof logs to the INFO log. Arguments are handled in the manner of fmt.Printf.
func Infof(format string, args ...interface{}) {
    logger.Infof(format, args...)
}

// Infoln logs to the INFO log. Arguments are handled in the manner of fmt.Println.
func Infoln(args ...interface{}) {
    logger.Infoln(args...)
}

// Warning logs to the WARNING log.
func Warning(args ...interface{}) {
    logger.Warning(args...)
}

// Warningf logs to the WARNING log. Arguments are handled in the manner of fmt.Printf.
func Warningf(format string, args ...interface{}) {
    logger.Warningf(format, args...)
}

// Warningln logs to the WARNING log. Arguments are handled in the manner of fmt.Println.
func Warningln(args ...interface{}) {
    logger.Warningln(args...)
}

// Error logs to the ERROR log.
func Error(args ...interface{}) {
    logger.Error(args...)
}

// Errorf logs to the ERROR log. Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, args ...interface{}) {
    logger.Errorf(format, args...)
}

// Errorln logs to the ERROR log. Arguments are handled in the manner of fmt.Println.
func Errorln(args ...interface{}) {
    logger.Errorln(args...)
}

// Fatal logs to the FATAL log. Arguments are handled in the manner of fmt.Print.
// It calls os.Exit() with exit code 1.
func Fatal(args ...interface{}) {
    logger.Fatal(args...)
    // Make sure fatal logs will exit.
    os.Exit(1)
}

// Fatalf logs to the FATAL log. Arguments are handled in the manner of fmt.Printf.
// It calls os.Exit() with exit code 1.
func Fatalf(format string, args ...interface{}) {
    logger.Fatalf(format, args...)
    // Make sure fatal logs will exit.
    os.Exit(1)
}

// Fatalln logs to the FATAL log. Arguments are handled in the manner of fmt.Println.
// It calle os.Exit()) with exit code 1.
func Fatalln(args ...interface{}) {
    logger.Fatalln(args...)
    // Make sure fatal logs will exit.
    os.Exit(1)
}

// Print prints to the logger. Arguments are handled in the manner of fmt.Print.
//
// Deprecated: use Info.
func Print(args ...interface{}) {
    logger.Info(args...)
}

// Printf prints to the logger. Arguments are handled in the manner of fmt.Printf.
//
// Deprecated: use Infof.
func Printf(format string, args ...interface{}) {
    logger.Infof(format, args...)
}

// Println prints to the logger. Arguments are handled in the manner of fmt.Println.
//
// Deprecated: use Infoln.
func Println(args ...interface{}) {
    logger.Infoln(args...)
}

// Exit logs to the FATAL log. Arguments are handled in the manner of fmt.Print.
// It calls os.Exit() with exit code 1.
func Exit(args ...interface{}) {
    logger.Exit(args...)
    // Make sure fatal logs will exit.
    os.Exit(1)
}

// Exitf logs to the FATAL log. Arguments are handled in the manner of fmt.Printf.
// It calls os.Exit() with exit code 1.
func Exitf(format string, args ...interface{}) {
    logger.Exitf(format, args...)
    // Make sure fatal logs will exit.
    os.Exit(1)
}

// Exitln logs to the FATAL log. Arguments are handled in the manner of fmt.Println.
// It calle os.Exit()) with exit code 1.
func Exitln(args ...interface{}) {
    logger.Exitln(args...)
    // Make sure fatal logs will exit.
    os.Exit(1)
}





