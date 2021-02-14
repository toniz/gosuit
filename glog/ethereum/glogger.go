package glogger

import (
    "github.com/toniz/gosuit/glog/ethereum/glog"
    l "github.com/toniz/gosuit/glog"
)

type LogV2 struct {
}

func init() {
    lv2 := &LogV2{}
    l.SetLogger(lv2)
}

func (*LogV2)Info(args ...interface{}) {
    glog.Info(args)
}

func (*LogV2) Infoln(args ...interface{}){
    glog.Infoln(args)
}

func (*LogV2) Infof(format string, args ...interface{}){
    glog.Infof(format, args)
}

func (*LogV2) Warning(args ...interface{}){
    glog.Warning(args)
}

func (*LogV2) Warningln(args ...interface{}){
    glog.Warningln(args)
}

func (*LogV2) Warningf(format string, args ...interface{}){
    glog.Warningf(format, args)
}

func (*LogV2) Error(args ...interface{}){
    glog.Error(args)
}

func (*LogV2) Errorln(args ...interface{}){
    glog.Errorln(args)
}

func (*LogV2) Errorf(format string, args ...interface{}){
    glog.Errorf(format, args)
}

func (*LogV2) Fatal(args ...interface{}){
    glog.Fatal(args)
}

func (*LogV2) Fatalln(args ...interface{}){
    glog.Fatalln(args)
}

func (*LogV2) Fatalf(format string, args ...interface{}){
    glog.Fatalf(format, args)
}

func (*LogV2) Exit(args ...interface{}){
    glog.Exit(args)
}

func (*LogV2) Exitln(args ...interface{}){
    glog.Exitln(args)
}

func (*LogV2) Exitf(format string, args ...interface{}){
    glog.Exitf(format, args)
}

func (*LogV2) V(level int) l.Verboser {
    return glog.V(glog.Level(level))
}

func (*LogV2) Flush() {
    glog.Flush()
}


