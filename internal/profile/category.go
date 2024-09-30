package profile

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/quinn-tao/hmis/v1/internal/amount"
	"github.com/quinn-tao/hmis/v1/internal/debug"
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
var freqFlagToEnum = map[string]Frequency {
    "": FreqNone,  // Default
    "dd": FreqDaily,
    "ww": FreqWeekly, 
    "2ww": FreqBiWeekly,
    "mm": FreqMonthly,
    "qq": FreqQuaterly,
    "yy": FreqYearly,
}
var freqEnumToFlag = map[Frequency]string {
    FreqNone: "",  // Default
    FreqDaily: "dd",
    FreqWeekly: "ww", 
    FreqBiWeekly: "2ww",
    FreqMonthly: "mm" ,
    FreqQuaterly: "qq",
    FreqYearly: "yy",
}

// TODO: logic handling duplicate categories 
type Category struct {
    Name string
    Recurr *Recurrence
    Sub map[string]*Category
}

type Recurrence struct {
    Freq Frequency
    Amount amount.RawAmountVal
    Date time.Time
}

// Find Category by searching particular path
func (c *Category) FindCategoryWithPath(path string) (retc *Category, exists bool) {
    tokens := strings.SplitN(path, "/",2)
    if c.Name != tokens[0] {
        return nil, false
    }
    
    if len(tokens) == 2 {
        for _, subcategory := range c.Sub {
            target, exists := subcategory.FindCategoryWithPath(tokens[1])
            if exists {
                return target, exists
            }
        }
    } else if len(tokens) == 1 {
        return c, true
    }

    return nil, false
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

func (this *Category) Equals(other *Category) bool {
    if this.Name != other.Name {
        return false
    }    
    for name, thisSub := range this.Sub {
        otherSub, exists := other.Sub[name]
        if !exists || !thisSub.Equals(otherSub) {
            return false
        }
    }
    return true 
}

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

var indent = 0
// Implement yaml unmarshalling
func (c *Category) MarshalYAML() (interface{}, error) {
    str := ""
    if c.Sub != nil && len(c.Sub) > 0 {
        for _, sc := range(c.Sub) {
            next, err := sc.marshalYaml(1)
            if err != nil {
                return nil, err
            } 
            str += fmt.Sprintf("\n%v", next)
        }
    }
    debug.Tracef("\n%v\n", str)
    return str, nil
}

func (c *Category) marshalYaml(indent int) (interface{}, error) {
    str := "- " + c.Name
    if c.Sub != nil && len(c.Sub) > 0 {
        str += ":"
        for _, sc := range(c.Sub) {
            prefix := ""
            for i := 0; i < indent; i++ {
                prefix += "    "
            }
            next, err := sc.marshalYaml(indent + 1)
            if err != nil {
                return nil, err
            }
            str += "\n" + prefix + next.(string)
        }
        return str, nil
    }
    if c.Recurr != nil {
        flag, exists := freqEnumToFlag[c.Recurr.Freq]
        if !exists {
            return nil, errors.New(fmt.Sprintf("frequency setting %v not exists", string(c.Recurr.Freq)))
        }
        amtStr := fmt.Sprintf(" %v",c.Recurr.Amount)
        str += amtStr + " " + flag  
    }
    return str, nil
}

