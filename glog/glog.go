package glog // import "google.golang.org/grpc/internal/logger"

import (
    "os"
)

// LoggerV2 does underlying logging work for logger.
// This is a copy of the LoggerV2 defined in the external logger package. It
// is defined here to avoid a circular dependency.
type TLogger interface {
    // Info logs to INFO log. Arguments are handled in the manner of fmt.Print.
    Info(args ...interface{})
    // Infoln logs to INFO log. Arguments are handled in the manner of fmt.Println.
    Infoln(args ...interface{})
    // Infof logs to INFO log. Arguments are handled in the manner of fmt.Printf.
    Infof(format string, args ...interface{})
    // Warning logs to WARNING log. Arguments are handled in the manner of fmt.Print.
    Warning(args ...interface{})
    // Warningln logs to WARNING log. Arguments are handled in the manner of fmt.Println.
    Warningln(args ...interface{})
    // Warningf logs to WARNING log. Arguments are handled in the manner of fmt.Printf.
    Warningf(format string, args ...interface{})
    // Error logs to ERROR log. Arguments are handled in the manner of fmt.Print.
    Error(args ...interface{})
    // Errorln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
    Errorln(args ...interface{})
    // Errorf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
    Errorf(format string, args ...interface{})
    // Fatal logs to ERROR log. Arguments are handled in the manner of fmt.Print.
    // gRPC ensures that all Fatal logs will exit with os.Exit(1).
    // Implementations may also call os.Exit() with a non-zero exit code.
    Fatal(args ...interface{})
    // Fatalln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
    // gRPC ensures that all Fatal logs will exit with os.Exit(1).
    // Implementations may also call os.Exit() with a non-zero exit code.
    Fatalln(args ...interface{})
    // Fatalf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
    // gRPC ensures that all Fatal logs will exit with os.Exit(1).
    // Implementations may also call os.Exit() with a non-zero exit code.
    Fatalf(format string, args ...interface{})
    // V reports whether verbosity level l is at least the requested verbose level.
    V(l int) bool
}

var logger TLogger

func SetLogger(l TLogger) {
    logger = l
}


// V reports whether verbosity level l is at least the requested verbose level.
func V(l int) bool {
    return logger.V(l)
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
