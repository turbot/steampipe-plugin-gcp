package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"google.golang.org/api/dns/v1"
)

type recordSetInfo = struct {
	RecordSet       *dns.ResourceRecordSet
	ManagedZoneName string
}

//// TABLE DEFINITION

func tableDnsRecordSet(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_dns_record_set",
		Description: "GCP DNS Record Set",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"managed_zone_name", "name", "type"}),
			Hydrate:    getDnsRecordSet,
		},
		List: &plugin.ListConfig{
			Hydrate:           listDnsRecordSets,
			ParentHydrate:     listDnsManagedZones,
			ShouldIgnoreError: isNotFoundError([]string{"403"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the record set.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RecordSet.Name"),
			},
			{
				Name:        "managed_zone_name",
				Description: "An user assigned, friendly name that identifies the managed zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The identifier of a supported record type. See the list of Supported DNS record types.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RecordSet.Type"),
			},
			{
				Name:        "kind",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RecordSet.Kind"),
			},
			{
				Name:        "routing_policy",
				Description: "Configures dynamic query responses based on geo location of querying user or a weighted round robin based routing policy. A ResourceRecordSet should only have either rrdata (static) or routing_policy(dynamic). An error is returned otherwise.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("RecordSet.RoutingPolicy"),
			},
			{
				Name:        "rrdatas",
				Description: "As defined in RFC 1035 (section 5) and RFC 1034 (section 3.6.1).",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("RecordSet.Rrdatas"),
			},
			{
				Name:        "signature_rrdatas",
				Description: "As defined in RFC 4034 (section 3.2).",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("RecordSet.SignatureRrdatas"),
			},
			{
				Name:        "ttl",
				Description: "Number of seconds that this ResourceRecordSet can be cached by resolvers.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("RecordSet.Ttl"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RecordSet.Name"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDnsRecordSetAka,
				Transform:   transform.FromValue(),
			},

			// GCP standard columns
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

func listDnsRecordSets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listDnsRecordSets")

	// Get the details of Cloud DNS Managed Zone
	managedZone := h.Item.(*dns.ManagedZone)

	// Create Service Connection
	service, err := DnsService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get Project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project

	resp, err := service.ResourceRecordSets.List(project, managedZone.Name).Do()
	if err != nil {
		return nil, err
	}

	for _, recordset := range resp.Rrsets {
		d.StreamListItem(ctx, recordSetInfo{recordset, managedZone.Name})
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getDnsRecordSet(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDnsRecordSet")

	// Create service connection
	service, err := DnsService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project

	name := d.KeyColumnQuals["name"].GetStringValue()
	rrset_type := d.KeyColumnQuals["type"].GetStringValue()
	managedZoneName := d.KeyColumnQuals["managedzone_name"].GetStringValue()

	resp, err := service.ResourceRecordSets.List(project, managedZoneName).Name(name).Type(rrset_type).Do()
	if err != nil {
		return nil, err
	}

	return recordSetInfo{resp.Rrsets[0], managedZoneName}, nil
}

func getDnsRecordSetAka(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get project details
	projectData, err := activeProject(ctx, d)
	if err != nil {
		return nil, err
	}
	project := projectData.Project

	data := h.Item.(recordSetInfo)

	akas := []string{"gcp://dns.googleapis.com/projects/" + project + "/managedZones/" + data.ManagedZoneName + "/rrsets/" + data.RecordSet.Name + "/" + data.RecordSet.Type}
	return akas, nil
}
