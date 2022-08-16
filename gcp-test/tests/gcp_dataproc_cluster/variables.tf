
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

resource "google_compute_network" "named_test_resource" {
  name = var.resource_name
}

resource "google_dataproc_cluster" "named_test_resource" {
  name   = var.resource_name
  region = var.gcp_region
  cluster_config {
    gce_cluster_config {
      zone    = var.gcp_zone
      network = google_compute_network.named_test_resource.name
    }
    master_config {
      num_instances = 1
      machine_type  = "n1-standard-1"
      disk_config {
        boot_disk_type = "pd-ssd"
        boot_disk_size_gb = 15
      }
    }
    worker_config {
      num_instances = 0
    }
    preemptible_worker_config {
      num_instances = 0
    }
    software_config {
      image_version = "1.4-debian10"
      override_properties = {
        "dataproc:dataproc.allow.zero.workers" = "true"
      }
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

