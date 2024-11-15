package utils

import (
	"strconv"
	"strings"

	"github.com/kballard/go-shellquote"
)

type EnvBuilder struct {
	val   map[string]string
	order []string
}

func (e *EnvBuilder) Set(k, v string) {
	v = shellquote.Join(v)

	_, ok := e.val[k]

	if !ok {
		e.order = append(e.order, k)
	}

	e.val[k] = v
}

func (e *EnvBuilder) SetMap(m map[string]string) {
	for k, v := range m {
		e.Set(k, v)
	}
}

func (e *EnvBuilder) Map() map[string]string {
	return e.val
}

func (e *EnvBuilder) Args() []string {
	res := []string{}

	for _, k := range e.order {
		res = append(res, k+"="+e.val[k])
	}

	return res
}

func (e *EnvBuilder) String() string {
	return strings.Join(e.Args(), " ")
}

func (e *EnvBuilder) SetP(k string, v *string) {
	if v == nil {
		return
	}

	e.Set(k, *v)
}

func (e *EnvBuilder) SetInt64P(k string, v *int64) {
	if v == nil {
		return
	}

	e.Set(k, strconv.FormatInt(*v, 10))
}

func (e *EnvBuilder) SetInt(k string, v *int) {
	if v == nil {
		return
	}

	e.Set(k, strconv.FormatInt(int64(*v), 10))
}

func (e *EnvBuilder) SetBoolP(k string, v *bool) {
	if v == nil {
		return
	}

	if *v {
		e.Set(k, "true")
	} else {
		e.Set(k, "false")
	}
}

func NewEnvBuilder() *EnvBuilder {
	b := &EnvBuilder{
		val: make(map[string]string),
	}

	return b
}
