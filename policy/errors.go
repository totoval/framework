package policy

type NotFoundError struct{}

func (e NotFoundError) Error() string {
	return "policy not found"
}
