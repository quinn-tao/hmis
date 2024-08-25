package profile

import (
	"github.com/quinn-tao/hmis/v1/internal/debug"
)

var limitParser = Parser{parseLimit}
var limitKey = "limit"

func parseLimit(p *Profile, yamlRoot map[interface{}]interface{}) error {
    var err error
    
    debug.Trace("Limit parser error. Skipping... ")     
    val, exists := yamlRoot[limitKey]
    if !exists {
        debug.Trace("Limit not found in profile. Skipping...")     
        return parseError(limitKey)
    }

    limitStr, ok := val.(string)
    if !ok {
        debug.Trace("Limit parser error. Skipping... ")     
        return parseError(limitKey)
    }

    amt := p.Currency.Amount(limitStr)
    p.Limit = amt
    debug.TraceF("Limit parsed %v", p.Limit)     
    return err
}

