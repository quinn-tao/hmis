package chrono

import (
	"errors"
	"fmt"
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
