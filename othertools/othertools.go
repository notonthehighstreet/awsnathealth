package othertools

// StringInSlice function checks if the slice contains the given string, return bool.
func StringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
