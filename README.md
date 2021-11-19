# CPA - Continuous Performance Analysis


This tool allows OpenShift users to run a watcher for Prometheus queries and define thresholds (using a yaml file) to observe the performance of the OpenShift cluster during performance testing.  It could be generalized to run constantly against a cluster and alert you when cluster is looking bad. It may sound like some of the other monitoring & alerting solutions but its supposed to be simple, scalable and user-friendly.


## Why use CPA:

- Runs external and can work with any "Prometheus"
- Can be extended to run queries other than Prometheus, such as ElasticSearch, or simple OC CLI commands
- History for each time you run - can be stored in log files


## Design:
```
                        ┌─────────────────────┐                             ┌───────────────────────────┐
                        │                     │                             │           OpenShift       │
                        │ Benchmark Job       │                             │                           │
                        │                     │                             │     ┌───────────────┐     │
                        │ (optional)          │                             │     │ Prometheus    │     │
                        └────────▲────────────┘                             │     │       ▲       │     │
                                 │                                          │     └───────┬───────┘     │ - At least one prometheus cluster info required
                         Ability to kill benchmark job                      │             │             │
                                 │                                          │             │             │
                                 │                                          └─────────────┼─────────────┘
                        ┌────────┴────────────┐                                           │
┌─────────────────┐     │                     │        Determines Url and Token           │
│                 │     │ Continuous Perf     ├───────────────────────────────────────────┘
│  Slack Notifs.  ◄─────┤   Analysis - CPA    │            Runs Queries
│                 │     │                     │
│                 │     │                     │                              ┌──────────────────────────┐
└─────────────────┘     └───────┬─────────────┘                              │                          │
                                │                                            │                          │
                                │         Requires Url and Token             │  Prometheus - external
                                └───────────────────────────────────────────►│                          │
                                              Runs Queries                   │                          │
                                                                             └──────────────────────────┘
```


## Features:

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
* [x] File logging the output
* [x] Print output to screen even when logging enabled - simultaneously
* [x] Let user decide query frequency
* [x] Slack Notification
* [x] Notify/Do Something(e.g. Pause/Kill benchmark jobs to preserve cluster) when results don't match conditions
* [x] Spawn goroutines to keep running queries and evaluating results to handle scale - e.g. when we have very large number of queries in the yaml file, we can divide and concurrently run queries
* [x] If slack config is not set, it is ignored and no attempts will be made to notify via slack
* [ ] debug mode
* [ ] use env vars
* [ ] Enhance log files to include uuid/time


## Usage:

* Then build the binary using make file: `make build` or update your binary using `make update`. You Can clean existin binary with `make clean` or do clean and update/build using `make all`.
* Set `KUBECONFIG` envvar, and make sure to review `config/queries.yaml`.
* You can then run the following command:
```sh

./bin/cpa -t 60s -h
Usage: cpa [--noclrscr] [--queries QUERIES] [--query-frequency QUERY-FREQUENCY] [--timeout TIMEOUT] [--log-output] [--terminate-benchmark TERMINATE-BENCHMARK]

Options:
  --noclrscr             Do not clear screen after each iteration. Clears screen by default. [default: false]
  --queries QUERIES, -q QUERIES
                         queries file to use [default: queries.yaml]
  --query-frequency QUERY-FREQUENCY, -f QUERY-FREQUENCY
                         How often do we run queries. You can pass values like 4h or 1h10m10s [default: 20s]
  --timeout TIMEOUT, -t TIMEOUT
                         Duration to run Continuous Performance Analysis. You can pass values like 4h or 1h10m10s [default: 4h]
  --log-output, -l       Output will be stored in a log file(cpa.log) in addition to stdout. [default: false]
  --terminate-benchmark TERMINATE-BENCHMARK, -k TERMINATE-BENCHMARK
                         When CPA is running in parallel with benchmark job, let CPA know to kill benchmark if any query fail. (E.g. -k <processID>) Helpful to preserve cluster for further analysis.
  --help, -h             display this help and exit
```

