# o11y

[![Test installation in Minikube](https://github.com/chazapp/o11y/actions/workflows/terraform_tests.yaml/badge.svg)](https://github.com/chazapp/o11y/actions/workflows/terraform_tests.yaml)
[![WallAPI](https://github.com/chazapp/o11y/actions/workflows/wall_api_tests.yaml/badge.svg)](https://github.com/chazapp/o11y/actions/workflows/wall_api_tests.yaml)
[![WallFront](https://github.com/chazapp/o11y/actions/workflows/wall_front_tests.yaml/badge.svg)](https://github.com/chazapp/o11y/actions/workflows/wall_front_tests.yaml)
[![End-to-End](https://github.com/chazapp/o11y/actions/workflows/e2e_tests.yaml/badge.svg)](https://github.com/chazapp/o11y/actions/workflows/e2e_tests.yaml)

A DevOps shop showcase of observability in action.  

[Navigate this repository with the blog article !](https://blog.chaz.pro/posts/2023-11-18/wall-experience)

## Contents

A Terraform project that deploys to a Minikube cluster the following tools:

- Grafana
- Alloy
- Loki
- Kube-Prometheus-Stack
- Tempo
- Minio
- Pyroscope

It also deploys the following example services:

- The WallAPI, a REST API written in Golang instrumented via OpenTelemetrySDK and profiled by Pyroscope via /pprof
- The WallClient, a front-end client written in React, auto-instrumented via Faro -> Grafana Agent -> Tempo & Loki

## Minikube

To test the complete stack, create a Minikube cluster then apply with Terraform.

The API Gateway requires `minikube tunnel` to be running in order for the LoadBalancer service to bind to an IP address.
If the tunnel is not up, the Chart will fail to install and the terraform process won't be able to finish.


``` bash
$ minikube start
$ minikube tunnel &> /dev/null & # runs tunnel in the background
...
$ terraform apply
...
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

