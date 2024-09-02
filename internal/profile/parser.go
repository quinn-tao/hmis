package profile

import (
	"errors"
	"fmt"
)

type Parser struct {
    parse func(* Profile, map[interface{}]interface{}) error
}

func parseError(key string) error {
    return errors.New(fmt.Sprintf("Error parsing %v", key))
}
