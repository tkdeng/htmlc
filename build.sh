#!/bin/bash

go get -u
go mod tidy

cd ./exec
go build -o ../htmlc
cd ../
