# Make sure that the `gcloud` is already installed, before running this test.

data "google_client_config" "current" {}

locals {
  path = "${path.cwd}/output.json"
}

resource "null_resource" "list_resource" {
  provisioner "local-exec" {
    command = "gcloud organizations list --format json > ${local.path}"
  }
}

data "local_file" "input" {
  depends_on = [null_resource.list_resource]
  filename   = local.path
}

output "display_name" {
  value = lookup(jsondecode(replace(replace(data.local_file.input.content, "[\n  ", ""), "\n]\n", "")), "displayName")
}

output "name" {
  value = lookup(jsondecode(replace(replace(data.local_file.input.content, "[\n  ", ""), "\n]\n", "")), "name")
}

output "lifecycle_state" {
  value = lookup(jsondecode(replace(replace(data.local_file.input.content, "[\n  ", ""), "\n]\n", "")), "lifecycleState")
}
