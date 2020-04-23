package mlog

import (
	"fmt"
	"os"
	"path"
	"time"
)

//	检查滚动标签
// 	标签类型："hour"、"day"、"month"、"year"、"size"
//	符合切割条件的，备份文件，重置该级别新日志文件句柄
func (l *Logger) checkTrail() {
	if !l.istail {
		for k, v := range l.fileObjs {
			switch l.trailType {
			case "hour":
				if time.Now().Format("04") == "00" { //0分切换
					//l.trailFile(file, lv)
					l.trailFile(v, k)
				}
			case "day":
				if time.Now().Format("15:04") == "00:00" { //0时0分切换
					l.trailFile(v, k)
				}
			case "month":
				if time.Now().Format("02 15:04") == "01 00:00" { //1日0时0分切换
					l.trailFile(v, k)
				}
			case "year":
				if time.Now().Format("01-02 15:04") == "01-01 00:00" { //1月1日0时0分切换
					l.trailFile(v, k)
				}
			case "size":
				fileInfo, err := v.Stat()
				if err != nil {
					fmt.Printf("get file info faild.err:%v\n", err)
					return
				}
				if fileInfo.Size() >= l.maxFileSize {
					l.trailFile(v, k)
				}
			}
		}
		l.istail = true
		time.Sleep(time.Second * 30)
	} else {
		l.istail = false
		time.Sleep(time.Second * 31)
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
