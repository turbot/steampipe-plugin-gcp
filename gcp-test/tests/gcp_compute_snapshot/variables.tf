
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

resource "google_compute_snapshot" "named_test_resource" {
  name        = var.resource_name
  description = "Test compute snapshot to verify table."
  source_disk = google_compute_disk.named_test_resource.name
  zone        = "us-central1-a"
  labels = {
    my_label = "value"
  }
  storage_locations = ["us-central1"]
}

data "google_compute_image" "named_test_resource" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_disk" "named_test_resource" {
  name  = var.resource_name
  image = data.google_compute_image.named_test_resource.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

output "project_id" {
  value = data.google_client_config.current.project
}

output "resource_aka" {
  value = "gcp://compute.googleapis.com/${google_compute_snapshot.named_test_resource.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_compute_snapshot.named_test_resource.id
}

output "self_link" {
  value = google_compute_snapshot.named_test_resource.self_link
}

output "source_disk" {
  value = google_compute_disk.named_test_resource.self_link
}
