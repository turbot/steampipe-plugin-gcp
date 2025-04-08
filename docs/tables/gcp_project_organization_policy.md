---
title: "Steampipe Table: gcp_project_organization_policy - Query GCP Project Organization Policies using SQL"
description: "Allows users to query Project Organization Policies in GCP, specifically the policy details, providing insights into organization level policies and their impact on the project."
folder: "Organization"
---

# Table: gcp_project_organization_policy - Query GCP Project Organization Policies using SQL

A Project Organization Policy in Google Cloud Platform (GCP) is a service that gives you the ability to manage and enforce consistent policy across your GCP resources. This service can be used to set fine-grained, resource-level policies anywhere in your resource hierarchy. It provides a simple and consistent way to manage and enforce organization-wide policies for your GCP resources.

## Table Usage Guide

The `gcp_project_organization_policy` table provides insights into Project Organization Policies within Google Cloud Platform (GCP). As a cloud engineer, explore policy-specific details through this table, including policy types, enforcement levels, and associated metadata. Utilize it to uncover information about policies, such as those with custom configurations, the hierarchical level of enforcement, and the verification of policy constraints.

## Examples

### Basic info
Explore which Google Cloud Platform (GCP) projects have recently been updated, including their unique identifiers and version numbers. This is useful for maintaining an overview of project changes and ensuring they are up-to-date.

```sql+postgres
select
  id,
  version,
  update_time
from
  gcp_project_organization_policy;
```

```sql+sqlite
select
  id,
  version,
  update_time
from
  gcp_project_organization_policy;
```

### Get organization policy constraints for each policy
Explore which organization policy constraints are applied to each policy within your Google Cloud Platform project. This can help in assessing the current policy configuration and ensure they align with your organization's security and compliance requirements.

```sql+postgres
select
  id,
  version,
  list_policy ->> 'allValues' as policy_value
from
  gcp_project_organization_policy;
```

```sql+sqlite
select
  id,
  version,
  json_extract(list_policy, '$.allValues') as policy_value
from
  gcp_project_organization_policy;
```