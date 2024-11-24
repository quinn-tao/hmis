package chronic

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/quinn-tao/hmis/v1/internal/debug"
)

// Date Range Formatter
// Used for fluent representation of date range formatting
//
// This formatter does not support time-zone, nor does it support
// selecting future date ranges
//
// This is primarily used for users who does not want
// to use precision timing, but rather want to specify a more
// abstract timing concept. For instance, select all data from
// **2 months ago**.

const (
	DateRangePresentToken = "Present"
)

func NewInvalidDateRangeErr(strDateRange string) error {
	errMsg := fmt.Sprintf("Invalid date range: %v because of %v",
		strDateRange, "")
	return errors.New(errMsg)
}

// Generate a parsing error message, returns a new error
func NewInvalidDateRangeErrWithReason(strDateRange string, reason string) error {
	errMsg := fmt.Sprintf("Invalid date range: %v because of %v",
		strDateRange, reason)
	return errors.New(errMsg)
}

// DateRange specifies a time range in dates in the past
// TODO: maybe compare this and that?
type DateRange struct {
	from time.Time
	to   time.Time
}

func (d DateRange) String() string {
	formatString := "2001-01-01"
	return fmt.Sprintf("[%v:%v]",
		d.to.Format(formatString),
		d.from.Format(formatString),
	)
}

// Create New DateRange from start and end tokens
func NewDateRangeFromString(fromDateToken string,
	toDateToken string) (DateRange, error) {
	fromDate, err := ParseDateToken(fromDateToken)
	if err != nil {
		return DateRange{}, nil
	}

	toDate, err := ParseDateToken(toDateToken)
	if err != nil {
		return DateRange{}, nil
	}

	dateRange := DateRange{
		from: fromDate,
		to:   toDate,
	}

	debug.Tracef("Date Range Parsed %v", dateRange)
	return dateRange, nil
}

// Parsing a single date specifier into time
// A date specifier is one of:
// 1. A special identifier
// 2. An actual date
// 3. A relative date
//
// # All of the specifier above are selects inclusive date ranges
//
// Specifier Formats:
// Special Identifiers:
// - "present", "today", and "now" translates to today
// - "mtd" and "month" translates to first day of current month
// - "ytd" and "year" translates to first day of current year.
//
// Date Formats:
// Supports a bit more than the standard date format parsing
// yyyy-mm-dd, yy-mm-dd, mm-dd, dd, mm/dd, mm/dd/yy, mm/dd/yyyy are
// all valid formats. If a year is specified using "yy", then it is
// treated as "20mm" year.
//
// Relative Dates
// Supported formats are: Xd, Xm, Xy, X where X is an positive integer
// Token Interpretation:
// - Xd, X - X days before current date
// - Xm - X months before current month
// - Xy - X years before current year
func ParseDateToken(strDateRangeToken string) (time.Time, error) {
	now := time.Now()

	// try Special Tokens
	t, err := parseDateUsingSpecialIdentifiers(strDateRangeToken)
	if err == nil {
		return t, nil
	}

	// try date formats
	t, err = parseDateUsingDateFmt(strDateRangeToken)
	if err == nil {
		return t, nil
	}

	// try Relative dates
	t, err = parseDateUsingRelativeDates(strDateRangeToken)
	if err == nil {
		return t, nil
	}

	return now, err
}

// Parsing string date token by trying to match them with special identifiers
func parseDateUsingSpecialIdentifiers(strDateRangeToken string) (time.Time, error) {
	now := time.Now()
	switch strings.ToLower(strDateRangeToken) {
	case "present":
		fallthrough
	case "today":
		fallthrough
	case "now":
		return time.Now(), nil
	case "mtd":
		fallthrough
	case "month":
		return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local), nil
	case "ytd":
		fallthrough
	case "year":
		return time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, time.Local), nil
	}
	return now, NewInvalidDateRangeErrWithReason(strDateRangeToken,
		"No matching special names")
}

// Parsing date string token using date format
func parseDateUsingDateFmt(strDateRangeToken string) (time.Time, error) {
	now := time.Now()

    // Specify what needs to be auto completed to current time value 
    // when user use date format
    const (
        complete = iota
        missingYY
        missingYYMM
    )
    type DateFmtAutoCompleteType int

    formatters := map[string]DateFmtAutoCompleteType{
        time.DateOnly: complete,
        "06-01-02": complete,
        "01-02": missingYY,
        "01/02": missingYY,
        "01/02/06": complete,
        "01/02/2006": complete,
        "01": missingYYMM,
    }

	for fm, autoCompletType := range formatters {
		t, err := time.Parse(fm, strDateRangeToken)
		if err == nil {
            switch autoCompletType {
            case missingYYMM:
                t = t.AddDate(0, int(now.Month()) - 1, 0) 
                fallthrough
            case missingYY:
                t = t.AddDate(now.Year(), 0, 0)
            } 

			return t, nil
		}
	}

	return now, NewInvalidDateRangeErrWithReason(strDateRangeToken,
		"Cannot parse date using date formats")
}

// Parsing date string token using relative dates formats:
func parseDateUsingRelativeDates(strDateRangeToken string) (time.Time, error) {
	now := time.Now()

	newInvalidValErr := func() error {
		return NewInvalidDateRangeErrWithReason(strDateRangeToken,
			"Invalid relative date value")
	}

	if strings.HasSuffix(strDateRangeToken, "y") {
		relativeDays, err :=
			strconv.Atoi(strDateRangeToken[:len(strDateRangeToken)-1])
		if err != nil {
			return now, newInvalidValErr()
		}

		return now.AddDate((0 - relativeDays), 0, 0), nil
	}

	if strings.HasSuffix(strDateRangeToken, "m") {
		relativeDays, err :=
			strconv.Atoi(strDateRangeToken[:len(strDateRangeToken)-1])
		if err != nil {
			return now, newInvalidValErr()
		}

		return now.AddDate(0, (0 - relativeDays), 0), nil
	}

	if strings.HasSuffix(strDateRangeToken, "d") {
		relativeDays, err :=
			strconv.Atoi(strDateRangeToken[:len(strDateRangeToken)-1])
		if err != nil {
			return now, newInvalidValErr()
		}
        
		return now.AddDate(0, 0, (0 - relativeDays)), nil
	}

	relativeDays, err := strconv.Atoi(strDateRangeToken)

	if err == nil {
		return now.AddDate(0, 0, (0 - relativeDays)), nil
	}

	return now, NewInvalidDateRangeErrWithReason(strDateRangeToken,
		"Cannot parse date using relative dates")
}
