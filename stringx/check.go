package stringx

// IsEmpty checks if the string is empty
func IsEmpty(str string) bool {
	return len(str) == 0
}

// HasEmpty checks if there are any empty strings
func HasEmpty(strings ...string) bool {
	for _, str := range strings {
		if IsEmpty(str) {
			return true
		}
	}

	return false
}
