// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package executors

import (
	"github.com/pedrogao/vectorsql/src/databases"
	"github.com/pedrogao/vectorsql/src/planners"
)

type CreateDatabaseExecutor struct {
	ctx  *ExecutorContext
	plan *planners.CreateDatabasePlan
}

func NewCreateDatabaseExecutor(ctx *ExecutorContext, plan planners.IPlan) IExecutor {
	return &CreateDatabaseExecutor{
		ctx:  ctx,
		plan: plan.(*planners.CreateDatabasePlan),
	}
}

func (executor *CreateDatabaseExecutor) Execute() (*Result, error) {
	ectx := executor.ctx
	ast := executor.plan.Ast

	databaseCtx := databases.NewDatabaseContext(ectx.log, ectx.conf)
	database, err := databases.DatabaseFactory(databaseCtx, ast)
	if err != nil {
		return nil, err
	}
	if err := database.Executor().CreateDatabase(); err != nil {
		return nil, err
	}

	result := NewResult()
	return result, nil
}

func (executor *CreateDatabaseExecutor) String() string {
	return ""
}
