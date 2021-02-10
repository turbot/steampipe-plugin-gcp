
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

resource "google_compute_backend_service" "named_test_resource" {
  name          = var.resource_name
  description   = "Test Backend Service to verify the table."
  health_checks = [google_compute_http_health_check.named_test_resource.id]
  enable_cdn    = true
  cdn_policy {
    signed_url_cache_max_age_sec = 7200
  }
}

resource "google_compute_http_health_check" "named_test_resource" {
  name               = var.resource_name
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}

output "resource_aka" {
  value = "gcp://compute.googleapis.com/${google_compute_backend_service.named_test_resource.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_compute_backend_service.named_test_resource.id
}

output "self_link" {
  value = google_compute_backend_service.named_test_resource.self_link
}

output "health_check" {
  value = google_compute_http_health_check.named_test_resource.self_link
}

output "project_id" {
  value = var.gcp_project
}
