package profile

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/quinn-tao/hmis/v1/internal/coins"
	"github.com/quinn-tao/hmis/v1/internal/debug"
	"github.com/quinn-tao/hmis/v1/internal/display"
)

var categoryParser = Parser{parseCategory}
var categoryKey = "category"

func parseCategory(p *Profile, yamlRoot map[interface{}]interface{}) error {
	var err error

	val, exists := yamlRoot[categoryKey]
	if !exists {
		debug.Tracef("Category not found in profile. Skipping...")
		return parseError(categoryKey)
	}

	categoryYamlRoot := make(map[string]interface{}, 1)
	categoryYamlRoot[categoryKey] = val
	p.Category = categoryTreeMake(p, categoryYamlRoot)

	debug.Tracef("Category parsed: %v", p.Category)
	return err
}

func categoryTreeMake(p *Profile, content interface{}) *Category {
	switch contentType := content.(type) {
	default:
		debug.Tracef("Category has unexpected type:%T", contentType)
		display.Errorf("Error parsing category {%v}", content)
		return nil
	case string:
		return categoryTreeMakeLeaf(content.(string))
	case map[string]interface{}:
		var c Category
		c.Sub = make(map[string]*Category)
		for key, subContentListIntf := range content.(map[string]interface{}) {
			subContentList, ok := subContentListIntf.([]interface{})
			if !ok {
				display.Errorf("Cannot parse category %v", key)
			}
			for _, subContent := range subContentList {
				subCategory := categoryTreeMake(p, subContent)
				if _, exists := c.Sub[subCategory.Name]; exists {
					display.Errorf("Duplicate category of %v", subCategory.Name)
				}
				c.Sub[subCategory.Name] = subCategory
			}
			c.Name = key
			return &c
		}
	}
	return nil
}

func categoryTreeMakeLeaf(content string) *Category {
	var category Category
	tokens := strings.Split(content, " ")

	getName := func(token string) error {
		category.Name = token
		return nil
	}

	getAmount := func(token string) error {
		amt, err := coins.NewFromString(token)
		if err != nil {
			return err
		}
		category.Recurr = &Recurrence{Amount: amt}
		return nil
	}

	getFreq := func(token string) error {
		freq, ok := freqFlagToEnum[token]
		if !ok {
			errMsg := fmt.Sprintf("Error parsing frequency setting in %v", content)
			return errors.New(errMsg)
		}
		category.Recurr.Freq = freq
		return nil
	}

	getDate := func(token string) error {
		date, err := time.Parse(time.DateOnly, token)
		if err != nil {
			return err
		}
		debug.Tracef("%v", date)
		return nil
	}

	getters := []func(string) error{
		getName,
		getAmount,
		getFreq,
		getDate,
	}

	for i, token := range tokens {
		if i >= len(getters) {
			display.Errorf("Too many arguments parsing %v", content)
		}
		err := getters[i](token)
		if err != nil {
			display.Error(err.Error())
		}
	}

	return &category
}
