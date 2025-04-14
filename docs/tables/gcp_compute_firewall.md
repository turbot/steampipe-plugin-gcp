---
title: "Steampipe Table: gcp_compute_firewall - Query Google Cloud Compute Engine Firewalls using SQL"
description: "Allows users to query Google Cloud Compute Engine Firewalls, providing insights into firewall rules and their configurations."
folder: "Compute"
---

# Table: gcp_compute_firewall - Query Google Cloud Compute Engine Firewalls using SQL

Google Cloud Compute Engine Firewalls are a networking security feature that allows you to control the traffic to your virtual machine instances. They provide a flexible and robust tool for securing your instances by defining what traffic is allowed to and from your instances. Firewalls are implemented at the network level and apply to all traffic that crosses the perimeter of the network.

## Table Usage Guide

The `gcp_compute_firewall` table can be used to gain insights into firewall rules within Google Cloud Compute Engine. As a network security administrator or a DevOps engineer, you can explore details about each firewall rule, including allowed and denied configurations, network associations, and priority. Utilize it to identify firewall rules that may be overly permissive or misconfigured, enhancing your network security posture.

## Examples

### Firewall rules basic info
Explore which firewall rules are in place for your Google Cloud Platform (GCP) compute instances. This allows you to understand the direction of traffic flow and assess the overall security configuration.

```sql+postgres
select
  name,
  id,
  description,
  direction
from
  gcp_compute_firewall;
```

```sql+sqlite
select
  name,
  id,
  description,
  direction
from
  gcp_compute_firewall;
```

### List of rules which are applied to TCP protocol
Explore which firewall rules are applied specifically to the TCP protocol in your Google Cloud Platform. This will help in assessing network security and identifying potential vulnerabilities.

```sql+postgres
select
  name,
  id,
  p ->> 'IPProtocol' as ip_protocol,
  p ->> 'ports' as ports
from
  gcp_compute_firewall,
  jsonb_array_elements(allowed) as p
where
  p ->> 'IPProtocol' = 'tcp';
```

```sql+sqlite
select
  f.name,
  f.id,
  json_extract(p.value, '$.IPProtocol') as ip_protocol,
  json_extract(p.value, '$.ports') as ports
from
  gcp_compute_firewall as f,
  json_each(allowed) as p
where
  json_extract(p.value, '$.IPProtocol') = 'tcp';
```

### List of disabled rules
Determine the areas in which firewall rules are disabled to strengthen your security posture in Google Cloud Platform. This can assist in identifying potential vulnerabilities and maintaining robust network security.

```sql+postgres
select
  name,
  id,
  description,
  disabled
from
  gcp_compute_firewall
where
  disabled;
```

```sql+sqlite
select
  name,
  id,
  description,
  disabled
from
  gcp_compute_firewall
where
  disabled = 1;
```

### List of Egress rules
Explore which firewall rules in your Google Cloud Platform are set to allow outbound traffic. This can help understand your network's security posture and identify potential vulnerabilities.

```sql+postgres
select
  name,
  id,
  direction,
  allowed,
  denied
from
  gcp_compute_firewall
where
  direction = 'EGRESS';
```

```sql+sqlite
select
  name,
  id,
  direction,
  allowed,
  denied
from
  gcp_compute_firewall
where
  direction = 'EGRESS';
```