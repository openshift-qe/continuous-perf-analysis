# CPA - Continuous Performance Analysis


This tool allows OpenShift users to run a watcher for Prometheus queries and define thresholds (using a yaml file) to observe the performance of the OpenShift cluster during performance testing.  It could be generalized to run constantly against a cluster and alert you when cluster is looking bad. It may sound like some of the other monitoring & alerting solutions but its supposed to be simple, scalable and user-friendly.

## To Do:

* [x] Create oc cli connection to OpenShift/Kubernetes using Kubeconfig
* [x] Determine Prometheus url, bearerToken for OpenShift
* [ ] If Prometheus url, bearerToken already included in the yaml, use that
* [ ] Create yaml format for queries, and expected outcomes (Use a struct to read that in)
* [ ] Spawn goroutines to keep running queries and evaluating results
* [ ] Notify/Do Something when results don't match conditions