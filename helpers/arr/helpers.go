package arr

func InSlice(needle interface{}, slice interface{}) bool {
	for _, value := range slice.([]interface{}) {
		if value == needle {
			return true
		}
	}
	return false
}