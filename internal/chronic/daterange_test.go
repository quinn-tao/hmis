package chronic_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/quinn-tao/hmis/v1/internal/chronic"
)

const (
	TypeSpecialIdentifier string = "Special Identifier"
	TypeDateFormat        string = "Date Format"
	TypeRelativeDate      string = "Relative Date"
)

var invalidDate = time.Now()

func TestParseDateToken(t *testing.T) {
	tcname := func(fmtType string, name string) string {
		return fmt.Sprintf("[TestNewDateRangeFromString] Format:%v: %v",
			fmtType, name)
	}

	// Depends on the timing, these tests may fail when ran
	// at time boundaries (end of year, end of month etc)
	tcs := []struct {
		Name        string
		InputTokens []string
		ExpErr      bool
		ExpDate     time.Time
	}{
		{
			Name:        tcname(TypeSpecialIdentifier, "present"),
			InputTokens: []string{"present", "Present", "PRESENT"},
			ExpErr:      false,
			ExpDate:     time.Now(),
		},
		{
			Name:        tcname(TypeSpecialIdentifier, "mtd"),
			InputTokens: []string{"mtd", "month"},
			ExpErr:      false,
			ExpDate: time.Date(time.Now().Year(), time.Now().Month(), 1,
				0, 0, 0, 0, time.Local),
		},
		{
			Name:        tcname(TypeSpecialIdentifier, "ytd"),
			InputTokens: []string{"ytd", "year"},
			ExpErr:      false,
			ExpDate: time.Date(time.Now().Year(), time.January,
				1, 0, 0, 0, 0, time.Local),
		},
		{
			Name: tcname(TypeDateFormat, "Valid year, month and date"),
			InputTokens: []string{"2000-02-01", "00-02-01", "02/01/2000",
				"02/01/00"},
			ExpErr: false,
			ExpDate: time.Date(2000, time.February, 1,
				0, 0, 0, 0, time.Local),
		},
		{
			Name:        tcname(TypeDateFormat, "Valid month and date"),
			InputTokens: []string{"02-01", "02/01"},
			ExpErr:      false,
			ExpDate: time.Date(time.Now().Year(), time.February,
				1, 0, 0, 0, 0, time.Local),
		},
		{
			Name:        tcname(TypeDateFormat, "Valid date"),
			InputTokens: []string{"01"},
			ExpErr:      false,
			ExpDate: time.Date(time.Now().Year(), time.Now().Month(),
				1, 0, 0, 0, 0, time.Local),
		},
		{
			Name:        tcname(TypeDateFormat, "Invalid Dates"),
			InputTokens: []string{"24-24", "13/1", "1/32", "2023-02-29"},
			ExpErr:      true,
			ExpDate:     invalidDate,
		},
		{
			Name:        tcname(TypeRelativeDate, "Valid date"),
			InputTokens: []string{"3d", "3"},
			ExpErr:      false,
			ExpDate:     time.Now().AddDate(0, 0, -3),
		},
		{
			Name:        tcname(TypeRelativeDate, "Valid month"),
			InputTokens: []string{"3m"},
			ExpErr:      false,
			ExpDate:     time.Now().AddDate(0, -3, 0),
		},
		{
			Name:        tcname(TypeRelativeDate, "Valid year"),
			InputTokens: []string{"3y"},
			ExpErr:      false,
			ExpDate:     time.Now().AddDate(-3, 0, 0),
		},
	}

	isSameDate := func(this time.Time, that time.Time) bool {
		return this.Day() == that.Day() &&
			this.Month() == that.Month() &&
			this.Year() == that.Year()
	}

	for _, tc := range tcs {
		t.Logf("Running %v", tc.Name)
		for _, token := range tc.InputTokens {
			actDate, actErr := chronic.ParseDateToken(token)
			if actErr != nil && !tc.ExpErr {
				t.Fatalf("Parsing %v. Expect no error, got %v",
                    token, actErr)
			}

            if actErr == nil && tc.ExpErr {
				t.Fatalf("Parsing %v. Expect error; got nil and parsed %v", 
                    token, actDate)
            }

			if !isSameDate(actDate, tc.ExpDate) {
				t.Fatalf("Parsing %v. Expecting %v, got %v",
					token, tc.ExpDate, actDate)
			}
		}
	}
}
