
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

resource "google_monitoring_alert_policy" "named_test_resource" {
  display_name = var.resource_name
  documentation {
    content   = "This is for testing purpose"
    mime_type = "text/markdown"
  }
  combiner = "OR"
  conditions {
    display_name = var.resource_name
    condition_threshold {
      filter     = "metric.type=\"compute.googleapis.com/instance/disk/write_bytes_count\" AND resource.type=\"gce_instance\""
      duration   = "60s"
      comparison = "COMPARISON_GT"
      aggregations {
        alignment_period   = "60s"
        per_series_aligner = "ALIGN_RATE"
      }
    }
  }
  user_labels = {
    foo = "bar"
  }
}

data "template_file" "resource_aka" {
  template = "gcp://monitoring.googleapis.com/${google_monitoring_alert_policy.named_test_resource.id}"
  vars = {
    resource_name = var.resource_name
    project       = data.google_client_config.current.project
    region        = data.google_client_config.current.region
    zone          = data.google_client_config.current.zone
  }
}

output "resource_aka" {
  depends_on = [google_monitoring_alert_policy.named_test_resource]
  value      = data.template_file.resource_aka.rendered
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_monitoring_alert_policy.named_test_resource.name
}

output "project_id" {
  value = var.gcp_project
}
