package profile

import (
	"github.com/quinn-tao/hmis/v1/internal/amount"
	"github.com/quinn-tao/hmis/v1/internal/debug"
	"github.com/quinn-tao/hmis/v1/internal/util"
	"golang.org/x/text/currency"
)

var currencyParser = Parser{parseCurrency}
var currencyKey = "currency"

func parseCurrency(p *Profile, yamlRoot map[interface{}]interface{}) error {
    var err error

    val, exists := yamlRoot[currencyKey]
    if !exists {
        debug.Tracef("Currency not found in profile. Skipping...")     
        return parseError(currencyKey)
    }

    code, ok := val.(string)
    if !ok {
        debug.Tracef("Currency parser error. Skipping... ")     
        return parseError(currencyKey)
    }
    
    unit, err := currency.ParseISO(code)
    util.CheckError(err) 
    p.Currency = amount.Currency(unit)
    debug.Tracef("Currency parsed %v", p.Currency) 
    return err
}

