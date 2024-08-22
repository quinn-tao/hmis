package debug

import "log"

var traceEnabled = false
var logger = log.Default()

func TraceF(msg string, val ...interface{}) {
    logger.Printf(msg, val)
}

func EnableTrace() {
    traceEnabled = true 
    logger.SetFlags(log.Ltime | log.Ldate | log.Lmicroseconds)
}


