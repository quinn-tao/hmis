package profile

import (
	"github.com/quinn-tao/hmis/v1/internal/debug"
)

func genStringFieldParser(retv *string, key string) Parser {
	parseFn := func(p *Profile, yamlRoot map[interface{}]interface{}) error {
		val, exists := yamlRoot[key]
		if !exists {
			debug.Tracef("key:%v not found in profile. Skipping...", key)
			return parseError(key)
		}

		strVal, ok := val.(string)
		if !ok {
			debug.Tracef("%v parser error. Skipping... ", strVal)
			return parseError(key)
		}

		*retv = strVal
		return nil
	}
	return Parser{parseFn}
}
