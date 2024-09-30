package profile

import (
	"fmt"

	"github.com/quinn-tao/hmis/v1/internal/amount"
	"github.com/quinn-tao/hmis/v1/internal/debug"
)

var limitParser = Parser{parseLimit}
var limitKey = "limit"

func parseLimit(p *Profile, yamlRoot map[interface{}]interface{}) error {
    var err error
    
    limitRaw, exists := yamlRoot[limitKey]
    if !exists {
        debug.Trace("Limit not found in profile. Skipping...")     
        return parseError(limitKey)
    }
    
    limitStr := fmt.Sprint(limitRaw)

    amt, err := amount.NewFromString(limitStr)
    if err != nil {
        return err
    }
    p.Limit = amt
    debug.Tracef("Limit parsed %v, raw %v", p.Limit, limitStr)     
    return err
}

