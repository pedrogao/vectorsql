// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package mocks

import (
	"context"
	"os"
	"sync"

	"github.com/pedrogao/vectorsql/src/base/xlog"
	"github.com/pedrogao/vectorsql/src/config"
	"github.com/pedrogao/vectorsql/src/databases"
	"github.com/pedrogao/vectorsql/src/sessions"
)

var once sync.Once

type Mock struct {
	Log     *xlog.Log
	Conf    *config.Config
	Session *sessions.Session
	Ctx     context.Context
	Cancel  context.CancelFunc
}

func NewMock() (*Mock, func()) {
	log := xlog.NewStdLog(xlog.Level(xlog.ERROR))
	conf := config.DefaultConfig()
	session := sessions.NewSession()
	ctx, cancel := context.WithCancel(context.Background())

	once.Do(func() {
		if err := databases.Load(log, conf); err != nil {
			log.Panic("%+v", err)
		}
	})

	mock := &Mock{
		Log:     log,
		Conf:    conf,
		Session: session,
		Ctx:     ctx,
		Cancel:  cancel,
	}
	return mock, func() {
		cancel()
		os.RemoveAll(conf.Server.Path)
	}
}
