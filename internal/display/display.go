package display

import (
	"bufio"
	"fmt"
	"os"

	"github.com/quinn-tao/hmis/v1/internal/util"
)


const (
    red = "\033[0;31m"
    green = "\033[0;33m"
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

// Ask user generic questions
func Dialog(msg string) string {
    reader := bufio.NewReader(os.Stdin)
    fmt.Println(msg)
    res, err := reader.ReadString('\n')
    util.CheckError(err)
    return res
}

// Ask user [Y/n] type of questions
func DialogYesNo(question string) bool {
    reader := bufio.NewReader(os.Stdin)
    fmt.Printf("%v [Y/n]\n", question)
    for {
        res, err := reader.ReadString('\n')
        util.CheckError(err)
        if res == "Y\n" {
            return true
        } else if res == "n\n" {
            return false
        } else {
            fmt.Println("[Y/n]\n")
        }
    }
}
