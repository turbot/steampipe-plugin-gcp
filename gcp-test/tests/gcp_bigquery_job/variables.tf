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

resource "google_bigquery_dataset" "named_test_resource" {
  dataset_id = var.resource_name
}

resource "google_bigquery_table" "named_test_resource" {
  depends_on          = [google_bigquery_dataset.named_test_resource]
  deletion_protection = false
  dataset_id          = google_bigquery_dataset.named_test_resource.dataset_id
  table_id            = var.resource_name
}

resource "google_bigquery_job" "named_test_resource" {
  depends_on = [google_bigquery_table.named_test_resource]
  job_id     = var.resource_name
  query {
    query = "SELECT state FROM ${google_bigquery_table.named_test_resource.id}"
  }
}

output "resource_aka" {
  depends_on = [google_bigquery_job.named_test_resource]
  value      = "gcp://bigquery.googleapis.com/${google_bigquery_job.named_test_resource.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  depends_on = [google_bigquery_job.named_test_resource]
  value      = "${var.gcp_project}:${google_bigquery_job.named_test_resource.location}.${var.resource_name}"
}

output "project_id" {
  value = var.gcp_project
}

output "region_id" {
  depends_on = [google_bigquery_job.named_test_resource]
  value      = google_bigquery_job.named_test_resource.location
}
