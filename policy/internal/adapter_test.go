package internal

import (
	"testing"

	"github.com/casbin/casbin"
)

func Test_KeyMatchRbac(t *testing.T) {
	conf := `
[request_definition]
r = sub, obj, act
[policy_definition]
p = sub, obj, act
[role_definition]
g = _ , _
[policy_effect]
e = some(where (p.eft == allow))
[matchers]
m = g(r.sub, p.sub)  && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)
`
	line := `
p, alice, /alice_data/*, (GET)|(POST)
p, alice, /alice_data/resource1, POST
p, data_group_admin, /admin/*, POST
p, data_group_admin, /bob_data/*, POST
g, alice, data_group_admin
`
	sa := NewAdapter(line)
	e := casbin.NewEnforcer(casbin.NewModel(conf), sa)
	sub := "alice"
	obj := "/alice_data/login"
	act := "POST"
	if e.Enforce(sub, obj, act) != true {
		t.Error("**error**")
	}
}

func Test_StringRbac(t *testing.T) {
	conf := `
[request_definition]
r = sub, obj, act
[policy_definition]
p = sub, obj, act
[role_definition]
g = _ , _
[policy_effect]
e = some(where (p.eft == allow))
[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`
	line := `
p, alice, data1, read
p, data_group_admin, data3, read
p, data_group_admin, data3, write
g, alice, data_group_admin
`
	sa := NewAdapter(line)
	e := casbin.NewEnforcer(casbin.NewModel(conf), sa)
	sub := "alice" // the user that wants to access a resource.
	obj := "data1" // the resource that is going to be accessed.
	act := "read"  // the operation that the user performs on the resource.
	if e.Enforce(sub, obj, act) != true {
		t.Error("**error**")
	}
}
