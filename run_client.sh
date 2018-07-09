#!/bin/sh

# omit GO_PATH if predefined
export GOPATH=~/Projects/go/
export MODE=development
export LOG_LEVEL=debug

go run client.go