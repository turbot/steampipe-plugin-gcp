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
  auto_create_subnetworks = false
  project = var.gcp_project
}

resource "google_compute_subnetwork" "named_test_resource" {
  name          = var.resource_name
  ip_cidr_range = "10.2.0.0/28"
  region        = var.gcp_region
  network       = google_compute_network.named_test_resource.id
}

resource "google_vpc_access_connector" "named_test_resource" {
  name          = var.resource_name
  subnet {
    name = google_compute_subnetwork.named_test_resource.name
  }
  machine_type = "e2-standard-4"
  min_instances = 2
  max_instances = 3
}

output "resource_aka" {
  value = "gcp://vpcaccess.googleapis.com/${google_vpc_access_connector.named_test_resource.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_vpc_access_connector.named_test_resource.id
}

output "project_id" {
  value = var.gcp_project
}

output "region_id" {
  value = var.gcp_region
}
