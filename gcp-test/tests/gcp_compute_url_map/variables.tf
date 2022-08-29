
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

resource "google_storage_bucket" "named_test_resource" {
  name     = var.resource_name
  location = "US"
}

resource "google_compute_backend_bucket" "static" {
  name        = var.resource_name
  bucket_name = google_storage_bucket.named_test_resource.name
  enable_cdn  = true
}

resource "google_compute_http_health_check" "default" {
  name               = var.resource_name
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}

resource "google_compute_backend_service" "login" {
  name        = "${var.resource_name}-login"
  port_name   = "http"
  protocol    = "HTTP"
  timeout_sec = 10

  health_checks = [google_compute_http_health_check.default.id]
}

resource "google_compute_backend_service" "home" {
  name        = "${var.resource_name}-home"
  port_name   = "http"
  protocol    = "HTTP"
  timeout_sec = 10

  health_checks = [google_compute_http_health_check.default.id]
}

resource "google_compute_url_map" "named_test_resource" {
  name        = var.resource_name
  description = "Test URL Map to validate integration test."

  default_service = google_compute_backend_service.home.id

  host_rule {
    hosts        = ["mysite.com"]
    path_matcher = "mysite"
  }

  host_rule {
    hosts        = ["myothersite.com"]
    path_matcher = "otherpaths"
  }

  path_matcher {
    name            = "mysite"
    default_service = google_compute_backend_service.home.id

    path_rule {
      paths   = ["/home"]
      service = google_compute_backend_service.home.id
    }

    path_rule {
      paths   = ["/login"]
      service = google_compute_backend_service.login.id
    }

    path_rule {
      paths   = ["/static"]
      service = google_compute_backend_bucket.static.id
    }
  }

  path_matcher {
    name            = "otherpaths"
    default_service = google_compute_backend_service.home.id
  }

  test {
    service = google_compute_backend_service.home.id
    host    = "hi.com"
    path    = "/home"
  }
}

output "resource_aka" {
  value = "gcp://compute.googleapis.com/${google_compute_url_map.named_test_resource.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "self_link" {
  value = google_compute_url_map.named_test_resource.self_link
}

output "default_service" {
  value = google_compute_backend_service.home.self_link
}

output "project_id" {
  value = var.gcp_project
}

output "backend_service_home" {
  value = google_compute_backend_service.home.self_link
}

output "backend_service_login" {
  value = google_compute_backend_service.login.self_link
}

output "backend_bucket" {
  value = google_compute_backend_bucket.static.self_link
}
