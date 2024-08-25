package display

import (
    "os"
    "fmt"
)


const (
    red = "\033[0;31m"
    end = "\033[0m"
)

func Error(msg string) {
    fmt.Printf("%s %s %s\n", red, msg, end)
    os.Exit(1)
}

func Errorf(msg string, args... interface{}) {
    errMsg := fmt.Sprintf(msg, args...)
    fmt.Printf("%s %s %s\n", red, errMsg, end)
    os.Exit(1)
}
