#!/bin/bash
export PATH=$PATH:$(go env GOPATH)/bin
rm -rf docs/ && swag init -g app/main.go  
go run app/main.go