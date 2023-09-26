package configs

func PriorityString(s ...string) string {
	for _, str := range s {
		if str != "" {
			return str
		}
	}
	return ""
}

func UniqueStrings(strings []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, s := range strings {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}

	return result
}
