#!/bin/sh
version=v2.3

GOOS=linux GOARCH=amd64 go build -o remote_storage_adapter
docker build -t registry.cn-hangzhou.aliyuncs.com/wise2c-test/prometheus-adapter:$version -t registry.cn-hangzhou.aliyuncs.com/wise2c-dev/prometheus-adapter:$version .

docker push registry.cn-hangzhou.aliyuncs.com/wise2c-test/prometheus-adapter:$version
docker push registry.cn-hangzhou.aliyuncs.com/wise2c-dev/prometheus-adapter:$version