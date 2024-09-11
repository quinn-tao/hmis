package profile

import (
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
        debug.TraceF("Currency not found in profile. Skipping...")     
        return parseError(currencyKey)
    }

    code, ok := val.(string)
    if !ok {
        debug.TraceF("Currency parser error. Skipping... ")     
        return parseError(currencyKey)
    }
    
    p.Currency, err = currency.ParseISO(code)
    util.CheckError(err) 
    debug.TraceF("Currency parsed %v", p.Currency) 
    return err
}

