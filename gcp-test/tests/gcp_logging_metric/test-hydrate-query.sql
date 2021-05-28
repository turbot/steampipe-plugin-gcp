select name, metric_descriptor_display_name, metric_descriptor_metric_kind, metric_descriptor_type, metric_descriptor_unit, metric_descriptor_value_type, metric_descriptor_labels, filter, description, value_extractor, label_extractors
from gcp.gcp_logging_metric
where name = '{{resourceName}}'