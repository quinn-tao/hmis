package amount

import "golang.org/x/text/currency"

type Currency currency.Unit

func (c Currency) MarshalYAML() (interface{}, error) {
    return currency.Unit(c).String(), nil
}
