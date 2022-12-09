// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package executors

import (
	"fmt"

	"github.com/pedrogao/vectorsql/src/planners"
	"github.com/pedrogao/vectorsql/src/processors"
	"github.com/pedrogao/vectorsql/src/transforms"
)

type ProjectionExecutor struct {
	ctx         *ExecutorContext
	plan        *planners.ProjectionPlan
	transformer processors.IProcessor
}

func NewProjectionExecutor(ctx *ExecutorContext, plan *planners.ProjectionPlan) IExecutor {
	return &ProjectionExecutor{
		ctx:  ctx,
		plan: plan,
	}
}

func (executor *ProjectionExecutor) Execute() (*Result, error) {
	log := executor.ctx.log
	conf := executor.ctx.conf

	transformCtx := transforms.NewTransformContext(executor.ctx.ctx, log, conf)
	transform := transforms.NewProjectionTransform(transformCtx, executor.plan)
	executor.transformer = transform

	result := NewResult()
	result.SetInput(transform)
	return result, nil
}

func (executor *ProjectionExecutor) String() string {
	transformer := executor.transformer.(*transforms.ProjectionTransform)
	return fmt.Sprintf("(%v, stats:%+v)", transformer.Name(), transformer.Stats())
}
