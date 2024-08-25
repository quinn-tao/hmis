package util

import (
	"fmt"
	"log"
	"os"
)

const (
    red = "\033[0;31m"
    end = "\033[0m"
)

// Fatals are used for user errors
func CheckFatalError(err error) {
    if err != nil {
        fmt.Printf("%s %v %s\n", red, err, end)
        os.Exit(1)
    }
}

func CheckFatalErrorf(err error, msg string, args ...interface{}) {
    if err != nil {
        errMsg := fmt.Sprintf(msg, args...)
        fmt.Printf("%s %s :%v %s\n", red, errMsg, err, end)
        os.Exit(1)
    }
}

// Panics are used for developer used errors
func CheckPanicError(err error) {
    if err != nil {
        panic(err)
    }
}

func CheckPanicErrorf(err error, msg string, args ...interface{}) {
    if err != nil {
        errMsg := fmt.Sprintf(msg, args...)
        log.Panicf("%s :%v", errMsg, err)
    }
}
