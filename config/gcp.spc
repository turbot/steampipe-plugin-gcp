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
  # impersonate_access_token = "ya29.c.c0ASRK0GZ7mv8lIV0iiudmiGBs9m1gqGfBYZzVRhHY9xQsu82jdCmZZcGe70CjkxcvsCsVXsPeCeGWE6aDl5K561WRcJi1T9l7pTB7fTuWgWDfAzTHye9Kg3z9dc66hVEct_i8seajX9f3WtdQBSqzYvZSenm_jdsuqfyWCSiiz_aVbx5y_MVgx3D_kT2Rz7ePbwSfuqnbsKfiByG0QI8YJlqB1_A6s5pyhITpvWmkNg1baWzqFW5IP9lzvLdGCaTAQ2Nf18IidExI50GLb5NdLQ1Njig_UnVNrEdc-Vke3X3J4PX0fsKY7y4tBcoGlq9Cc9bGEzPHrElEENaJ_iP7WCbZf-b8aUqoqkVQl8Wl15AoJwbrcksjyuXVTg28QpC3cUb6a_aXUnrQcu9Ru1ULRHjVA9JKU0ayu1tpYNLSRxCeM5cMCtGet3f34xCaV7kT1wut6vD4hlpKRfOBHrnd6bTw2dF8Q89m9AA2jHR-X2v68KngKFkvmkTuUEhPJPOXl1VXUAkN1W1oI0ccW8Y04TN_YSreMh56dmeziAZ9O4S_WB9eYcBbgNvW6ghkrQ83UWJ15DeD8hiWPTKFppITzjgaFCQDE4Q1aF5yKlVGFFdpS1Fe9UZvIY4bDKgE645DsSewy-Swiw4wRIiF6wgdwp0cdavl5BoodtI4OSomxxbh3n_u-wJmO12xBXUMiaRMlu72a_ilXow5ynU9U-wdlqScr2uf4bZwZyfUrUX1xXpUJmd8kXka8Skpv6wtOnywQmkWeupbMQXRO504S76u-cXekOdcUSR4RJlZ9s9geB5aWFJ5SmFkwcYXS3Ijh66m_Mq-s2JmI51wd4F-ZY0U75pRw-OmiRX2xtBk5c2mS7gZfoae88MmmU4J2aBwsOwcedX9fUrBl-4QSBxhSRcdsRyFp1eXf0-whBd8mQ9WyJOOb1v9zd1qrmBZpXma2i5ltst6FsizQrxmSb98xROqjY6iqtmqIyshWjydY3RzFYcdOlhb5aMYJd"

  # `impersonate_service_account` (optional) - The GCP service account (string) which should be impersonated.
  # If not set, no impersonation is done.
  #impersonate_service_account = "YOUR_SERVICE_ACCOUNT"

  # `quota_project` (optional) - The project ID used for billing and quota. When set,
  # this project ID is used to track quota usage and billing for the operations performed with the GCP connection.
  # If `quota_project` is not specified directly, the system will look for the `GOOGLE_CLOUD_QUOTA_PROJECT`
  # environment variable to determine which project to use for billing and quota.
  # If neither is specified, billing and quota are tracked against the project associated with the credentials used for authentication.
  # quota_project = "YOUR_QUOTA_PROJECT_ID"

  # `ignore_error_codes` (optional) - List of additional GCP error codes to ignore for all queries.
  # By default, common not found error codes are ignored and will still be ignored even if this argument is not set.
  # Refer https://cloud.google.com/resource-manager/docs/core_errors#Global_Errors for more information on GCP error codes
  #ignore_error_codes = ["401", "403"]
}
