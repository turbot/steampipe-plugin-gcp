---
title: "Steampipe Table: gcp_dns_managed_zone - Query Google Cloud DNS Managed Zones using SQL"
description: "Allows users to query Google Cloud DNS Managed Zones, specifically the configuration and status details, providing insights into DNS management and potential configuration issues."
folder: "Cloud DNS"
---

# Table: gcp_dns_managed_zone - Query Google Cloud DNS Managed Zones using SQL

Google Cloud DNS is a scalable, reliable, and managed authoritative Domain Name System (DNS) service running on the same infrastructure as Google. It provides a simple, cost-effective way to make your applications and services available to your users. This service translates requests for domain names like www.google.com into IP addresses like 74.125.29.101.

## Table Usage Guide

The `gcp_dns_managed_zone` table provides insights into DNS Managed Zones within Google Cloud DNS. As a network engineer, explore zone-specific details through this table, including DNS configuration, visibility, and associated metadata. Utilize it to uncover information about zones, such as those with private visibility, DNSSEC state, and the verification of DNS configurations.

## Examples

### Basic info
Explore the basic information about Google Cloud Platform's DNS managed zones, such as their names, identifiers, DNS names, creation times, and visibility settings. This query can help you gain insights into the configuration and status of your DNS managed zones to ensure they are set up as expected.

```sql+postgres
select
  name,
  id,
  dns_name,
  creation_time,
  visibility
from
  gcp_dns_managed_zone;
```

```sql+sqlite
select
  name,
  id,
  dns_name,
  creation_time,
  visibility
from
  gcp_dns_managed_zone;
```

### List public zones with DNSSEC disabled
Explore which public zones have the DNSSEC feature disabled. This can be used to identify potential security vulnerabilities in your DNS configuration.

```sql+postgres
select
  name,
  id,
  dns_name,
  dnssec_config_state,
  visibility
from
  gcp_dns_managed_zone
where 
  visibility = 'public'
  and 
  (
    dnssec_config_state is null
    or dnssec_config_state = 'off'
  );
```

```sql+sqlite
select
  name,
  id,
  dns_name,
  dnssec_config_state,
  visibility
from
  gcp_dns_managed_zone
where 
  visibility = 'public'
  and 
  (
    dnssec_config_state is null
    or dnssec_config_state = 'off'
  );
```