package policy

import (
	"github.com/totoval/framework/model"
)

type key = string
type value = string
type Policier interface {
	Before() *bool
	Create(userIF model.IUser, routeParamMap map[key]value) bool
	Update(userIF model.IUser, routeParamMap map[key]value) bool
	Delete(userIF model.IUser, routeParamMap map[key]value) bool
	ForceDelete(userIF model.IUser, routeParamMap map[key]value) bool
	View(userIF model.IUser, routeParamMap map[key]value) bool
	Restore(userIF model.IUser, routeParamMap map[key]value) bool
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

func policyValidate(user model.IUser, polices Policier, action Action, routeParamMap map[key]value) bool {
	if beforeResult := polices.Before(); beforeResult != nil {
		return *beforeResult
	}

	switch action {
	case ActionCreate:
		return polices.Create(user, routeParamMap)
	case ActionUpdate:
		return polices.Update(user, routeParamMap)
	case ActionDelete:
		return polices.Delete(user, routeParamMap)
	case ActionForceDelete:
		return polices.ForceDelete(user, routeParamMap)
	case ActionView:
		return polices.View(user, routeParamMap)
	case ActionRestore:
		return polices.Restore(user, routeParamMap)
	default:
		return false
	}
}
