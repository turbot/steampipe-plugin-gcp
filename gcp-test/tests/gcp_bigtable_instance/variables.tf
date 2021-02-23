
variable "resource_name" {
  type        = string
  default     = "turbottest20200125createupdate"
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

locals {
  resource = {
    scope = "gcp://cloudresourcemanager.googleapis.com/projects/${data.google_client_config.current.project}"
  }
}

resource "google_service_account" "named_test_resource" {
  account_id   = var.resource_name
  display_name = var.resource_name
}

resource "google_bigtable_instance" "named_test_resource" {
  name                = var.resource_name
  instance_type       = "DEVELOPMENT"
  deletion_protection = false

  cluster {
    cluster_id   = var.resource_name
    zone         = data.google_client_config.current.zone
    storage_type = "HDD"
  }

  labels = {
    name = var.resource_name
  }
}

resource "google_bigtable_instance_iam_member" "member" {
  instance = var.resource_name
  role     = "roles/bigtable.user"
  member   = "serviceAccount:${google_service_account.named_test_resource.email}"

  depends_on = [
    google_bigtable_instance.named_test_resource
  ]
}

output "resource_aka" {
  value = "gcp://bigtableadmin.googleapis.com/${google_bigtable_instance.named_test_resource.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_bigtable_instance.named_test_resource.id
}

output "project_id" {
  value = var.gcp_project
}

output "etag" {
  value = google_bigtable_instance_iam_member.member.etag
}

output "email" {
  value = google_service_account.named_test_resource.email
}
