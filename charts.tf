resource "kubernetes_namespace" "tools_namespace" {
  metadata {
    name = var.tools_namespace
  }
}

resource "helm_release" "kube-prometheus-stack" {
  name       = "kube-prometheus-stack"
  repository = "https://prometheus-community.github.io/helm-charts"
  chart      = "kube-prometheus-stack"
  version    = "51.8.1"

  namespace  = var.tools_namespace

  values = [
    "${file("configs/kube-prometheus-stack.yaml")}"
  ]
}


resource "helm_release" "grafana" {
  name       = "grafana"
  repository = "https://grafana.github.io/helm-charts/"
  chart      = "grafana"
  version    = "6.60.6"

  namespace  = var.tools_namespace

  values = [
    "${file("configs/grafana.yaml")}"
  ]
  depends_on = [
    helm_release.kube-prometheus-stack
  ]
}

resource "helm_release" "grafana-agent" {
  name       = "grafana-agent"
  repository = "https://grafana.github.io/helm-charts/"
  chart      = "grafana-agent"
  version    = "0.26.0"

  namespace  = var.tools_namespace

  values = [
    "${file("configs/grafana-agent.yaml")}"
  ]
  
  depends_on = [
    helm_release.kube-prometheus-stack
  ]
}

resource "helm_release" "loki" {
  name       = "loki"
  repository = "https://grafana.github.io/helm-charts/"
  chart      = "loki"
  version    = "5.31.0"

  namespace  = var.tools_namespace

  values = [
    "${file("configs/loki.yaml")}"
  ]
  depends_on = [
    helm_release.kube-prometheus-stack
  ]
}

resource "helm_release" "minio" {
  name       = "minio"
  repository = "https://charts.min.io/"
  chart      = "minio"
  version    = "5.0.14"

  namespace  = var.tools_namespace

  values = [
    "${file("configs/minio.yaml")}"
  ]
  depends_on = [
    helm_release.kube-prometheus-stack
  ]
}

resource "helm_release" "promtail" {
  name       = "promtail"
  repository = "https://grafana.github.io/helm-charts/"
  chart      = "promtail"
  version    = "6.15.2"

  namespace  = var.tools_namespace

  values = [
    "${file("configs/promtail.yaml")}"
  ]

  depends_on = [
    helm_release.kube-prometheus-stack
  ]
}

resource "helm_release" "tempo" {
  name       = "tempo"
  repository = "https://grafana.github.io/helm-charts/"
  chart      = "tempo"
  version    = " 1.6.2"

  namespace  = var.tools_namespace

  values = [
    "${file("configs/tempo.yaml")}"
  ]

  depends_on = [
    helm_release.kube-prometheus-stack
  ]
}

resource "helm_release" "pyroscope" {
  name       = "pyroscope"
  repository = "https://grafana.github.io/helm-charts/"
  chart      = "pyroscope"
  version    = " 1.1.0"

  namespace  = var.tools_namespace

  values = [
    "${file("configs/pyroscope.yaml")}"
  ]

  depends_on = [
    helm_release.kube-prometheus-stack
  ]
}