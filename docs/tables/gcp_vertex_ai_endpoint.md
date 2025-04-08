---
title: "Steampipe Table: gcp_vertex_ai_endpoint - Query GCP Vertex AI Endpoints using SQL"
description: "Allows users to query GCP Vertex AI Endpoints, specifically the detailed information about each endpoint in the Google Cloud project."
folder: "Vertex AI"
---

# Table: gcp_vertex_ai_endpoint - Query GCP Vertex AI Endpoints using SQL

Google Cloud's Vertex AI is a unified ML platform for building, deploying, and scaling AI applications. It offers a suite of tools and services for data scientists and ML engineers, which includes the ability to manage endpoints. An endpoint in Vertex AI is a resource that serves predictions from one or more deployed models.

## Table Usage Guide

The `gcp_vertex_ai_endpoint` table provides insights into Vertex AI endpoints within Google Cloud Platform. As a data scientist or ML engineer, explore endpoint-specific details through this table, including deployed models, traffic split, and associated metadata. Utilize it to manage and monitor your AI application endpoints, such as those serving high traffic, the distribution of traffic among deployed models, and the status of each endpoint.

## Examples

### Basic info
Explore which AI endpoints are active within your Google Cloud Platform, gaining insights into their creation time and associated networks. This can be particularly useful for managing and auditing your AI resources.

```sql+postgres
select
  name,
  display_name,
  create_time,
  network
from
  gcp_vertex_ai_endpoint;
```

```sql+sqlite
select
  name,
  display_name,
  create_time,
  network
from
  gcp_vertex_ai_endpoint;
```

### List endpoints that are exposed via private service connect
Explore which endpoints are made accessible through private service connect to enhance your understanding of the network's security and accessibility. This can help in assessing potential vulnerabilities and managing access controls.

```sql+postgres
select
  name,
  display_name,
  create_time,
  enable_private_service_connect
from
  gcp_vertex_ai_endpoint
where
  enable_private_service_connect;
```

```sql+sqlite
select
  name,
  display_name,
  create_time,
  enable_private_service_connect
from
  gcp_vertex_ai_endpoint
where
  enable_private_service_connect;
```

### List endpoints created in the last 30 days
Determine the areas in which new endpoints have been established within the past month. This can be useful for tracking recent changes and developments in your network.

```sql+postgres
select
  name,
  display_name,
  network,
  create_time,
  update_time
from
  gcp_vertex_ai_endpoint
where
  create_time >= now() - interval '30' day;
```

```sql+sqlite
select
  name,
  display_name,
  network,
  create_time,
  update_time
from
  gcp_vertex_ai_endpoint
where
  create_time >= datetime('now', '-30 day');
```

### Get customer-managed key details of endpoints
Explore the encryption specifics of your AI endpoints to understand their security setup and creation time. This information could be useful in assessing the security measures in place and identifying potential areas for improvement.

```sql+postgres
select
  name,
  create_time,
  encryption_spec ->> 'KmsKeyName' as kms_key_name
from
  gcp_vertex_ai_endpoint;
```

```sql+sqlite
select
  name,
  create_time,
  json_extract(encryption_spec, '$.KmsKeyName') as kms_key_name
from
  gcp_vertex_ai_endpoint;
```

### Get prediction request response config of endpoints
Explore the configuration of prediction request responses in AI endpoints to understand if logging is enabled, the sampling rate, and whether the destination is BigQuery. This can help optimize your data analysis process by ensuring the right logs are being captured and stored at the correct location.

```sql+postgres
select
  name,
  network,
  predict_request_response_logging_config ->> 'Enabled' as enabled,
  predict_request_response_logging_config ->> 'SamplingRate' as sampling_rate,
  predict_request_response_logging_config ->> 'BigqueryDestination' as bigquery_destination
from
  gcp_vertex_ai_endpoint;
```

```sql+sqlite
select
  name,
  network,
  json_extract(predict_request_response_logging_config, '$.Enabled') as enabled,
  json_extract(predict_request_response_logging_config, '$.SamplingRate') as sampling_rate,
  json_extract(predict_request_response_logging_config, '$.BigqueryDestination') as bigquery_destination
from
  gcp_vertex_ai_endpoint;
```