package profile

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/quinn-tao/hmis/v1/internal/coins"
	"github.com/quinn-tao/hmis/v1/internal/debug"
)

var (
	ErrInvalidCategoryPath = errors.New("Invalid category path")
	ErrAlreadyExists       = errors.New("Category path already exists")
)

// Recurrent expense cycle settings
type Frequency string

const (
	FreqNone     = "None"
	FreqDaily    = "Daily"
	FreqWeekly   = "Weekly"
	FreqBiWeekly = "BiWeekly"
	FreqMonthly  = "Monthly"
	FreqQuaterly = "Quaterly"
	FreqYearly   = "Yearly"
)

var freqFlagToEnum = map[string]Frequency{
	"":    FreqNone, // Default
	"dd":  FreqDaily,
	"ww":  FreqWeekly,
	"2ww": FreqBiWeekly,
	"mm":  FreqMonthly,
	"qq":  FreqQuaterly,
	"yy":  FreqYearly,
}
var freqEnumToFlag = map[Frequency]string{
	FreqNone:     "", // Default
	FreqDaily:    "dd",
	FreqWeekly:   "ww",
	FreqBiWeekly: "2ww",
	FreqMonthly:  "mm",
	FreqQuaterly: "qq",
	FreqYearly:   "yy",
}

func (f Frequency) String() string {
	return freqEnumToFlag[f]
}

type Recurrence struct {
	Freq   Frequency
	Amount coins.RawAmountVal
	Date   time.Time
}

func (r Recurrence) String() string {
	return fmt.Sprintf("%v %v", r.Freq, r.Amount)
}

// TODO: logic handling duplicate categories
type Category struct {
	Name   string
	Recurr *Recurrence
	Sub    map[string]*Category
}

// Category APIs
// ================================================================================

// Insert the category into category tree
// New category is defined by using a unix-like path, where
// the last element of the path is the name of the new category
func (c *Category) AddCategory(path string) (*Category, error) {
	pathTokens := strings.Split(path, "/")
	if len(pathTokens) < 1 {
		return nil, ErrInvalidCategoryPath
	}

	newCategoryName := pathTokens[len(pathTokens)-1]
	newCategory := Category{Name: newCategoryName}

	if len(pathTokens) == 1 {
		err := c.insertSub(&newCategory)
		if err != nil {
			return nil, err
		}
		return &newCategory, nil
	}

	// Break down new path into <prefix>/<new category name>
	prefixPath := strings.Join(pathTokens[:len(pathTokens)-1], "/")

	// Locate parent category
	parent, exists := c.FindCategoryWithPath(prefixPath)
	if !exists {
		return nil, ErrInvalidCategoryPath
	}

	// Insert new child
	err := parent.insertSub(&newCategory)
	if err != nil {
		return nil, err
	}

	return &newCategory, nil
}

// Find Category by searching particular path
func (c *Category) FindCategoryWithPath(path string) (retc *Category, exists bool) {
	tokens := strings.SplitN(path, "/", 2)

	if c.Sub == nil || len(c.Sub) == 0 {
		return nil, false
	}

	next, exists := c.Sub[tokens[0]]
	if !exists {
		return nil, false
	}

	if len(tokens) == 1 {
		return next, true
	}
	return next.FindCategoryWithPath(tokens[1])
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

type CategorySelector func(*Category) bool

var CategorySelectAll CategorySelector = func(c *Category) bool { return true }

type CategoryAccumulator func(*Category, interface{}) interface{}

// Visit Category Tree and apply selector and/or accumulator
// Selected categories and Accumulated Values are returned
// If selector is nil, then all Categories are selected
// If accumulator is nil, then nil will be returned as accumulated value
func (c *Category) Visit(selector CategorySelector,
	accumulator CategoryAccumulator, currAcc interface{}) ([]*Category, interface{}) {
	var sel []*Category

	if selector == nil || selector(c) {
		sel = []*Category{c}
	}

	acc := accumulator(c, currAcc)
	for _, sub := range c.Sub {
		subSel, subAcc := sub.Visit(selector, accumulator, acc)
		sel = append(sel, subSel...)
		acc = subAcc
	}

	debug.Tracef("visiting %v acc %v", c, acc)

	return sel, acc
}

// Utilities
// ============================================================================

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
	for name, sc := range c.Sub {
		prefix := ""
		for i := 0; i < indent; i++ {
			prefix += "        "
		}
		str += "\n" + prefix + name + sc.String()
	}
	indent -= 1
	if c.Recurr != nil {
		str += fmt.Sprintf(": %v", c.Recurr)
	}
	return str
}

func (c *Category) MarshalYAML() (interface{}, error) {
	yml, err := c.marshalYAML()
	if err != nil {
		return nil, err
	}
	content, exists := yml.(map[string][]interface{})["category"]

	debug.Tracef("%v", c)
	if !exists {
		return nil, errors.New("Marshalled category tree does not contain 'category' key")
	}
	return content, nil
}

// Marshalling a category into yaml produces one of:
//	a string, represents a leaf node in the category tree, or
//	a mapping from name of this category to an array of children yaml values
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
	debug.Tracef("%v == %v\n", c.Name, yml[c.Name])
	return yml, nil
}

func (c *Category) insertSub(sub *Category) error {
	if c.Sub == nil {
		c.Sub = make(map[string]*Category, 1)
	}
	_, exists := c.Sub[sub.Name]
	if exists {
		return ErrAlreadyExists
	}
	c.Sub[sub.Name] = sub
	return nil
}
