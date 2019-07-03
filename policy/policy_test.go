package policy

import (
	"strings"
	"testing"

	"github.com/totoval/framework/model"
	"github.com/totoval/framework/policy/internal"
)

type testUser struct {
}

func (tu *testUser) Scan(userId uint) error {
	*tu = testUser{}
	return nil
}

type UserPolicy struct {
}

func newUserPolicy() *UserPolicy {
	return &UserPolicy{}
}
func (up *UserPolicy) Before() *bool                       { return nil }
func (up *UserPolicy) Create(userIF model.IUser) bool      { panic("need implements create ") }
func (up *UserPolicy) Update(userIF model.IUser) bool      { panic("need implements update ") }
func (up *UserPolicy) Delete(userIF model.IUser) bool      { panic("need implements delete ") }
func (up *UserPolicy) ForceDelete(userIF model.IUser) bool { panic("need implements force delete") }
func (up *UserPolicy) View(userIF model.IUser) bool        { panic("need implements view ") }
func (up *UserPolicy) Restore(userIF model.IUser) bool     { panic("need implements restore ") }

func Test_Enforcer(t *testing.T) {
	enfc := internal.NewEnforcer()
	enfc.AddFunction("totovalPolicyValidate", func(args ...interface{}) (i interface{}, e error) {
		user := args[0].(model.IUser)
		resource := args[1].(Policier)
		operation := strings.ToUpper(args[2].(string)[0:1]) + strings.ToLower(args[2].(string)[1:])

		if beforeResult := resource.Before(); beforeResult != nil {
			return *beforeResult, nil
		}

		switch operation {
		case "Create":
			return resource.Create(user), nil
		case "Update":
			return resource.Update(user), nil
		case "Delete":
			return resource.Delete(user), nil
		case "ForceDelete":
			return resource.ForceDelete(user), nil
		case "View":
			return resource.View(user), nil
		case "Restore":
			return resource.Restore(user), nil
		default:
			return false, nil
		}

	})

	sub := new(testUser)   // the user that wants to access a resource.
	obj := newUserPolicy() // the resource that is going to be accessed.
	act := "create"        // the operation that the user performs on the resource.
	if enfc.Enforce(sub, obj, act) != true {
		t.Error("**error**")
	}
}
