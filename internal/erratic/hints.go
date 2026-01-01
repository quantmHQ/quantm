package erratic

type (
	// Hints represents a map of key-value pairs providing additional information about an error.
	//
	// Example:
	//
	//  info := Hints{"field": "invalid value"}
	//  fmt.Println(info) // Output: map[string]any{"field": "invalid value"}
	Hints map[string]any
)

func NewHints(args ...string) Hints {
	odd := len(args)%2 != 0

	limit := len(args)
	if odd {
		limit--
	}

	details := make(Hints)

	for i := 0; i < limit; i += 2 {
		details[args[i]] = args[i+1]
	}

	if odd {
		details["unknown"] = args[len(args)-1]
	}

	return details
}
