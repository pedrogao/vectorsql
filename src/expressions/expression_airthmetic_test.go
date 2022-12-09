// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package expressions

import (
	"testing"

	"github.com/pedrogao/vectorsql/src/datavalues"

	"github.com/stretchr/testify/assert"
)

func TestAirthmeticsExpression(t *testing.T) {
	tests := []struct {
		name      string
		expr      IExpression
		expect    datavalues.IDataValue
		errstring string
	}{
		{
			name:   "(1+2)",
			expr:   ADD(1, 2),
			expect: datavalues.ToValue(3),
		},
		{
			name:   "(a+b)",
			expr:   ADD("a", "b"),
			expect: datavalues.ToValue(3),
		},
		{
			name:   "(a+3)",
			expr:   ADD("a", 3),
			expect: datavalues.ToValue(4),
		},
		{
			name:   "(a+3)",
			expr:   ADD("a", CONST(3)),
			expect: datavalues.ToValue(4),
		},
		{
			name:   "(1+3)",
			expr:   ADD(CONST(1), CONST(3)),
			expect: datavalues.ToValue(4),
		},
		{
			name:   "(1-3)",
			expr:   SUB(CONST(1), CONST(3)),
			expect: datavalues.ToValue(-2),
		},
		{
			name:   "a+(1-3)",
			expr:   ADD("a", SUB(CONST(1), CONST(3))),
			expect: datavalues.ToValue(-1),
		},
		{
			name:   "a+(b*3)",
			expr:   ADD("a", MUL("b", 3)),
			expect: datavalues.ToValue(7),
		},
		{
			name:   "a/b",
			expr:   DIV("a", "b"),
			expect: datavalues.ToValue(0.5),
		},
		{
			name:      "a+c",
			expr:      ADD("a", "c"),
			errstring: "not-ok",
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			params := Map{
				"a": datavalues.ToValue(1),
				"b": datavalues.ToValue(2),
				"c": datavalues.MakeString("c"),
			}
			actual, err := test.expr.Update(params)
			if test.errstring != "" {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, test.expect, actual)

				err = test.expr.Walk(func(e IExpression) (bool, error) {
					return true, nil
				})
				assert.Nil(t, err)
			}
		})
	}
}

func TestAirthmeticsParamsExpression(t *testing.T) {
	tests := []struct {
		name      string
		expr      IExpression
		expect    datavalues.IDataValue
		errstring string
	}{
		{
			name:   "(1+2)",
			expr:   ADD(1, 2),
			expect: datavalues.ToValue(3),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.expr.Eval()
			if test.errstring != "" {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				actual := test.expr.Result()
				assert.Equal(t, test.expect, actual)
			}
		})
	}
}
