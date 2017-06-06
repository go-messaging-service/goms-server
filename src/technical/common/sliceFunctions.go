package technicalCommon

// ContainsString removes the element e from the slice s.
func ContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
