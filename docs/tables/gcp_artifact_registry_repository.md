---
title: "Steampipe Table: gcp_artifact_registry_repository - Query Artifact Registry Repositories using SQL"
description: "Allows users to query Artifact Registry Repositories in Google Cloud Platform (GCP), specifically the details about the repositories, including their name, location, format, and creation time."
folder: "Artifact Registry"
---

# Table: gcp_artifact_registry_repository - Query Artifact Registry Repositories using SQL

Artifact Registry is a scalable and managed repository provided by Google Cloud Platform (GCP) that allows teams to store, manage, and secure software packages. It supports various formats including Docker, Maven, and npm. Artifact Registry offers fine-grained access control, detailed audit logging, and can be integrated with Google Cloud’s other services.

## Table Usage Guide

The `gcp_artifact_registry_repository` table provides insights into the repositories within GCP's Artifact Registry. As a DevOps engineer, you can explore repository-specific details through this table, including their name, location, format, and creation time. Utilize it to manage and secure your software packages, ensuring they are stored in the correct format and location.

## Examples

### Basic info
Explore the basic details of your Google Cloud Platform's Artifact Registry repositories such as their names, cleanup policies, creation times, formats, key names, and modes. This information can help you manage and monitor your repositories more effectively.

```sql+postgres
select
  name,
  cleanup_policy_dry_run,
  create_time,
  format,
  kms_key_name,
  mode
from
  gcp_artifact_registry_repository;
```

```sql+sqlite
select
  name,
  cleanup_policy_dry_run,
  create_time,
  format,
  kms_key_name,
  mode
from
  gcp_artifact_registry_repository;
```

### List unencrypted repositories
Analyze the settings to understand which artifact repositories in your Google Cloud Platform do not have encryption enabled. This helps in identifying potential security risks and taking necessary actions to secure your repositories.

```sql+postgres
select
  name,
  cleanup_policy_dry_run,
  create_time,
  kms_key_name
from
  gcp_artifact_registry_repository
where
  kms_key_name = '';
```

```sql+sqlite
select
  name,
  cleanup_policy_dry_run,
  create_time,
  kms_key_name
from
  gcp_artifact_registry_repository
where
  kms_key_name is null;
```

### List docker format package repositories
Explore the GCP Artifact Registry to identify repositories that store Docker format packages. This can help you understand your usage patterns and manage your resources effectively.

```sql+postgres
select
  name,
  create_time,
  description,
  size_bytes,
  format
from
  gcp_artifact_registry_repository
where
  format = 'DOCKER';
```

```sql+sqlite
select
  name,
  create_time,
  description,
  size_bytes,
  format
from
  gcp_artifact_registry_repository
where
  format = 'DOCKER';
```

### List standard repositories
Explore which artifact repositories in GCP have been set to the 'standard' mode. This can help in assessing the configuration for optimal resource utilization and management.

```sql+postgres
select
  name,
  format,
  mode,
  create_time
from
  gcp_artifact_registry_repository
where
  mode = 'STANDARD_REPOSITORY';
```

```sql+sqlite
select
  name,
  format,
  mode,
  create_time
from
  gcp_artifact_registry_repository
where
  mode = 'STANDARD_REPOSITORY';
```

### List repositories that satisfies physical zone separation
Uncover the details of repositories that meet the physical zone separation requirements. This can be useful in assessing compliance with specific data residency and redundancy policies.

```sql+postgres
select
  name,
  mode,
  format,
  satisfies_pzs,
  description,
  create_time
from
  gcp_artifact_registry_repository
where
  satisfies_pzs;
```

```sql+sqlite
select
  name,
  mode,
  format,
  satisfies_pzs,
  description,
  create_time
from
  gcp_artifact_registry_repository
where
  satisfies_pzs = 1;
```

### Get docker configuration of repositories
Explore the configuration settings of Docker repositories to understand the immutability of tags and the handling of force send fields and null fields. This can be useful to review and manage your Docker repositories effectively.

```sql+postgres
select
  name,
  docker_config -> 'ImmutableTags' as immutable_tags,
  docker_config ->> 'ForceSendFields' as force_send_fields,
  docker_config ->> 'NullFields' as null_fields
from
  gcp_artifact_registry_repository;
```

```sql+sqlite
select
  name,
  json_extract(docker_config, '$.ImmutableTags') as immutable_tags,
  json_extract(docker_config, '$.ForceSendFields') as force_send_fields,
  json_extract(docker_config, '$.NullFields') as null_fields
from
  gcp_artifact_registry_repository;
```

### Get remote repository config details of repositories
Uncover the details of remote repositories' configurations to better understand the types of repositories being used, such as Apt, Docker, Maven, Npm, Python, and Yum. This can be useful for managing and optimizing the usage of different repository types in your Google Cloud Platform Artifact Registry.

```sql+postgres
select
  name,
  remote_repository_config ->> 'AptRepository' as apt_repository,
  remote_repository_config ->> 'DockerRepository' as docker_repository,
  remote_repository_config ->> 'MavenRepository' as maven_repository,
  remote_repository_config ->> 'NpmRepository' as npm_repository,
  remote_repository_config ->> 'PythonRepository' as python_repository,
  remote_repository_config ->> 'YumRepository' as yum_repository
from
  gcp_artifact_registry_repository;
```

```sql+sqlite
select
  name,
  json_extract(remote_repository_config, '$.AptRepository') as apt_repository,
  json_extract(remote_repository_config, '$.DockerRepository') as docker_repository,
  json_extract(remote_repository_config, '$.MavenRepository') as maven_repository,
  json_extract(remote_repository_config, '$.NpmRepository') as npm_repository,
  json_extract(remote_repository_config, '$.PythonRepository') as python_repository,
  json_extract(remote_repository_config, '$.YumRepository') as yum_repository
from
  gcp_artifact_registry_repository;
```