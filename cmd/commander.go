package cmd

import "errors"

type Commander interface {
	Command() string
	Description() string
	Handler(arg *Arg) error
}

type Arg struct {
	argList *map[string]string
}

func newArg(argList *map[string]string) *Arg {
	return &Arg{argList: argList}
}
func (a *Arg) Get(argName string) (data *string, err error) {
	if a.argList == nil {
		return nil, errors.New("arg <" + argName + "> is not set")
	}

	al := *a.argList
	_data, ok := al[argName]
	if !ok {
		return nil, errors.New("arg is not set")
	}
	return &_data, nil
}
