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