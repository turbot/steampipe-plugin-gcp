connection "gcp" {
  plugin    = "gcp"

  # `project` (optional) - The project ID to connect to. This is the project id (string), not the
  # project number. If the `project` argument is not specified for a connection,
  # the project will be determined in the following order:
  #   - The standard gcloud SDK `CLOUDSDK_CORE_PROJECT` environment variable, if set; otherwise
  #   - The `GCP_PROJECT` environment variable, if set (this is deprecated); otherwise
  #   - The current active project project, as returned by the `gcloud config get-value project` command
  #project  = "YOUR_PROJECT_NAME"

  # `credential_file` (optional) -  - The path to a JSON credential file that contains
  # Google application credentials.  If `credential_file` is not specified in a connection,
  # credentials will be loaded from:
  #   - The path specified in the `GOOGLE_APPLICATION_CREDENTIALS` environment variable, if set; otherwise
  #   - The standard location (`~/.config/gcloud/application_default_credentials.json`)
  #credential_file    = "~/.config/gcloud/application_default_credentials.json"

  # `impersonate_service_account` (optional) - The GCP service account (string) which should be impersonated.
  # If not set, no impersonation is done.
  #impersonate_service_account = "YOUR_SERVICE_ACCOUNT"
}
