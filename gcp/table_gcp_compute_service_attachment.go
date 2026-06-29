package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"

	compute "google.golang.org/api/compute/v0.beta"
)

//// TABLE DEFINITION

func tableGcpComputeServiceAttachment(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_compute_service_attachment",
		Description: "GCP Compute Service Attachment",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getComputeServiceAttachment,
			Tags:       map[string]string{"service": "compute", "action": "serviceAttachments.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listComputeServiceAttachments,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "connection_preference", Require: plugin.Optional, Operators: []string{"<>", "="}},
				{Name: "enable_proxy_protocol", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
			Tags: map[string]string{"service": "compute", "action": "serviceAttachments.list"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "Name of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique identifier for the resource.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "description",
				Description: "An optional description of this resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "region",
				Description: "URL of the region where the service attachment resides.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_timestamp",
				Description: "Creation timestamp in RFC3339 text format.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "connection_preference",
				Description: "The connection preference of the service attachment (ACCEPT_AUTOMATIC, ACCEPT_MANUAL, or CONNECTION_PREFERENCE_UNSPECIFIED).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "enable_proxy_protocol",
				Description: "If true, enable the proxy protocol for client IP forwarding through proxies.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "fingerprint",
				Description: "A hash of the contents stored in this object, used for optimistic locking.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "producer_forwarding_rule",
				Description: "URL of the internal load balancer forwarding rule serving the endpoint identified by this service attachment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "propagated_connection_limit",
				Description: "The number of consumer spokes that connected PSC endpoints can be propagated to via Network Connectivity Center.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "psc_service_attachment_id",
				Description: "An 128-bit global unique ID of the PSC service attachment.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "reconcile_connections",
				Description: "Whether a consumer accept/reject list change can reconcile statuses of existing PSC endpoints.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "target_service",
				Description: "URL of a service serving the endpoint identified by this service attachment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "connected_endpoints",
				Description: "An array of connections for all the consumers connected to this service attachment.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "consumer_accept_lists",
				Description: "Specifies which consumer projects or networks are allowed to connect to the service attachment.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "consumer_reject_lists",
				Description: "Specifies a list of projects or networks that are not allowed to connect to this service attachment.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "domain_names",
				Description: "Domain names used during integration between PSC connected endpoints and Cloud DNS.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "nat_subnets",
				Description: "URLs of subnets provided by the service producer to use for NAT in this service attachment.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tunneling_config",
				Description: "Tunneling configuration that, when set, encapsulates traffic between consumer and producer.",
				Type:        proto.ColumnType_JSON,
			},

			// standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromP(serviceAttachmentSelfLinkToTurbotData, "Akas"),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Region").Transform(lastPathElement),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromP(serviceAttachmentSelfLinkToTurbotData, "Project"),
			},
		},
	}
}

//// LIST FUNCTION

func listComputeServiceAttachments(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	service, err := ComputeBetaService(ctx, d)
	if err != nil {
		return nil, err
	}

	filterQuals := []filterQualMap{
		{"connection_preference", "connectionPreference", "string"},
		{"enable_proxy_protocol", "enableProxyProtocol", "boolean"},
	}

	filters := buildQueryFilterFromQuals(filterQuals, d.Quals)
	filterString := ""
	if len(filters) > 0 {
		filterString = strings.Join(filters, " ")
	}

	pageSize := types.Int64(500)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *pageSize {
			pageSize = limit
		}
	}

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	resp := service.ServiceAttachments.AggregatedList(project).Filter(filterString).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *compute.ServiceAttachmentAggregatedList) error {
		d.WaitForListRateLimit(ctx)

		for _, item := range page.Items {
			for _, sa := range item.ServiceAttachments {
				d.StreamListItem(ctx, sa)

				if d.RowsRemaining(ctx) == 0 {
					page.NextPageToken = ""
					return nil
				}
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getComputeServiceAttachment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	service, err := ComputeBetaService(ctx, d)
	if err != nil {
		return nil, err
	}

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	var found compute.ServiceAttachment
	name := d.EqualsQuals["name"].GetStringValue()

	resp := service.ServiceAttachments.AggregatedList(project).Filter("name=" + name)
	if err := resp.Pages(
		ctx,
		func(page *compute.ServiceAttachmentAggregatedList) error {
			for _, item := range page.Items {
				for _, sa := range item.ServiceAttachments {
					found = *sa
				}
			}
			return nil
		},
	); err != nil {
		return nil, err
	}

	if len(found.Name) < 1 {
		return nil, nil
	}

	return &found, nil
}

//// TRANSFORM FUNCTIONS

func serviceAttachmentSelfLinkToTurbotData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	sa := d.HydrateItem.(*compute.ServiceAttachment)
	param := d.Param.(string)

	project := strings.Split(sa.SelfLink, "/")[6]
	region := getLastPathElement(types.SafeString(sa.Region))

	turbotData := map[string]interface{}{
		"Project": project,
		"Akas":    []string{"gcp://compute.googleapis.com/projects/" + project + "/regions/" + region + "/serviceAttachments/" + sa.Name},
	}

	return turbotData[param], nil
}
