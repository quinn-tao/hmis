package profile

import (
	"github.com/quinn-tao/hmis/v1/internal/debug"
	"github.com/quinn-tao/hmis/v1/internal/display"
)

var modeParser = Parser{parseMode}
var modeKey = "mode"

func parseMode(p *Profile, yamlRoot map[interface{}]interface{}) error {
    var err error

    val, exists := yamlRoot[modeKey]
    if !exists {
        debug.Trace("Mode not found in profile. Skipping...")     
        return parseError(modeKey)
    }

    modeStr, ok := val.(string)
    if !ok {
        debug.Trace("Mode parser error. Skipping... ")     
        return parseError(modeKey)
    }
    
    mode, ok := modeEnumMap[modeStr]
    if !ok {
        display.Errorf("Error parsing mode {%v}", modeStr)
    }
    p.Mode = mode

    debug.TraceF("Mode parsed %v", p.Mode) 
    return err
}

