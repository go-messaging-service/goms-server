package util

// ContainsString removes the element e from the slice s.
func ContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// RemoveStrings will remove the given strings from the array s.
func RemoveStrings(s []string, e []string) []string {
	result := []string{}

	for _, a := range s {
		if !ContainsString(e, a) {
			result = append(result, a)
		}
	}

	return result
}

// RemoveString will remove the given string e from the slice s.
func RemoveString(s []string, e string) []string {
	result := []string{}

	for _, a := range s {
		if a != e {
			result = append(result, a)
		}
	}

	return result
}
