package cmd

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	exutil "github.com/openshift/openshift-tests/test/extended/util"
	"github.com/prometheus/common/model"
	v1 "k8s.io/api/core/v1"
	kapierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	e2e "k8s.io/kubernetes/test/e2e/framework"
	yaml "sigs.k8s.io/yaml"
)

const configPath = "./config/"

type prometheusConfig struct {
	Url         string `json:"URL"`
	BearerToken string `json:"BearerToken"`
}

func (c *prometheusConfig) Parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}

func readPrometheusConfig() (url, bearerToken string, ok bool) {
	data, err := ioutil.ReadFile(configPath + "prometheus.yaml")
	if err != nil {
		log.Printf("Cound't read %s/prometheus.yaml", configPath)
		return "", "", false
	}
	var config prometheusConfig
	if err := config.Parse(data); err != nil {
		log.Fatal(err)
	}
	if (config != prometheusConfig{}) {
		return config.Url, config.BearerToken, true
	}
	return "", "", false
}

func LocatePrometheus(oc *exutil.CLI) (url, bearerToken string, ok bool) {
	url, bearerToken, ok = readPrometheusConfig()
	if ok {
		return
	}

	_, err := oc.AdminKubeClient().CoreV1().Services("openshift-monitoring").Get("prometheus-k8s", metav1.GetOptions{})
	if kapierrs.IsNotFound(err) {
		return "", "", false
	}
	for i := 0; i < 30; i++ {
		secrets, err := oc.AdminKubeClient().CoreV1().Secrets("openshift-monitoring").List(metav1.ListOptions{})
		if err != nil {
			log.Printf("An Error has occured %s", err)
			return
		}
		for _, secret := range secrets.Items {
			if secret.Type != v1.SecretTypeServiceAccountToken {
				continue
			}
			if !strings.HasPrefix(secret.Name, "prometheus-") {
				continue
			}
			bearerToken = string(secret.Data[v1.ServiceAccountTokenKey])
			break
		}
		if len(bearerToken) == 0 {
			log.Println("Waiting for prometheus service account secret to show up")
			time.Sleep(time.Second)
			continue
		}
	}
	route, err := oc.AdminRouteClient().RouteV1().Routes("openshift-monitoring").Get("prometheus-k8s", metav1.GetOptions{})
	if kapierrs.IsNotFound(err) {
		return "", "", false
	}
	return "https://" + route.Spec.Host, bearerToken, true
}

type prometheusResponse struct {
	Status string                 `json:"status"`
	Data   prometheusResponseData `json:"data"`
}

type prometheusResponseData struct {
	ResultType string       `json:"resultType"`
	Result     model.Vector `json:"result"`
}

const (
	maxPrometheusQueryAttempts = 5
	prometheusQueryRetrySleep  = 10 * time.Second
)

func runQueryViaHTTP(url, bearer string) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearer))

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Printf("Prometheus request returned HTTP Code: %d\n", resp.StatusCode)
	return string(contents), nil
}

func RunQueries(promQueries map[string]bool, oc *exutil.CLI, ns, execPodName, baseURL, bearerToken string) {
	// expect all correct metrics within a reasonable time period
	queryErrors := make(map[string]error)
	passed := make(map[string]struct{})
	for i := 0; i < maxPrometheusQueryAttempts; i++ {
		for query, expected := range promQueries {
			if _, ok := passed[query]; ok {
				continue
			}
			//TODO when the http/query apis discussed at https://github.com/prometheus/client_golang#client-for-the-prometheus-http-api
			// and introduced at https://github.com/prometheus/client_golang/blob/master/api/prometheus/v1/api.go are vendored into
			// openshift/origin, look to replace this homegrown http request / query param with that API
			url := fmt.Sprintf("%s/api/v1/query?%s", baseURL, (url.Values{"query": []string{query}}).Encode())
			contents, err := runQueryViaHTTP(url, bearerToken)
			if err != nil {
				log.Fatal(err)
			}
			// check query result, if this is a new error log it, otherwise remain silent
			var result prometheusResponse
			if err := json.Unmarshal([]byte(contents), &result); err != nil {
				e2e.Logf("unable to parse query response for %s: %v", query, err)
				continue
			}
			metrics := result.Data.Result
			if result.Status != "success" {
				data, _ := json.Marshal(metrics)
				msg := fmt.Sprintf("promQL query: %s had reported incorrect status:\n%s", query, data)
				if prev, ok := queryErrors[query]; !ok || prev.Error() != msg {
					e2e.Logf("%s", msg)
				}
				queryErrors[query] = fmt.Errorf(msg)
				continue
			}
			if (len(metrics) > 0 && !expected) || (len(metrics) == 0 && expected) {
				data, _ := json.Marshal(metrics)
				msg := fmt.Sprintf("promQL query: %s had reported incorrect results:\n%s", query, data)
				if prev, ok := queryErrors[query]; !ok || prev.Error() != msg {
					e2e.Logf("%s", msg)
				}
				queryErrors[query] = fmt.Errorf(msg)
				continue
			}

			// query successful
			passed[query] = struct{}{}
			delete(queryErrors, query)
		}

		if len(queryErrors) == 0 {
			break
		}
		time.Sleep(prometheusQueryRetrySleep)
	}

	if len(queryErrors) != 0 {
		exutil.DumpPodLogsStartingWith("prometheus-0", oc)
	}
}
