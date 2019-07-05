package policy

import (
	"github.com/gin-gonic/gin"

	"github.com/totoval/framework/auth"
	"github.com/totoval/framework/model"
)

type key = string
type value = string
type Policier interface {
	Before(IUser model.IUser, routeParamMap map[key]value) *bool
	Create(IUser model.IUser, routeParamMap map[key]value) bool
	Update(IUser model.IUser, routeParamMap map[key]value) bool
	Delete(IUser model.IUser, routeParamMap map[key]value) bool
	ForceDelete(IUser model.IUser, routeParamMap map[key]value) bool
	View(IUser model.IUser, routeParamMap map[key]value) bool
	Restore(IUser model.IUser, routeParamMap map[key]value) bool
}

type Action byte

const (
	ActionCreate Action = iota
	ActionUpdate
	ActionDelete
	ActionForceDelete
	ActionView
	ActionRestore
)

type Authorization struct {
	auth.AuthUser
}

func (a *Authorization) Authorize(c *gin.Context, policies Policier, action Action) (isAbort bool, user model.IUser) {
	if a.AuthUser.Scan(c) {
		return true, nil
	}

	rpm := make(map[key]value)
	if !policyValidate(user, policies, action, rpm) {
		forbid(c)
		return true, a.AuthUser.User()
	}

	return false, a.AuthUser.User()
}

func policyValidate(user model.IUser, policies Policier, action Action, routeParamMap map[key]value) bool {
	if user == nil {
		return true
	}
	if policies == nil {
		return true
	}

	if beforeResult := policies.Before(); beforeResult != nil {
		return *beforeResult
	}

	switch action {
	case ActionCreate:
		return policies.Create(user, routeParamMap)
	case ActionUpdate:
		return policies.Update(user, routeParamMap)
	case ActionDelete:
		return policies.Delete(user, routeParamMap)
	case ActionForceDelete:
		return policies.ForceDelete(user, routeParamMap)
	case ActionView:
		return policies.View(user, routeParamMap)
	case ActionRestore:
		return policies.Restore(user, routeParamMap)
	default:
		return false
	}
}
