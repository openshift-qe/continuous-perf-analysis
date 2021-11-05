# CPA - Continuous Performance Analysis


This tool allows OpenShift users to run a watcher for Prometheus queries and define thresholds (using a yaml file) to observe the performance of the OpenShift cluster during performance testing.  It could be generalized to run constantly against a cluster and alert you when cluster is looking bad. It may sound like some of the other monitoring & alerting solutions but its supposed to be simple, scalable and user-friendly.

## To Do:

* [x] Create oc cli connection to OpenShift/Kubernetes using Kubeconfig
* [x] Determine Prometheus url, bearerToken for OpenShift
* [x] If Prometheus url, bearerToken already included in the yaml, use that
* [x] Create yaml format for queries, and expected outcomes (Use a struct to read that in)
* [x] Spwan go routine to run queries and analyze results
* [x] Spwan goroutine to receive notification when a query yields "False" value
* [x] Update to latest go and recompile
* [x] Add CLI to the program
  * [x] Add a parameter to read different query files in config dir
  * [x] Add parameter for clearing/not-clearing screen
  * [x] Add Parameter for timeout
* [ ] Notify/Do Something when results don't match conditions
* [ ] Spawn goroutines to keep running queries and evaluating results to handle scale - e.g. when we have very large number of queries in the yaml file, we can divide and concurrently run queries
* [ ] File logging the output, screen will give current status, look at Prometheus alerts
