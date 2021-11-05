package cmd

/*
Program to read the prometheus config, and use it to run queries against OCP Cluster prometheus.
You can either set the Prometheus URL and BearerToken in config/prometheus.yaml or with the
KUBECONFIG env var being set, the program can use oc get routes and oc get secrets to find the URL
and bearer Token.
*/

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
	Url         string `json:"url"`
	BearerToken string `json:"bearerToken"`
}

func (c *prometheusConfig) Parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}

func readPrometheusConfig() (url, bearerToken string, err error) {
	data, err := ioutil.ReadFile(configPath + "prometheus.yaml")
	msg := fmt.Sprintf("Cound't read %s/prometheus.yaml", configPath)
	if err != nil {
		log.Println(msg)
		return "", "", fmt.Errorf(msg)
	}
	var config prometheusConfig
	if err := config.Parse(data); err != nil {
		log.Fatal(err)
	}
	if (config != prometheusConfig{}) {
		return config.Url, config.BearerToken, nil
	}
	return "", "", fmt.Errorf(msg)
}

func LocatePrometheus(oc *exutil.CLI) (url, bearerToken string, err error) {
	url, bearerToken, err = readPrometheusConfig()
	if err != nil {
		return "", "", err
	}

	_, err = oc.AdminKubeClient().CoreV1().Services("openshift-monitoring").Get("prometheus-k8s", metav1.GetOptions{})
	if kapierrs.IsNotFound(err) {
		return "", "", err
	}
	for i := 0; i < 30; i++ {
		secrets, secretErr := oc.AdminKubeClient().CoreV1().Secrets("openshift-monitoring").List(metav1.ListOptions{})
		if secretErr != nil {
			log.Printf("An Error has occured %s", secretErr)
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
		return "", "", err
	}
	return "https://" + route.Spec.Host, bearerToken, nil
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

	return string(contents), nil
}

func RunQuery(promQuery string, oc *exutil.CLI, baseURL, bearerToken string) (prometheusResponse, error) {
	// expect all correct metrics within a reasonable time period
	queryErrors := make(map[string]error)
	passed := make(map[string]struct{})
	var result prometheusResponse
	for i := 0; i < maxPrometheusQueryAttempts; i++ {
		if _, ok := passed[promQuery]; ok {
			continue
		}

		url := fmt.Sprintf("%s/api/v1/query?%s", baseURL, (url.Values{"query": []string{promQuery}}).Encode())
		contents, err := runQueryViaHTTP(url, bearerToken)
		if err != nil {
			log.Fatal(err)
		}
		// check query result, if this is a new error log it, otherwise remain silent
		if err := json.Unmarshal([]byte(contents), &result); err != nil {
			e2e.Logf("unable to parse query response for %s: %v", promQuery, err)
			continue
		}
		metrics := result.Data.Result
		if result.Status != "success" {
			data, _ := json.Marshal(metrics)
			msg := fmt.Sprintf("promQL query: %s had reported incorrect status:\n%s", promQuery, data)
			if prev, ok := queryErrors[promQuery]; !ok || prev.Error() != msg {
				e2e.Logf("%s", msg)
			}
			queryErrors[promQuery] = fmt.Errorf(msg)
			continue
		}

		// for _, r := range result.Data.Result {
		// 	log.Printf("Type is %[1]T \n Metric is: %[1]s\n\n", r.Metric["phase"])
		// 	log.Printf("Type is %[1]T \n Value is: %[1]s\n\n", r.Value)
		// 	if r.Metric["phase"] == "Running" {
		// 		log.Printf("We have %v pod in %s phase", r.Value, r.Metric["phase"])
		// 	}
		// }
		// query successful
		passed[promQuery] = struct{}{}
		delete(queryErrors, promQuery)
		// if there were no errors let's break out of the loop
		if len(queryErrors) == 0 {
			return result, nil
		}
		// else sleep for 10 sec
		time.Sleep(prometheusQueryRetrySleep)
	}

	// if there were errors, dump logs from prometheus pod
	if len(queryErrors) != 0 {
		exutil.DumpPodLogsStartingWith("prometheus-0", oc)
	}
	return result, queryErrors[promQuery]
}
