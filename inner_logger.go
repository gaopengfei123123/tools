package tools

// 用于打印包内日志, 如果需要调用日志方法, 需要注入指定的日志格式

type Logger func(level int, traceData map[string]string, msg string, v ...interface{})

var innerLogger Logger

func printInnerLog(level int, traceData map[string]string, msg string, v ...interface{}) {
	if innerLogger != nil {
		innerLogger(level, traceData, msg, v...)
		return
	}

	//msg = fmt.Sprintf(msg, v...)
	//logs.Trace("level: %v, traceData: %+v, msg: %v", level, traceData, msg)
	return
}

func SetLogger(logger Logger) {
	if logger == nil {
		return
	}
	innerLogger = logger
}
