#!/bin/bash

echo "Clear old builds"
rm -rf ./build

echo "Build for windows..."
export GOOS=windows
export GOARCH=386
go get -v
go build -o ./build/webm-api.windows.exe ServerEntry.go

echo "Build for linux..."
export GOOS=linux
export GOARCH=386
go get -v
go build -o ./build/webm-api.linux ServerEntry.go