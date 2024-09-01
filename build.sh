#!/bin/bash

go build -ldflags "-w -s" ./
#upx --best gobuilder

GOOS=windows go build -ldflags "-w -s" ./
#upx --best gobuilder.exe