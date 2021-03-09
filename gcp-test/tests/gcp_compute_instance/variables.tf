
variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125-create-update"
  description = "Name of the resource used throughout the test."
}

variable "gcp_project" {
  type        = string
  default     = "niteowl-aaa"
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

resource "google_compute_network" "named_test_resource" {
  name = var.resource_name
}

resource "google_service_account" "named_test_resource" {
  account_id   = var.resource_name
  display_name = var.resource_name
}

resource "google_compute_instance" "named_test_resource" {
  name         = var.resource_name
  machine_type = "f1-micro"
  zone         = "us-east1-b"
  description  = "Test VM instance to verify the table."

  tags = ["foo", "bar"]

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

  network_interface {
    network = google_compute_network.named_test_resource.name

    access_config {
      // Ephemeral IP
    }
  }

  metadata = {
    foo = "bar"
  }

  labels = {
    name = var.resource_name
  }

  metadata_startup_script = "echo hi > /test.txt"

  service_account {
    # Google recommends custom service accounts that have cloud-platform scope and permissions granted via IAM Roles.
    email  = google_service_account.named_test_resource.email
    scopes = ["cloud-platform"]
  }
}

resource "google_compute_instance_iam_member" "member" {
  project       = google_compute_instance.named_test_resource.project
  zone          = google_compute_instance.named_test_resource.zone
  instance_name = google_compute_instance.named_test_resource.name
  role          = "roles/compute.osLogin"
  member        = "serviceAccount:${google_service_account.named_test_resource.email}"
}

output "resource_aka" {
  value = "gcp://compute.googleapis.com/${google_compute_instance.named_test_resource.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "self_link" {
  value = google_compute_instance.named_test_resource.self_link
}

output "network" {
  value = google_compute_network.named_test_resource.self_link
}

output "service_account" {
  value = google_service_account.named_test_resource.email
}

output "metadata_fingerprint" {
  value = google_compute_instance.named_test_resource.metadata_fingerprint
}

output "cpu_platform" {
  value = google_compute_instance.named_test_resource.cpu_platform
}

output "etag" {
  value = google_compute_instance_iam_member.member.etag
}

output "label_fingerprint" {
  value = google_compute_instance.named_test_resource.label_fingerprint
}

output "project_id" {
  value = var.gcp_project
}
