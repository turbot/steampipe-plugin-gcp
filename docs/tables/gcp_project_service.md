---
title: "Steampipe Table: gcp_project_service - Query Google Cloud Platform Project Services using SQL"
description: "Allows users to query Project Services in Google Cloud Platform, specifically to get details about the enabled APIs and services for each project."
folder: "Project"
---

# Table: gcp_project_service - Query Google Cloud Platform Project Services using SQL

A Project Service in Google Cloud Platform is a service that has been enabled for a specific project. These services include various APIs and cloud-based solutions provided by Google, such as Compute Engine, App Engine, Cloud Storage, BigQuery, and more. Having access to project service information allows users to understand and manage the services and APIs that are currently in use for a given project.

## Table Usage Guide

The `gcp_project_service` table provides insights into the enabled services and APIs within a Google Cloud Platform project. As a cloud architect or administrator, you can explore service-specific details through this table, including service names, project IDs, and states of services. Use it to manage and monitor the enabled services and APIs, ensuring optimal utilization and compliance with organizational policies.

## Examples

### Basic info
Explore the various services associated with your Google Cloud Platform project. This query is useful for understanding what services are currently active within your project, aiding in project management and resource allocation.

```sql+postgres
select
  *
from
  gcp_project_service;
```

```sql+sqlite
select
  *
from
  gcp_project_service;
```

### List of services which are enabled
Explore which services are currently active within your Google Cloud Platform project. This is useful for maintaining awareness of your operational services and ensuring that only necessary services are enabled, enhancing security and cost-efficiency.

```sql+postgres
select
  name,
  state
from
  gcp_project_service
where
  state = 'ENABLED';
```

```sql+sqlite
select
  name,
  state
from
  gcp_project_service
where
  state = 'ENABLED';
```