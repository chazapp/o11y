variable "tools_namespace" {
    description = "Namespace to install monitoring tools"
    type = string
    default = "monitoring"
}

variable "kube_config" {
    description = "Path to .kube/config file"
    type = string
}

variable "kube_context" {
    description = "Kube Context entry"
    type = string
}

variable "kube-prometheus-stack-override" {
    description = "Override values.yaml for Kube-Prometeus-Stack Helm Chart"
    type        = any
    default     = null
}

variable "grafana-override" {
    description = "Override values.yaml for Grafana Helm Chart"
    type        = any
    default     = null
}

variable "grafana-agent-override" {
    description = "Override values.yaml for Grafana-Agent Helm Chart"
    type        = any
    default     = null
}

variable "loki-override" {
    description = "Override values.yaml for Loki Helm Chart"
    type        = any
    default     = null
}

variable "minio-override" {
    description = "Override values.yaml for Minio Helm Chart"
    type        = any
    default     = null
}

variable "promtail-override" {
    description = "Override values.yaml for Promtail Helm Chart"
    type        = any
    default     = null
}

variable "tempo-override" {
    description = "Override values.yaml for Tempo Helm Chart"
    type        = any
    default     = null
}

variable "pyroscope-override" {
    description = "Override values.yaml for Pyroscope Helm Chart"
    type        = any
    default     = null
}
