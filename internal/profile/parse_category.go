package profile

import (
	"strconv"
	"strings"
	"time"

	"github.com/quinn-tao/hmis/v1/internal/debug"
	"github.com/quinn-tao/hmis/v1/internal/display"
	"github.com/quinn-tao/hmis/v1/internal/util"
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
        return categoryTreeMakeLeaf(p, content.(string))
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
                c.Sub[subCategory.Name] = subCategory
            } 
            c.Name = key 
            return &c
        }
    }
    return nil
}

func categoryTreeMakeLeaf(p *Profile, content string) *Category {
    var category Category
    tokens := strings.Split(content, " ")
    
    getName := func(token string) {
        category.Name = token
    }

    getAmount := func(token string) {
        // TODO: convert dollar amount to cents
        intval, err := strconv.Atoi(token)
        util.CheckError(err)
        amt := p.Currency.Amount(intval)
        category.Recurr = &Recurrence{Amount: amt}
    }

    getFreq := func(token string) {
        freq, ok := freqEnumMap[token]
        if !ok {
            display.Errorf("Error parsing frequency setting in %v", content)
        }
        category.Recurr.Freq = freq
    }

    getDate := func(token string) {
        date, err := time.Parse(time.DateOnly, token)
        util.CheckError(err)
        debug.Tracef("%v", date)
    }

    getters := []func(string) {
        getName,
        getAmount,
        getFreq,
        getDate,
    }

    for i, token := range tokens {
        if i >= len(getters) {
            display.Errorf("Too many arguments parsing %v", content)
        } 
        getters[i](token)  
    }
    return &category
}

