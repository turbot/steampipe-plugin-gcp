select name, parent, tag_value, title
from gcp.gcp_tag_binding
where parent = '{{ output.parent.value }}' and tag_value = '{{ output.tag_value.value }}';