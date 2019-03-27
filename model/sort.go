package model

import "errors"

type Sort struct {
	Key       string
	Direction sortDirection
}

const (
	ASC sortDirection = iota
	DESC
)

type sortDirection byte

func (sd sortDirection) String() string {
	switch sd {
	case ASC:
		return "asc"
	case DESC:
		return "desc"
	}
	panic(errors.New("type sortDirection parsed error"))
}
