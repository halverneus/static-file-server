package cli

// Args parsed from the command-line.
type Args []string

// Parse command-line arguments into Args. Value is returned to support daisy
// chaining.
func Parse(values []string) Args {
	args := Args(values)
	return args
}

// Matches is used to determine if the arguments match the provided pattern.
func (args Args) Matches(pattern ...string) bool {
	// If lengths don't match then nothing does.
	if len(pattern) != len(args) {
		return false
	}

	// Compare slices using '*' as a wildcard.
	for index, value := range pattern {
		if "*" != value && value != args[index] {
			return false
		}
	}
	return true
}
