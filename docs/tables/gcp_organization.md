---
title: "Steampipe Table: gcp_organization - Query GCP Organizations using SQL"
description: "Allows users to query GCP Organizations, specifically their metadata, providing insights into their configuration and governance."
folder: "Organization"
---

# Table: gcp_organization - Query GCP Organizations using SQL

A GCP Organization represents a collection of GCP resources that share common IAM policies. It is the root node in the GCP resource hierarchy and is associated with a domain that has a Google Workspace or Cloud Identity account. The Organization resource provides centralized control and oversight of all GCP resources.

## Table Usage Guide

The `gcp_organization` table provides insights into GCP Organizations within Google Cloud Platform. As a cloud architect or administrator, explore organization-specific details through this table, including their associated metadata, lifecycle state, directory customer ID, and more. Utilize it to uncover information about organizations, such as their creation time, owner details, and the verification of IAM policies.

**Important Notes**
- This table requires the `resourcemanager.organizations.get` permission to retrieve organization details.

## Examples

### Basic info
Explore the general details of your Google Cloud Platform organizations, such as its display name, associated organization ID, lifecycle state, and creation time. This information can help you assess the status and history of your organizations, which can be useful for administrative and auditing purposes.

```sql+postgres
select
  display_name,
  organization_id,
  lifecycle_state,
  creation_time
from
  gcp_organization;
```

```sql+sqlite
select
  display_name,
  organization_id,
  lifecycle_state,
  creation_time
from
  gcp_organization;
```

### Get essential contacts for organizations
Explore which essential contacts are associated with specific organizations. This is useful for quickly identifying key contacts within each organization, which can streamline communication and improve operational efficiency.

```sql+postgres
select
  organization_id,
  jsonb_pretty(essential_contacts) as essential_contacts
from
  gcp_organization;
```

```sql+sqlite
select
  organization_id,
  essential_contacts
from
  gcp_organization;
```