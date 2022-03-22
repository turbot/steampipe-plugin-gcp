
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

data "google_compute_default_service_account" "default" {
}

resource "google_compute_project_metadata" "named_test_resource" {
  metadata = {
    oslogin = "TRUE"
    fizz    = "buzz"
  }
}

output "project_aka" {
  value = "gcp://cloudresourcemanager.googleapis.com/projects/${data.google_client_config.current.project}"
}

output "project_id" {
  value = data.google_client_config.current.project
}

output "service_account" {
  value = data.google_compute_default_service_account.default.email
}
