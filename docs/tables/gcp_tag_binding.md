---
title: "Steampipe Table: gcp_tag_binding - Query GCP Tag Bindings using SQL"
description: "Allows users to query GCP Tag Bindings, specifically the detailed information about each tag binding in the Google Cloud project."
folder: "Resource Tags"
---

# Table: gcp_tag_binding - Query GCP Tag Bindings using SQL

Google Cloud's Tag Manager allows users to attach metadata to GCP resources in the form of tags. Tags are key-value pairs that help in organizing and managing resources, such as VM instances, storage buckets, and networks. Tag bindings are the association between a tag and a resource, which allows users to filter and manage resources based on their metadata.

## Examples

### Basic info

Explore which resources are tagged within your Google Cloud Platform, gaining insights into the associated tags and resource types. This can be particularly useful for managing and organizing your resources.

```sql+postgres
select
  name,
  parent,
  tag_value,
  title
from
  gcp_tag_binding
where 
  parent='//cloudresourcemanager.googleapis.com/projects/your-project-id';
```

```sql+sqlite
select
  name,
  parent,
  tag_value,
  title
from
  gcp_tag_binding
where 
  parent='//cloudresourcemanager.googleapis.com/projects/your-project-id';
```
