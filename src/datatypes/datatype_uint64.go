// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package datatypes

import (
	"fmt"
	"io"

	"github.com/pedrogao/vectorsql/src/base/binary"
	"github.com/pedrogao/vectorsql/src/base/errors"
	"github.com/pedrogao/vectorsql/src/datavalues"
)

const (
	DataTypeUInt64Name = "UInt64"
)

type UInt64DataType struct {
}

func NewUInt64DataType() IDataType {
	return &UInt64DataType{}
}

func (datatype *UInt64DataType) Name() string {
	return DataTypeUInt64Name
}

func (datatype *UInt64DataType) Serialize(writer *binary.Writer, v datavalues.IDataValue) error {
	return writer.UInt64(uint64(datavalues.AsInt(v)))
}

func (datatype *UInt64DataType) SerializeText(writer io.Writer, v datavalues.IDataValue) error {
	_, err := writer.Write([]byte(fmt.Sprintf("%d", uint64(datavalues.AsInt(v)))))
	return err
}

func (datatype *UInt64DataType) Deserialize(reader *binary.Reader) (datavalues.IDataValue, error) {
	if res, err := reader.UInt64(); err != nil {
		return nil, errors.Wrap(err)
	} else {
		return datavalues.ToValue(res), nil
	}
}
