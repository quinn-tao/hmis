package profile

import (
	"time"

	"golang.org/x/text/currency"
)

// Recurrent expense cycle settings
type Frequency string
const (
    FreqNone = "None" 
    FreqDaily = "Daily"
    FreqWeekly = "Weekly"
    FreqBiWeekly = "BiWeekly"
    FreqMonthly = "Monthly"
    FreqQuaterly = "Quaterly"
    FreqYearly = "Yearly"
)
var freqEnumMap = map[string]Frequency {
    "": FreqNone,  // Default
    "dd": FreqDaily,
    "ww": FreqWeekly, 
    "2ww": FreqBiWeekly,
    "mm": FreqMonthly,
    "qq": FreqQuaterly,
    "yy": FreqYearly,
}

type Category struct {
    Name string
    Recurr *Recurrence
    Sub map[string]*Category
}

type Recurrence struct {
    Freq Frequency
    Amount currency.Amount
    Date time.Time
}

var indent = 0
func (c Category) String() string {
    str := ""
    indent += 1
    for name, sc := range(c.Sub) {
        prefix := ""
        for i := 0; i < indent; i++ {
            prefix += "  "
        }
        str += "\n" + prefix + name + sc.String()
    }
    indent -= 1
    if c.Recurr != nil {
        str += ": " + string(c.Recurr.Freq)
    }
    return str
}
