// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package datavalues

import (
	"unsafe"

	"github.com/pedrogao/vectorsql/src/base/docs"
	"github.com/pedrogao/vectorsql/src/base/errors"
)

type ValueString string

func MakeString(v string) IDataValue {
	r := ValueString(v)
	return &r
}

func ZeroString() IDataValue {
	r := ValueString("")
	return &r
}

func (v *ValueString) Size() uintptr {
	return unsafe.Sizeof(v) + uintptr(len(*v))
}

func (v *ValueString) String() string {
	return string(*v)
}

func (v *ValueString) Type() Type {
	return TypeString
}

func (v *ValueString) Family() Family {
	return FamilyString
}

func (v *ValueString) AsString() string {
	return string(*v)
}

func (v *ValueString) Compare(other IDataValue) (Comparison, error) {
	if other.Type() != TypeString {
		return 0, errors.Errorf("type mismatch between values")
	}

	a := string(*v)
	b := AsString(other)
	switch {
	case a > b:
		return 1, nil
	case b > a:
		return -1, nil
	default:
		return 0, nil
	}
}

func (v *ValueString) Document() docs.Documentation {
	return docs.Text("String")
}

func AsString(v IDataValue) string {
	if t, ok := v.(*ValueString); ok {
		return string(*t)
	}
	return ""
}
