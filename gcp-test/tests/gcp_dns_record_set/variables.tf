
variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125-create-update"
  description = "Name of the resource used throughout the test."
}

variable "record_set_name" {
  type = string
  default = "test-record.turbot.com."
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

variable "record_set_type" {
  type = string
  default = "A"
}

variable "ttl" {
  type = number
  default = 86400
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

resource "google_dns_managed_zone" "named_test_resource" {
  name        = var.resource_name
  dns_name    = "turbot.com."
  description = "Test managed zone to validate the table outcome."
  labels = {
    name = var.resource_name
  }

  visibility = "private"
}

resource "google_dns_record_set" "named_test_resource" {
  managed_zone = google_dns_managed_zone.named_test_resource.name
  name         = var.record_set_name
  type         = var.record_set_type
  rrdatas      = ["10.0.0.1", "10.1.0.1"]
  ttl          = var.ttl
}

output "resource_aka" {
  value = "gcp://dns.googleapis.com/${google_dns_record_set.named_test_resource.id}"
}

output "record_set_name" {
  value = var.record_set_name
}

output "resource_type" {
  value = var.record_set_type
}

output "resource_id" {
  value = google_dns_record_set.named_test_resource.id
}

output "resource_name" {
  value = var.resource_name
}

output "project_id" {
  value = var.gcp_project
}
