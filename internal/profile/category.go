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

// TODO: logic handling duplicate categories 
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

// Recursively find catgory in category tree
func (c *Category) FindCategoryRecursive(name string) (retc *Category, exists bool) {
    if c.Name == name {
        return c, true
    }     
    for _, subcategory := range c.Sub {
        target, exists := subcategory.FindCategoryRecursive(name)
        if exists {
            return target, exists
        }
    }
    return nil, false 
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
