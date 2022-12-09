// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package datavalues

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pedrogao/vectorsql/src/base/lru"
)

var (
	re      = regexp.MustCompile(`([^\\]?|[\\]{2})[%_]`)
	reCache = lru.New(65536)
)

func replacer(s string) string {
	if strings.HasPrefix(s, `\\`) {
		return s[2:]
	}

	result := strings.Replace(s, "%", ".*", -1)
	result = strings.Replace(result, "_", ".", -1)
	return result
}

func LikeToRegexp(likeExpr string) *regexp.Regexp {
	if likeExpr == "" {
		return regexp.MustCompile("^.*$") // Can never fail
	}

	keyPattern := regexp.QuoteMeta(likeExpr)
	keyPattern = re.ReplaceAllStringFunc(keyPattern, replacer)
	keyPattern = fmt.Sprintf("^%s$", keyPattern)
	return regexp.MustCompile(keyPattern) // Can never fail
}

func Like(likeExpr string, x IDataValue) bool {
	var re *regexp.Regexp

	if rex, ok := reCache.Get(likeExpr); ok {
		re = rex.(*regexp.Regexp)
	} else {
		re = LikeToRegexp(likeExpr)
		reCache.Add(likeExpr, re)
	}
	return re.Match([]byte(x.String()))
}
