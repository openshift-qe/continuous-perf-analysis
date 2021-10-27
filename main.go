package main

import (
	"log"

	prometheus "github.com/kedark3/cpa/cmd"
	exutil "github.com/openshift/openshift-tests/test/extended/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	oc := exutil.NewCLI("prometheus-cpa", exutil.KubeConfigPath())
	secrets, err := oc.AdminKubeClient().CoreV1().Secrets("openshift-monitoring").List(metav1.ListOptions{})

	if err != nil {
		log.Printf("An Error has occured %s", err)
		return
	}
	log.Printf("Found following secrets %d", secrets.Size())
	prometheus.LocatePrometheus()

}
