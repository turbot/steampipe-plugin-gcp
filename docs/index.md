---
organization: Turbot
category: ["public cloud"]
icon_url: "/images/plugins/turbot/gcp.svg"
brand_color: "#ea4335"
display_name: "GCP"
name: "gcp"
description: "Steampipe plugin for Google Cloud Platform (GCP) services and resource types"
og_description: Query GCP with SQL! Open source CLI. No DB required. 
og_image: "/images/plugins/turbot/gcp-social-graphic.png"
---

# GCP + Steampipe

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

[GCP](https://gcp.amazon.com/) provides on-demand cloud computing platforms and APIs to authenticated customers on a metered pay-as-you-go basis. 

For example:

```sql
select
  name,
  location,
  versioning_enabled
from
  gcp_storage_bucket
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

- **[Table definitions & examples →](gcp/tables)**

## Get started

### Install

Download and install the latest GCP plugin:

```bash
steampipe plugin install gcp
```

### Credentials

| Item | Description |
| - | - |
| Credentials | configure your [Application Default Credentials](https://cloud.google.com/sdk/gcloud/reference/auth/application-default) |
| Permissions | Grant the `ReadOnlyAccess` policy to your user or role. |
| Radius | Each connection represents a single GCP project. |
| Resolution |  1. Credentials from the json file specified by the `credential_file` paramater in your steampipe config.<br />2. Credentials from the json file specified by the `GOOGLE_APPLICATION_CREDENTIALS` environment variable.<br />3. Credentials from the default json file location (~/.config/gcloud/application_default_credentials.json). |

### Configuration

Installing the latest gcp plugin will create a config file (`~/.steampipe/config/gcp.spc`) with a single connection named `gcp`:

```hcl
connection "gcp" {
  plugin    = "gcp"
  project   = "my-project"
  credential_file = "~/my-service-account-creds.json"
}
```

## Get involved

* Open source: https://github.com/turbot/steampipe-plugin-gcp
* Community: [Discussion forums](https://github.com/turbot/steampipe/discussions)