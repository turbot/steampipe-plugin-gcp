# Table: gcp_artifact_registry_repository

Google Cloud Artifact Registry is a managed artifact repository service provided by Google Cloud Platform (GCP). It is designed to store and manage software packages, container images, and other development artifacts. Artifact Registry helps organizations manage their dependencies, share software artifacts, and ensure the reliability and security of their software supply chain.

### Basic info

```sql
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

```sql
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

### List docker format package repositories

```sql
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

```sql
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

```sql
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

### Get docker configuration of repositories

```sql
select
  name,
  docker_config -> 'ImmutableTags' as immutable_tags,
  docker_config ->> 'ForceSendFields' as force_send_fields,
  docker_config ->> 'NullFields' as null_fields
from
  gcp_artifact_registry_repository;
```

### Get remote repository config details of repositories

```sql
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