
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

resource "google_service_account" "named_test_resource" {
  account_id   = var.resource_name
  display_name = var.resource_name
  description  = "Test service account to verify the table"
}

resource "google_service_account_key" "named_test_resource" {
  service_account_id = google_service_account.named_test_resource.name
  public_key_type    = "TYPE_X509_PEM_FILE"
}

output "resource_aka" {
  value = "gcp://iam.googleapis.com/${google_service_account_key.named_test_resource.id}"
}

output "resource_id" {
  value = google_service_account_key.named_test_resource.id
}

output "service_account_name" {
  value = split("/", google_service_account_key.named_test_resource.id)[3]
}

output "name" {
  value = split("/", google_service_account_key.named_test_resource.id)[5]
}

output "project_id" {
  value = var.gcp_project
}
