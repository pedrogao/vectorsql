// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package transforms

import (
	"sync"
	"time"

	"github.com/pedrogao/vectorsql/src/base/collections"
	"github.com/pedrogao/vectorsql/src/datablocks"
	"github.com/pedrogao/vectorsql/src/expressions"
	"github.com/pedrogao/vectorsql/src/planners"
	"github.com/pedrogao/vectorsql/src/processors"
	"github.com/pedrogao/vectorsql/src/sessions"

	"github.com/gammazero/workerpool"
)

type GroupBySelectionTransform struct {
	ctx            *TransformContext
	plan           *planners.SelectionPlan
	progressValues sessions.ProgressValues
	processors.BaseProcessor
}

func NewGroupBySelectionTransform(ctx *TransformContext, plan *planners.SelectionPlan) processors.IProcessor {
	return &GroupBySelectionTransform{
		ctx:           ctx,
		plan:          plan,
		BaseProcessor: processors.NewBaseProcessor("transform_groupby_selection"),
	}
}

func (t *GroupBySelectionTransform) Execute() {
	ctx := t.ctx
	plan := t.plan
	out := t.Out()
	defer out.Close()

	var mu sync.Mutex
	groupers := make([]*collections.HashMap, 0, 32)
	workerPool := workerpool.New(ctx.conf.Runtime.ParallelWorkerNumber)

	onNext := func(x interface{}) {
		switch y := x.(type) {
		case *datablocks.DataBlock:
			workerPool.Submit(func() {
				start := time.Now()
				grouper, err := y.GroupBySelectionByPlan(plan)
				if err != nil {
					out.Send(err)
					return
				}
				mu.Lock()
				groupers = append(groupers, grouper)
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
		final := collections.NewHashMap()
		for _, grouper := range groupers {
			iter := grouper.GetIterator()
			for {
				curKey, curVal, ok := iter.Next()
				if !ok {
					break
				}

				// Check.
				mergeVal, mergeHash, ok, err := final.Get(curKey)
				if err != nil {
					out.Send(err)
					return
				}

				// Merge state.
				if ok {
					curVal := curVal.([]expressions.IExpression)
					mergeVal := mergeVal.([]expressions.IExpression)
					for i := range mergeVal {
						if _, err := mergeVal[i].Merge(curVal[i]); err != nil {
							out.Send(err)
							return
						}
					}

				} else {
					if err := final.SetByHash(curKey, mergeHash, curVal); err != nil {
						out.Send(err)
						return
					}
				}
			}
		}

		// Final state.
		iter := final.GetIterator()
		for {
			_, val, ok := iter.Next()
			if !ok {
				break
			}
			if finalBlock, err := datablocks.BuildOneBlockFromExpressions(val.([]expressions.IExpression)); err != nil {
				out.Send(err)
			} else {
				out.Send(finalBlock)
			}
		}
	}
	t.Subscribe(onNext, onDone)
}

func (t *GroupBySelectionTransform) Stats() sessions.ProgressValues {
	return t.progressValues
}
