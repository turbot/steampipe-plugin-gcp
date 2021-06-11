
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

resource "google_service_account" "default" {
  account_id   = var.resource_name
  display_name = var.resource_name
}

resource "google_compute_network" "named_test_resource" {
  name = var.resource_name
}

resource "google_container_cluster" "named_test_resource" {
  name     = var.resource_name
  location = var.gcp_region
  network  = google_compute_network.named_test_resource.name

  # We can't create a cluster with no node pool defined, but we want to only use
  # separately managed node pools. So we create the smallest possible default
  # node pool and immediately delete it.
  remove_default_node_pool = true
  initial_node_count       = 1
}

resource "google_container_node_pool" "named_test_resource" {
  name       = var.resource_name
  location   = var.gcp_region
  cluster    = google_container_cluster.named_test_resource.name
  node_count = 1

  node_config {
    preemptible  = true
    machine_type = "n1-standard-1"

    # Google recommends custom service accounts that have cloud-platform scope and permissions granted via IAM Roles.
    service_account = google_service_account.default.email
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
  }
}


output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_container_node_pool.named_test_resource.id
}

output "resource_aka" {
  value = "gcp://container.googleapis.com/v1/${google_container_node_pool.named_test_resource.id}"
}

output "cluster_name" {
  value = google_container_cluster.named_test_resource.name
}

output "project_id" {
  value = var.gcp_project
}

output "location" {
  value = var.gcp_region
}
