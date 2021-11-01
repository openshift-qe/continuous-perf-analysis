package main

import (
	"fmt"
	"log"

	analyze "github.com/kedark3/cpa/cmd/analyze"

	g "github.com/onsi/ginkgo"
	o "github.com/onsi/gomega"
)

func main() {
	o.RegisterFailHandler(g.Fail)

	// oc := exutil.NewCLI("prometheus-cpa", exutil.KubeConfigPath())
	// secrets, err := oc.AdminKubeClient().CoreV1().Secrets("openshift-monitoring").List(metav1.ListOptions{})

	// if err != nil {
	// 	log.Printf("An Error has occured %s", err)
	// 	return
	// }
	// log.Printf("Found following secrets %d", secrets.Size())
	// url, bearerToken, ok := prometheus.LocatePrometheus(oc)
	// if !ok {
	// 	log.Printf("Oops something went wrong while trying to fetch Prometheus url and bearerToken")
	// 	return
	// }

	// queries := []string{
	// `sum(kube_pod_status_phase{}) by (phase) > 0`, // pod count by phase
	// `sum(kube_namespace_status_phase) by (phase)`, // namespace count by phase
	// `sum(kube_node_status_condition{status="true"}) by (condition) > 0`,                                                                                   // node condition by status
	// `sum by (instance) (rate(ovnkube_master_pod_creation_latency_seconds_sum[2m]))`,                                                                       // OVN pod creation latency
	// `sum by (instance) (rate(ovnkube_node_cni_request_duration_seconds_sum{command="ADD"}[2m]))`,                                                          // CNI Request duration for "ADD" command over 2m interval
	// `sum by (instance) (rate(ovnkube_node_cni_request_duration_seconds_sum{command="DEL"}[2m]))`,                                                          // CNI Request duration for "DEL" command over 2m interval
	// `sum(container_memory_working_set_bytes{pod=~"ovnkube-master-.*",namespace="openshift-ovn-kubernetes",container=""}) by (pod, node)`,                  // ovnkube-master Memory Usage
	// `sum(container_memory_working_set_bytes{pod=~"ovnkube-master-.*",namespace="openshift-ovn-kubernetes",container!=""}) by (pod, node)`,                 // ovnkube-master Memory Usage
	// `topk(10, rate(container_cpu_usage_seconds_total{pod=~"ovnkube-.*",namespace="openshift-ovn-kubernetes",container="ovn-controller"}[2m])*100)`,        // top 10 - ovn-controller cpu usage
	// `topk(10, sum(container_memory_working_set_bytes{pod=~"ovnkube-node-.*",namespace="openshift-ovn-kubernetes",container="ovn-controller"}) by (node))`, // top 10 - ovn-controller memory usage
	// `sum(container_memory_rss{pod="prometheus-k8s-0",namespace!="",name!="",container="prometheus"}) by (pod)`,                                            // Prometheus replica 0 rss memory
	// `sum(container_memory_rss{pod="prometheus-k8s-1",namespace!="",name!="",container="prometheus"}) by (pod)`,                                            // Prometheus replica 1 rss memory
	// `rate(container_cpu_usage_seconds_total{pod=~"ovnkube-master.*",namespace="openshift-ovn-kubernetes",container!=""}[2m])*100`,                         // CPU usage ovnkube-master components over 2m interval
	// `sum by (condition)(cluster_operator_conditions{condition!=""})`,
	// }
	// log.Printf("URL is %s and bearerToken is %s", url, bearerToken)
	// for _, query := range queries {
	// 	fmt.Println(prometheus.RunQuery(query, oc, url, bearerToken))
	// 	fmt.Println()
	// }
	queryList, err := analyze.ReadPrometheusQueries()
	if err != nil {
		log.Println(err)
	}
	for _, items := range queryList {
		fmt.Println(items.Query)
		for _, watchItems := range items.WatchFor {
			fmt.Println(watchItems.Key, watchItems.Val, watchItems.Threshold)
		}
	}
}
