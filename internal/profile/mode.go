package profile

// Setting Mode for recurrent expenses
type Mode string

const (
	// Any recurrent expenses would be treated as spent on
	// the first day of the cycle. This overrides user config on expense date
	ModePessimistic = "pessimistic"
	// Any recurrent expenses would be treated as spent on
	// the last day of the cycle. This overrides user config on expense date
	ModeOptimistic = "optimistic"
	// Any recurrent expenses would be treated as spent on
	// the first day of the cycle unless specified by user
	ModePrecisePessimistic = "precise-pessimistic"
	// Any recurrent expenses would be treated as spent on
	// the last day of the cycle unless specified by user
	ModePreciseOptimistic = "precise-optimistic"
)

var modeEnumMap = map[string]Mode{
	"":                    ModePessimistic, // Default
	"pessimistic":         ModePessimistic,
	"optimistic":          ModeOptimistic,
	"precise-pessimistic": ModePrecisePessimistic,
	"precise-optimistic":  ModePreciseOptimistic,
}
