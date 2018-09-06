# Requests per Second Demo

## Install Prometheus Operator

We will use the Prometheus Operator to create the Prometheus cluster and to create the relevant configuration to monitor our app. So first, install the operator:

```bash
kubectl apply -f https://raw.githubusercontent.com/coreos/prometheus-operator/master/bundle.yaml
```

This will create a number of CRDs, Service Accounts and RBAC things. Once it's finished, 

## Create a Prometheus instance

```bash
kubectl apply -f prometheus-config/prometheus
```

This will create a single replica Prometheus instance.

## Expose a Service for Prometheus instance

```bash
kubectl expose pod prometheus-prometheus-0 --port 9090 --target-port 9090
```

## Deploy The Adopt Dog App

```bash
helm install ./charts/adoptdog --name rps-prom
```

This will deploy with an ingress and should create the HPA, Prometheus ServiceMonitor and everything else needed, except the adapter. Do that next.

Change the values.yaml as needed (especially for ingress)

## Deploy the Prometheus Metric Adapter

```bash
helm install stable/prometheus-adapter --name prometheus-adaptor -f prometheus-config/prometheus-adapter/values.yaml
```

NOTE: if you have the Azure application insights adapter installed, you'll need to remove that first.

There might be some lag time between when you create the adapter and when the metrics are available.

```bash
kubectl get --raw /apis/custom.metrics.k8s.io/v1beta1/namespaces/default/pod/*/requests_per_second | jq .
```

This should show metrics if everything is setup correctly. Example:

```console
$ kubectl get --raw /apis/custom.metrics.k8s.io/v1beta1/namespaces/default/pod/*/requests_per_second | jq .
{
  "kind": "MetricValueList",
  "apiVersion": "custom.metrics.k8s.io/v1beta1",
  "metadata": {
    "selfLink": "/apis/custom.metrics.k8s.io/v1beta1/namespaces/default/pod/%2A/requests_per_second"
  },
  "items": [
    {
      "describedObject": {
        "kind": "Pod",
        "namespace": "default",
        "name": "rps-prom-adoptdog-8684976576-7hvc9",
        "apiVersion": "/__internal"
      },
      "metricName": "requests_per_second",
      "timestamp": "2018-09-05T03:57:44Z",
      "value": "0"
    },
    {
      "describedObject": {
        "kind": "Pod",
        "namespace": "default",
        "name": "rps-prom-adoptdog-8684976576-p6wm7",
        "apiVersion": "/__internal"
      },
      "metricName": "requests_per_second",
      "timestamp": "2018-09-05T03:57:44Z",
      "value": "0"
    }
  ]
}
```

## Hit it with some Load

I've been using [Hey](https://github.com/rakyll/hey)

```
export GOPATH=~/go
export PATH=$GOPATH/bin:$PATH
go get -u github.com/rakyll/hey
hey -z 20m http://<whatever-the-ingress-url-is>
```

## Watch it scale

```
$ kubectl get hpa rps-prom-adoptdog -w
NAME                REFERENCE                      TARGETS   MINPODS   MAXPODS   REPLICAS   AGE
rps-prom-adoptdog   Deployment/rps-prom-adoptdog   0 / 10    2         10        2          4m
rps-prom-adoptdog   Deployment/rps-prom-adoptdog   128500m / 10   2         10        2         4m
rps-prom-adoptdog   Deployment/rps-prom-adoptdog   170500m / 10   2         10        4         5m
rps-prom-adoptdog   Deployment/rps-prom-adoptdog   111500m / 10   2         10        4         5m
rps-prom-adoptdog   Deployment/rps-prom-adoptdog   95250m / 10   2         10        4         6m
rps-prom-adoptdog   Deployment/rps-prom-adoptdog   141 / 10   2         10        4         6m
```

Overtime, this should go up.

## Some Prometheus Queries

Rounded average requests per second per container

```
round(avg(irate(request_durations_histogram_secs_count[1m])))
```

Rounded total requests per second

```
round(sum(irate(request_durations_histogram_secs_count[1m])))
```

Individual number of data points, should correspond to number of containers

```
count(irate(request_durations_histogram_secs_count[1m]))
```

The number of containers running, per the counter

```
running_containers
```

Average response time in seconds:

```
avg(request_durations_histogram_secs_sum / request_durations_histogram_secs_count)
```
