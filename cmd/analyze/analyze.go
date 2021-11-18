package analyze

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"

	prometheus "github.com/kedark3/cpa/cmd/prometheus"
	exutil "github.com/openshift/openshift-tests/test/extended/util"

	"github.com/prometheus/common/model"

	"sigs.k8s.io/yaml"
)

/*
This package will read the queries from relevant config file under `config/` dir
and run relevant queries using prometheus package. Then retrieve the result and
analyze it against the threshold.
*/
var wg sync.WaitGroup

const configPath = "./config/"

type queryList []queries

type queries struct {
	Query    string `json:"query"`
	WatchFor []watchList
}

type watchList struct {
	Key       string `json:"Key"`
	Val       string `json:"Val"`
	Threshold string `json:"threshold"`
	Operator  string `json:"operator"`
}

func (q *queryList) Parse(data []byte) error {
	return yaml.Unmarshal(data, q)
}

func ReadPrometheusQueries(queriesFile string) (queriesList queryList, err error) {
	data, err := ioutil.ReadFile(configPath + queriesFile)
	if err != nil {
		log.Printf("Cound't read %s/queries.yaml", configPath)
		return queriesList, err
	}
	if err := queriesList.Parse(data); err != nil {
		log.Fatal(err)
	}
	// log.Println(queriesList)
	if len(queriesList) == 0 {
		return queriesList, fmt.Errorf("query list is empty: %v", queriesList)
	}

	return queriesList, nil
}

func Queries(queryList queryList, oc *exutil.CLI, baseURL, bearerToken string, c chan string, tb chan bool, terminateBenchmark string) {
	// start := time.Now()
	for _, item := range queryList {
		go runQuery(item, oc, baseURL, bearerToken, c, tb, terminateBenchmark)
	}
	wg.Wait()
	// end := time.Since(start)
	// log.Printf("\n It takes %s time to run queries", end)
}

func runQuery(q queries, oc *exutil.CLI, baseURL, bearerToken string, c chan string, tb chan bool, terminateBenchmark string) {
	wg.Add(1)
	defer wg.Done()
	result, err := prometheus.RunQuery(q.Query, oc, baseURL, bearerToken)
	if err != nil {
		log.Println(err)
		return
	}
	opMap := map[string]string{"eq": "==", "lt": "<", "gt": ">", "lte": "<=", "gte": ">="}
	for _, metric := range result.Data.Result {
		for _, watchItems := range q.WatchFor {
			// log.Println(watchItems.Key, watchItems.Val, watchItems.Threshold)
			// log.Println(metric.Metric[model.LabelName(watchItems.Key)], model.LabelValue(watchItems.Val), metric.Value)
			// e.g. if "metric.Metric[model.LabelName(watchItems.Key)]" --> metric.Metric["phase"] ==  model.LabelValue(watchItems.Val)  --> "Running"
			// or watchItems key is nil - meaning its a numerical query such as max()
			if metric.Metric[model.LabelName(watchItems.Key)] == model.LabelValue(watchItems.Val) || watchItems.Key == "nil" {
				// log.Println(metric.Metric[model.LabelName(watchItems.Key)], metric.Value, watchItems.Threshold, watchItems.Operator)
				v1, _ := strconv.ParseFloat(metric.Value.String(), 64)
				v2, _ := strconv.ParseFloat(watchItems.Threshold, 64)
				b := true // if this becomes false we send message on go channel
				switch watchItems.Operator {
				case "eq":
					b = v1 == v2
				case "gt":
					b = v1 > v2
				case "lt":
					b = v1 < v2
				case "lte":
					b = v1 <= v2
				case "gte":
					b = v1 >= v2
				}
				log.Printf(`
Query:%s
Value: %.4f %s Threshold: %.4f is %t
`, q.Query, v1, opMap[watchItems.Operator], v2, b)
				if !b {
					log.Printf("\n%[2]s\n Comparison of Value and Threshold is %[1]t. Notifying...\n%[2]s\n", b, strings.Repeat("~", 80))
					c <- fmt.Sprintf("\nQuery: %s\nValue: %.4f %s Threshold: %.4f is %t for key: %q and val: %q\n", q.Query, v1, opMap[watchItems.Operator], v2, b, watchItems.Key, watchItems.Val)
					if terminateBenchmark != "" {
						tb <- true // send signal to terminate benchmark channel
					}
				}
			}
		}
	}
}
