# prom-exporter
Repo for learning how to implement Prometheus exporter. This contains the complete code of the snippets in [this dev.to article](https://dev.to/metonymicsmokey/custom-prometheus-metrics-with-go-520n). 

### Instructions to run:    
* Have Prometheus running on port 9090 and make sure port 2112 is not in use.   
* Async exporter: `go run async.go temp.go`.     
* Sync exporter: `go run sync.go temp.go`.     
