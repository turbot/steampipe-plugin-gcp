
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

data "google_project" "project" {}

resource "google_dataplex_lake" "named_test_resource" {
  location     = "us-central1"
  name         = var.resource_name
  description  = "Lake for DCL"
  display_name = var.resource_name
  project      = data.google_client_config.current.project

  labels = {
    name = var.resource_name
  }
}

resource "google_dataplex_task" "named_test_resource" {

    task_id      = var.resource_name
    location     = "us-central1"
    lake         = google_dataplex_lake.named_test_resource.name

    description = "Test Task Basic"

    display_name = var.resource_name

    labels = { "count": "3" }

    trigger_spec  {
        type = "RECURRING"
        disabled = false
        max_retries = 3
        start_time = "2023-10-02T15:01:23Z"
        schedule = "1 * * * *"
    }

    execution_spec {
        service_account = "${data.google_project.project.number}-compute@developer.gserviceaccount.com"
        project = data.google_client_config.current.project
        max_job_execution_lifetime = "100s"
        kms_key = "234jn2kjn42k3n423"
    }

    spark {
        python_script_file = "gs://dataproc-examples/pyspark/hello-world/hello-world.py"

    }

    project = data.google_client_config.current.project

}

output "project_id" {
  value = data.google_client_config.current.project
}

output "resource_aka" {
  value = "gcp://dataplex.googleapis.com/${google_dataplex_task.named_test_resource.id}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_dataplex_task.named_test_resource.id
}

output "lake_name" {
  value = google_dataplex_lake.named_test_resource.id
}

