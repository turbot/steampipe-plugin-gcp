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

resource "google_service_account" "named_test_resource" {
  account_id = var.resource_name
}

resource "google_bigquery_dataset" "named_test_resource" {
  dataset_id                  = var.resource_name
  friendly_name               = "turbot_${var.resource_name}"
  description                 = "This is a test dataset to validate the table outcome."
  location                    = var.gcp_region
  default_table_expiration_ms = 3600000

  labels = {
    name = var.resource_name
  }

  access {
    role          = "OWNER"
    user_by_email = google_service_account.named_test_resource.email
  }

  access {
    role   = "READER"
    domain = "hashicorp.com"
  }
}

output "resource_aka" {
  value = "gcp://bigquery.googleapis.com/${google_bigquery_dataset.named_test_resource.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_bigquery_dataset.named_test_resource.id
}

output "service_account_email" {
  value = google_service_account.named_test_resource.email
}

output "self_link" {
  value = google_bigquery_dataset.named_test_resource.self_link
}

output "etag" {
  value = google_bigquery_dataset.named_test_resource.etag
}

output "project_id" {
  value = var.gcp_project
}

output "region_id" {
  value = var.gcp_region
}
