package profile

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/quinn-tao/hmis/v1/internal/amount"
)

var (
    ErrInvalidCategoryPath = errors.New("Invalid category path")
    ErrAlreadyExists = errors.New("Category path already exists")
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

func (f Frequency) String() string {
    return freqEnumToFlag[f]
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

func (c *Category) AddCategory(path string) (*Category, error) {
    pathTokens := strings.Split(path, "/")
    if len(pathTokens) < 1 {
        return nil, ErrInvalidCategoryPath
    }
    prefixPath := strings.Join(pathTokens[:len(pathTokens)-1], "/") 
    newCategoryName := pathTokens[len(pathTokens)-1]
    newCategory := Category{Name:newCategoryName}
    
    parent, exists := c.FindCategoryWithPath(prefixPath)
    if !exists {
        return nil, ErrInvalidCategoryPath
    }
    
    if parent.Sub == nil {
        parent.Sub = make(map[string]*Category)
    }
    _, exists = parent.Sub[newCategoryName]
    if exists {
        return nil, ErrAlreadyExists
    }

    parent.Sub[newCategoryName] = &newCategory
    return &newCategory, nil
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

func (c *Category) MarshalYAML() (interface{}, error) {
    yml, err := c.marshalYAML()
    if err != nil {
        return  nil, err
    }
    content, exists := yml.(map[string][]interface{})["category"]
    if !exists {
        return nil, errors.New("Marshalled category tree does not contain 'category' key")
    } 
    return content, nil
}

// Marshalling a category into yaml produces one of:
//   a string, represents a leaf node in the category tree, or
//   a mapping from name of this category to an array of children yaml values
func (c *Category) marshalYAML() (interface{}, error) {
    if c.Sub == nil || len(c.Sub) == 0 {
        if c.Recurr == nil {
            return fmt.Sprintf("%v", c.Name), nil
        }
        return fmt.Sprintf("%v %v %v", c.Name, c.Recurr.Amount, c.Recurr.Freq), nil
    }
    yml := map[string][]interface{}{
        c.Name: make([]interface{}, 0, len(c.Sub)),
    }
    for _, sub := range c.Sub {
        subYml, err := sub.marshalYAML()
        if err != nil {
            return nil, err
        }
        yml[c.Name] = append(yml[c.Name], subYml)
    }  
    return yml, nil
}
