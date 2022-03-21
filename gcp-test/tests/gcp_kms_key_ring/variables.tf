
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

data "google_iam_policy" "resource" {
  binding {
    role = "roles/editor"

    members = [
      "user:sourav@turbot.com",
    ]
  }
}
resource "google_kms_key_ring" "named_test_resource" {
  name     = var.resource_name
  location = "global"
}

resource "google_kms_key_ring_iam_policy" "named_test_resource" {
  key_ring_id = google_kms_key_ring.named_test_resource.id
  policy_data = data.google_iam_policy.resource.policy_data
}

output "etag" {
  value = google_kms_key_ring_iam_policy.named_test_resource.etag
}

output "resource_aka" {
  value = "gcp://cloudkms.googleapis.com/${google_kms_key_ring.named_test_resource.id}"
}

output "resource_name" {
  value = var.resource_name
}
