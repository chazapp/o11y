resource "tls_private_key" "o11y" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "o11y" {
  private_key_pem = tls_private_key.o11y.private_key_pem

  subject {
    common_name  = "o11y.local"
    organization = "o11y Inc."
  }

  validity_period_hours = 24 * 365

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]

  dns_names = [
    "api.o11y.local",
    "o11y.local",
  ]
}

resource "kubernetes_secret" "gateway-certs" {
  metadata {
    name      = "gateway-certs"
    namespace = "apps"
  }

  type = "kubernetes.io/tls"
  
  data = {
    "tls.crt" = tls_self_signed_cert.o11y.cert_pem
    "tls.key" = tls_private_key.o11y.private_key_pem
  }
}

resource "local_file" "o11y-cert" {
  content  = tls_self_signed_cert.o11y.cert_pem
  filename = "${path.module}/.runtime/o11y.crt"
}