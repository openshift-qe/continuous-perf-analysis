package analyze

import (
	"io/ioutil"
	"log"

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
