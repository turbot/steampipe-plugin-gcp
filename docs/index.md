---
organization: Turbot
category: ["public cloud"]
icon_url: "/images/plugins/turbot/gcp.svg"
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

Installing the latest gcp plugin will create a default connection named `gcp`. This connection will dynamically determine the scope and credentials using the same mechanism as the CLI. In effect, this means that by default Steampipe will execute with the same credentials and against the same project as the `gcloud` command would - The GCP plugin uses the standard sdk environment variables and credential files as used by the CLI.  (Of course this also  implies that the `gcloud` cli needs to be configured with the proper credentials before the steampipe gcp plugin can be used).

Note that there is nothing special about the default connection, other than that it is created by default on plugin install - You can delete or rename this connection, or modify its configuration options (via the configuration file).


## Connection Configuration
Connection configurations are defined using HCL in one or more Steampipe config files.  Steampipe will load ALL configuration files from `~/.steampipe/config` that have a `.spc` extension. A config file may contain multiple connections.


### Scope

A GCP connection is scoped to a single GCP project, with a single set of credentials. 

The Google APIs are scoped at multiple levels (global, regional, zonal) depending on the resource/service, however the `gcloud` cli list commands tend to be global regardless of the API scope. For example, `gcloud compute instances list` returns ALL instances, even though the API scope is zonal. Steampipe will behave the same way as the cli -- A GCP Steampipe connection will query all regions, locations, and zones in the project.


### Credentials

By default, the GCP plugin uses your [Application Default Credentials](https://cloud.google.com/sdk/gcloud/reference/auth/application-default) to connect to GCP.  If you have not set up ADC, simply run `gcloud auth application-default login`.  This command will prompt you to log in, and then will download the application default credentials to `~/.config/gcloud/application_default_credentials.json`.

Alternatively, Steampipe can use a service account to connect to GCP
1. In the Cloud Console, go to the Create [service account key page](https://console.cloud.google.com/apis/credentials/serviceaccountkey).
2. From the Service account list, select New service account.
3. In the Service account name field, enter a name.
4. From the Role list, select Project > Viewer.
5. Click Create. A JSON file that contains your key downloads to your computer.

### Configuration Arguments

The GCP plugin allows you set credentials static credentials with the following arguments:
- `project` - The project ID to connect to. This is the project id (string), not the project number. If the `project` argument is not specified for a connection, the project will be determined in the following order:
   - The standard gcloud SDK `CLOUDSDK_CORE_PROJECT` environment variable, if set; otherwise
   - The `GCP_PROJECT` environment variable, if set (this is deprecated); otherwise
   - The current active project project, as returned by the `gcloud config get-value project` command
- `credential_file` - The path to a JSON credential file that contains Google application credentials.  If `credential_file` is not specified in a connection, credentials will be loaded from:
   - The path specified in the `GOOGLE_APPLICATION_CREDENTIALS` environment variable, if set; otherwise
   - The standard location (`~/.config/gcloud/application_default_credentials.json`)

#### Example configurations

- The default connection.  This uses standard Application Default Credentials (ADC) against the active project as configured for `gcloud`
   ```hcl
   connection "gcp" {
   plugin    = "gcp"                 
   }
   ```

- A connection to a specific project, using standard ADC Credentials.
   ```hcl
   connection "gcp_my_project" {
   plugin    = "gcp"   
   project   = "my-project"              
   }
   ```

- A connection to a specific project, using non-default credentials.
   ```hcl
   connection "gcp_my_other_project" {
   plugin             = "gcp"   
   project            = "my-other-project"
   credential_file    = "/home/me/my-service-account-creds.json"        
   }
   ```

