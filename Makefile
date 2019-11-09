SOURCEDIR = .
REPO=gitlab.com/krinklesaurus
DOCKER_REGISTRY=registry.${REPO}
NAME ?= go-p2pb2b
GO_VERSION=1.13.4
GO_RUN=docker run --rm -v ${PWD}:/usr/src/myapp -w /usr/src/myapp golang:${GO_VERSION}
VERSION=$(shell git rev-parse --short HEAD)

default: clean test build

.PHONY: init
init:
	export GO111MODULE=on &&\
		${GO_RUN} go mod init ${REPO}/${NAME}


.PHONY: clean
clean:
	@if [ -f ${NAME} ] ; then rm ${NAME}; fi &&\
		rm -rf vendor/


# validate the project is correct and all necessary information is available
.PHONY: validate
validate:
	export GO111MODULE=on &&\
		${GO_RUN} go mod tidy &&\
		${GO_RUN} go mod vendor &&\
		${GO_RUN} go mod verify &&\
		${GO_RUN} golint $$(${GO_RUN} go list ./... | grep -v /vendor/)


.PHONY: test
test:
	${GO_RUN} go vet -v $$(${GO_RUN} go list ./... | grep -v /vendor/) &&\
	${GO_RUN} go test -v -race -cover $$(${GO_RUN} go list ./... | grep -v /vendor/)


.PHONY: build
build:
	${GO_RUN} go build -v