package policy

type RoutePolicier interface {
	Can(policy Policier, action Action)
}
