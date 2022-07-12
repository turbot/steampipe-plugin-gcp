# Make sure that the `gcloud` is already installed, before running this test.

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

locals {
  path = "${path.cwd}/output.json"
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

# Create GCP > SQL > Backup
resource "null_resource" "named_test_resource" {
  depends_on = [google_sql_database_instance.named_test_resource]
  provisioner "local-exec" {
    command = "gcloud sql backups create --instance ${var.resource_name}"
  }
}

# List GCP > SQL > Backup and store the output in a local file
resource "null_resource" "list_resource" {
  depends_on = [null_resource.named_test_resource]
  provisioner "local-exec" {
    command = "gcloud sql backups list --instance ${var.resource_name} --format json > ${local.path}"
  }
}

data "local_file" "input" {
  depends_on = [null_resource.list_resource]
  filename   = local.path
}

output "backup_id" {
  value = lookup(jsondecode(replace(replace(data.local_file.input.content, "[\n  ", ""), "\n]\n", "")), "id", "what_1")
}

output "self_link" {
  value = lookup(jsondecode(replace(replace(data.local_file.input.content, "[\n  ", ""), "\n]\n", "")), "selfLink", "what_1")
}

output "project_id" {
  value = var.gcp_project
}

output "resource_name" {
  value = var.resource_name
}
