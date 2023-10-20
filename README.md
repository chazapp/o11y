# o11y
[![Test installation in Minikube](https://github.com/chazapp/o11y/actions/workflows/tests.yml/badge.svg)](https://github.com/chazapp/o11y/actions/workflows/tests.yml)  

A fully configured observability stack based on Grafana's tools.

## Contents

A Terraform project that applies to a Minikube cluster the following Helm Charts:
  - Grafana
  - Grafana-Agent
  - Kube-Prometheus-Stack
  - Loki
  - Promtail
  - Tempo
  - Minio


## Usage

Create a Minikube cluster, enable the `ingress` addon, then apply with Terraform

```
$ minikube start
$ minikube addons enable ingress
$ terraform apply
```

Once applied, you may add the Ingress for the various services to your `/etc/hosts` file

```
$ kubectl get ingress -n monitoring
NAME                               CLASS    HOSTS                 ADDRESS        PORTS   AGE
grafana                            <none>   grafana.local         192.168.49.2   80      10m
grafana-agent                      <none>   grafana-agent.local   192.168.49.2   80      10m
kube-prometheus-stack-prometheus   <none>   prometheus.local      192.168.49.2   80      10m
loki-gateway                       <none>   loki.local            192.168.49.2   80      10m
minio-console                      <none>   console.minio.local   192.168.49.2   80      10m
```

*Note: On Windows, you may need to use the `minikube tunnel` command and use 127.0.0.1 instead !*

You may now access the various services via your browser.

