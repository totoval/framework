package internal

import (
	"bytes"
	"errors"
	"strings"

	"github.com/casbin/casbin/model"
	"github.com/casbin/casbin/persist"
	"github.com/casbin/casbin/util"
)

/*type Adapter interface {
    // LoadPolicy loads all policy rules from the storage.
    LoadPolicy(model model.Model) error
    // SavePolicy saves all policy rules to the storage.
    SavePolicy(model model.Model) error
    // AddPolicy adds a policy rule to the storage.
    // This is part of the Auto-Save feature.
    AddPolicy(sec string, ptype string, rule []string) error
    // RemovePolicy removes a policy rule from the storage.
    // This is part of the Auto-Save feature.
    RemovePolicy(sec string, ptype string, rule []string) error
    // RemoveFilteredPolicy removes policy rules that match the filter from the storage.
    // This is part of the Auto-Save feature.
    RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error
}*/

type Adapter struct {
	Line string
}

func NewAdapter(line string) persist.Adapter {
	return &Adapter{
		Line: line,
	}
}

func (sa *Adapter) LoadPolicy(model model.Model) error {
	if sa.Line == "" {
		return errors.New("invalid line, line cannot be empty")
	}
	strs := strings.Split(sa.Line, "\n")
	for _, str := range strs {
		if str == "" {
			continue
		}
		persist.LoadPolicyLine(str, model)
	}

	return nil
}

func (sa *Adapter) SavePolicy(model model.Model) error {
	var tmp bytes.Buffer
	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			tmp.WriteString(ptype + ", ")
			tmp.WriteString(util.ArrayToString(rule))
			tmp.WriteString("\n")
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			tmp.WriteString(ptype + ", ")
			tmp.WriteString(util.ArrayToString(rule))
			tmp.WriteString("\n")
		}
	}
	sa.Line = strings.TrimRight(tmp.String(), "\n")
	return nil
}

func (sa *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

func (sa *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	sa.Line = ""
	return nil
}

func (sa *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return errors.New("not implemented")
}
