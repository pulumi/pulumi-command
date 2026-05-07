pulumi {
  required_providers {
    eks = {
      source  = "pulumi/eks"
      version = ">= 3.0.0"
    }
  }
}

resource "eks_cluster" "cluster" {
}

resource "command_local_command" "cleanup_kubernetes_namespaces" {
  delete = "kubectl --kubeconfig <(echo \"$KUBECONFIG_DATA\") delete namespace nginx\n"

  interpreter = ["/bin/bash", "-c"]

  environment = {
    KUBECONFIG_DATA = eks_cluster.cluster.kubeconfig_json
  }
}
