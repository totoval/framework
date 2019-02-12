package helpers

func InSlice(needle interface{}, slice []interface{}) bool {
	for _, value := range slice {
		if value == needle {
			return true
		}
	}
	return false
}
