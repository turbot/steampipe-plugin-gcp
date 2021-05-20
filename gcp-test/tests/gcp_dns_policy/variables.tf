
variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125-create-update"
  description = "Name of the resource used throughout the test."
}

variable "gcp_project" {
  type        = string
  default     = "niteowl-aaa"
  description = "GCP project used for the test."
}

variable "gcp_region" {
  type        = string
  default     = "us-east1"
  description = "GCP region used for the test."
}

variable "gcp_zone" {
  type    = string
  default = "us-east1-b"
}

variable "enable_inbound_forwarding" {
  type = bool
  default = true
}

variable "enable_logging" {
  type = bool
  default = true
}

provider "google" {
  project = var.gcp_project
  region  = var.gcp_region
  zone    = var.gcp_zone
}


data "google_client_config" "current" {}

data "null_data_source" "resource" {
  inputs = {
    scope = "gcp://cloudresourcemanager.googleapis.com/projects/${data.google_client_config.current.project}"
  }
}

resource "google_dns_policy" "named_test_resource" {
  name        = var.resource_name
  enable_inbound_forwarding = var.enable_inbound_forwarding

  enable_logging = var.enable_logging

  alternative_name_server_config {
    target_name_servers {
      ipv4_address    = "172.16.1.10"
      forwarding_path = "private"
    }
    target_name_servers {
      ipv4_address = "172.16.1.20"
    }
  }


  networks {
    network_url = google_compute_network.named_test_resource.id
  }
}

resource "google_compute_network" "named_test_resource" {
  name                    = var.resource_name
  auto_create_subnetworks = false
}

output "resource_aka" {
  value = "gcp://dns.googleapis.com/${google_dns_policy.named_test_resource.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "enable_inbound_forwarding" {
  value = var.enable_inbound_forwarding
}

output "enable_logging" {
  value = var.enable_logging
}

output "network" {
  value = google_compute_network.named_test_resource.self_link
}

output "project_id" {
  value = var.gcp_project
}
