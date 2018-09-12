REPO=mojo
VERSION=latest

build:
	GOOS=linux GOARCH=amd64 go build -v
docker:
	docker build -t registry.cn-hangzhou.aliyuncs.com/$(REPO)/prometheus-adapter:$(VERSION) .

docker-push:
	docker push registry.cn-hangzhou.aliyuncs.com/$(REPO)/prometheus-adapter:$(VERSION)
run:
	docker run -it prometheus-adapter -influxdb-url=http://localhost:8080
all: build docker docker-push