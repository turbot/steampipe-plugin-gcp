
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

provider "google" {
  project = var.gcp_project
  region  = var.gcp_region
}

data "google_client_config" "current" {}

data "null_data_source" "resource" {
  inputs = {
    scope = "gcp://cloudresourcemanager.googleapis.com/projects/${data.google_client_config.current.project}"
  }
}

resource "google_compute_node_template" "named_test_resource" {
  name      = var.resource_name
  region    = "us-central1"
  node_type = "c2-node-60-240"
}

resource "google_compute_node_group" "named_test_resource" {
  name               = var.resource_name
  zone               = "us-central1-a"
  description        = "example google_compute_node_group for Terraform Google Provider"
  maintenance_policy = "RESTART_IN_PLACE"
  size               = 1
  node_template      = google_compute_node_template.named_test_resource.id
  autoscaling_policy {
    mode      = "ONLY_SCALE_OUT"
    min_nodes = 1
    max_nodes = 10
  }
}

output "resource_aka" {
  value = "gcp://compute.googleapis.com/${google_compute_node_group.named_test_resource.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_compute_node_group.named_test_resource.id
}

output "self_link" {
  value = google_compute_node_group.named_test_resource.self_link
}

output "node_template" {
  value = google_compute_node_template.named_test_resource.self_link
}

output "project_id" {
  value = var.gcp_project
}
