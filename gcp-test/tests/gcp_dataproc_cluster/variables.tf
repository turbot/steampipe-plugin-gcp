
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

provider "google" {
  project = var.gcp_project
  region  = var.gcp_region
}

data "google_client_config" "current" {}

data "null_data_source" "resource" {
  inputs = {
    scope = "gcp://dataproc.googleapis.com/projects/${data.google_client_config.current.project}"
  }
}

# resource "google_compute_network" "named_test_resource" {
#   name                    = var.resource_name
#   mtu                     = 1500
#   auto_create_subnetworks = true
#   description             = "Test network to validate integration test."
# }

resource "google_dataproc_cluster" "named_test_resource" {
  name   = var.resource_name
  region = var.gcp_region
  cluster_config {
    gce_cluster_config {
      zone = "us-east1-b"

      # One of the below to hook into a custom network / subnetwork
      # network    = google_compute_network.named_test_resource.name
      subnetwork = "projects/parker-aaa/regions/us-east1/subnetworks/test21"

      tags = ["foo", "bar"]
    }
  }
}


output "project_id" {
  value = data.google_client_config.current.project
}

output "resource_aka" {
  value = "gcp://dataproc.googleapis.com/${google_dataproc_cluster.named_test_resource.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_dataproc_cluster.named_test_resource.id
}

