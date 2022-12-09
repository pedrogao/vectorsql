// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package executors

import (
	"github.com/pedrogao/vectorsql/src/planners"
	"github.com/pedrogao/vectorsql/src/processors"
)

type SinkExecutor struct {
	ctx  *ExecutorContext
	plan *planners.SinkPlan
}

func NewSinkExecutor(ctx *ExecutorContext, plan *planners.SinkPlan) IExecutor {
	return &SinkExecutor{
		ctx:  ctx,
		plan: plan,
	}
}

func (executor *SinkExecutor) Execute() (*Result, error) {
	proc := processors.NewSink("transforms_sink")

	result := NewResult()
	result.SetInput(proc)
	return result, nil
}

func (executor *SinkExecutor) String() string {
	return ""
}
