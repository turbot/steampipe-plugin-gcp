
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

resource "google_pubsub_topic" "named_test_resource" {
  name = var.resource_name
}

resource "google_pubsub_subscription" "named_test_resource" {
  name  = var.resource_name
  topic = google_pubsub_topic.named_test_resource.name

  ack_deadline_seconds = 20

  labels = {
    name = var.resource_name
  }

  push_config {
    push_endpoint = "https://example.com/push"

    attributes = {
      x-goog-version = "v1"
    }
  }
}

resource "google_pubsub_subscription_iam_binding" "editor" {
  subscription = google_pubsub_subscription.named_test_resource.id
  role         = "roles/editor"
  members = [
    "allUsers",
  ]
}

output "resource_aka" {
  value = "gcp://pubsub.googleapis.com/${google_pubsub_subscription.named_test_resource.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_pubsub_subscription.named_test_resource.id
}

output "topic_id" {
  value = google_pubsub_topic.named_test_resource.id
}

output "etag" {
  value = google_pubsub_subscription_iam_binding.editor.etag
}

output "project_id" {
  value = var.gcp_project
}
