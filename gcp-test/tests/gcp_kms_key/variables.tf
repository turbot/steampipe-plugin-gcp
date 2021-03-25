
variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125-create-update"
  description = "Name of the resource used throughout the test."
}

variable "gcp_project" {
  type        = string
  default     = "pikachu-aaa"
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

resource "google_kms_key_ring" "named_test_resource" {
  name     = var.resource_name
  location = "global"
}

resource "google_kms_crypto_key" "named_test_resource" {
  name            = var.resource_name
  key_ring        = google_kms_key_ring.named_test_resource.id
  rotation_period = "100000s"
}

output "resource_aka" {
  value = "gcp://cloudkms.googleapis.com/${google_kms_crypto_key.named_test_resource.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_kms_crypto_key.named_test_resource.id
}

output "project_id" {
  value = var.gcp_project
}
