// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package transforms

import (
	"context"
	"testing"

	"github.com/pedrogao/vectorsql/src/columns"
	"github.com/pedrogao/vectorsql/src/datablocks"
	"github.com/pedrogao/vectorsql/src/datatypes"
	"github.com/pedrogao/vectorsql/src/mocks"
	"github.com/pedrogao/vectorsql/src/planners"
	"github.com/pedrogao/vectorsql/src/processors"

	"github.com/stretchr/testify/assert"
)

func TestLimitTransfrom(t *testing.T) {
	tests := []struct {
		name   string
		plan   planners.IPlan
		source []interface{}
		expect *datablocks.DataBlock
	}{
		{
			name: "simple",
			plan: planners.NewLimitPlan(
				planners.NewConstantPlan(1),
				planners.NewConstantPlan(2),
			),
			source: mocks.NewSourceFromSlice(
				mocks.NewBlockFromSlice(
					[]*columns.Column{
						{Name: "name", DataType: datatypes.NewStringDataType()},
						{Name: "age", DataType: datatypes.NewInt32DataType()},
					},
					[]interface{}{"x", 11},
					[]interface{}{"z", 13},
					[]interface{}{"y", 12},
					[]interface{}{"y", 13},
				)),
			expect: mocks.NewBlockFromSlice(
				[]*columns.Column{
					{Name: "name", DataType: datatypes.NewStringDataType()},
					{Name: "age", DataType: datatypes.NewInt32DataType()},
				},
				[]interface{}{"z", 13},
				[]interface{}{"y", 12},
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, cleanup := mocks.NewMock()
			defer cleanup()
			ctx := NewTransformContext(mock.Ctx, mock.Log, mock.Conf)

			stream := mocks.NewMockBlockInputStream(test.source)
			datasource := NewDataSourceTransform(ctx, stream)

			limit := NewLimitransform(ctx, test.plan.(*planners.LimitPlan))

			sink := processors.NewSink("sink")
			pipeline := processors.NewPipeline(context.Background())
			pipeline.Add(datasource)
			pipeline.Add(limit)
			pipeline.Add(sink)
			pipeline.Run()

			err := pipeline.Wait(func(x interface{}) error {
				actual := x.(*datablocks.DataBlock)
				expect := test.expect
				assert.True(t, mocks.DataBlockEqual(actual, expect))
				return nil
			})
			assert.Nil(t, err)
			stats := limit.(*Limitransform).Stats()
			assert.True(t, stats.TotalRowsToRead.Get() > 0)
		})
	}
}
