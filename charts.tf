resource "kubernetes_namespace" "tools_namespace" {
  metadata {
    name = var.tools_namespace
  }
}

resource "helm_release" "kube-prometheus-stack" {
  name       = "kube-prometheus-stack"
  repository = "https://prometheus-community.github.io/helm-charts"
  chart      = "kube-prometheus-stack"
  version    = "58.7.2"
  
  namespace  = var.tools_namespace

  values = [
    var.kube-prometheus-stack-override != null ? var.kube-prometheus-stack-override : "${file("${path.module}/configs/kube-prometheus-stack.yaml")}"
  ]
}


resource "helm_release" "grafana" {
  name       = "grafana"
  repository = "https://grafana.github.io/helm-charts/"
  chart      = "grafana"
  version    = "9.2.10"

  namespace  = var.tools_namespace

  values = [
    var.grafana-override != null ? var.grafana-override : "${file("${path.module}/configs/grafana.yaml")}"
  ]
  depends_on = [
    helm_release.kube-prometheus-stack
  ]
}

resource "helm_release" "loki" {
  name       = "loki"
  repository = "https://grafana.github.io/helm-charts/"
  chart      = "loki"
  version    = "6.31.0"

  namespace  = var.tools_namespace

  values = [
    var.loki-override != null ? var.loki-override : "${file("${path.module}/configs/loki.yaml")}"
  ]
  depends_on = [
    helm_release.kube-prometheus-stack
  ]
}

resource "helm_release" "minio" {
  name       = "minio"
  repository = "https://charts.min.io/"
  chart      = "minio"
  version    = "5.2.0"

  namespace  = var.tools_namespace

  values = [
    var.minio-override != null ? var.minio-override : "${file("${path.module}/configs/minio.yaml")}"
  ]
  depends_on = [
    helm_release.kube-prometheus-stack
  ]
}

resource "helm_release" "tempo" {
  name       = "tempo"
  repository = "https://grafana.github.io/helm-charts/"
  chart      = "tempo"
  version    = "1.9.0"

  namespace  = var.tools_namespace

  values = [
    var.tempo-override != null ? var.tempo-override : "${file("${path.module}/configs/tempo.yaml")}"
  ]

  depends_on = [
    helm_release.kube-prometheus-stack
  ]
}

# resource "helm_release" "pyroscope" {
#   name       = "pyroscope"
#   repository = "https://grafana.github.io/helm-charts/"
#   chart      = "pyroscope"
#   version    = "1.6.1"

#   namespace  = var.tools_namespace

#   values = [
#     var.pyroscope-override != null ? var.pyroscope-override : "${file("${path.module}/configs/pyroscope.yaml")}"
#   ]

#   depends_on = [
#     helm_release.kube-prometheus-stack
#   ]
# }

resource "helm_release" "alloy" {
  name    = "alloy"
  repository = "https://grafana.github.io/helm-charts/"
  chart      = "alloy"
  version    = "1.1.2"
  namespace = var.tools_namespace

  values = [
    "${file("${path.module}/configs/alloy/alloy.yaml")}",
  ]
  depends_on = [ 
    kubernetes_config_map.alloy-config
   ]
}

resource "kubernetes_config_map" "alloy-config" {
  metadata {
    name = "alloy-config"
    namespace = var.tools_namespace
  }
  data = {
    config = "${file("${path.module}/configs/alloy/config.alloy")}"
  }
}

resource "kubernetes_namespace" "apps_namespace" {
  metadata {
    name = "apps"
  }
}

resource "helm_release" "wall_api" {
  name      = "wall-api"
  chart = "${path.module}/apps/wall_api/chart"
  version = "2.4.0"
  namespace = "apps"
  depends_on = [
    helm_release.kube-prometheus-stack
  ]
}

resource "helm_release" "wall_front" {
  name      = "wall-front"
  chart = "${path.module}/apps/wall_front/chart"
  version = "2.4.0"
  namespace = "apps"
}
