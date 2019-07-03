package internal

import (
	"github.com/casbin/casbin"
	"github.com/casbin/casbin/model"
	"github.com/casbin/casbin/persist"
)

func m() *model.Model {
	modelText := `[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && (r.act == p.act || p.act == "*")`
	modelText = `[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = totovalPolicyValidate(r.sub, r.obj, r.act)`
	m := casbin.NewModel(modelText)
	return &m
}
func p() persist.Adapter {
	basicPolicy := ``
	p := NewAdapter(basicPolicy)
	return p
}

func NewEnforcer() *casbin.Enforcer {
	return casbin.NewEnforcer(*m(), p())
}
