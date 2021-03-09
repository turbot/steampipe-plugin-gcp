
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

data "google_iam_policy" "admin" {
  binding {
    role = "roles/storage.admin"
    members = [
      "user:bob@turbot.com",
    ]
  }
}


resource "google_storage_bucket" "named_test_resource" {
  name                        = var.resource_name
  location                    = "EU"
  force_destroy               = true
  uniform_bucket_level_access = true
  storage_class               = "STANDARD"
  labels = {
    "name" = "test"
  }

  website {
    main_page_suffix = "index.html"
    not_found_page   = "404.html"
  }

  cors {
    origin          = ["http://image-store.com"]
    method          = ["GET", "HEAD", "PUT", "POST", "DELETE"]
    response_header = ["*"]
    max_age_seconds = 3600
  }

  lifecycle_rule {
    condition {
      age = 3
    }
    action {
      type = "Delete"
    }
  }
}

# resource "google_storage_bucket_acl" "image-store-acl" {
#   bucket = google_storage_bucket.named_test_resource.name

#   role_entity = [
#     "OWNER:bob@turbot.com",
#   ]
# }

resource "google_storage_bucket_iam_policy" "policy" {
  bucket      = google_storage_bucket.named_test_resource.name
  policy_data = data.google_iam_policy.admin.policy_data
}

output "resource_aka" {
  value = "gcp://storage.googleapis.com/projects/${data.google_client_config.current.project}/buckets/${google_storage_bucket.named_test_resource.name}"
}

output "resource_name" {
  value = var.resource_name
}

output "resource_id" {
  value = google_storage_bucket.named_test_resource.id
}

output "self_link" {
  value = google_storage_bucket.named_test_resource.self_link
}

