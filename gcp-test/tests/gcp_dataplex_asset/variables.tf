
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

resource "google_service_account" "named_test_resource" {
  account_id = var.resource_name
}

resource "google_bigquery_dataset" "named_test_resource" {
  dataset_id                  = var.resource_name
  friendly_name               = "turbot_${var.resource_name}"
  description                 = "This is a test dataset to validate the table outcome."
  location                    = "us-central1"
  default_table_expiration_ms = 3600000

  labels = {
    name = var.resource_name
  }

  access {
    role          = "OWNER"
    user_by_email = google_service_account.named_test_resource.email
  }

  access {
    role   = "READER"
    domain = "hashicorp.com"
  }
}

resource "google_dataplex_lake" "basic_lake" {
  name         = var.resource_name
  location     = "us-central1"
  project = var.gcp_project
}


resource "google_dataplex_zone" "basic_zone" {
  name         = "${var.resource_name}-dp-zone-bq"
  location     = "us-central1"
  lake = google_dataplex_lake.basic_lake.name
  type = "RAW"

  discovery_spec {
    enabled = false
  }


  resource_spec {
    location_type = "SINGLE_REGION"
  }

  project = var.gcp_project
}


resource "google_dataplex_asset" "primary" {
  display_name  = var.resource_name
  name          = var.resource_name
  location      = "us-central1"

  lake = google_dataplex_lake.basic_lake.name
  dataplex_zone = google_dataplex_zone.basic_zone.name

  discovery_spec {
    enabled = false
  }

  resource_spec {
    name = "projects/${var.gcp_project}/datasets/${google_bigquery_dataset.named_test_resource.dataset_id}"
    type = "BIGQUERY_DATASET"
  }

  labels = {
    env     = "foo"
    my-asset = "exists"
  }


  project = var.gcp_project
}

output "project_id" {
  value = data.google_client_config.current.project
}

output "resource_aka" {
  value = "gcp://dataplex.googleapis.com/${google_dataplex_asset.primary.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_dataplex_asset.primary.id
}

output "lake_name" {
  value = google_dataplex_lake.basic_lake.id
}

output "zone_name" {
  value = google_dataplex_zone.basic_zone.id
}

