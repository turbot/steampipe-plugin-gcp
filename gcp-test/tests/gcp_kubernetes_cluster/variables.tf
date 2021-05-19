
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
  
provider "google" {
  project = var.gcp_project
  region  = var.gcp_region
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

resource "google_container_cluster" "named_test_resource" {
  name               = var.resource_name
  location           = var.gcp_region
  network            = google_compute_network.named_test_resource.name
  initial_node_count = 1

  logging_service    = "none"
  monitoring_service = "none"
}


# resource "google_container_cluster" "named_test_resource" {
#   name     = var.resource_name
#   location = var.gcp_region
#   network = google_compute_network.named_test_resource.name
#   remove_default_node_pool = false
#   initial_node_count       = 1
#   resource_labels = {
#     name = var.resource_name
#   }
# }

# resource "google_container_node_pool" "named_test_resource" {
#   name       = var.resource_name
#   location   = var.gcp_region
#   cluster    = google_container_cluster.named_test_resource.name
#   node_count = 1

#   node_config {
#     preemptible  = true
#     machine_type = "e2-medium"

#     # Google recommends custom service accounts that have cloud-platform scope and permissions granted via IAM Roles.
#     service_account = google_service_account.named_test_resource.email
#     oauth_scopes    = [
#       "https://www.googleapis.com/auth/cloud-platform"
#     ]
#   }
# }

# output "resource_aka" {
#   value = "gcp://compute.googleapis.com/projects/${var.gcp_project}${google_compute_firewall.named_test_resource.id}"
# }

//container.googleapis.com/projects/pikachu-aaa/zones/us-east1-b/clusters/cluster-1

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_container_cluster.named_test_resource.id
}

output "self_link" {
  value = google_container_cluster.named_test_resource.self_link
}

output "services_ipv4_cidr" {
  value = google_container_cluster.named_test_resource.services_ipv4_cidr
}

output "project_id" {
  value = var.gcp_project
}

output "location" {
  value = var.gcp_region
}
