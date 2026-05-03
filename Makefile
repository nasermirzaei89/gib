#@IgnoreInspection BashAddShebang
ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))

build:
	go build -o $(ROOT)/bin/engine $(ROOT)/cmd/engine/main.go