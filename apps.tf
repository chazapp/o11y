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


resource "helm_release" "gateway" {
  name = "gateway"
  chart = "${path.module}/apps/gateway/charts"
  version = "1.0.0"
  namespace = "apps"
  depends_on = [ helm_release.istio-gateway, kubernetes_secret.gateway-certs ]
}

resource "helm_release" "landing" {
  name      = "landing"
  chart     = "${path.module}/apps/landing/charts"
  version   = "0.1.0"
  namespace = "apps"
  depends_on = [ helm_release.gateway ]
}

resource "helm_release" "auth" {
  name      = "auth"
  chart     = "${path.module}/apps/auth/charts"
  version   = "0.1.0"
  namespace = "apps"
  depends_on = [
    helm_release.gateway,
    data.external.keygen
  ]

  set_sensitive {
    name  = "secrets.jwt.privateKey"
    value = data.external.keygen.result.privateKey
  }

  set_sensitive {
    name  = "secrets.jwt.publicKey"
    value = data.external.keygen.result.publicKey
  }
}

resource "null_resource" "keygen" {
  provisioner "local-exec" {
    command = <<EOT
      mkdir -p .runtime/ && cd .runtime
      jose-util generate-key -use sig -alg RS256
      mv jwk-*-priv.json privateKey.json
      mv jwk-*-pub.json publicKey.json
    EOT
  }
}

data "external" "keygen" {
  depends_on = [ null_resource.keygen ]

  program = ["bash", "-c", "./scripts/tf-get-jwt-keys.sh"]
}
