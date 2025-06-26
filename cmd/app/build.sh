#!/bin/sh

rm linux-amd64/sd-push-linux-amd64

env GOOS=linux GOARCH=amd64 go build -o linux-amd64/sd-push-linux-amd64 -ldflags '-s -w'
#env GOOS=linux GOARCH=arm64 go build -o bin/sd-push-linux-arm64
#env GOOS=netbsd GOARCH=amd64 go build -o bin/sd-push-netbsd-amd64
#env GOOS=netbsd GOARCH=arm64 go build -o bin/sd-push-netbsd-arm64
#env GOOS=openbsd GOARCH=amd64 go build -o bin/sd-push-openbsd-amd64
#env GOOS=openbsd GOARCH=arm64 go build -o bin/sd-push-openbsd-arm64


scp linux-amd64/sd-push-linux-amd64 bjones@10.0.0.205:/home/bjones/easy-diffusion/

