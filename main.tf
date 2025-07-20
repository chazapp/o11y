terraform {
    required_providers {
        kubernetes = {
            source = "hashicorp/kubernetes"
            version = "2.30.0"
        }
        helm = {
            source = "hashicorp/helm"
            version = "2.13.2"
        }
        tls = {
          source = "hashicorp/tls"
          version = "4.1.0"
        }
    }
}

provider "helm" {
    kubernetes {
        config_path = "~/.kube/config"
    }
}

provider "kubernetes" {
    config_path = "~/.kube/config"
    config_context = "minikube"
}

provider "tls" {
  
}