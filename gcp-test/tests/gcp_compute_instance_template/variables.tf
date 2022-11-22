
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

resource "google_compute_network" "named_test_resource" {
  name = var.resource_name
}

resource "google_service_account" "named_test_resource" {
  account_id   = var.resource_name
  display_name = var.resource_name
}

resource "google_compute_image" "named_test_resource" {
  name = var.resource_name

  raw_disk {
    source = "https://storage.googleapis.com/bosh-gce-raw-stemcells/bosh-stemcell-97.98-google-kvm-ubuntu-xenial-go_agent-raw-1557960142.tar.gz"
  }
}

data "google_compute_image" "named_test_resource" {
  depends_on = [
    google_compute_image.named_test_resource
  ]
  name = var.resource_name
}

resource "google_compute_instance_template" "named_test_resource" {
  name        = var.resource_name
  description = "Test instance template to verify the table."

  tags = ["foo", "bar"]

  labels = {
    name = var.resource_name
  }

  instance_description = "A dummy description"
  machine_type         = "f1-micro"
  can_ip_forward       = false

  scheduling {
    automatic_restart   = true
    on_host_maintenance = "MIGRATE"
  }

  // Create a new boot disk from an image
  disk {
    source_image = data.google_compute_image.named_test_resource.name
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = google_compute_network.named_test_resource.name
  }

  metadata = {
    foo = "bar"
  }

  service_account {
    # Google recommends custom service accounts that have cloud-platform scope and permissions granted via IAM Roles.
    email  = google_service_account.named_test_resource.email
    scopes = ["cloud-platform"]
  }
}

output "resource_aka" {
  value = "gcp://compute.googleapis.com/${google_compute_instance_template.named_test_resource.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_compute_instance_template.named_test_resource.id
}

output "self_link" {
  value = google_compute_instance_template.named_test_resource.self_link
}

output "network" {
  value = google_compute_network.named_test_resource.self_link
}

output "service_account" {
  value = google_service_account.named_test_resource.email
}

output "metadata_fingerprint" {
  value = google_compute_instance_template.named_test_resource.metadata_fingerprint
}

output "project_id" {
  value = var.gcp_project
}
