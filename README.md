# o11y
[![Test installation in Minikube](https://github.com/chazapp/o11y/actions/workflows/terraform_tests.yaml/badge.svg)](https://github.com/chazapp/o11y/actions/workflows/terraform_tests.yaml)



A fully configured observability stack based on Grafana's tools.

## Contents

A Terraform project that deploys to a Minikube cluster the following tools:
  - Grafana
  - Grafana-Agent
  - Kube-Prometheus-Stack
  - Loki
  - Promtail
  - Tempo
  - Minio
  - Pyroscope

It also deploys the following example services:
  - The WallAPI, a REST API written in Golang instrumented via OpenTelemetrySDK and profiled by Pyroscope via /pprof
  - The WallClient, a front-end client written in React, auto-instrumented via Faro -> Grafana Agent -> Tempo & Loki



## Minikube

To test the complete stack, create a Minikube cluster then apply with Terraform

```
$ minikube start
$ terraform apply
```

Once applied, you may add the Ingresses for the various services to your `/etc/hosts` file

```
$ kubectl get ingress -n monitoring
NAME                               CLASS    HOSTS                 ADDRESS        PORTS   AGE
grafana                            <none>   grafana.local         192.168.49.2   80      10m
...
```

*Note: On Windows, you may need to use the `minikube tunnel` command and use 127.0.0.1 instead !*

You may now access services via your browser.

