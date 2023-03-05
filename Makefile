VERSION=`git describe --tags --dirty --always`
COMMIT=`git rev-parse HEAD`
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION} -X main.commit=${COMMIT}"

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: lint
lint:
	golangci-lint run -v ./...

.PHONY: build
build:
	CGO_ENABLED=0 go build ${LDFLAGS} -o bin/unkatan cmd/unkatan/main.go

.PHONY: build_linux
build_linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ${LDFLAGS} -o bin/unkatan cmd/unkatan/main.go

.PHONY: run_tests
run_tests: env_up
	GOPATH=`go env GOPATH` docker-compose -f docker-compose.yml -f docker-compose.ci.yml run tests

.PHONY: cover
cover: run_tests
	go tool cover -html=profile.cov -o oc_cover.html
	rm profile.cov

.PHONY: run
run: build
	bin/unkatan --config=./unkatan.yml

.PHONY: to_dev
to_dev: build_linux
	cat bin/unkatan | ssh "${DEV_SSH}" "cat - > unkatan && chmod +x unkatan;"