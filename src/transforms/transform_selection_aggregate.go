// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package transforms

import (
	"sync"
	"time"

	"github.com/pedrogao/vectorsql/src/datablocks"
	"github.com/pedrogao/vectorsql/src/expressions"
	"github.com/pedrogao/vectorsql/src/planners"
	"github.com/pedrogao/vectorsql/src/processors"
	"github.com/pedrogao/vectorsql/src/sessions"

	"github.com/gammazero/workerpool"
)

type AggregateSelectionTransform struct {
	ctx            *TransformContext
	plan           *planners.SelectionPlan
	progressValues sessions.ProgressValues
	processors.BaseProcessor
}

func NewAggregateSelectionTransform(ctx *TransformContext, plan *planners.SelectionPlan) processors.IProcessor {
	return &AggregateSelectionTransform{
		ctx:           ctx,
		plan:          plan,
		BaseProcessor: processors.NewBaseProcessor("transform_aggregate_selection"),
	}
}

func (t *AggregateSelectionTransform) Execute() {
	ctx := t.ctx
	out := t.Out()
	defer out.Close()
	plan := t.plan

	// Get all base fields by the expression.
	fields, err := planners.BuildVariableValues(plan.Projects)
	if err != nil {
		out.Send(err)
		return
	}

	var mu sync.Mutex
	var exprs [][]expressions.IExpression
	workerPool := workerpool.New(ctx.conf.Runtime.ParallelWorkerNumber)

	onNext := func(x interface{}) {
		switch y := x.(type) {
		case *datablocks.DataBlock:
			workerPool.Submit(func() {
				start := time.Now()
				expr, err := y.AggregateSelectionByPlan(fields, plan)
				if err != nil {
					out.Send(err)
					return
				}
				mu.Lock()
				exprs = append(exprs, expr)
				mu.Unlock()

				cost := time.Since(start)
				t.progressValues.Cost.Add(cost)
				t.progressValues.ReadBytes.Add(int64(y.TotalBytes()))
				t.progressValues.ReadRows.Add(int64(y.NumRows()))
				t.progressValues.TotalRowsToRead.Add(int64(y.NumRows()))
			})
		case error:
			out.Send(y)
		}
	}
	onDone := func() {
		workerPool.StopWait()
		if len(exprs) > 0 {
			var mergeExpr []expressions.IExpression
			// Do merge.
			for i, expr := range exprs {
				if i == 0 {
					mergeExpr = expr
					continue
				}
				for i := range mergeExpr {
					if _, err := mergeExpr[i].Merge(expr[i]); err != nil {
						out.Send(err)
						return
					}
				}
			}
			if merger, err := datablocks.BuildOneBlockFromExpressions(mergeExpr); err != nil {
				out.Send(err)
			} else {
				out.Send(merger)
			}
		}
	}
	t.Subscribe(onNext, onDone)
}

func (t *AggregateSelectionTransform) Stats() sessions.ProgressValues {
	return t.progressValues
}
