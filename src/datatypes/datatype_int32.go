// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package datatypes

import (
	"fmt"
	"io"
	"reflect"

	"base/binary"
	"base/errors"
	"datavalues"
)

const (
	DataTypeInt32Name = "Int32"
)

type Int32DataType struct {
	DataTypeBase
}

func NewInt32DataType() IDataType {
	return &Int32DataType{}
}

func (datatype *Int32DataType) Type() reflect.Type {
	return reflect.ValueOf(int32(0)).Type()
}

func (datatype *Int32DataType) Name() string {
	return DataTypeInt32Name
}

func (datatype *Int32DataType) Serialize(writer *binary.Writer, v datavalues.IDataValue) error {
	if err := writer.Int32(int32(datavalues.AsInt(v))); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func (datatype *Int32DataType) SerializeText(writer io.Writer, v datavalues.IDataValue) error {
	if _, err := writer.Write([]byte(fmt.Sprintf("%d", int32(datavalues.AsInt(v))))); err != nil {
		return errors.Wrap(err)
	}
	return nil
}
