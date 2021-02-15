package glogger

import (
    "fmt"
    "github.com/toniz/gosuit/glog/ethereum/glog"
    l "github.com/toniz/gosuit/glog"
)

type glogger struct{}

func init() {
    l.SetLogger(&glogger{})
}

func (g *glogger) Info(args ...interface{}) {
    glog.InfoDepth(2, args...)
}

func (g *glogger) Infoln(args ...interface{}) {
    glog.InfoDepth(2, fmt.Sprintln(args...))
}

func (g *glogger) Infof(format string, args ...interface{}) {
    glog.InfoDepth(2, fmt.Sprintf(format, args...))
}

func (g *glogger) InfoDepth(depth int, args ...interface{}) {
    glog.InfoDepth(depth+2, args...)
}

func (g *glogger) Warning(args ...interface{}) {
    glog.WarningDepth(2, args...)
}

func (g *glogger) Warningln(args ...interface{}) {
    glog.WarningDepth(2, fmt.Sprintln(args...))
}

func (g *glogger) Warningf(format string, args ...interface{}) {
    glog.WarningDepth(2, fmt.Sprintf(format, args...))
}

func (g *glogger) WarningDepth(depth int, args ...interface{}) {
    glog.WarningDepth(depth+2, args...)
}

func (g *glogger) Error(args ...interface{}) {
    glog.ErrorDepth(2, args...)
}

func (g *glogger) Errorln(args ...interface{}) {
    glog.ErrorDepth(2, fmt.Sprintln(args...))
}

func (g *glogger) Errorf(format string, args ...interface{}) {
    glog.ErrorDepth(2, fmt.Sprintf(format, args...))
}

func (g *glogger) ErrorDepth(depth int, args ...interface{}) {
    glog.ErrorDepth(depth+2, args...)
}

func (g *glogger) Fatal(args ...interface{}) {
    glog.FatalDepth(2, args...)
}

func (g *glogger) Fatalln(args ...interface{}) {
    glog.FatalDepth(2, fmt.Sprintln(args...))
}

func (g *glogger) Fatalf(format string, args ...interface{}) {
    glog.FatalDepth(2, fmt.Sprintf(format, args...))
}

func (g *glogger) FatalDepth(depth int, args ...interface{}) {
    glog.FatalDepth(depth+2, args...)
}

func (g *glogger) Exit(args ...interface{}) {
    glog.ExitDepth(2, args...)
}

func (g *glogger) Exitln(args ...interface{}) {
    glog.ExitDepth(2, fmt.Sprintln(args...))
}

func (g *glogger) Exitf(format string, args ...interface{}) {
    glog.ExitDepth(2, fmt.Sprintf(format, args...))
}

func (g *glogger) ExitDepth(depth int, args ...interface{}) {
    glog.ExitDepth(depth+2, args...)
}

func (*LogV2) V(level int) l.Verboser {
    return glog.V(glog.Level(level))
}

func (g *glogger) Flush() {
    glog.Flush()
}

func (g *glogger) CopyStandardLogTo(name string) {
    glog.CopyStandardLogTo(name)
}


