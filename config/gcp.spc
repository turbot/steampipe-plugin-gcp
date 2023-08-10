connection "gcp" {
  plugin    = "gcp"

  # `project` (optional) - The project ID to connect to. This is the project ID (string), not the
  # project name or number. If the `project` argument is not specified for a connection,
  # the project will be determined in the following order:
  #   - The standard gcloud SDK `CLOUDSDK_CORE_PROJECT` environment variable, if set; otherwise
  #   - The `GCP_PROJECT` environment variable, if set (this is deprecated); otherwise
  #   - The current active project, as returned by the `gcloud config get-value project` command
  #project = "YOUR_PROJECT_ID"

  # `credentials` (optional) - Either the path to a JSON credential file that contains Google application credentials,
  # or the contents of a service account key file in JSON format. If `credentials` is not specified in a connection,
  # credentials will be loaded from:
  #   - The path specified in the `GOOGLE_APPLICATION_CREDENTIALS` environment variable, if set; otherwise
  #   - The standard location (`~/.config/gcloud/application_default_credentials.json`)
  #credentials = "~/.config/gcloud/application_default_credentials.json"

  # `impersonate_service_account` (optional) - The GCP service account (string) which should be impersonated.
  # If not set, no impersonation is done.
  #impersonate_service_account = "YOUR_SERVICE_ACCOUNT"

  # `ignore_error_messages` (optional) - List of additional GCP error message pattern to ignore for all queries.
  #  ignore_error_messages = ["API has not been used", "other regex"]

  # `ignore_error_codes` (optional) - List of additional GCP error codes to ignore for all queries.
  # By default, common not found error codes are ignored and will still be ignored even if this argument is not set.
  # Refer https://cloud.google.com/resource-manager/docs/core_errors#Global_Errors for more information on GCP error codes
  #ignore_error_codes = ["401", "403"]
}
