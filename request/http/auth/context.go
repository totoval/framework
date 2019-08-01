package auth

type Context interface {
	AuthClaimID() (ID uint, exist bool)
	IUserModel() IUser
}
