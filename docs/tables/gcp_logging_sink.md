# Table:  table_gcp_logging_sink

Sinks represent filtered exports for log entries.

## Examples

### List writer identity that writes the export logs of logging sink

```sql
select
  name,
  unique_writer_identity
from
  gcp_logging_sink;
```


### List the destination path for each sink

```sql
select
  name,
  destination
from
  gcp_logging_sink;
```