go log
Cut log by configuration

//	lv 级别：DEBUG INFO WARN ERROR FATAL；
// 	fp 日志文件路径；
// 	pf 日志文件前缀；
// 	ts 日志切割类型及参数："minite=10"、"hour=2"、"day=1"、"size=1024";

//	1、选时间模式时，以每个方式的起始值切割，如：hour,即每小时0分时切割; day,即每日0时切割；
//  2、选size模式时，值以MB为单位
//  3、异步写入磁盘
//  4、输出到控制台


package main

import (
	"time"

	"github.com/lzy3240/mlog"
)

var log *mlog.Logger

// init log
func init() {
	log = mlog.Newlog("info", "./logs/", "Test", "hour=1")
	// log = mlog.Newlog("error", "./logs/", "Test", "hour=3")
	// log = mlog.Newlog("info", "./logs/", "Test", "day=1")
	// log = mlog.Newlog("info", "./logs/", "Test", "size=100")
}

func main() {
	a := "just a test msg"
	for {
		log.Info("this is a info test,%v", a)
		log.Debug("this is a Debug test,%v", a)
		log.Fatal("this is a Fatal test,%v", a)
		log.Error("this is a Error test,%v", a)
		log.Warn("this is a Warn test,%v", a)
		time.Sleep(time.Second)
	}
}
