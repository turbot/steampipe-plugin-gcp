
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

resource "google_compute_disk" "named_test_resource" {
  name = var.resource_name
  labels = {
    name = var.resource_name
  }
}

resource "google_service_account" "named_test_resource" {
  account_id   = var.resource_name
  display_name = var.resource_name
  description  = "Test service account to verify the table"
}

resource "google_compute_image" "named_test_resource" {
  name        = var.resource_name
  description = "Test image to verify the table."
  labels = {
    name = var.resource_name
  }
  source_disk = google_compute_disk.named_test_resource.id
}

resource "google_compute_image_iam_member" "member" {
  project = google_compute_image.named_test_resource.project
  image   = google_compute_image.named_test_resource.name
  role    = "roles/compute.imageUser"
  member  = "serviceAccount:${google_service_account.named_test_resource.email}"
}

output "resource_aka" {
  value = "gcp://compute.googleapis.com/${google_compute_image.named_test_resource.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_compute_image.named_test_resource.id
}

output "self_link" {
  value = google_compute_image.named_test_resource.self_link
}

output "source_disk" {
  value = google_compute_disk.named_test_resource.self_link
}

output "etag" {
  value = google_compute_image_iam_member.member.etag
}

output "email" {
  value = google_service_account.named_test_resource.email
}

output "project_id" {
  value = var.gcp_project
}
