---
organization: Turbot
category: ["public cloud"]
icon_url: "/images/plugins/turbot/gcp.svg"
brand_color: "#ea4335"
display_name: "GCP"
name: "gcp"
description: "Steampipe plugin for querying buckets, instances, functions and more from GCP."
og_description: Query GCP with SQL! Open source CLI. No DB required.
og_image: "/images/plugins/turbot/gcp-social-graphic.png"
engines: ["steampipe", "sqlite", "postgres", "export"]
---

# GCP + Steampipe

[Steampipe](https://steampipe.io) is an open-source zero-ETL engine to instantly query cloud APIs using SQL.

[GCP](https://cloud.google.com) provides on-demand cloud computing platforms and APIs to authenticated customers on a metered pay-as-you-go basis.

For example:

```sql
select
  name,
  location,
  versioning_enabled
from
  gcp_storage_bucket;
```

```
+--------------------+----------+--------------------+
| name               | location | versioning_enabled |
+--------------------+----------+--------------------+
| steampipe-io-dev   | us-east1 | false              |
| steampipe-io-stage | us       | false              |
| steampipe-io-prod  | us       | true               |
+--------------------+----------+--------------------+
```

## Documentation

- **[Table definitions & examples →](/plugins/turbot/gcp/tables)**

## Get started

### Install

Download and install the latest GCP plugin:

```bash
steampipe plugin install gcp
```

### Credentials

| Item        | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| ----------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Credentials | When running locally, you must configure your [Application Default Credentials](https://cloud.google.com/sdk/gcloud/reference/auth/application-default). If you are running in Cloud Shell or Cloud Code, [the tool uses the credentials you provided when you logged in, and manages any authorizations required](https://cloud.google.com/docs/authentication/provide-credentials-adc#cloud-based-dev). |
| Permissions | Assign the `Viewer` role to your user or service account. You may also need additional permissions related to IAM policies, like `pubsub.subscriptions.getIamPolicy`, `pubsub.topics.getIamPolicy`, `storage.buckets.getIamPolicy`, since these are not included in the `Viewer` role. You can grant these by creating a custom role in your project. |
| Radius      | Each connection represents a single GCP project, except for some tables like `gcp_organization` and `gcp_organization_project` which return all resources the credentials attached to the connection have access to. |
| Resolution  | 1. Credentials from the JSON file specified by the `credentials` parameter in your steampipe config.<br />2. Credentials from the JSON file specified by the `GOOGLE_APPLICATION_CREDENTIALS` environment variable.<br />3. Credentials from the default JSON file location (~/.config/gcloud/application_default_credentials.json). <br />4. Credentials from [the metadata server](https://cloud.google.com/docs/authentication/application-default-credentials#attached-sa) |

### Configuration

Installing the latest gcp plugin will create a config file (`~/.steampipe/config/gcp.spc`) with a single connection named `gcp`:

```hcl
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

  # `impersonate_access_token` (optional) - You can generate an OAuth 2.0 access token by using the gcloud CLI, the REST API, or the Cloud Client Libraries and Google API Client Libraries.
  # Refer https://cloud.google.com/iam/docs/create-short-lived-credentials-direct#gcloud_2 for generating the access token.
  # impersonate_access_token = "ya29.c.c0ASRK0GZ7mv8lIV0iiudmiGBs9m1gqGfBYZzV...aMYJd"

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
  # By default, the common not found error codes are ignored and will still be ignored even if this argument is not set.
  # Refer https://cloud.google.com/resource-manager/docs/core_errors#Global_Errors for more information on GCP error codes
  #ignore_error_codes = ["401", "403"]
}
```

## Advanced configuration options

By default, the GCP plugin uses your [Application Default Credentials](https://cloud.google.com/sdk/gcloud/reference/auth/application-default) to connect to GCP. If you have not set up ADC, simply run `gcloud auth application-default login`. This command will prompt you to log in, and then will download the application default credentials to ~/.config/gcloud/application_default_credentials.json.

For users with multiple GCP projects and more complex authentication use cases, here are some examples of advanced configuration options:

### Use a service account

Generate and download a JSON key for an existing service account using: [create service account key page](https://console.cloud.google.com/apis/credentials/serviceaccountkey).

```hcl
connection "gcp_my_other_project" {
  plugin      = "gcp"
  project     = "my-other-project"
  credentials = "/home/me/my-service-account-creds.json"
}
```

### Use impersonation access token

Generate an impersonate access token using: [gcloud CLI command](https://cloud.google.com/iam/docs/create-short-lived-credentials-direct#gcloud_2).

```hcl
connection "gcp_my_other_project" {
  plugin                   = "gcp"
  project                  = "my-other-project"
  impersonate_access_token = "ya29.c.c0ASRK0GZ7mv8lIV0iiudmiGBs9m1gqGfBYZzV...aMYJd"
}
```

## Multi-Project Connections

You may create multiple gcp connections:

```hcl
connection "gcp_all" {
  type        = "aggregator"
  plugin      = "gcp"
  connections = ["gcp_project_*"]
}

connection "gcp_project_aaa" {
  plugin  = "gcp"
  project = "project-aaa"
}

connection "gcp_project_bbb" {
  plugin  = "gcp"
  project = "project-bbb"
}

connection "gcp_project_ccc" {
  plugin  = "gcp"
  project = "project-ccc"
}
```

Depending on the mode of authentication, a multi-project configuration can also look like:

```hcl
connection "gcp_all" {
  type        = "aggregator"
  plugin      = "gcp"
  connections = ["gcp_project_*"]
}

connection "gcp_project_aaa" {
  plugin      = "gcp"
  project     = "project-aaa"
  credentials = "/home/me/my-service-account-creds-for-project-aaa.json"
}

connection "gcp_project_bbb" {
  plugin      = "gcp"
  project     = "project-bbb"
  credentials = "/home/me/my-service-account-creds-for-project-bbb.json"
}

connection "gcp_project_ccc" {
  plugin      = "gcp"
  project     = "project-ccc"
  credentials = "/home/me/my-service-account-creds-for-project-ccc.json"
}
```

Each connection is implemented as a distinct [Postgres schema](https://www.postgresql.org/docs/current/ddl-schemas.html). As such, you can use qualified table names to query a specific connection:

```sql
select * from gcp_project_aaa.gcp_project
```

Alternatively, you can use an unqualified name and it will be resolved according to the [Search Path](https://steampipe.io/docs/using-steampipe/managing-connections#setting-the-search-path):

```sql
select * from gcp_project
```

You can create multi-project connections by using an [**aggregator** connection](https://steampipe.io/docs/using-steampipe/managing-connections#using-aggregators). Aggregators allow you to query data from multiple connections for a plugin as if they are a single connection:

```hcl
connection "gcp_all" {
  plugin      = "gcp"
  type        = "aggregator"
  connections = ["gcp_project_aaa", "gcp_project_bbb", "gcp_project_ccc"]
}
```

Querying tables from this connection will return results from the `gcp_project_aaa`, `gcp_project_bbb`, and `gcp_project_ccc` connections:

```sql
select * from gcp_all.gcp_project
```

Steampipe supports the `*` wildcard in the connection names. For example, to aggregate all the GCP plugin connections whose names begin with `gcp_`:

```hcl
connection "gcp_all" {
  type        = "aggregator"
  plugin      = "gcp"
  connections = ["gcp_*"]
}
```

### Specify static credentials using environment variables

```sh
export CLOUDSDK_CORE_PROJECT=myproject
export GOOGLE_CLOUD_QUOTA_PROJECT=billingproject
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/my/creds.json
```
