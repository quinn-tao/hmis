package coins_test

import (
	"testing"

	"github.com/quinn-tao/hmis/v1/internal/coins"
)

func TestStringToCents(t *testing.T) {
	tcs := []struct {
		Name           string
		AmountStr      string
		ExpAmountCents int64
		ExpErr         error
	}{
		{
			Name:           "Dollar amount without cents",
			AmountStr:      "100",
			ExpAmountCents: 10000,
			ExpErr:         nil,
		}, {
			Name:           "Dollar amount with one digit cents",
			AmountStr:      "100.5",
			ExpAmountCents: 10050,
			ExpErr:         nil,
		}, {
			Name:           "Dollar amount with one zero and one digit cents",
			AmountStr:      "100.05",
			ExpAmountCents: 10005,
			ExpErr:         nil,
		}, {
			Name:           "Dollar amount with two digit cents",
			AmountStr:      "100.55",
			ExpAmountCents: 10055,
			ExpErr:         nil,
		}, {
			Name:           "Dollar amount with two digits of zero",
			AmountStr:      "100.00",
			ExpAmountCents: 10000,
			ExpErr:         nil,
		}, {
			Name:           "Dollar amount with one digit of zero",
			AmountStr:      "100.0",
			ExpAmountCents: 10000,
			ExpErr:         nil,
		}, {
			Name:           "Dollar amount with one ',' no cents",
			AmountStr:      "10,000",
			ExpAmountCents: 1000000,
			ExpErr:         nil,
		}, {
			Name:           "Dollar amount with one ',' and cents",
			AmountStr:      "10,000.23",
			ExpAmountCents: 1000023,
			ExpErr:         nil,
		}, {
			Name:           "Dollar amount with more ',' and cents",
			AmountStr:      "10,000,230.23",
			ExpAmountCents: 1000023023,
			ExpErr:         nil,
		}, {
			Name:           "Invalid amount with too many cent digits",
			AmountStr:      "100.023",
			ExpAmountCents: 10000,
			ExpErr:         coins.ErrInvalidAmount,
		}, {
			Name:           "Invalid amount with wrong place of ','",
			AmountStr:      "1,00.023",
			ExpAmountCents: 0,
			ExpErr:         coins.ErrInvalidAmount,
		}, {
			Name:           "Invalid amount with 0 and ',' at beginning of string",
			AmountStr:      "0,100.023",
			ExpAmountCents: 0,
			ExpErr:         coins.ErrInvalidAmount,
		}, {
			Name:           "Invalid amount with ',' at beginning of string",
			AmountStr:      ",100.023",
			ExpAmountCents: 0,
			ExpErr:         coins.ErrInvalidAmount,
		},
	}

	for _, tc := range tcs {
		t.Logf("[TestStringToCents] running %v", tc.Name)
		actAmountCents, actErr := coins.NewFromString(tc.AmountStr)
		if tc.ExpErr != nil {
			if actErr == nil {
				t.Fatalf("Expected error %v, got %v", tc.ExpErr, actErr)
			}
			return
		}
		if actAmountCents != coins.RawAmountVal(tc.ExpAmountCents) {
			t.Fatalf("Expected %v, got %v", coins.RawAmountVal(tc.ExpAmountCents), actAmountCents)
		}
	}
}

func TestCentsToString(t *testing.T) {
	tcs := []struct {
		Name         string
		AmountCents  int
		ExpAmountStr string
	}{
		{
			Name:         "Dollar amount without cents",
			AmountCents:  10000,
			ExpAmountStr: "100.00",
		}, {
			Name:         "Dollar amount with one digit cents",
			AmountCents:  10050,
			ExpAmountStr: "100.50",
		}, {
			Name:         "Dollar amount with one zero and one digit cents",
			AmountCents:  10005,
			ExpAmountStr: "100.05",
		}, {
			Name:         "Dollar amount with two digit cents",
			AmountCents:  10055,
			ExpAmountStr: "100.55",
		}, {
			Name:         "Dollar amount with more digits",
			AmountCents:  10055055,
			ExpAmountStr: "100,550.55",
		}, {
			Name:         "Dollar amount with even more digits. Pattern xx,000,xx",
			AmountCents:  10100055,
			ExpAmountStr: "101,000.55",
		}, {
			Name:         "Dollar amount with even more digits. Pattern xx,x00,0xx",
			AmountCents:  110001055,
			ExpAmountStr: "1,100,010.55",
		}, {
			Name:         "Dollar amount with even more digits. Negative value",
			AmountCents:  -110001055,
			ExpAmountStr: "-1,100,010.55",
		},
	}

	for _, tc := range tcs {
		t.Logf("[TestCentsToString] Running %v", tc.Name)

		ActualAmountStr := coins.RawAmountVal(tc.AmountCents).String()
		if ActualAmountStr != tc.ExpAmountStr {
			t.Fatalf("Expected %v, Got %v", tc.ExpAmountStr, ActualAmountStr)
		}
	}
}
