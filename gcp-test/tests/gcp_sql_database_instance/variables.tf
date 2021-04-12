
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

resource "google_sql_database_instance" "named_test_resource" {
  name                = var.resource_name
  database_version    = "MYSQL_5_6"
  region              = var.gcp_region
  deletion_protection = false

  settings {
    tier      = "db-f1-micro"
    disk_size = 10
    disk_type = "PD_HDD"
    backup_configuration {
      enabled = true
    }

    user_labels = {
      name = var.resource_name
    }
  }
}

resource "google_sql_user" "users" {
  name     = var.resource_name
  instance = google_sql_database_instance.named_test_resource.name
  host     = "me.com"
  password = "changeme"
}

output "resource_aka" {
  value = "gcp://cloudsql.googleapis.com/projects/${var.gcp_project}/regions/${var.gcp_region}/instances/${var.resource_name}"
}

output "resource_name" {
  value = var.resource_name
}

output "self_link" {
  value = google_sql_database_instance.named_test_resource.self_link
}

output "project_id" {
  value = var.gcp_project
}
