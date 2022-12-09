// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package datastreams

import (
	"github.com/pedrogao/vectorsql/src/datablocks"
)

type IteratorFunc func() (*datablocks.DataBlock, error)

type IteratorBlockInputStream struct {
	iteratorFunc IteratorFunc
}

func NewIteratorBlockInputStream(iterator IteratorFunc) IDataBlockInputStream {
	return &IteratorBlockInputStream{
		iteratorFunc: iterator,
	}
}

func (stream *IteratorBlockInputStream) Name() string {
	return "IteratorBlockInputStream"
}

func (stream *IteratorBlockInputStream) Read() (*datablocks.DataBlock, error) {
	return stream.iteratorFunc()
}

func (stream *IteratorBlockInputStream) Close() {}
