
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
  default     = "us-central1"
  description = "GCP region used for the test."
}

variable "gcp_zone" {
  type    = string
  default = "us-central1-f"
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

resource "google_service_account" "named_test_resource" {
  account_id   = var.resource_name
  display_name = var.resource_name
}

resource "google_compute_autoscaler" "named_test_resource" {
  name        = var.resource_name
  description = var.resource_name
  zone        = var.gcp_zone

  target = google_compute_instance_group_manager.named_test_resource.id

  autoscaling_policy {
    max_replicas    = 1
    min_replicas    = 0
    cooldown_period = 60
  }
}

resource "google_compute_instance_template" "named_test_resource" {
  name           = var.resource_name
  machine_type   = "e2-micro"
  can_ip_forward = false

  tags = ["foo", "bar"]

  disk {
    source_image = data.google_compute_image.debian_9.id
  }

  network_interface {
    network = "default"
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

resource "google_compute_target_pool" "named_test_resource" {
  name = var.resource_name
}

resource "google_compute_instance_group_manager" "named_test_resource" {
  name = var.resource_name
  zone = var.gcp_zone

  version {
    instance_template = google_compute_instance_template.named_test_resource.id
    name              = "primary"
  }

  target_pools       = [google_compute_target_pool.named_test_resource.id]
  base_instance_name = var.resource_name
}

data "google_compute_image" "debian_9" {
  family  = "debian-11"
  project = "debian-cloud"
}

output "resource_aka" {
  value = "gcp://compute.googleapis.com/${google_compute_autoscaler.named_test_resource.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_compute_autoscaler.named_test_resource.id
}

output "self_link" {
  value = google_compute_autoscaler.named_test_resource.self_link
}

output "project_id" {
  value = var.gcp_project
}
