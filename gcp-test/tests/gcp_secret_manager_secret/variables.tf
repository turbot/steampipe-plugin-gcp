
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

data "google_project" "current" {}

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

resource "google_secret_manager_secret" "named_test_resource" {
  secret_id = var.resource_name

  labels = {
    label = var.resource_name
  }

  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
      replicas {
        location = "us-east1"
      }
    }
  }
}

output "resource_aka" {
  value = "gcp://secretmanager.googleapis.com/projects/${data.google_project.current.number}/secrets/${var.resource_name}"
}

output "resource_id" {
  value = "projects/${data.google_project.current.number}/secrets/${var.resource_name}"
}

output "resource_name" {
  value = var.resource_name
}

output "project_id" {
  value = var.gcp_project
}

output "project_number" {
  value = data.google_project.current.number
}
