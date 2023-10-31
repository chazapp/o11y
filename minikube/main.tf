terraform {
    backend "local" {
    }
}
module "chazapp_o11y" {
    source = "../"

    kube_config = "~/.kube/config"
    kube_context = "minikube"
}