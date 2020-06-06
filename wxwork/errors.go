package wxwork

import (
	"errors"
)

var (
	ErrEmptyCorp  = errors.New("empty corpID or corpSecret")
	ErrEmptyArg   = errors.New("empty argument")
	ErrNotFound   = errors.New("not found")
	ErrOutofRange = errors.New("out of range")
)
