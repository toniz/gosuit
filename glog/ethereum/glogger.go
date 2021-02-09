package logger

import (
    "github.com/toniz/gosuit/ethereum/glog"
    l "github.com/toniz/gosuit/glog"
)

type LogV2 struct {
}

func init() {
    l.SetLogger(new(LogV2))
}

func (LogV2 *)Info(args ...interface{}) {
    glog.Info(args ...interface{})
}

func (LogV2 *) Infoln(args ...interface{}){
    glog.Infoln(args ...interface{})
}

func (LogV2 *) Infof(format string, args ...interface{}){
    glog.Infof(format string, args ...interface{})
}

func (LogV2 *) Warning(args ...interface{}){
    glog.Warning(args ...interface{})
}

func (LogV2 *) Warningln(args ...interface{}){
    glog.Warningln(args ...interface{})
}

func (LogV2 *) Warningf(format string, args ...interface{}){
    glog.Warningf(format string, args ...interface{})
}

func (LogV2 *) Error(args ...interface{}){
    glog.Error(args ...interface{})
}

func (LogV2 *) Errorln(args ...interface{}){
    glog.Errorln(args ...interface{})
}

func (LogV2 *) Errorf(format string, args ...interface{}){
    glog.Errorf(format string, args ...interface{})
}

func (LogV2 *) Fatal(args ...interface{}){
    glog.Fatal(args ...interface{})
}

func (LogV2 *) Fatalln(args ...interface{}){
    glog.Fatalln(args ...interface{})
}

func (LogV2 *) Fatalf(format string, args ...interface{}){
    glog.Fatalf(format string, args ...interface{})
}

func (LogV2 *) V(l int) glog.Verbose {
    return glog.V(l int)
}



