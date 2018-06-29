build:
	GOOS=linux GOARCH=amd64 go build -v
docker:
	docker build -t registry.cn-hangzhou.aliyuncs.com/mojo/prometheus-adapter .

docker-push:
	docker push registry.cn-hangzhou.aliyuncs.com/mojo/prometheus-adapter
run:
	docker run -it prometheus-adapter -influxdb-url=http://localhost:8080
all: build docker docker-push