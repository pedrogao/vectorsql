// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package datatypes

import (
	"github.com/pedrogao/vectorsql/src/base/errors"
)

type dataTypeCreator func() IDataType

var (
	table = map[string]dataTypeCreator{
		NewStringDataType().Name():  NewStringDataType,
		NewInt32DataType().Name():   NewInt32DataType,
		NewUInt32DataType().Name():  NewUInt32DataType,
		NewInt64DataType().Name():   NewInt64DataType,
		NewUInt64DataType().Name():  NewUInt64DataType,
		NewFloat64DataType().Name(): NewFloat64DataType,
	}
)

func DataTypeFactory(name string) (IDataType, error) {
	dt, ok := table[name]
	if !ok {
		return nil, errors.Errorf("Unsupported data type:%s", name)
	}
	return dt(), nil
}
