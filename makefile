build:
	@echo "--> Building..."
	go build -v -o bin/vectorsql-server src/cmd/server.go

goyacc:
	goyacc -o src/parsers/sqlparser/sql.go src/parsers/sqlparser/sql.y

clean:
	@echo "--> Cleaning..."
	@go clean
	@rm -f bin/*


test:
	@echo "--> Testing..."
	@$(MAKE) testbase
	@$(MAKE) testconfig
	@$(MAKE) testsessions
	@$(MAKE) testexpressions
	@$(MAKE) testprocessors
	@$(MAKE) testdatatypes
	@$(MAKE) testdatastreams
	@$(MAKE) testparsers
	@$(MAKE) testplanners
	@$(MAKE) testoptimizers
	@$(MAKE) testexecutors
	@$(MAKE) testtransforms

testbase:
	go test -v -race github.com/pedrogao/vectorsql/src/base/xlog
	go test -v -race github.com/pedrogao/vectorsql/src/base/lru
	go test -v -race github.com/pedrogao/vectorsql/src/base/metric

testconfig:
	go test -v -race github.com/pedrogao/vectorsql/src/config

testsessions:
	go test -v -race github.com/pedrogao/vectorsql/src/sessions

testprocessors:
	go test -v -race github.com/pedrogao/vectorsql/src/processors

testdatatypes:
	go test -v -race github.com/pedrogao/vectorsql/src/datatypes

testdatastreams:
	go test -v -race github.com/pedrogao/vectorsql/src/datastreams

testexpressions:
	go test -v -race github.com/pedrogao/vectorsql/src/expressions

testparsers:
	go test -v -race github.com/pedrogao/vectorsql/src/parsers/...

testplanners:
	go test -v -race github.com/pedrogao/vectorsql/src/planners

testoptimizers:
	go test -v -race github.com/pedrogao/vectorsql/src/optimizers

testexecutors:
	go test -v -race github.com/pedrogao/vectorsql/src/executors

testtransforms:
	go test -v -race github.com/pedrogao/vectorsql/src/transforms


pkgs =	github.com/pedrogao/vectorsql/src/config		\
		github.com/pedrogao/vectorsql/src/sessions	\
		github.com/pedrogao/vectorsql/src/processors	\
		github.com/pedrogao/vectorsql/src/parsers/...	\
		github.com/pedrogao/vectorsql/src/datatypes	\
		github.com/pedrogao/vectorsql/src/datastreams	\
		github.com/pedrogao/vectorsql/src/expressions	\
		github.com/pedrogao/vectorsql/src/planners	\
		github.com/pedrogao/vectorsql/src/optimizers	\
		github.com/pedrogao/vectorsql/src/executors	\
		github.com/pedrogao/vectorsql/src/transforms

coverage:
	go build -v -o bin/gotestcover \
	src/vendor/github.com/pierrre/gotestcover/*.go;
	bin/gotestcover -coverprofile=coverage.out -v $(pkgs)
	go tool cover -html=coverage.out

check:
	go get -v github.com/golangci/golangci-lint/cmd/golangci-lint
	bin/golangci-lint --skip-dirs github run src/... --skip-files sql.go

fmt:
	go fmt $(pkgs)


.PHONY: build clean install fmt test coverage
