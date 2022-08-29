
variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125-create-update"
  description = "Name of the resource used throughout the test."
}

variable "gcp_project" {
  type        = string
  default     = "parker-aaa"
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

provider "google-beta" {
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

resource "google_compute_node_template" "named_test_resource" {
  provider = google-beta

  name        = var.resource_name
  description = "Test Node Template resource to verify the table."
  region      = var.gcp_region
  node_type   = "c2-node-60-240"

  node_affinity_labels = {
    name = var.resource_name
  }

  server_binding {
    type = "RESTART_NODE_ON_MINIMAL_SERVERS"
  }
}

output "resource_aka" {
  value = "gcp://compute.googleapis.com/${google_compute_node_template.named_test_resource.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_compute_node_template.named_test_resource.id
}

output "self_link" {
  value = google_compute_node_template.named_test_resource.self_link
}

output "project_id" {
  value = var.gcp_project
}
