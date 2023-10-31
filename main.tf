terraform {
    required_providers {
        kubernetes = {
            source = "hashicorp/kubernetes"
            version = "2.23.0"
        }
        helm = {
            source = "hashicorp/helm"
            version = "2.11.0"
        }
    }
}

provider "helm" {
    kubernetes {
        config_path = var.kube_config
    }
}

provider "kubernetes" {
    config_path = var.kube_config
    config_context = var.kube_context
}
