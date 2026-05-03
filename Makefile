#@IgnoreInspection BashAddShebang
ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))

build:
	go build -o $(ROOT)/bin/gib $(ROOT)/cmd/gib/main.go