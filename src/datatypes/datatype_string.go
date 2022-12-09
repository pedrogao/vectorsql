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

const (
	DataTypeStringName = "String"
)

type StringDataType struct {
}

func NewStringDataType() IDataType {
	return &StringDataType{}
}

func (datatype *StringDataType) Name() string {
	return DataTypeStringName
}

func (datatype *StringDataType) Serialize(writer *binary.Writer, v datavalues.IDataValue) error {
	return writer.String(datavalues.AsString(v))
}

func (datatype *StringDataType) SerializeText(writer io.Writer, v datavalues.IDataValue) error {
	_, err := writer.Write([]byte(datavalues.AsString(v)))
	return err
}

func (datatype *StringDataType) Deserialize(reader *binary.Reader) (datavalues.IDataValue, error) {
	if res, err := reader.String(); err != nil {
		return nil, errors.Wrap(err)
	} else {
		return datavalues.MakeString(res), nil
	}
}
