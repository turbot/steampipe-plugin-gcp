---
title: "Steampipe Table: gcp_compute_ssl_policy - Query GCP Compute SSL Policies using SQL"
description: "Allows users to query SSL Policies in GCP Compute, providing insights into the SSL policies and their configurations."
folder: "Compute"
---

# Table: gcp_compute_ssl_policy - Query GCP Compute SSL Policies using SQL

A GCP Compute SSL Policy is a resource in Google Cloud Platform's Compute Engine that allows for the flexible configuration of SSL features. SSL Policies are used to control the features of SSL connections that proxy or load balancers negotiate, with the SSL policy being associated with a TargetHttpsProxy or TargetSslProxy resource. They are primarily used to control the minimum version of SSL/TLS protocol, as well as the SSL features that the proxy or load balancer negotiates.

## Table Usage Guide

The `gcp_compute_ssl_policy` table provides comprehensive insights into SSL Policies within Google Cloud Platform's Compute Engine. As a security analyst, you can explore policy-specific details through this table, including minimum SSL version, profile, and custom features. Use this table to uncover information about SSL policies, such as their configurations, associated resources, and any potential security vulnerabilities due to outdated SSL versions or weak ciphers.

## Examples

### Basic info
Explore the basic information of your SSL policies in Google Cloud Platform to understand their configurations and ensure they are using the most secure version of TLS. This can help in maintaining the security standards and compliance of your infrastructure.

```sql+postgres
select
  name,
  id,
  self_link,
  min_tls_version
from
  gcp_compute_ssl_policy;
```

```sql+sqlite
select
  name,
  id,
  self_link,
  min_tls_version
from
  gcp_compute_ssl_policy;
```

### List SSL policies with minimum TLS version 1.2 and the MODERN profile
Determine the areas in which SSL policies are utilizing minimum TLS version 1.2 and the modern profile. This is useful to ensure that your network security is up to date and adheres to modern standards.

```sql+postgres
select
  name,
  id,
  min_tls_version
from
  gcp_compute_ssl_policy
where
  min_tls_version = 'TLS_1_2'
  and profile = 'MODERN';
```

```sql+sqlite
select
  name,
  id,
  min_tls_version
from
  gcp_compute_ssl_policy
where
  min_tls_version = 'TLS_1_2'
  and profile = 'MODERN';
```

### List SSL policies with the RESTRICTED profile
Determine the areas in which SSL policies adhere to a 'RESTRICTED' profile. This can be useful for maintaining security standards and ensuring compliance within your Google Cloud Platform environment.

```sql+postgres
select
  name,
  id,
  profile
from
  gcp_compute_ssl_policy
where
  profile = 'RESTRICTED';
```

```sql+sqlite
select
  name,
  id,
  profile
from
  gcp_compute_ssl_policy
where
  profile = 'RESTRICTED';
```

### List SSL policies with weak cipher suites
Discover the segments that have weak SSL policies enabled. This is particularly useful for identifying potential security vulnerabilities within your system.

```sql+postgres
select
  name,
  id,
  enabled_feature
from
  gcp_compute_ssl_policy,
  jsonb_array_elements_text(enabled_features) as enabled_feature
where
  profile = 'CUSTOM'
  and enabled_feature in('TLS_RSA_WITH_AES_128_GCM_SHA256', 'TLS_RSA_WITH_AES_256_GCM_SHA384', 'TLS_RSA_WITH_AES_128_CBC_SHA', 'TLS_RSA_WITH_AES_256_CBC_SHA', 'TLS_RSA_WITH_3DES_EDE_CBC_SHA');
```

```sql+sqlite
select
  p.name,
  p.id,
  enabled_feature
from
  gcp_compute_ssl_policy as p,
  json_each(enabled_features) as enabled_feature
where
  profile = 'CUSTOM'
  and enabled_feature.value in('TLS_RSA_WITH_AES_128_GCM_SHA256', 'TLS_RSA_WITH_AES_256_GCM_SHA384', 'TLS_RSA_WITH_AES_128_CBC_SHA', 'TLS_RSA_WITH_AES_256_CBC_SHA', 'TLS_RSA_WITH_3DES_EDE_CBC_SHA');
```