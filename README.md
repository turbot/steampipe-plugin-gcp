![image](https://hub.steampipe.io/images/plugins/turbot/gcp-social-graphic.png)

# GCP Plugin for Steampipe

Use SQL to query infrastructure including servers, networks, databases and more from GCP.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/gcp)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/gcp/tables)
- Community: [Slack Channel](https://join.slack.com/t/steampipe/shared_invite/zt-oij778tv-lYyRTWOTMQYBVAbtPSWs3g)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-gcp/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install gcp
```

Run a query:

```sql
select
  name,
  role_id
from
  gcp_iam_role;
```

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-gcp.git
cd steampipe-plugin-gcp
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```sh
make
```

Configure the plugin:

```sh
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/gcp.spc
```

Try it!

```shell
steampipe query
> .inspect gcp
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-gcp/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [GCP Plugin](https://github.com/turbot/steampipe-plugin-gcp/labels/help%20wanted)
