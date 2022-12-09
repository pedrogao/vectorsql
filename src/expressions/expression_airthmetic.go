// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package expressions

import (
	"github.com/pedrogao/vectorsql/src/base/docs"
	"github.com/pedrogao/vectorsql/src/datavalues"
)

func ADD(left interface{}, right interface{}) IExpression {
	exprs := expressionsFor(left, right)
	return &BinaryExpression{
		name: "+",
		argumentNames: [][]string{
			{"left", "right"},
		},
		description: docs.Text("Returns the sum of the two arguments."),
		validate: OneOf(
			SameFamily(datavalues.FamilyInt),
			SameFamily(datavalues.FamilyFloat),
		),
		left:  exprs[0],
		right: exprs[1],
		updateFn: func(left datavalues.IDataValue, right datavalues.IDataValue) (datavalues.IDataValue, error) {
			return datavalues.Add(left, right)
		},
	}
}

func SUB(left interface{}, right interface{}) IExpression {
	exprs := expressionsFor(left, right)
	return &BinaryExpression{
		name: "-",
		argumentNames: [][]string{
			{"left", "right"},
		},
		description: docs.Text("Returns the difference between the two arguments."),
		validate: OneOf(
			SameFamily(datavalues.FamilyInt),
			SameFamily(datavalues.FamilyFloat),
		),
		left:  exprs[0],
		right: exprs[1],
		updateFn: func(left datavalues.IDataValue, right datavalues.IDataValue) (datavalues.IDataValue, error) {
			return datavalues.Sub(left, right)
		},
	}
}

func MUL(left interface{}, right interface{}) IExpression {
	exprs := expressionsFor(left, right)
	return &BinaryExpression{
		name: "*",
		argumentNames: [][]string{
			{"left", "right"},
		},
		description: docs.Text("Returns the dot product of the two arguments."),
		validate: OneOf(
			SameFamily(datavalues.FamilyInt),
			SameFamily(datavalues.FamilyFloat),
		),
		left:  exprs[0],
		right: exprs[1],
		updateFn: func(left datavalues.IDataValue, right datavalues.IDataValue) (datavalues.IDataValue, error) {
			return datavalues.Mul(left, right)
		},
	}
}

func DIV(left interface{}, right interface{}) IExpression {
	exprs := expressionsFor(left, right)
	return &BinaryExpression{
		name: "/",
		argumentNames: [][]string{
			{"left", "right"},
		},
		description: docs.Text("Returns the division of the two arguments."),
		validate: OneOf(
			SameFamily(datavalues.FamilyInt),
			SameFamily(datavalues.FamilyFloat),
		),
		left:  exprs[0],
		right: exprs[1],
		updateFn: func(left datavalues.IDataValue, right datavalues.IDataValue) (datavalues.IDataValue, error) {
			return datavalues.Div(left, right)
		},
	}
}
