go log
Cut log by configuration

1、可日志级别：DEBUG INFO WARN ERROR FATAL；
2、大于设置的级别的日志，单独落地文件；
3、可添加落地文件前缀；
4、可设置日志切割类型："default"、"minite"、"hour"、"day"、"month"、"year"、"size";
    4.1、选时间模式时，以每个方式的起始值切割，如：hour,即每小时0分0秒时切割; month,即每月1日0时切割；
    4.2、选size时，按文件大小切割，参数值为ms，单位byte;不选size时，ms值不生效，可写0；
