package record

import (
	"fmt"

	"github.com/quinn-tao/hmis/v1/internal/chrono"
	"github.com/quinn-tao/hmis/v1/internal/coins"
)

type Record struct {
	Id       int
	Amount   coins.RawAmountVal
	Name     string
	Category string
	Date     chrono.Date
}

func (r *Record) String() string {
	return fmt.Sprintf("%v: id=%v cents=%v [%v] on %v",
		r.Name, r.Id, r.Amount, r.Category, r.Date)
}
