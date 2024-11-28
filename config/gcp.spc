connection "gcp" {
  plugin = "gcp"

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

  # `impersonate_access_token` (optional) - You can generate an OAuth 2.0 access token by using the gcloud CLI, the REST API, or the Cloud Client Libraries and Google API Client Libraries.
  # Refer https://cloud.google.com/iam/docs/create-short-lived-credentials-direct#gcloud_2 for generating the access token.
  # impersonate_access_token = "ya29.c.c0ASRK0GZ7mv8lIV0iiudmiGBs...hb5aMYJd"

  # `impersonate_service_account` (optional) - The GCP service account (string) which should be impersonated.
  # If not set, no impersonation is done.
  #impersonate_service_account = "YOUR_SERVICE_ACCOUNT"

  # `quota_project` (optional) - The project ID used for billing and quota. When set,
  # this project ID is used to track quota usage and billing for the operations performed with the GCP connection.
  # If `quota_project` is not specified directly, the system will look for the `GOOGLE_CLOUD_QUOTA_PROJECT`
  # environment variable to determine which project to use for billing and quota.
  # If neither is specified, billing and quota are tracked against the project associated with the credentials used for authentication.
  # quota_project = "YOUR_QUOTA_PROJECT_ID"

  # `ignore_error_messages` (optional) - List of additional GCP error message pattern to ignore for all queries.
  #  ignore_error_messages = ["^.*API has not been used.*$"]

  # `ignore_error_codes` (optional) - List of additional GCP error codes to ignore for all queries.
  # By default, common not found error codes are ignored and will still be ignored even if this argument is not set.
  # Refer https://cloud.google.com/resource-manager/docs/core_errors#Global_Errors for more information on GCP error codes
  #ignore_error_codes = ["401", "403"]
}
