
variable "resource_name" {
  type        = string
  default     = "turbot-test-20200125-create-update"
  description = "Name of the resource used throughout the test."
}

variable "gcp_project" {
  type        = string
  default     = "pikachu-aaa"
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

resource "google_compute_vpn_gateway" "named_test_resource" {
  name    = var.resource_name
  network = google_compute_network.named_test_resource.id
  description = "Test to verify compute vpn gateway table."
}

resource "google_compute_network" "named_test_resource" {
  name = var.resource_name
}

resource "google_compute_address" "named_test_resource" {
  name = var.resource_name
}

resource "google_compute_forwarding_rule" "named_test_resource" {
  name        = var.resource_name
  ip_protocol = "ESP"
  ip_address  = google_compute_address.named_test_resource.address
  target      = google_compute_vpn_gateway.named_test_resource.id
}

resource "google_compute_forwarding_rule" "named_test_resource1" {
  name        = var.resource_name
  ip_protocol = "UDP"
  port_range  = "500"
  ip_address  = google_compute_address.named_test_resource.address
  target      = google_compute_vpn_gateway.named_test_resource.id
}

resource "google_compute_forwarding_rule" "named_test_resource2" {
  name        = var.resource_name
  ip_protocol = "UDP"
  port_range  = "4500"
  ip_address  = google_compute_address.named_test_resource.address
  target      = google_compute_vpn_gateway.named_test_resource.id
}

resource "google_compute_vpn_tunnel" "named_test_resource" {
  name          = var.resource_name
  peer_ip       = "15.0.0.120"
  shared_secret = "a secret message"

  target_vpn_gateway = google_compute_vpn_gateway.named_test_resource.id

  depends_on = [
    google_compute_forwarding_rule.named_test_resource,
    google_compute_forwarding_rule.named_test_resource1,
    google_compute_forwarding_rule.named_test_resource2,
  ]
}

resource "google_compute_route" "named_test_resource" {
  name       = var.resource_name
  network    = google_compute_network.named_test_resource.name
  dest_range = "15.0.0.0/24"
  priority   = 1000

  next_hop_vpn_tunnel = google_compute_vpn_tunnel.named_test_resource.id
}

output "resource_aka" {
  value = "gcp://compute.googleapis.com/${google_compute_vpn_gateway.named_test_resource.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_compute_vpn_gateway.named_test_resource.id
}

output "description" {
  value = google_compute_vpn_gateway.named_test_resource.description
}

output "self_link" {
  value = google_compute_vpn_gateway.named_test_resource.self_link
}

output "forwarding_rules" {
  value = google_compute_forwarding_rule.named_test_resource.self_link
}

output "region" {
  value = var.gcp_region
}

output "project_id" {
  value = var.gcp_project
}
