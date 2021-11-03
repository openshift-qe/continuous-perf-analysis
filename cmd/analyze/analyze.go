package analyze

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"time"

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

func ReadPrometheusQueries() (queriesList queryList, err error) {
	data, err := ioutil.ReadFile(configPath + "queries.yaml")
	if err != nil {
		log.Printf("Cound't read %s/queries.yaml", configPath)
		return queriesList, err
	}
	if err := queriesList.Parse(data); err != nil {
		log.Fatal(err)
	}
	// fmt.Println(queriesList)
	if len(queriesList) == 0 {
		return queriesList, nil
	}

	return queriesList, nil
}

func Queries(queryList queryList, oc *exutil.CLI, baseURL, bearerToken string, c chan string) {
	for _, items := range queryList {
		fmt.Printf("\nQuery:%s\n", items.Query)
		result, err := prometheus.RunQuery(items.Query, oc, baseURL, bearerToken)
		if err != nil {
			fmt.Println(err)
			continue
		}
		opMap := map[string]string{"eq": "==", "lt": "<", "gt": ">", "lte": "<=", "gte": ">="}
		for _, metric := range result.Data.Result {
			for _, watchItems := range items.WatchFor {
				// fmt.Println(watchItems.Key, watchItems.Val, watchItems.Threshold)
				// fmt.Println(metric.Metric[model.LabelName(watchItems.Key)], model.LabelValue(watchItems.Val), metric.Value)
				// e.g. if "metric.Metric[model.LabelName(watchItems.Key)]" --> metric.Metric["phase"] ==  model.LabelValue(watchItems.Val)  --> "Running"
				// or watchItems key is nil - meaning its a numerical query such as max()
				if metric.Metric[model.LabelName(watchItems.Key)] == model.LabelValue(watchItems.Val) || watchItems.Key == "nil" {
					// fmt.Println(metric.Metric[model.LabelName(watchItems.Key)], metric.Value, watchItems.Threshold, watchItems.Operator)
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
					fmt.Printf("\nValue: %.4f %s Threshold: %.4f is %t\n", v1, opMap[watchItems.Operator], v2, b)
					if !b {
						fmt.Printf("\n Comparison of Value and Threshold is %t. Notifying...\n", b)
						c <- fmt.Sprintf("\nValue: %.4f %s Threshold: %.4f is %t\n", v1, opMap[watchItems.Operator], v2, b)
					}
				}
			}
		}
	}
}

func Notify(c chan string) {
	waitChars := []string{"/", "-", "\\", "|"}
	for {
		select {
		case msg := <-c:
			fmt.Println("***************************************")
			fmt.Println("Received following on the channel:", msg)
			fmt.Println("***************************************")
		default:
			fmt.Printf("\r%s Please Wait. No new message received on the channel....", waitChars[rand.Intn(4)])
			time.Sleep(time.Millisecond * 500)
		}
	}

}
