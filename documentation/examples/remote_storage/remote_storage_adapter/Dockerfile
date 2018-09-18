FROM busybox
MAINTAINER excessivespeed@126.com

COPY ./remote_storage_adapter /go/bin/remote_storage_adapter

WORKDIR /go/bin/
ENTRYPOINT ["/go/bin/remote_storage_adapter"]
EXPOSE 9201