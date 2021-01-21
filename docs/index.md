---
organization: Turbot
category: ["public cloud"]
icon_url: "/images/plugins/turbot/google.svg"
brand_color: "#e7483a"
display_name: "GCP"
name: "gcp"
description: "Steampipe plugin for Google Cloud Platform (GCP) services and resource types"
---

# GCP

The Google Cloud Platform (GCP) plugin is used to interact with the many resources supported by GCP.

### Installation

To download and install the latest gcp plugin:

```bash
$ steampipe plugin install gcp
Installing plugin gcp...
$
```

### Scope

An GCP connection is scoped to a single gcp project, with a single set of credentials. Currently, a connection is limited to a single GCP region.

### Configure Environment Variables

This topic explains how to authenticate an application as a service account. To use Steampipe, you need to create an ervice account in GCP with the appropriate permissions.

1. In the Cloud Console, go to the Create [service account key page](https://console.cloud.google.com/apis/credentials/serviceaccountkey).
2. From the Service account list, select New service account.
3. In the Service account name field, enter a name.
4. From the Role list, select Project > Viewer.
5. Click Create. A JSON file that contains your key downloads to your computer.
6. Set `GOOGLE_APPLICATION_CREDENTIALS` environment variable with path to downloaded JSON file.
7. Set `GCP_PROJECT` environment variable to name of the project

Set GCP project and credential as environment variable (Mac, Linux):

```bash
export GOOGLE_APPLICATION_CREDENTIALS=/Users/abc/Downloads/project-aaa.json
export GCP_PROJECT=project-aaa
```

Run a query:

```bash
$ steampipe query
Welcome to Steampipe v0.0.11
Type ".inspect" for more information.
> select name, expire_time, project, topic from gcp.gcp_pubsub_snapshot;
+--------+--------------------------+-------------+------------------------------------------------------------------+
|  name  |       expire_time        |   project   |                              topic                               |
+--------+--------------------------+-------------+------------------------------------------------------------------+
| test12 | 2021-01-28T09:17:21.986Z | project-aaa | projects/project-aaa/topics/event_handler                        |
+--------+--------------------------+-------------+------------------------------------------------------------------+
```
