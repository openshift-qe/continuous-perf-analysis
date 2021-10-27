package main

import (
	"log"

	prometheus "github.com/kedark3/cpa/cmd"
	exutil "github.com/openshift/openshift-tests/test/extended/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	g "github.com/onsi/ginkgo"
	o "github.com/onsi/gomega"
)

func main() {
	o.RegisterFailHandler(g.Fail)

	oc := exutil.NewCLI("prometheus-cpa", exutil.KubeConfigPath())
	secrets, err := oc.AdminKubeClient().CoreV1().Secrets("openshift-monitoring").List(metav1.ListOptions{})

	if err != nil {
		log.Printf("An Error has occured %s", err)
		return
	}
	log.Printf("Found following secrets %d", secrets.Size())
	url, bearerToken, ok := prometheus.LocatePrometheus(oc)
	if !ok {
		log.Printf("Oops something went wrong while trying to fetch Prometheus url and bearerToken")
		return
	}

	oc.SetupProject()
	ns := oc.Namespace()

	execPod := exutil.CreateCentosExecPodOrFail(oc.AdminKubeClient(), ns, "execpod", nil)
	defer func() { oc.AdminKubeClient().CoreV1().Pods(ns).Delete(execPod.Name, metav1.NewDeleteOptions(1)) }()

	tests := map[string]bool{
		// Should have successfully sent at least some metrics to remote write endpoint
		// uncomment this once https://github.com/openshift/cluster-monitoring-operator/pull/434
		// is merged, and remove the other two checks.
		// `prometheus_remote_storage_succeeded_samples_total{job="prometheus-k8s"} >= 1`: true,

		// should have successfully sent at least once to remote
		`metricsclient_request_send{client="federate_to",job="telemeter-client",status_code="200"} >= 1`: true,
		// should have scraped some metrics from prometheus
		`federate_samples{job="telemeter-client"} >= 10`: true,
	}
	prometheus.RunQueries(tests, oc, ns, execPod.Name, url, bearerToken)
	log.Printf("URL is %s and bearerToken is %s", url, bearerToken)

}
