#!/usr/bin/env sh

set -e 

go fmt ./... && \
    goimports -w . && \
    go vet ./... 
