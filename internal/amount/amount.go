package amount

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Internal Amount Representation
// Amount is in terms of "cents", and is detached from the currency
// so it represents a raw unit-less value
// 'hi' is the dollor part; 'lo' is the cent part
type RawAmountVal int64

var (
    ErrInvalidAmount = errors.New("Invalid amount string")
    ErrAmountTooLarge = errors.New("Amount too large")
)

/// Parsers a new RawAmountVal from string 
func NewFromString(amountStr string) (RawAmountVal, error) {
    tokens := strings.Split(amountStr, ".")

    if len(tokens) > 2 {
        return 0, ErrInvalidAmount
    }

    hi := int64(0)
    if len(tokens) > 0 {
        val, err := getHi(tokens[0])
        if err != nil {
            return 0, err
        }
        hi = val
    }
    
    lo := int64(0)
    if len(tokens) > 1 {
        val, err := getLo(tokens[1]) 
        if err != nil {
            return 0, err
        }
        lo = val
    }

    return RawAmountVal(hi * 100 + lo), nil
}

func (r RawAmountVal) String() string {
    hi := int64(r) / 100
    lo := int64(r) % 100
    sb := strings.Builder{}

    if hi < 0 {
        hi = hi * -1
        lo = lo * -1
        sb.WriteRune('-')
    }

    st := []int64{lo}
    for ; hi > 0; hi /= 1000 {
        st = append(st, hi % 1000) 
    }

    for i := len(st) - 1; i >= 0; i-- {
        v := st[i]
        if i == 0 {
            sb.WriteString(fmt.Sprintf(".%02v", v)) 
            return sb.String()
        } 
        if i == len(st) - 1 {
            sb.WriteString(fmt.Sprintf("%v", v)) 
        } else {
            sb.WriteString(fmt.Sprintf("%03v", v)) 
        }
        if i != 1 {
            sb.WriteRune(',')
        }
    }
    return sb.String()
}

func (r RawAmountVal) MarshalYAML() (interface{}, error){
    return r.String(), nil
}

// Parsers hi (dollar) parts of an amount value from string
func getHi(hiStr string) (int64, error) {
    tokens := strings.Split(hiStr, ",")
    
    amount := int64(0)
    for i, token := range tokens {
        if intval, err := strconv.ParseInt(token, 10, 64); err == nil {
            if i == 0 {
                if intval == 0 {
                    return 0, ErrInvalidAmount
                }
            } else {
                if len(token) != 3 {
                    return 0, ErrInvalidAmount
                }
                if amount > (math.MaxInt64 - intval) / 1000 {
                    return 0, ErrAmountTooLarge
                }
            }
            amount = amount * 1000 + intval
        } else {
            return 0, ErrInvalidAmount
        }
    }

    return amount, nil
}

func getLo(loStr string) (int64, error) {
    if len(loStr) > 2 {
        return 0, ErrInvalidAmount
    }
    intval, err := strconv.ParseInt(loStr, 10, 64)
    if err != nil {
        return 0, ErrInvalidAmount
    }
    if len(loStr) == 1 {
        intval = intval * 10
    }
    return intval, nil
}

