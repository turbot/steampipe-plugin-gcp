select name, self_link
from gcp.gcp_vpc_access_connector
where self_link = 'https://vpcaccess.googleapis.com/v1/{{ output.resource_id.value }}';