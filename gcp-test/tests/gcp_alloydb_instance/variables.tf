
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
    scope = "gcp://alloydb.googleapis.com/projects/${data.google_client_config.current.project}"
  }
}

resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = var.resource_name
  instance_type = "PRIMARY"

  machine_config {
    cpu_count = 2
  }

  depends_on = [google_service_networking_connection.vpc_connection]
}

resource "google_alloydb_cluster" "default" {
  cluster_id = var.resource_name
  location   = "us-central1"
  network_config {
    network = google_compute_network.default.id
  }

  initial_user {
    password = "alloydb-cluster"
  }
}

resource "google_compute_network" "default" {
  name = var.resource_name
}

resource "google_compute_global_address" "private_ip_alloc" {
  name          =  var.resource_name
  address_type  = "INTERNAL"
  purpose       = "VPC_PEERING"
  prefix_length = 16
  network       = google_compute_network.default.id
}

resource "google_service_networking_connection" "vpc_connection" {
  network                 = google_compute_network.default.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
}

output "project_id" {
  value = data.google_client_config.current.project
}

output "resource_aka" {
  value = "gcp://alloydb.googleapis.com/${google_alloydb_instance.default.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_alloydb_instance.default.id
}

