// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package datatypes

import (
	"io"

	"github.com/pedrogao/vectorsql/src/base/binary"
	"github.com/pedrogao/vectorsql/src/base/errors"
	"github.com/pedrogao/vectorsql/src/datavalues"
)

type IDataType interface {
	Name() string
	Serialize(*binary.Writer, datavalues.IDataValue) error
	SerializeText(io.Writer, datavalues.IDataValue) error
	Deserialize(*binary.Reader) (datavalues.IDataValue, error)
}

func GetDataTypeByValue(val datavalues.IDataValue) (IDataType, error) {
	switch val.Type() {
	case datavalues.TypeString:
		return NewStringDataType(), nil
	case datavalues.TypeInt:
		return NewInt64DataType(), nil
	case datavalues.TypeInt32:
		return NewInt32DataType(), nil
	case datavalues.TypeFloat:
		return NewFloat64DataType(), nil
	default:
		return nil, errors.Errorf("Unsupported value type:%v", val.Type())
	}
}
