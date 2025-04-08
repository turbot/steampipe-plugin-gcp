---
title: "Steampipe Table: gcp_dns_policy - Query Google Cloud DNS Policies using SQL"
description: "Allows users to query Google Cloud DNS Policies, providing insights into policy configurations and settings."
folder: "Cloud DNS"
---

# Table: gcp_dns_policy - Query Google Cloud DNS Policies using SQL

Google Cloud DNS Policies are a resource that allows users to configure how DNS queries are handled in Google Cloud. These policies can be used to control DNS behavior in a flexible and granular way, such as by configuring DNS forwarding, alternative name servers, or enabling private DNS zones. Google Cloud DNS Policies provide a way to manage DNS settings across multiple networks, improving network security and reliability.

## Table Usage Guide

The `gcp_dns_policy` table provides insights into DNS Policies within Google Cloud DNS. As a network engineer or a security analyst, explore policy-specific details through this table, including configurations, settings, and associated metadata. Utilize it to uncover information about policies, such as those with specific forwarding paths, the alternative name servers, and the status of private DNS zones.

## Examples

### Basic info
Explore the configuration settings of your Google Cloud Platform's DNS policies to understand their current setup. This can help in identifying instances where inbound forwarding or logging is enabled, which can be crucial for security and network management.

```sql+postgres
select
  name,
  id,
  kind,
  enable_inbound_forwarding,
  enable_logging,
  target_name_servers
from
  gcp_dns_policy;
```

```sql+sqlite
select
  name,
  id,
  kind,
  enable_inbound_forwarding,
  enable_logging,
  target_name_servers
from
  gcp_dns_policy;
```

### List DNS policies with logging disabled
Discover the DNS policies that have logging disabled. This can be useful to identify potential security risks or compliance issues related to lack of logging in your GCP environment.

```sql+postgres
select
  name,
  id,
  enable_logging
from
  gcp_dns_policy
where
  not enable_logging;
```

```sql+sqlite
select
  name,
  id,
  enable_logging
from
  gcp_dns_policy
where
  enable_logging = 0;
```

### List DNS policies not associated with any network
Discover policies in Google Cloud Platform's DNS service that aren't linked to any network. This can help identify unused resources or potential configuration issues.

```sql+postgres
select
  name,
  id,
  networks
from
  gcp_dns_policy
where
  networks = '[]';
```

```sql+sqlite
select
  name,
  id,
  networks
from
  gcp_dns_policy
where
  networks = '[]';
```