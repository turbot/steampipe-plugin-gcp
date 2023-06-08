
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
  default = "us-central1-a"
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

resource "google_compute_image" "named_test_resource" {
  name = var.resource_name

  raw_disk {
    source = "https://storage.googleapis.com/bosh-gce-raw-stemcells/bosh-stemcell-97.98-google-kvm-ubuntu-xenial-go_agent-raw-1557960142.tar.gz"
  }
}

data "google_compute_image" "my_image" {
  depends_on = [
    google_compute_image.named_test_resource
  ]
  name = var.resource_name
}

resource "google_compute_resource_policy" "policy" {
  name   = "my-resource-policy"
  region = "us-central1"
  snapshot_schedule_policy {
    schedule {
      daily_schedule {
        days_in_cycle = 1
        start_time    = "04:00"
      }
    }
  }
}

resource "google_compute_disk" "named_test_resource" {
  name        = var.resource_name
  description = var.resource_name
  type        = "pd-ssd"
  zone        = "us-central1-a"
  size        = 50
  image       = data.google_compute_image.my_image.self_link
  labels = {
    name = var.resource_name
  }
}

resource "google_compute_disk_resource_policy_attachment" "attachment" {
  name = google_compute_resource_policy.policy.name
  disk = google_compute_disk.named_test_resource.name
  zone = "us-central1-a"
}

data "google_iam_policy" "admin" {
  binding {
    role = "roles/compute.instanceAdmin"
    members = [
      "serviceAccount:test-685@parker-aaa.iam.gserviceaccount.com",
    ]
  }
}

resource "google_compute_disk_iam_policy" "policy" {
  project     = google_compute_disk.named_test_resource.project
  zone        = google_compute_disk.named_test_resource.zone
  name        = google_compute_disk.named_test_resource.name
  policy_data = data.google_iam_policy.admin.policy_data
}

output "resource_aka" {
  value = "gcp://compute.googleapis.com/${google_compute_disk.named_test_resource.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_compute_disk.named_test_resource.id
}

output "self_link" {
  value = google_compute_disk.named_test_resource.self_link
}

output "project" {
  value = var.gcp_project
}

output "etag" {
  value = google_compute_disk_iam_policy.policy.etag
}

