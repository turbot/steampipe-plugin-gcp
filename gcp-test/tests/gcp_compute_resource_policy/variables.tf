
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

resource "google_compute_resource_policy" "named_test_resource" {
  name        = var.resource_name
  region      = var.gcp_region
  description = "Start and stop instances"
  instance_schedule_policy {
    vm_start_schedule {
      schedule = "0 * * * *"
    }
    vm_stop_schedule {
      schedule = "15 * * * *"
    }
    time_zone = "US/Central"
  }
}

output "resource_aka" {
  value = "gcp://compute.googleapis.com/${google_compute_resource_policy.named_test_resource.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "self_link" {
  value = google_compute_resource_policy.named_test_resource.self_link
}

output "project_id" {
  value = var.gcp_project
}
