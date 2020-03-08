package mlog

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"time"
)

// 	日志滚动类型初始化，
// 	分类计算滚动类型的切割label，
// 	以map记录当前级别以上的切割label，
// 	滚动类型："minite"、"hour"、"day"、"week"、"month"、"year"、"size"、"default"。
func (l *Logger) logTrail() {
	l.timeLabel = make(map[LogLevel]int64)
	var x = [5]LogLevel{DEBUG, INFO, WARN, ERROR, FATAL}
	for _, v := range x {
		if v >= l.Level {
			switch l.trailType {
			case "minite":
				timeStr := time.Now().Format("2006-01-02 15:04:")
				t, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr+"59", time.Local)
				label := t.Unix() + 1
				l.timeLabel[v] = label
			case "hour":
				timeStr := time.Now().Format("2006-01-02 15:")
				t, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr+"59:59", time.Local)
				label := t.Unix() + 1
				l.timeLabel[v] = label
			case "day":
				timeStr := time.Now().Format("2006-01-02")
				t, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr+" 23:59:59", time.Local)
				label := t.Unix() + 1
				l.timeLabel[v] = label
			case "week":
				w := int(time.Now().Weekday())
				timeStr := time.Now().Format("2006-01-02")
				t, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr+" 23:59:59", time.Local)
				label := t.Unix() + int64(86400*(6-w)+1)
				l.timeLabel[v] = label
			case "month":
				m := int(time.Now().Month())
				switch m {
				case 1, 2, 3, 4, 5, 6, 7, 8:
					yres := strconv.Itoa(time.Now().Year())
					mres := "0" + strconv.Itoa(m+1)
					t, _ := time.ParseInLocation("2006-01-02 15:04:05", yres+"-"+mres+"-01 00:00:00", time.Local)
					label := t.Unix()
					l.timeLabel[v] = label
				case 9, 10, 11:
					yres := strconv.Itoa(time.Now().Year())
					mres := strconv.Itoa(m + 1)
					t, _ := time.ParseInLocation("2006-01-02 15:04:05", yres+"-"+mres+"-01 00:00:00", time.Local)
					label := t.Unix()
					l.timeLabel[v] = label
				default: //12
					yres := strconv.Itoa(time.Now().Year() + 1)
					mres := "01"
					t, _ := time.ParseInLocation("2006-01-02 15:04:05", yres+"-"+mres+"-01 00:00:00", time.Local)
					label := t.Unix()
					l.timeLabel[v] = label
				}
			case "year":
				res := strconv.Itoa(time.Now().Year() + 1)
				t, _ := time.ParseInLocation("2006-01-02 15:04:05", res+"-01-01 00:00:00", time.Local)
				label := t.Unix()
				l.timeLabel[v] = label
			case "size":
				l.timeLabel[v] = l.maxFileSize
			case "default":
				l.timeLabel[v] = 0
			default:
				panic("日志切割参数不正确.")
			}
		}
	}
}

//	检查滚动标签，
//	分类检查当前级别的label，
//	符合切割条件的重置该级别的切割label，备份文件，重置该级别新日志文件句柄。
func (l *Logger) checkLabel(file *os.File, lv LogLevel) {
	now := time.Now().Unix()
	if now >= l.timeLabel[lv] {
		//修改下次检查label
		switch l.trailType {
		case "minite":
			l.timeLabel[lv] += 60
		case "hour":
			l.timeLabel[lv] += 3600
		case "day":
			l.timeLabel[lv] += 86400
		case "week":
			l.timeLabel[lv] += 604800
		case "month":
			//tl := time.Unix(l.timeLabel[0], 0).Format("2006-01-02 15:04:05") //int64 -> string
			tl := time.Unix(l.timeLabel[0], 0) //int64 ->time.time
			m := int(tl.Month())
			switch m {
			case 1, 2, 3, 4, 5, 6, 7, 8:
				yres := strconv.Itoa(time.Now().Year())
				mres := "0" + strconv.Itoa(m+1)
				t, _ := time.ParseInLocation("2006-01-02 15:04:05", yres+"-"+mres+"-01 00:00:00", time.Local)
				label := t.Unix()
				l.timeLabel[lv] = label
			case 9, 10, 11:
				yres := strconv.Itoa(time.Now().Year())
				mres := strconv.Itoa(m + 1)
				t, _ := time.ParseInLocation("2006-01-02 15:04:05", yres+"-"+mres+"-01 00:00:00", time.Local)
				label := t.Unix()
				l.timeLabel[lv] = label
			default: //12
				yres := strconv.Itoa(time.Now().Year() + 1)
				mres := "01"
				t, _ := time.ParseInLocation("2006-01-02 15:04:05", yres+"-"+mres+"-01 00:00:00", time.Local)
				label := t.Unix()
				l.timeLabel[lv] = label
			}
		case "year":
			l.timeLabel[lv] += 31536000
		case "size":
			l.timeLabel[lv] += 0
		}
		//备份及打开文件
		nowStr := time.Now().Format("20060102150405")
		fileInfo, err := file.Stat()
		if err != nil {
			fmt.Printf("get file info faild.err:%v\n", err)
			return
		}
		oldName := path.Join(l.filePath, fileInfo.Name())
		newName := fmt.Sprintf("%s%s", oldName, nowStr)
		//1,关闭文件
		file.Close()
		//2,备份原文件
		os.Rename(oldName, newName)
		//3,重新打开文件
		fileObj, err := os.OpenFile(oldName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("open new file faild.err:%v\n", err)
			return
		}
		l.fileObjs[lv] = fileObj
	}
}

//	检查滚动标签，
//	分类检查当前级别的label，
//	符合切割条件的，备份文件，重置该级别新日志文件句柄。
func (l *Logger) checkSize(file *os.File, lv LogLevel) {
	nowStr := time.Now().Format("20060102150405")
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("get file info faild.err:%v\n", err)
		return
	}
	oldName := path.Join(l.filePath, fileInfo.Name())
	newName := fmt.Sprintf("%s%s", oldName, nowStr)
	if fileInfo.Size() >= l.timeLabel[lv] {
		//1,关闭文件
		file.Close()
		//2,备份原文件
		os.Rename(oldName, newName)
		//3,重新打开文件
		fileObj, err := os.OpenFile(oldName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("open new file faild.err:%v\n", err)
			return
		}
		l.fileObjs[lv] = fileObj
	}
}
