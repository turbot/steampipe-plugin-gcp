
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
    scope = "gcp://dataproc.googleapis.com/projects/${data.google_client_config.current.project}"
  }
}

resource "google_dataproc_metastore_service" "named_test_resource" {
  service_id = var.resource_name
  location   = var.gcp_region
  port       = 9080
  tier       = "DEVELOPER"

  maintenance_window {
    hour_of_day = 2
    day_of_week = "SUNDAY"
  }

  hive_metastore_config {
    version = "2.3.6"
  }

  labels = {
    env = "test"
  }
}


output "project_id" {
  value = data.google_client_config.current.project
}

output "resource_aka" {
  value = "gcp://metastore.googleapis.com/${google_dataproc_metastore_service.named_test_resource.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_dataproc_metastore_service.named_test_resource.id
}

output "state" {
  value = google_dataproc_metastore_service.named_test_resource.state
}

output "uid" {
  value = google_dataproc_metastore_service.named_test_resource.uid
}

