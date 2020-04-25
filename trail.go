package mlog

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"
)

//  初始化滚动标签
func (l *Logger) initTag() error {
	switch l.trailType {
	case "minite":
		m, _ := time.ParseDuration("1m")
		t, err := time.ParseInLocation("2006-01-02 15:04", time.Now().Format("2006-01-02 15:04"), time.Local)
		if err != nil {
			fmt.Printf("err:%v\n", err)
		}
		l.trailTag = t.Add(time.Duration(l.trailCap) * m)
	case "hour":
		h, _ := time.ParseDuration("1h")
		t, err := time.ParseInLocation("2006-01-02 15", time.Now().Format("2006-01-02 15"), time.Local)
		if err != nil {
			fmt.Printf("err:%v\n", err)
		}
		l.trailTag = t.Add(time.Duration(l.trailCap) * h)
	case "day":
		d, _ := time.ParseDuration("24h")
		t, err := time.ParseInLocation("2006-01-02", time.Now().Format("2006-01-02"), time.Local)
		if err != nil {
			fmt.Printf("err:%v\n", err)
		}
		l.trailTag = t.Add(time.Duration(l.trailCap) * d)
	case "size":
	default:
		return errors.New("cut type attribute is wrong")
	}
	return nil
}

//	检查滚动标签
// 	标签类型："minite"、"hour"、"day"、"size"
//	符合切割条件的，备份文件，重置该级别新日志文件句柄
func (l *Logger) checkTrail() {
	switch l.trailType {
	case "minite":
		m, _ := time.ParseDuration("1m")
		if time.Now().After(l.trailTag) { //当前时间是否大于标签时间
			l.trailFile()
			l.trailTag = time.Now().Add(time.Duration(l.trailCap) * m) //标签时间加1个单位
		}
	case "hour":
		h, _ := time.ParseDuration("1h")
		if time.Now().After(l.trailTag) {
			l.trailFile()
			l.trailTag = time.Now().Add(time.Duration(l.trailCap) * h)
		}
	case "day":
		d, _ := time.ParseDuration("24h")
		if time.Now().After(l.trailTag) {
			l.trailFile()
			l.trailTag = time.Now().Add(time.Duration(l.trailCap) * d)
		}
	case "size":
		fileInfo, err := l.fileObj.Stat()
		if err != nil {
			fmt.Printf("get file info faild.err:%v\n", err)
			return
		}
		if fileInfo.Size() >= l.trailCap*1024*1024 { //以Byte为单位
			l.trailFile()
		}
	}
}

// 文件滚动
func (l *Logger) trailFile() {
	//备份及打开文件
	nowStr := time.Now().Format("20060102150405")
	fileInfo, err := l.fileObj.Stat()
	if err != nil {
		fmt.Printf("get file info faild.err:%v\n", err)
		return
	}
	oldName := path.Join(l.filePath, fileInfo.Name())
	newName := fmt.Sprintf("%s%s", oldName, nowStr)
	l.fileObj.Close()
	os.Rename(oldName, newName)
	fobj, err := os.OpenFile(oldName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open new file faild.err:%v\n", err)
		return
	}
	l.fileObj = fobj
}
