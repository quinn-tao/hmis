package profile

import (
	"github.com/quinn-tao/hmis/v1/internal/debug"
)

var limitParser = Parser{parseLimit}
var limitKey = "limit"

func parseLimit(p *Profile, yamlRoot map[interface{}]interface{}) error {
    var err error
    
    limitStr, exists := yamlRoot[limitKey]
    if !exists {
        debug.Trace("Limit not found in profile. Skipping...")     
        return parseError(limitKey)
    }

    amt := p.Currency.Amount(limitStr)
    p.Limit = amt
    debug.Tracef("Limit parsed %v", p.Limit)     
    return err
}

