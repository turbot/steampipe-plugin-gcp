variable "project_id" {
  type        = string
  default     = "1111111111"
  description = "The project ID."
}

variable "project_name" {
  type        = string
  default     = "parker-aaa"
  description = "The project name."
}

variable "tag_key_name" {
  type        = string
  default     = "env"
  description = "Short name for the tag key."
}

variable "tag_value_name" {
  type        = string
  default     = "dev"
  description = "Short name for the tag value."
}

resource "google_tags_tag_key" "key" {
  parent      = "projects/${var.project_id}"
  short_name  = var.tag_key_name
  description = "Tag Key for ${var.tag_key_name} resources."
}

resource "google_tags_tag_value" "value" {
  parent      = "tagKeys/${google_tags_tag_key.key.name}"
  short_name  = var.tag_value_name
  description = "Tag Value for ${var.tag_value_name} resources."
}

resource "google_tags_tag_binding" "binding" {
  parent    = "//cloudresourcemanager.googleapis.com/projects/${var.project_id}"
  tag_value = "tagValues/${google_tags_tag_value.value.name}"
}

output "parent" {
  value = google_tags_tag_binding.binding.parent
}

output "name" {
  value = google_tags_tag_binding.binding.id
}

output "tag_value" {
  value = google_tags_tag_binding.binding.tag_value
}
