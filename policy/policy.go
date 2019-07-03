package policy

import (
	"github.com/casbin/casbin"

	"github.com/totoval/framework/model"
	"github.com/totoval/framework/policy/internal"
)

type Policier interface {
	Before() *bool
	Create(userIF model.IUser) bool
	Update(userIF model.IUser) bool
	Delete(userIF model.IUser) bool
	ForceDelete(userIF model.IUser) bool
	View(userIF model.IUser) bool
	Restore(userIF model.IUser) bool
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

var enfc *casbin.Enforcer

func Initialize() {
	enfc = initEnforcer()
}

func initEnforcer() *casbin.Enforcer {
	enfc := internal.NewEnforcer()
	enfc.AddFunction("totovalPolicyValidate", func(args ...interface{}) (i interface{}, e error) {
		user := args[0].(model.IUser)
		resource := args[1].(Policier)
		operation := args[2].(Action)
		
		if beforeResult := resource.Before(); beforeResult != nil {
			return *beforeResult, nil
		}

		switch operation {
		case ActionCreate:
			return resource.Create(user), nil
		case ActionUpdate:
			return resource.Update(user), nil
		case ActionDelete:
			return resource.Delete(user), nil
		case ActionForceDelete:
			return resource.ForceDelete(user), nil
		case ActionView:
			return resource.View(user), nil
		case ActionRestore:
			return resource.Restore(user), nil
		default:
			return false, nil
		}
	})

	return enfc
}
