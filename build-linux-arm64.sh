#!/bin/bash

GOOS=linux GOARCH=arm64 go build -o info-server-arm64 main.go
