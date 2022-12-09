// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package datablocks

import (
	"github.com/pedrogao/vectorsql/src/datavalues"
	"github.com/pedrogao/vectorsql/src/expressions"
	"github.com/pedrogao/vectorsql/src/planners"
)

func (block *DataBlock) FilterByPlan(fields []string, plan *planners.FilterPlan) error {
	expr, err := planners.BuildExpression(plan.SubPlan)
	if err != nil {
		return err
	}

	i := 0
	params := make(expressions.Map)
	checks := make([]datavalues.IDataValue, block.NumRows())
	it, err := block.MixsIterator(fields)
	if err != nil {
		return err
	}
	for it.Next() {
		row := it.Value()
		for k := range row {
			params[it.Column(k).Name] = row[k]
		}
		v, err := expr.Update(params)
		if err != nil {
			return err
		}
		checks[i] = v
		i++
	}

	// In place filter.
	n := 0
	seqs := block.seqs
	for i, check := range checks {
		if datavalues.AsBool(check) {
			seqs[n] = seqs[i]
			n++
		}
	}
	block.mu.Lock()
	block.seqs = seqs[:n]
	block.mu.Unlock()
	return nil
}
