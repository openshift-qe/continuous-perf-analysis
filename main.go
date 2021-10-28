package main

import (
	"log"

	prometheus "github.com/kedark3/cpa/cmd"
	exutil "github.com/openshift/openshift-tests/test/extended/util"

	g "github.com/onsi/ginkgo"
	o "github.com/onsi/gomega"
)

func main() {
	o.RegisterFailHandler(g.Fail)

	oc := exutil.NewCLI("prometheus-cpa", exutil.KubeConfigPath())
	// secrets, err := oc.AdminKubeClient().CoreV1().Secrets("openshift-monitoring").List(metav1.ListOptions{})

	// if err != nil {
	// 	log.Printf("An Error has occured %s", err)
	// 	return
	// }
	// log.Printf("Found following secrets %d", secrets.Size())
	url, bearerToken, ok := prometheus.LocatePrometheus(oc)
	if !ok {
		log.Printf("Oops something went wrong while trying to fetch Prometheus url and bearerToken")
		return
	}

	tests := map[string]bool{
		`sum(kube_pod_status_phase{}) by (phase) > 0`:                       true,
		`sum(kube_namespace_status_phase) by (phase)`:                       true,
		`sum(kube_node_status_condition{status="true"}) by (condition) > 0`: true,
	}
	prometheus.RunQueries(tests, oc, url, bearerToken)
	// log.Printf("URL is %s and bearerToken is %s", url, bearerToken)

}
