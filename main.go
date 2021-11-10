// +build go1.17

package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/alexflint/go-arg"
	analyze "github.com/kedark3/cpa/cmd/analyze"
	notify "github.com/kedark3/cpa/cmd/notify"
	prometheus "github.com/kedark3/cpa/cmd/prometheus"
	exutil "github.com/openshift/openshift-tests/test/extended/util"

	g "github.com/onsi/ginkgo"
	o "github.com/onsi/gomega"
)

func main() {
	var args struct {
		NoClrscr           bool          `arg:"--noclrscr" help:"Do not clear screen after each iteration. Clears screen by default." default:"false"`
		Queries            string        `arg:"-q,--queries" help:"queries file to use" default:"queries.yaml"`
		QueryFrequency     time.Duration `arg:"-f,--query-frequency" help:"How often do we run queries. You can pass values like 4h or 1h10m10s" default:"20s"`
		Timeout            time.Duration `arg:"-t,--timeout" help:"Duration to run Continuous Performance Analysis. You can pass values like 4h or 1h10m10s" default:"4h"`
		LogOutput          bool          `arg:"-l,--log-output" help:"Output will be stored in a log file(cpa.log) in addition to stdout." default:"false"`
		TerminateBenchmark string        `arg:"-k,--terminate-benchmark" help:"When CPA is running in parallel with benchmark job, let CPA know to kill benchmark if any query fail. (E.g. -k <processID>) Helpful to preserve cluster for further analysis." default:""`
	}
	arg.MustParse(&args)

	o.RegisterFailHandler(g.Fail)
	if args.LogOutput {
		f, err := os.OpenFile("cpa.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		multiWriter := io.MultiWriter(os.Stdout, f)
		if err != nil {
			log.Fatal(err)
		}

		//defer to close when you're done with it, not because you think it's idiomatic!
		defer f.Close()

		//set output of logs to f
		log.SetOutput(multiWriter)
	}

	oc := exutil.NewCLI("prometheus-cpa", exutil.KubeConfigPath())
	// secrets, err := oc.AdminKubeClient().CoreV1().Secrets("openshift-monitoring").List(metav1.ListOptions{})

	// if err != nil {
	// 	log.Printf("An Error has occured %s", err)
	// 	return
	// }
	// log.Printf("Found following secrets %d", secrets.Size())
	url, bearerToken, err := prometheus.LocatePrometheus(oc)
	if err != nil {
		log.Printf("Oops something went wrong while trying to fetch Prometheus url and bearerToken")
		log.Println(err)
		return
	}

	slackConfig, err := notify.ReadslackConfig()
	if err != nil {
		log.Printf("Oops something went wrong while trying to fetch Slack Config")
		log.Println(err)
		return
	}
	// fmt.Println("UserID, Channel ID, slackToken: ", slackConfig)
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
	tb := make(chan bool)
	c := make(chan string)

	thread_ts := slackConfig.SlackNotify("New benchmark started, we will monitor it for performance and notify here with the issues.", "")
	go func(c chan string) {
		for i := 1; ; i++ {
			log.Printf("Iteration no. %d\n", i)
			queryList, err := analyze.ReadPrometheusQueries(args.Queries)
			if err != nil {
				log.Println(err)
				return
			}
			analyze.Queries(queryList, oc, url, bearerToken, c, tb, args.TerminateBenchmark)
			time.Sleep(args.QueryFrequency)
			if !args.NoClrscr {
				log.Print("\033[H\033[2J") // clears screen before printing next iteration
			}
		}
	}(c)
	go slackConfig.Notify(c, thread_ts)

	if args.TerminateBenchmark != "" {
		go notify.TerminateBenchmark(tb, args.TerminateBenchmark)
	}

	d, err := time.ParseDuration(args.Timeout.String())
	if err != nil {
		log.Println(err)
	}
	time.Sleep(d)
}
