# o11y
[![Test installation in Minikube](https://github.com/chazapp/o11y/actions/workflows/tests.yml/badge.svg)](https://github.com/chazapp/o11y/actions/workflows/tests.yml)  

A fully configured observability stack based on Grafana's tools.

## Contents

A Terraform module that deploys to a Kubernetes cluster the following tools:
  - Grafana
  - Grafana-Agent
  - Kube-Prometheus-Stack
  - Loki
  - Promtail
  - Tempo
  - Minio
  - Pyroscope

It also deploys the following example services:
  - The WallAPI, a REST API written in Golang, auto-instrumented via Bayla and profiled by Pyroscope via /pprof
  - The WallClient, a front-end client written in React, auto-instrumented client-side via Faro -> Grafana Agent -> Tempo & Loki



## Minikube

The Terraform module provides default configuration values for all tools that are deployed within it. These
are tailored for usage in a Minikube cluster, of which an application of the module is available in it own folder.


To test the complete stack, create a Minikube cluster, enable the `ingress` addon, then apply with Terraform

```
$ minikube start
$ minikube addons enable ingress
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

