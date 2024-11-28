
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
    scope = "gcp://alloydb.googleapis.com/projects/${data.google_client_config.current.project}"
  }
}

resource "google_dataplex_lake" "named_test_resource" {
  location     = "us-central1"
  name         = var.resource_name
  description  = "Lake for DCL"
  display_name = var.resource_name
  project      = data.google_client_config.current.project

  labels = {
    name = var.resource_name
  }
}

output "project_id" {
  value = data.google_client_config.current.project
}

output "resource_aka" {
  value = "gcp://dataplex.googleapis.com/${google_dataplex_lake.named_test_resource.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_dataplex_lake.named_test_resource.id
}

