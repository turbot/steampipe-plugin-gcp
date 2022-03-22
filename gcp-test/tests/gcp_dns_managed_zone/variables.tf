
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

resource "google_dns_managed_zone" "named_test_resource" {
  name        = var.resource_name
  dns_name    = "turbot.com."
  description = "Test managed zone to validate the table outcome."
  labels = {
    name = var.resource_name
  }

  visibility = "private"

  private_visibility_config {
    networks {
      network_url = google_compute_network.named_test_resource.id
    }
  }
}

resource "google_compute_network" "named_test_resource" {
  name                    = var.resource_name
  auto_create_subnetworks = false
}

output "resource_aka" {
  value = "gcp://dns.googleapis.com/${google_dns_managed_zone.named_test_resource.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_dns_managed_zone.named_test_resource.id
}

output "network" {
  value = google_compute_network.named_test_resource.self_link
}

output "project_id" {
  value = var.gcp_project
}
