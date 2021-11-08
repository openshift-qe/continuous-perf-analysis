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
* [x] Add a Makefile
* [ ] Notify/Do Something when results don't match conditions
* [ ] Spawn goroutines to keep running queries and evaluating results to handle scale - e.g. when we have very large number of queries in the yaml file, we can divide and concurrently run queries
* [ ] File logging the output, screen will give current status, look at Prometheus alerts



## Usage:

* Then build the binary using make file: `make build` or update your binary using `make update`. You Can clean existin binary with `make clean` or do clean and update/build using `make all`.
* Set `KUBECONFIG` envvar, and make sure to review `config/queries.yaml`.
* You can then run the following command:
```sh

./bin/cpa --help
Usage: cpa [--noclrscr] [--queries QUERIES] [--timeout TIMEOUT]

Options:
  --noclrscr             Do not clear screen after each iteration. [default: false]
  --queries QUERIES      queries file to use [default: queries.yaml]
  --timeout TIMEOUT      Duration to run Continuous Performance Analysis. You can pass values like 4h or 1h10m10s [default: 4h]
  --help, -h             display this help and exit
```