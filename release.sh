#!/bin/sh
# linux static + UPX
CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -mod=vendor -v -trimpath -ldflags "-s" -a -o chill && upx chill
# rpi 2/3
GOARCH=arm GOARM=6 GOOS=linux go build -mod=vendor -o chill.rpi2
