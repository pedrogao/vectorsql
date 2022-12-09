// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package transforms

import (
	"context"

	"github.com/pedrogao/vectorsql/src/base/xlog"
	"github.com/pedrogao/vectorsql/src/config"
	"github.com/pedrogao/vectorsql/src/sessions"
)

type TransformContext struct {
	ctx              context.Context
	log              *xlog.Log
	conf             *config.Config
	progressCallback func(values *sessions.ProgressValues)
}

func NewTransformContext(ctx context.Context, log *xlog.Log, conf *config.Config) *TransformContext {
	return &TransformContext{
		ctx:  ctx,
		log:  log,
		conf: conf,
	}
}

func (ctx *TransformContext) SetProgressCallback(fn func(pv *sessions.ProgressValues)) {
	ctx.progressCallback = fn
}
