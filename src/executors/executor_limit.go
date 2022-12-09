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

type LimitExecutor struct {
	ctx         *ExecutorContext
	plan        *planners.LimitPlan
	transformer processors.IProcessor
}

func NewLimitExecutor(ctx *ExecutorContext, plan *planners.LimitPlan) IExecutor {
	return &LimitExecutor{
		ctx:  ctx,
		plan: plan,
	}
}

func (executor *LimitExecutor) Execute() (*Result, error) {
	log := executor.ctx.log
	conf := executor.ctx.conf

	transformCtx := transforms.NewTransformContext(executor.ctx.ctx, log, conf)
	transform := transforms.NewLimitransform(transformCtx, executor.plan)
	executor.transformer = transform

	result := NewResult()
	result.SetInput(transform)
	return result, nil
}

func (executor *LimitExecutor) String() string {
	transformer := executor.transformer.(*transforms.Limitransform)
	return fmt.Sprintf("(%v, stats:%+v)", transformer.Name(), transformer.Stats())
}
