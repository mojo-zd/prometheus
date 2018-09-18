# Remote storage adapter

This is a write adapter that receives samples via Prometheus's remote write
protocol and stores them in Graphite, InfluxDB, or OpenTSDB. It is meant as a
replacement for the built-in specific remote storage implementations that have
been removed from Prometheus.

For InfluxDB, this binary is also a read adapter that supports reading back
data through Prometheus via Prometheus's remote read protocol.

## Building

```
go build
```

## Running

Graphite example:

```
./remote_storage_adapter -graphite-address=localhost:8080
```

OpenTSDB example:

```
./remote_storage_adapter -opentsdb-url=http://localhost:8081/
```

InfluxDB example:

```
./remote_storage_adapter -influxdb-url=http://localhost:8086/ -influxdb.database=prometheus -influxdb.retention-policy=autogen
```

To show all flags:

```
./remote_storage_adapter -h
```

## Configuring Prometheus

To configure Prometheus to send samples to this binary, add the following to your `prometheus.yml`:

```yaml
# Remote write configuration (for Graphite, OpenTSDB, or InfluxDB).
remote_write:
  - url: "http://localhost:9201/write"

# Remote read configuration (for InfluxDB only at the moment).
remote_read:
  - url: "http://localhost:9201/read"
```

# Build docker image
```
cd prometheus/documentation/examples/remote_storage/remote_storage_adapter
sh publish.sh
```

# K8S Deployment with Opentsdb Adapter
```
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: prometheus-adapter
  labels:
    app: prometheus-adapter
spec:
  replicas: 1
  selector:
    matchLabels:
      app: prometheus-adapter
  template:
    metadata:
      labels:
        app: prometheus-adapter
    spec:
      containers:
      - name: adapter
        image: registry.cn-hangzhou.aliyuncs.com/wise2c-test/prometheus-adapter:v2.3
        args:
        - -opentsdb-url=http://10.0.0.131:4242
```

```
influxdb args

- -influxdb-url=http://localhost:8086
- -influxdb.username=prom
- -influxdb.database=prometheus
- -influxdb.retention-policy=autogen
```