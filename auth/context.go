package auth

import "github.com/totoval/framework/context"

type Context interface {
	AuthClaimID() (ID uint, exist bool)
	IUserModel() IUser
	context.ResponseContextor
	context.DataContextor
}
