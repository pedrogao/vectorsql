// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package executors

import (
	"testing"

	"github.com/pedrogao/vectorsql/src/mocks"
	"github.com/pedrogao/vectorsql/src/planners"

	"github.com/stretchr/testify/assert"
)

func TestDropDatabaseExecutor(t *testing.T) {
	tests := []struct {
		name  string
		query string
		err   string
	}{
		{
			name:  "create",
			query: "create database db1",
		},
		{
			name:  "drop",
			query: "drop database db1",
		},
		{
			name:  "drop-not-exists",
			query: "drop database xxdb1",
			err:   "database:xxdb1 doesn't exists",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, cleanup := mocks.NewMock()
			defer cleanup()

			plan, err := planners.PlanFactory(test.query)
			assert.Nil(t, err)

			ctx := NewExecutorContext(mock.Ctx, mock.Log, mock.Conf, mock.Session)
			executor, err := ExecutorFactory(ctx, plan)
			assert.Nil(t, err)

			result, err := executor.Execute()
			if test.err != "" {
				assert.Equal(t, test.err, err.Error())
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}
