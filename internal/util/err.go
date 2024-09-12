package util

import (
	"fmt"
	"os"

	"github.com/quinn-tao/hmis/v1/internal/debug"
)

const (
    red = "\033[0;31m"
    end = "\033[0m"
)

func CheckError(err error) {
    if err != nil {
        debug.Tracef("%v", err)
        fmt.Printf("%s %v %s\n", red, err, end)
        os.Exit(1)
    }
}

func CheckErrorf(err error, msg string, args ...interface{}) {
    if err != nil {
        errMsg := fmt.Sprintf(msg, args...)
        debug.Tracef("%v", errMsg)
        fmt.Printf("%s %s :%v %s\n", red, errMsg, err, end)
        os.Exit(1)
    }
}
