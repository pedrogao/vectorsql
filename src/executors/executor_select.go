// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package executors

import (
	"github.com/pedrogao/vectorsql/src/base/errors"
	"github.com/pedrogao/vectorsql/src/planners"
)

type SelectExecutor struct {
	ctx  *ExecutorContext
	plan *planners.SelectPlan
	tree *ExecutorTree
}

func NewSelectExecutor(ctx *ExecutorContext, plan planners.IPlan) IExecutor {
	return &SelectExecutor{
		ctx:  ctx,
		tree: NewExecutorTree(ctx),
		plan: plan.(*planners.SelectPlan),
	}
}

func (executor *SelectExecutor) Execute() (*Result, error) {
	ectx := executor.ctx
	tree := executor.tree

	children := executor.plan.SubPlan.SubPlans

	for _, plan := range children {
		switch plan := plan.(type) {
		case *planners.TableValuedFunctionPlan:
			executor := NewTableValuedFunctionExecutor(ectx, plan)
			tree.Add(executor)
		case *planners.ScanPlan:
			executor := NewScanExecutor(ectx, plan)
			tree.Add(executor)
		case *planners.FilterPlan:
			executor := NewFilterExecutor(ectx, plan)
			tree.Add(executor)
		case *planners.SelectionPlan:
			executor := NewSelectionExecutor(ectx, plan)
			tree.Add(executor)
		case *planners.OrderByPlan:
			executor := NewOrderByExecutor(ectx, plan)
			tree.Add(executor)
		case *planners.LimitPlan:
			executor := NewLimitExecutor(ectx, plan)
			tree.Add(executor)
		case *planners.ProjectionPlan:
			executor := NewProjectionExecutor(ectx, plan)
			tree.Add(executor)
		case *planners.SinkPlan:
			executor := NewSinkExecutor(ectx, plan)
			tree.Add(executor)
		default:
			return nil, errors.Errorf("Unsupported plan:%T", plan)
		}
	}
	pipeline, err := tree.BuildPipeline()
	if err != nil {
		return nil, err
	}
	pipeline.Run()

	result := NewResult()
	result.SetInput(pipeline.Last())
	return result, nil
}

func (executor *SelectExecutor) String() string {
	res := ""
	for _, t := range executor.tree.subExecutors {
		res += t.String()
		res += " -> "
	}
	return res
}
