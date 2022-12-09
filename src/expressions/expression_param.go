// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package expressions

import (
	"github.com/pedrogao/vectorsql/src/datavalues"
)

type IParams interface {
	Get(name string) (datavalues.IDataValue, bool)
}

type Map map[string]datavalues.IDataValue

func (p Map) Get(name string) (datavalues.IDataValue, bool) {
	v, ok := p[name]
	return v, ok
}
