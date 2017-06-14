#!/bin/bash
protoc -I=. --gofast_out=plugins=grpc:. hello.proto