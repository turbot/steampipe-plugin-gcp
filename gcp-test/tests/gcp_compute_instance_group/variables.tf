
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
  default     = "us-central1"
  description = "GCP region used for the test."
}

variable "gcp_zone" {
  type    = string
  default = "us-central1-a"
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

resource "google_compute_network" "named_test_resource" {
  name = var.resource_name
}

resource "google_compute_instance_group" "named_test_resource" {
  name        = var.resource_name
  description = var.resource_name
  zone        = var.gcp_zone
  network     = google_compute_network.named_test_resource.id
}

output "resource_aka" {
  value = "gcp://compute.googleapis.com/projects/${var.gcp_project}/zones/${var.gcp_zone}/instanceGroups/${google_compute_instance_group.named_test_resource.name}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_compute_instance_group.named_test_resource.id
}

output "self_link" {
  value = google_compute_instance_group.named_test_resource.self_link
}

output "size" {
  value = google_compute_instance_group.named_test_resource.size
}

output "project" {
  value = var.gcp_project
}
