package mlog

import (
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
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
	level     LogLevel
	filePath  string
	perFix    string
	trailType string
	trailTag  time.Time
	trailCap  int64
	fileObj   *os.File
	logChan   chan *logMsg //
}

// logMsg 结构体 //
type logMsg struct {
	level    LogLevel
	msg      string
	funcName string
	fileName string
	timeStmp string
	line     int
}

// Newlog 构造函数
//	lv 级别：DEBUG INFO WARN ERROR FATAL；
// 	fp 日志文件路径；
// 	pf 日志文件前缀；
// 	ts 日志切割类型及参数："minite=10"、"hour=2"、"day=1"、"size=1024";
//	1、选时间模式时，以每个方式的起始值切割，如：hour,即每小时0分时切割; day,即每日0时切割；
//  2、选size模式时，值以MB为单位
func Newlog(lv, fp, pf, ts string) *Logger {
	level, err := parseLogLevel(lv)
	if err != nil {
		panic(err)
	}
	tt, ms, err := parseTrailAttr(ts)
	if err != nil {
		panic(err)
	}
	if ms <= 0 {
		err = errors.New("cut interval attribute must bigger than 0")
		panic(err)
	}
	fl := &Logger{
		level:     level,
		filePath:  fp,
		perFix:    pf,
		trailType: tt,
		trailCap:  ms,
		logChan:   make(chan *logMsg, 50000),
	}
	//初始化滚动标签
	err = fl.initTag()
	if err != nil {
		panic(err)
	}
	//初始化文件句柄
	err = fl.creatFile()
	if err != nil {
		panic(err)
	}
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

	fileName := l.perFix + ".log"
	fullFileName := path.Join(l.filePath, fileName)
	fobj, err := os.OpenFile(fullFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open log file faild. err:%v\n", err)
		return err
	}
	l.fileObj = fobj
	//开启1个后台goroutine落盘
	go l.writeBackGround()

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

// 格式化日志切割参数
func parseTrailAttr(str string) (tt string, tc int64, err error) {
	tmpStr := strings.Split(str, "=")
	tt = tmpStr[0]
	tc, err = strconv.ParseInt(tmpStr[1], 10, 64)
	return
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
		err := errors.New("unknown log level")
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

// 后台从通道取出落盘
func (l *Logger) writeBackGround() {
	for {
		//日志文件切割
		l.checkTrail()
		select {
		case logTmp := <-l.logChan:
			logInfo := fmt.Sprintf("[%s] [%s] [%s:%s:%d] %s\n", logTmp.timeStmp, unparseLogLevel(logTmp.level), logTmp.fileName, logTmp.funcName, logTmp.line, logTmp.msg)
			fmt.Printf(logInfo)
			fmt.Fprintf(l.fileObj, logInfo)
		default:
			time.Sleep(time.Millisecond * 10)
		}
	}
}

// 按格式写日志，切割日志文件
func (l *Logger) writeLog(lv LogLevel, format string, args ...interface{}) {
	if l.enable(lv) {
		msg := fmt.Sprintf(format, args...) //拼接接口类变量
		now := time.Now()
		funcName, fileName, line := getRunInfo(3)
		//构造日志信息结构体
		logTmp := &logMsg{
			level:    lv,
			msg:      msg,
			funcName: funcName,
			fileName: fileName,
			timeStmp: now.Format("2006-01-02 15:04:05.000"),
			line:     line,
		}
		//发送到通道
		select {
		case l.logChan <- logTmp:
		default:
		}
	}
}

// 取得运行信息
func getRunInfo(skip int) (funcName, fileName string, line int) {
	pc, file, line, ok := runtime.Caller(skip) //runtime.Caller()
	if !ok {
		fmt.Printf("runtime.Caller() faild.\n")
		return
	}
	funcName = strings.Split(runtime.FuncForPC(pc).Name(), ".")[1] //返回主调用函数名
	fileName = path.Base(file)                                     //返回主调用文件名
	return
}

// CloseFile 关闭文件句柄，内部不使用
func (l *Logger) closeFile(file *os.File) {
	file.Close()
}
