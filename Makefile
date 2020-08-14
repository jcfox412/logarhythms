SHELL := /bin/bash

build:
	go build github.com/jcfox412/logarhythms/cmd/logarhythms

run: build
	./logarhythms

lint: build
	# Lint code
	gometalinter ./...

test:
	# Running tests
	go test -race ./... -coverprofile cover.out

mocks:
	# Generate mocks for tests
	mockery -note @generated -case=underscore -name Manager -output internal/audio/mocks -dir internal/audio
