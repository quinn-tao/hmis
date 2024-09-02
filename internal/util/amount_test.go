package util_test

import (
	"errors"
	"testing"

	"github.com/quinn-tao/hmis/v1/internal/util"
)

func TestStringToCents(t *testing.T) {
    tcs := []struct{
        Name string 
        AmountStr string
        ExpAmountCents int 
        ExpErr error
    }{{
        Name: "Dollar amount without cents",
        AmountStr: "100",
        ExpAmountCents: 10000,
        ExpErr: nil,
    },{
        Name: "Dollar amount with one digit cents",
        AmountStr: "100.5",
        ExpAmountCents: 10050,
        ExpErr: nil,
    }, {
        Name: "Dollar amount with one zero and one digit cents",
        AmountStr: "100.05",
        ExpAmountCents: 10005,
        ExpErr: nil,
    }, {
        Name: "Dollar amount with two digit cents",
        AmountStr: "100.55",
        ExpAmountCents: 10055,
        ExpErr: nil,
    }, {
        Name: "Dollar amount with two digits of zero",
        AmountStr: "100.00",
        ExpAmountCents: 10000,
        ExpErr: nil,
    }, {
        Name: "Dollar amount with one digit of zero",
        AmountStr: "100.0",
        ExpAmountCents: 10000,
        ExpErr: nil,
    }, {
        Name: "Invalid amount",
        AmountStr: "100.023",
        ExpAmountCents: 10000,
        ExpErr: errors.New("Invalid amount 100.023"),
    }}

    for _, tc := range tcs {
        t.Logf("[TestStringToCents] Running %v", tc.Name)
        ActualAmountCents, ActualErr := util.StringToCents(tc.AmountStr)

        if tc.ExpErr != nil {
            if ActualErr == nil {
                t.Fatalf("Expected error, got None")
            }
            return
        }

        if ActualAmountCents != tc.ExpAmountCents {
            t.Fatalf("Expected %v, got %v", tc.ExpAmountCents, ActualAmountCents)
        }
    }
}

func TestCentsToString(t *testing.T) {
    tcs := []struct{
        Name string 
        AmountCents int 
        ExpAmountStr string
    }{{
        Name: "Dollar amount without cents",
        AmountCents: 10000,
        ExpAmountStr: "100.00",
    },{
        Name: "Dollar amount with one digit cents",
        AmountCents: 10050,
        ExpAmountStr: "100.50",
    }, {
        Name: "Dollar amount with one zero and one digit cents",
        AmountCents: 10005,
        ExpAmountStr: "100.05",
    }, {
        Name: "Dollar amount with two digit cents",
        AmountCents: 10055,
        ExpAmountStr: "100.55",
    }}

    for _, tc := range tcs {
        t.Logf("[TestCentsToString] Running %v", tc.Name)

        ActualAmountStr := util.CentsToString(tc.AmountCents)
        if ActualAmountStr != tc.ExpAmountStr {
            t.Fatalf("Expected %v, Got %v", tc.ExpAmountStr, ActualAmountStr)
        }
    }
}

