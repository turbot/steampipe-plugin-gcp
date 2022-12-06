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
---

# GCP + Steampipe

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

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

| Item | Description |
| - | - |
| Credentials | When running locally, you must configure your [Application Default Credentials](https://cloud.google.com/sdk/gcloud/reference/auth/application-default). If you are running in Cloud Shell or Cloud Code, [the tool uses the credentials you provided when you logged in, and manages any authorizations required](https://cloud.google.com/docs/authentication/provide-credentials-adc#cloud-based-dev). |
| Permissions | Assign the `Viewer` role to your user or service account. |
| Radius | Each connection represents a single GCP project. |
| Resolution |  1. Credentials from the json file specified by the `credentials` parameter in your steampipe config.<br />2. Credentials from the json file specified by the `GOOGLE_APPLICATION_CREDENTIALS` environment variable.<br />3. Credentials from the default json file location (~/.config/gcloud/application_default_credentials.json). <br />4. Credentials from [the metadata server](https://cloud.google.com/docs/authentication/application-default-credentials#attached-sa)|

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
  #   - The current active project project, as returned by the `gcloud config get-value project` command
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
}
```

**NOTE:** The `credential_file` property has been deprecated and will be removed in the next major version. Please use `credentials` instead.

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-gcp
- Community: [Slack Channel](https://steampipe.io/community/join)

## Advanced configuration options

By default, the GCP plugin uses your [Application Default Credentials](https://cloud.google.com/sdk/gcloud/reference/auth/application-default) to connect to GCP. If you have not set up ADC, simply run `gcloud auth application-default login`. This command will prompt you to log in, and then will download the application default credentials to ~/.config/gcloud/application_default_credentials.json.

For users with multiple GCP project and more complex authentication use cases, here are some examples of advanced configuration options:

### Use a service account

Generate and download a JSON key for an existing service account using: [create service account key page](https://console.cloud.google.com/apis/credentials/serviceaccountkey).

```hcl
connection "gcp_my_other_project" {
  plugin      = "gcp"
  project     = "my-other-project"
  credentials = "/home/me/my-service-account-creds.json"
}
```

### Specify multiple projects

A common configuration is to have multiple connections to different projects, using the same standard ADC Credentials for all connections:

```hcl
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

### Specify static credentials using environment variables

```sh
export CLOUDSDK_CORE_PROJECT=myproject
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/my/creds.json
```
