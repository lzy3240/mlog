package mlog

import (
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

// LogLevel 级别类型
type LogLevel uint16

// 日志级别常量
const (
	DEFAULT LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

// Lever 结构体
type Lever struct {
	Level LogLevel
}

// Logger 结构体
type Logger struct {
	level       LogLevel
	filePath    string
	perFix      string
	trailType   string
	maxFileSize int64
	istail      bool
	fileObjs    map[LogLevel]*os.File
}

// Newlog 构造函数
//	lv 级别：DEBUG INFO WARN ERROR FATAL；
// 	fp 日志文件路径；
// 	pf 日志文件前缀；
// 	tt 日志切割类型："hour"、"day"、"month"、"year"、"size";
//	1、选时间模式时，以每个方式的起始值切割，如：hour,即每小时0分时切割; month,即每月1日0时切割；
//	2、选size时，按文件大小切割，参数值为ms，单位byte;不选size时，ms值不生效，可写0
func Newlog(lv, fp, pf, tt string, ms int64) *Logger {
	level, err := parseLogLevel(lv)
	if err != nil {
		panic(err)
	}
	fl := &Logger{
		level:       level,
		filePath:    fp,
		perFix:      pf,
		trailType:   tt,
		maxFileSize: ms,
		istail:      false,
	}
	err = fl.creatFile() //创建文件
	if err != nil {
		panic(err)
	}
	go fl.checkTrail()
	return fl
}

// 初始化日志文件
func (l *Logger) creatFile() error {
	// 判断filePath是否存在及创建
	exists, err := isExists(l.filePath)
	if err != nil {
		fmt.Printf("get filepath err. err:%v\n", err)
	}
	if !exists {
		err := os.Mkdir(l.filePath, os.ModePerm)
		if err != nil {
			fmt.Printf("create file path faild. err:%v\n", err)
		}
	}
	//判断日志级别，生成高于该级别的日志句柄
	var tmp = [5]LogLevel{DEBUG, INFO, WARN, ERROR, FATAL}
	l.fileObjs = make(map[LogLevel]*os.File, 5)
	for _, v := range tmp {
		if v >= l.level {
			fileName := l.perFix + unparseLogLevel(v) + ".log"
			fullFileName := path.Join(l.filePath, fileName)
			fileObj, err := os.OpenFile(fullFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Printf("open log file faild. err:%v\n", err)
				return err
			}
			l.fileObjs[v] = fileObj
		}
	}

	//fmt.Println(l.fileObjs)
	return nil
}

// 检查目录是否存在
func isExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 格式化日志级别
func parseLogLevel(s string) (LogLevel, error) {
	s = strings.ToLower(s)
	switch s {
	case "debug":
		return DEBUG, nil
	case "info":
		return INFO, nil
	case "warn":
		return WARN, nil
	case "error":
		return ERROR, nil
	case "fatal":
		return FATAL, nil
	default:
		err := errors.New("未知的日志级别")
		return DEBUG, err
	}
}

// 反格式化日志级别
func unparseLogLevel(lv LogLevel) string {
	switch lv {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	}
	return "DEBUG"
}

// 按格式写日志，切割日志文件
func (l *Logger) writeLog(lv LogLevel, format string, args ...interface{}) {
	if l.enable(lv) {
		msg := fmt.Sprintf(format, args...) //拼接接口类变量
		now := time.Now()
		funcName, fileName, line := getRunInfo(3)
		fmt.Printf("[%s] [%s] [%s:%s:%d] %s\n", now.Format("2006-01-02 15:04:05.000"), unparseLogLevel(lv), fileName, funcName, line, msg)
		//日志文件切割
		//l.checkTrail(l.fileObjs[lv], lv)
		fmt.Fprintf(l.fileObjs[lv], "[%s] [%s] [%s:%s:%d] %s\n", now.Format("2006-01-02 15:04:05.000"), unparseLogLevel(lv), fileName, funcName, line, msg)
	}
}

// 取得运行信息
func getRunInfo(skip int) (funcName, fileName string, line int) {
	pc, file, line, ok := runtime.Caller(skip) //runtime.Caller()
	if !ok {
		fmt.Printf("runtime.Caller() faild.\n")
		return
	}
	//funcName = runtime.FuncForPC(pc).Name()
	funcName = strings.Split(runtime.FuncForPC(pc).Name(), ".")[1] //返回主调用函数名
	fileName = path.Base(file)                                     //返回主调用文件名
	return
}

// CloseFile 关闭文件句柄，内部不使用
func (l *Logger) CloseFile(file *os.File) {
	file.Close()
}
