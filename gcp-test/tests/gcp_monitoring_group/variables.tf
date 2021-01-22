
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

resource "google_monitoring_group" "parent" {
  display_name = "parent-${var.resource_name}"
  filter       = "resource.metadata.region=\"europe-west2\""
}

resource "google_monitoring_group" "named_test_resource" {
  display_name = var.resource_name
  filter       = "resource.metadata.region=\"europe-west2\""
  parent_name  = google_monitoring_group.parent.name
}

output "resource_aka" {
  value = "gcp://monitoring.googleapis.com/${google_monitoring_group.named_test_resource.name}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_monitoring_group.named_test_resource.name
}

output "parent_resource_id" {
  value = google_monitoring_group.parent.name
}

output "project_id" {
  value = var.gcp_project
}
