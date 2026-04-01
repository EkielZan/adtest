#!/bin/env bash
if [[ $1 == 1 ]];then
force="-a"
fi

export CGO_ENABLED=1
export GOOS=linux
export GOARCH=amd64

cd src
echo "Build Main Binaries"
go build -installsuffix cgo $force -v -o ../bin/adtest .
cd ..
