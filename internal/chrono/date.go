package chrono

import "time"

type Date struct {
	time.Time
}

var invalidDate = Date{time.Now()}

// Get new date from internal string format
func NewDate(raw string) (Date, error) {
	date, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return invalidDate, err
	}
	return Date{date}, nil
}

// Parse new date from user input token
func ParseDate(token string) (Date, error) {
	date, err := ParseDateToken(token)
	if err != nil {
		return Date{time.Now()}, err
	}
	return Date{date}, nil
}

func Today() Date {
	return Date{time.Now()}
}

func (d Date) String() string {
	return d.Format("2006-01-02")
}

func (d Date) Equal(other Date) bool {
	yy, mm, dd := d.Date()
	return yy == other.Year() &&
		mm == other.Month() &&
		dd == other.Day()
}
