package mlog

import (
	"fmt"
	"os"
	"path"
	"time"
)

//	检查滚动标签
// 	标签类型："minite"、"hour"、"day"、"month"、"year"、"size"
//	符合切割条件的，备份文件，重置该级别新日志文件句柄
func (l *Logger) checkTrail(file *os.File, lv LogLevel) {
	switch l.trailType {
	case "default":
	case "minite":
		if time.Now().Format("05") == "00" {
			l.trailFile(file, lv)
		}
	case "hour":
		if time.Now().Format("04:05") == "00:00" {
			l.trailFile(file, lv)
		}
	case "day":
		if time.Now().Format("15:04:05") == "00:00:00" {
			l.trailFile(file, lv)
		}
	case "month":
		if time.Now().Format("02 15:04:05") == "01 00:00:00" {
			l.trailFile(file, lv)
		}
	case "year":
		if time.Now().Format("01-02 15:04:05") == "01-01 00:00:00" {
			l.trailFile(file, lv)
		}
	case "size":
		fileInfo, err := file.Stat()
		if err != nil {
			fmt.Printf("get file info faild.err:%v\n", err)
			return
		}
		if fileInfo.Size() >= l.maxFileSize {
			l.trailFile(file, lv)
		}
	}
}

// 文件滚动
func (l *Logger) trailFile(file *os.File, lv LogLevel) {
	//备份及打开文件
	nowStr := time.Now().Format("20060102150405")
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("get file info faild.err:%v\n", err)
		return
	}
	oldName := path.Join(l.filePath, fileInfo.Name())
	newName := fmt.Sprintf("%s%s", oldName, nowStr)
	file.Close()
	os.Rename(oldName, newName)
	fileObj, err := os.OpenFile(oldName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open new file faild.err:%v\n", err)
		return
	}
	l.fileObjs[lv] = fileObj
}
