package gcp

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"google.golang.org/api/dns/v1"
)

//// TABLE DEFINITION

func tableGcpDnsManagedZone(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_dns_managed_zone",
		Description: "GCP DNS Managed Zone",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getDnsManagedZone,
			Tags:       map[string]string{"service": "dns", "action": "managedZones.get"},
		},
		List: &plugin.ListConfig{
			Hydrate: listDnsManagedZones,
			Tags:    map[string]string{"service": "dns", "action": "managedZones.list"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "An user assigned, friendly name that identifies the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique identifier for the resource, defined by the server.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "kind",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "dns_name",
				Description: "The DNS name of this managed zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "visibility",
				Description: "Specifies the zone's visibility. public zones are exposed to the Internet, while private zones are visible only to Virtual Private Cloud resources.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_time",
				Description: "The time that this resource was created on the server.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "A user-specified, human-readable description of the managed zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "dnssec_config_non_existence",
				Description: "Specifies the mechanism for authenticated denial-of-existence responses.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DnssecConfig.NonExistence"),
			},
			{
				Name:        "dnssec_config_state",
				Description: "Specifies whether DNSSEC is enabled, and what mode it is in.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DnssecConfig.State"),
			},
			{
				Name:        "name_server_set",
				Description: "Specifies the NameServerSet for this ManagedZone. A NameServerSet is a set of DNS name servers that all host the same ManagedZones.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_directory_config_namespace_deletion_time",
				Description: "The time that the namespace backing this zone was deleted.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ServiceDirectoryConfig.Namespace.DeletionTime").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "self_link",
				Description: "Server-defined URL for the managed zone.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDnsZoneSelfLink,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "dnssec_config_default_key_specs",
				Description: "Specifies parameters for generating initial DnsKeys for this ManagedZone.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DnssecConfig.DefaultKeySpecs"),
			},
			{
				Name:        "forwarding_config_target_name_servers",
				Description: "A list of target name servers to forward to.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ForwardingConfig.TargetNameServers"),
			},
			{
				Name:        "peering_config_target_network",
				Description: "Specifies the configuration of the network with which to peer.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PeeringConfig.TargetNetwork"),
			},
			{
				Name:        "private_visibility_config_networks",
				Description: "A set of Virtual Private Cloud resources that the zone is visible from.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("PrivateVisibilityConfig.Networks"),
			},
			{
				Name:        "labels",
				Description: "A set labels attached with the resource.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "name_servers",
				Description: "Delegate your managed_zone to these virtual name servers; defined by the server.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Labels"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDnsManagedZoneAka,
				Transform:   transform.FromValue(),
			},

			// Standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromConstant("global"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getProject,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listDnsManagedZones(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listDnsManagedZones")

	// Create Service Connection
	service, err := DnsService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Max limit isn't mentioned in the documentation
	// Default limit is set as 1000
	pageSize := types.Int64(1000)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *pageSize {
			pageSize = limit
		}
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	resp := service.ManagedZones.List(project).MaxResults(*pageSize)
	if err := resp.Pages(ctx, func(page *dns.ManagedZonesListResponse) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, managedZone := range page.ManagedZones {
			d.StreamListItem(ctx, managedZone)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				page.NextPageToken = ""
				return nil
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getDnsManagedZone(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDnsManagedZone")

	// Create Service Connection
	service, err := DnsService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)
	name := d.EqualsQuals["name"].GetStringValue()

	resp, err := service.ManagedZones.Get(project, name).Do()
	if err != nil {
		return nil, err
	}

	// The API doesn't return any error, if we pass any invalid parameter
	if len(resp.Name) > 0 {
		return resp, nil
	}

	return nil, nil
}

func getDnsManagedZoneAka(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(*dns.ManagedZone)

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	akas := []string{"gcp://dns.googleapis.com/projects/" + project + "/managedZones/" + data.Name}

	return akas, nil
}

func getDnsZoneSelfLink(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	zone := h.Item.(*dns.ManagedZone)

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	selfLink := "https://www.googleapis.com/dns/v1/projects/" + project + "/managedZones/" + zone.Name

	return selfLink, nil
}
