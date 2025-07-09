
variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125-create-update"
  description = "Name of the resource used throughout the test."
}

variable "gcp_project" {
  type        = string
  default     = "nagraj-aaa"
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

resource "google_compute_security_policy" "named_test_resource" {
  name        = var.resource_name
  description = "Test security policy for validation."
  rule {
    action   = "allow"
    priority = 2147483647
    match {
      versioned_expr = "SRC_IPS_V1"
      config {
        src_ip_ranges = ["*"]
      }
    }
    description = "default rule"
  }
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_compute_security_policy.named_test_resource.id
}

output "self_link" {
  value = google_compute_security_policy.named_test_resource.self_link
}

output "fingerprint" {
  value = google_compute_security_policy.named_test_resource.fingerprint
}

output "project_id" {
  value = var.gcp_project
}
