
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

resource "google_compute_instance" "names_test_resource" {
  provider     = google-beta
  name         = var.resource_name
  machine_type = "f1-micro"
  zone         = "us-east1-b"
  project      = var.gcp_project

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    network = "default"
  }
}

resource "google_compute_machine_image" "names_test_resource" {
  provider        = google-beta
  project         = var.gcp_project
  name            = var.resource_name
  source_instance = google_compute_instance.names_test_resource.self_link
}

output "machine_type" {
  value = "f1-micro"
}

output "resource_name" {
  value = var.resource_name
}

output "self_link" {
  value = google_compute_machine_image.names_test_resource.self_link
}

output "resource_aka" {
  value = "gcp://compute.googleapis.com/projects/${var.gcp_project}/machineImages/${var.resource_name}"
}
