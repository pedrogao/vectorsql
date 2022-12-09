// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package storages

import (
	"github.com/pedrogao/vectorsql/src/columns"
	"github.com/pedrogao/vectorsql/src/datastreams"
	"github.com/pedrogao/vectorsql/src/sessions"
)

type IStorage interface {
	Name() string
	Columns() []*columns.Column
	GetInputStream(*sessions.Session) (datastreams.IDataBlockInputStream, error)
	GetOutputStream(*sessions.Session) (datastreams.IDataBlockOutputStream, error)
	Close()
}
