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
			Tags:      map[string]TagValue{"job": TagValue("prometheus"), "monitor": TagValue("codelab-monitor"), "quantile": TagValue("0"), "instance": TagValue("localhost:9099")},
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
