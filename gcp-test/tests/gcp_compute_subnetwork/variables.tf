
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

resource "google_compute_network" "named_test_resource" {
  name                    = var.resource_name
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "named_test_resource" {
  name          = var.resource_name
  description   = "Test subnetwork to validate integration test."
  ip_cidr_range = "10.2.0.0/16"
  region        = "us-east1"
  network       = google_compute_network.named_test_resource.id
  secondary_ip_range {
    range_name    = var.resource_name
    ip_cidr_range = "192.168.10.0/24"
  }
}

data "google_iam_policy" "admin" {
  binding {
    role = "roles/compute.networkUser"
    members = [
      "user:test@example.com",
    ]
  }
}

resource "google_compute_subnetwork_iam_policy" "policy" {
  project     = google_compute_subnetwork.named_test_resource.project
  region      = google_compute_subnetwork.named_test_resource.region
  subnetwork  = google_compute_subnetwork.named_test_resource.name
  policy_data = data.google_iam_policy.admin.policy_data
}
output "project_id" {
  value = data.google_client_config.current.project
}

output "resource_aka" {
  value = "gcp://compute.googleapis.com/${google_compute_subnetwork.named_test_resource.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_compute_subnetwork.named_test_resource.id
}

output "self_link" {
  value = google_compute_subnetwork.named_test_resource.self_link
}

output "network" {
  value = google_compute_subnetwork.named_test_resource.network
}

output "etag" {
  value = google_compute_subnetwork_iam_policy.policy.etag
}
