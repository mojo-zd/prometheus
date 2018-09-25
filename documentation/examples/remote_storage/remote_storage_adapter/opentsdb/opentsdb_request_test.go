package opentsdb

import (
	"testing"
	"encoding/json"
	"context"
	"time"
	"bytes"
	"golang.org/x/net/context/ctxhttp"
	"net/http"
	"github.com/kr/pretty"
	"io/ioutil"
	"github.com/prometheus/prometheus/prompb"
	"github.com/go-kit/kit/log"
	"github.com/prometheus/common/promlog"
)

var (
	putURL   = "http://10.0.0.131:4242/api/put"
	queryURL = "http://10.0.0.131:4242//api/query"
)

type Query struct {
	Start   int64                  `json:"start"`
	End     int64                  `json:"end"`
	Queries []*StoreSamplesRequest `json:"queries"`
}

func TestOpentsdbWrite(t *testing.T) {
	by := Unmarshal()
	pretty.Println("post data", string(by))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	resp, err := ctxhttp.Post(ctx, http.DefaultClient, putURL, contentTypeJSON, bytes.NewBuffer(by))
	if err != nil {
		t.Error(err.Error())
		return
	}

	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	pretty.Println("http status code", string(b))
}

func Unmarshal() []byte {
	req := []StoreSamplesRequest{
		{
			Metric:    "go.gc.duration.seconds",
			Timestamp: 1537170046,
			Value:     17,
			Tags: map[string]TagValue{
				"job":       TagValue("prometheus"),
				"monitor":   TagValue("codelab-monitor"),
				"quantile":  TagValue("0"),
				"instance":  TagValue("localhost:9099"),
				"job1":      TagValue("prometheus"),
				"monitor1":  TagValue("codelab-monitor"),
				"quantile1": TagValue("0"),
				"instance1": TagValue("localhost:9099"),
				"instance2": TagValue("localhost:9099"),
			},
		},
	}

	b, _ := json.Marshal(&req)
	return b
}

func TestOpentsdbReader(t *testing.T) {
	//query := &otdbQueryReq{
	//	Start: 1536832936,
	//	End:   1537006036,
	//	Queries: []otdbQuery{
	//		{
	//			Metric: "go_gc_duration_seconds",
	//			Filters: []otdbFilter{
	//				{Type: "literal_or", Tagk: "monitor", Filter: "codelab-monitor", GroupBy: true},
	//			},
	//			Aggregator: "sum",
	//		},
	//	},
	//}
	//q, _ := json.Marshal(query)
	q := `{
  "start": 1536832936,
  "end": 1537006036,
  "queries": [
    {
      "metric": "go__gc__duration__seconds",
      "filters": [
        {
          "type": "literal_or",
          "tagk": "monitor",
          "filter": "codelab-monitor",
          "groupBy": true
        }
      ],
      "aggregator": "sum"
    }
  ]
}`
	//pretty.Println(string(q))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	resp, err := ctxhttp.Post(ctx, http.DefaultClient, queryURL, contentTypeJSON, bytes.NewBuffer([]byte(q)))
	if err != nil {
		pretty.Errorf("error info is %s", err.Error())
		return
	}

	o, _ := ioutil.ReadAll(resp.Body)
	pretty.Println("response is ", string(o))
}

func TestReq(t *testing.T) {
	readReq := &prompb.ReadRequest{
		Queries: []*prompb.Query{{
			StartTimestampMs: 1537496087552,
			EndTimestampMs:   1537496387552,
			Matchers: []*prompb.LabelMatcher{
				{Type: 0, Name: "__name__", Value: "namespace_name:kube_pod_container_resource_requests_cpu_cores:sum"},
				{Type: 0, Name: "prometheus", Value: "monitoring/k8s"},
				{Type: 0, Name: "prometheus_replica", Value: "prometheus-k8s-0"},
			},
			Hints: &prompb.ReadHints{StepMs: 0, Func: "sum"},
		}},
	}

	logLevel := promlog.AllowedLevel{}
	logLevel.Set("debug")
	logger := promlog.New(logLevel)
	client := NewClient(log.With(logger, "storage", "OpenTSDB"),
		"http://10.0.0.131:4242",
		time.Second*30, )

	client.Read(readReq)
}

// request format
//`
//readReq := &prompb.ReadRequest{
//		Queries: []*prompb.Query{{
//			StartTimestampMs: 1537496087552,
//			EndTimestampMs:   1537496387552,
//			Matchers: []*prompb.LabelMatcher{
//				{Type: 0, Name: "__name__", Value: "namespace_name:kube_pod_container_resource_requests_cpu_cores:sum"},
//				{Type: 0, Name: "prometheus", Value: "monitoring/k8s"},
//				{Type: 0, Name: "prometheus_replica", Value: "prometheus-k8s-0"},
//			},
//			Hints: &prompb.ReadHints{StepMs: 0, Func: "sum"},
//		}},
//	}
//`

// converted format: {
//`
//{
//  "start": 1537496087,
//  "end": 1537496387,
//  "queries": [
//    {
//      "metric": "namespace__name_.kube__pod__container__resource__requests__cpu__cores_.sum",
//      "filters": [
//        {
//          "type": "literal_or",
//          "tagk": "prometheus",
//          "filter": "monitoring/k8s",
//          "groupBy": true
//        },
//        {
//          "type": "literal_or",
//          "tagk": "prometheus_replica",
//          "filter": "prometheus-k8s-0",
//          "groupBy": true
//        }
//      ],
//      "aggregator": "sum"
//    }
//  ]
//}
//`
