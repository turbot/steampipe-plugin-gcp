---
title: "Steampipe Table: gcp_dns_record_set - Query GCP DNS Record Sets using SQL"
description: "Allows users to query DNS Record Sets in Google Cloud Platform (GCP), specifically the details of each record set, providing insights into DNS configurations and potential anomalies."
folder: "Cloud DNS"
---

# Table: gcp_dns_record_set - Query GCP DNS Record Sets using SQL

A DNS Record Set in Google Cloud Platform (GCP) is a collection of DNS records of the same type that share the same domain name. They are used to map a domain name to an IP address or other data. DNS Record Sets are essential for the functioning of the internet, enabling the translation of human-readable domain names into numerical IP addresses that computers can understand.

## Table Usage Guide

The `gcp_dns_record_set` table provides insights into DNS Record Sets within Google Cloud Platform (GCP). As a network engineer, explore details of each record set through this table, including record types, record data, and associated metadata. Utilize it to uncover information about DNS configurations, such as those with potential misconfigurations, the mapping of domain names to IP addresses, and the verification of DNS configurations.

## Examples

### Basic info
Explore the configuration of your Google Cloud Platform DNS record sets to gain insights into their management zones, types, routing policies, and time-to-live (TTL) settings. This is useful for understanding your DNS infrastructure and making necessary changes for optimal performance.

```sql+postgres
select
  name,
  managed_zone_name,
  type,
  kind,
  routing_policy,
  rrdatas,
  signature_rrdatas,
  ttl
from
  gcp_dns_record_set;
```

```sql+sqlite
select
  name,
  managed_zone_name,
  type,
  kind,
  routing_policy,
  rrdatas,
  signature_rrdatas,
  ttl
from
  gcp_dns_record_set;
```

### List CNAME record sets
Explore which alias domain names are associated with actual domain names in your Google Cloud DNS. This is useful for managing and maintaining DNS records, ensuring proper redirections, and troubleshooting potential issues.

```sql+postgres
select
  name,
  managed_zone_name,
  type,
  ttl
from
  gcp_dns_record_set
where
 type = 'CNAME';
```

```sql+sqlite
select
  name,
  managed_zone_name,
  type,
  ttl
from
  gcp_dns_record_set
where
 type = 'CNAME';
```