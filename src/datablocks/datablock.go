// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package datablocks

import (
	"fmt"
	"sync"

	"github.com/pedrogao/vectorsql/src/base/errors"
	"github.com/pedrogao/vectorsql/src/columns"
	"github.com/pedrogao/vectorsql/src/datavalues"
)

type DataBlock struct {
	mu         sync.RWMutex
	seqs       []int
	info       *DataBlockInfo
	values     []*DataBlockValue
	totalBytes uint64
}

func NewDataBlock(cols []*columns.Column) *DataBlock {
	block := &DataBlock{
		info:   &DataBlockInfo{},
		seqs:   make([]int, 0),
		values: make([]*DataBlockValue, len(cols)),
	}

	for i, col := range cols {
		cv := NewDataBlockValue(col)
		block.values[i] = cv
	}
	return block
}

func newDataBlock(seqs []int, values []*DataBlockValue) *DataBlock {
	return &DataBlock{
		info:   &DataBlockInfo{},
		seqs:   seqs,
		values: values,
	}
}

func (block *DataBlock) Clone() *DataBlock {
	return NewDataBlock(block.Columns())
}

func (block *DataBlock) DeepClone() *DataBlock {
	clone := NewDataBlock(block.Columns())
	clone.totalBytes = block.totalBytes
	clone.seqs = make([]int, len(block.seqs))
	copy(clone.seqs, block.seqs)

	for i, value := range block.values {
		clone.values[i] = value.DeepClone()
	}
	return clone
}

func (block *DataBlock) Info() *DataBlockInfo {
	return block.info
}

func (block *DataBlock) TotalBytes() uint64 {
	block.mu.RLock()
	defer block.mu.RUnlock()
	return block.totalBytes
}

func (block *DataBlock) NumRows() int {
	block.mu.RLock()
	defer block.mu.RUnlock()
	return len(block.seqs)
}

func (block *DataBlock) NumColumns() int {
	return len(block.values)
}

func (block *DataBlock) Columns() []*columns.Column {
	cols := make([]*columns.Column, len(block.values))
	for i, cv := range block.values {
		cols[i] = cv.column
	}
	return cols
}

func (block *DataBlock) Column(name string) (*columns.Column, error) {
	for _, cv := range block.values {
		if cv.column.Name == name {
			return cv.column, nil
		}
	}
	return nil, errors.Errorf("Can't find column:%v", name)
}

func (block *DataBlock) DataBlockValue(name string) (*DataBlockValue, error) {
	for _, cv := range block.values {
		if cv.column.Name == name {
			return cv, nil
		}
	}
	return nil, errors.Errorf("Can't find column:%v", name)
}

func (block *DataBlock) RowIterator() *DataBlockRowIterator {
	return newDataBlockRowIterator(block)
}

func (block *DataBlock) ColumnIterator(name string) (*DataBlockColumnIterator, error) {
	for i, v := range block.values {
		if v.column.Name == name {
			return newDataBlockColumnIterator(block, i), nil
		}
	}
	return nil, errors.Errorf("Can't find column:%v", name)
}

func (block *DataBlock) ColumnIterators() []*DataBlockColumnIterator {
	var iterators []*DataBlockColumnIterator

	for i := range block.values {
		iter := newDataBlockColumnIterator(block, i)
		iterators = append(iterators, iter)
	}
	return iterators
}

func (block *DataBlock) MixsIterator(columns []string) (*DataBlockMixsIterator, error) {
	return newDataBlockColsRowIterator(block, columns)
}

func (block *DataBlock) WriteRow(values []datavalues.IDataValue) error {
	cols := block.NumColumns()
	if len(values) != cols {
		return errors.Errorf("Can't append row, expect column length:%v", cols)
	}

	offset := len(block.values[0].values)
	for i := 0; i < cols; i++ {
		block.totalBytes += uint64(values[i].Size())
		block.values[i].values = append(block.values[i].values, values[i])
	}
	block.seqs = append(block.seqs, offset)
	return nil
}

func (block *DataBlock) Close() {
	block.seqs = nil
	block.values = nil
	block.totalBytes = 0
}

func (block *DataBlock) Dump() {
	header := "\n| "
	for _, cv := range block.values {
		header += fmt.Sprintf("%v(%T)", cv.column.Name, cv.column.DataType)
		header += " | "
	}
	fmt.Print(header)

	body := "\n"
	it := block.RowIterator()
	for it.Next() {
		row := it.Value()
		for _, v := range row {
			body += v.String()
			body += " "
		}
		body += "\n"
	}
	fmt.Print(body)
}
