
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

resource "google_project_organization_policy" "named_test_resource" {
  project    = var.gcp_project
  constraint = "serviceuser.services"

  list_policy {
    allow {
      all = true
    }
  }
}

output "project_aka" {
  value = "gcp://cloudresourcemanager.googleapis.com/projects/${var.gcp_project}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_title" {
  value = google_project_organization_policy.named_test_resource.constraint
}

output "resource_id" {
  value = split(":", google_project_organization_policy.named_test_resource.id)[1]
}

output "project_id" {
  value = var.gcp_project
}
