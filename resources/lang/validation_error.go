package lang

type ValidationError map[string]string

func (ve *ValidationError) First() string {
	for _, value := range *ve {
		return value
	}
	return ""
}
