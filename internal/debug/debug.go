package debug

import "log"

var traceEnabled *bool
var logger = log.Default()

func Trace(msg string, val ...interface{}) {
    if traceEnabled != nil && *traceEnabled {
        logger.Println(msg)
    }
}

func TraceF(msg string, val ...interface{}) {
    if traceEnabled != nil && *traceEnabled {
        logger.Printf(msg, val...)
    }
}

func SetTraceFlag(flag *bool) {
    traceEnabled = flag
}
