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

resource "google_compute_network" "network_test" {
  name                    = var.resource_name
  provider                = google-beta
  project                 = var.gcp_project
  auto_create_subnetworks = false
}

resource "google_compute_vpn_gateway" "target_gateway" {
  name    = var.resource_name
  description   = "Test VPN target gateway to validate integration test."
  network = google_compute_network.network_test.id
}

resource "google_compute_address" "test_address" {
  name = var.resource_name
  network_tier = "PREMIUM"
}

resource "google_compute_forwarding_rule" "fr_esp" {
  name        = "fr-esp"
  ip_protocol = "ESP"
  ip_address  = google_compute_address.test_address.address
  target      = google_compute_vpn_gateway.target_gateway.id
  network_tier = "PREMIUM"
}

resource "google_compute_forwarding_rule" "fr_udp500" {
  name        = "fr-udp500"
  ip_protocol = "UDP"
  port_range  = "500"
  ip_address  = google_compute_address.test_address.address
  target      = google_compute_vpn_gateway.target_gateway.id
  network_tier = "PREMIUM"
}

resource "google_compute_forwarding_rule" "fr_udp4500" {
  name        = "fr-udp4500"
  ip_protocol = "UDP"
  port_range  = "4500"
  ip_address  = google_compute_address.test_address.address
  target      = google_compute_vpn_gateway.target_gateway.id
  network_tier = "PREMIUM"
}

output "resource_aka" {
  value = "gcp://compute.googleapis.com/${google_compute_vpn_gateway.target_gateway.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_compute_vpn_gateway.target_gateway.id
}

output "description" {
  value = google_compute_vpn_gateway.target_gateway.description
}

output "self_link" {
  value = google_compute_vpn_gateway.target_gateway.self_link
}

output "region" {
  value = var.gcp_region
}

output "project_id" {
  value = var.gcp_project
}
