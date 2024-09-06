package util

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func StringToCents(amountStr string) (int, error){
    tokens := strings.Split(amountStr, ".")

    if len(tokens) > 2 {
        return 0, errors.New(fmt.Sprintf("Invalid amount %v", amountStr))
    }
    
    centsHi := 0
    if len(tokens) > 0 {
        hi, err := strconv.Atoi(tokens[0]) 
        if err != nil {
            return 0, errors.New(fmt.Sprintf("Invalid amount %v", amountStr))
        }
        centsHi = hi
    } 
    
    centsLo := 0
    if len(tokens) > 1 {
        lo, err := strconv.Atoi(tokens[1]) 
        if err != nil || len(tokens[1]) > 2 {
            return 0, errors.New(fmt.Sprintf("Invalid amount %v", amountStr))
        }
        
        if len(tokens[1]) == 1 {
            centsLo = lo * 10
        } else {
            centsLo = lo 
        }
    }
    
    return centsHi * 100 + centsLo, nil
}

func CentsToString(amountCents int) string {
    centsHi := amountCents / 100
    centsLo := amountCents % 100
    return fmt.Sprintf("%v.%02v", centsHi, centsLo)
}
