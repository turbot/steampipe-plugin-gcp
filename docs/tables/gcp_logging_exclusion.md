---
title: "Steampipe Table: gcp_logging_exclusion - Query Google Cloud Platform Logging Exclusions using SQL"
description: "Allows users to query Logging Exclusions in Google Cloud Platform, specifically the filters and details of the excluded logs, providing insights into log management and potential security implications."
folder: "Cloud Logging"
---

# Table: gcp_logging_exclusion - Query Google Cloud Platform Logging Exclusions using SQL

Google Cloud Platform Logging Exclusions allow users to exclude certain logs from being stored, which can be critical for managing costs, avoiding unnecessary data retention, and adhering to privacy requirements. This service provides a way to set filters on logs based on resource type, log severity, and other attributes. It is an essential part of Google Cloud's logging and monitoring framework.

## Table Usage Guide

The `gcp_logging_exclusion` table provides insights into Logging Exclusions within Google Cloud Platform. As a security analyst or cloud administrator, explore exclusion-specific details through this table, including filters, descriptions, and associated metadata. Utilize it to uncover information about exclusions, such as those with broad filters, the resources affected by exclusions, and the verification of exclusion settings.

## Examples

### Basic info
Explore which logging exclusions are currently disabled in your Google Cloud Platform (GCP) system. This allows you to identify potential gaps in your logging coverage and rectify them for better system monitoring and security.

```sql+postgres
select
  name,
  disabled,
  filter,
  description
from
  gcp_logging_exclusion;
```

```sql+sqlite
select
  name,
  disabled,
  filter,
  description
from
  gcp_logging_exclusion;
```

### List of exclusions which are disabled
Explore which logging exclusions are currently disabled in your Google Cloud Platform (GCP) setup. This can help ensure you're capturing all necessary logs for audit and compliance purposes.

```sql+postgres
select
  name,
  disabled
from
  gcp_logging_exclusion
where
  disabled;
```

```sql+sqlite
select
  name,
  disabled
from
  gcp_logging_exclusion
where
  disabled = 1;
```