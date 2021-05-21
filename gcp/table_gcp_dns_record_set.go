package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"google.golang.org/api/dns/v1"
)

type recordSetInfo = struct {
	RecordSet *dns.ResourceRecordSet
	ManagedZoneName   string
}

//// TABLE DEFINITION

func tableDnsRecordSet(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_dns_resord_set",
		Description: "GCP DNS Record Set",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"managedzone_name", "name", "type"}),
			Hydrate:    getDnsRecordSet,
		},
		List: &plugin.ListConfig{
			Hydrate:       listDnsRecordSet,
			ParentHydrate: listDnsManagedZones,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The record set name",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField("RecordSet.Name"),
			},
			{
				Name: "managedzone_name",
				Description: "An user assigned, friendly name that identifies the managed zone.",
				Type: proto.ColumnType_STRING,
				Transform: transform.FromField("ManagedZoneName"),
			},
			{
				Name:        "type",
				Description: "The identifier of a supported record type. See the list of Supported DNS record types.",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField("RecordSet.Type"),
			},
			{
				Name:        "kind",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform: transform.FromField("RecordSet.Kind"),
			},
			{
				Name:        "routing_policy",
				Description: "Configures dynamic query responses based on geo location of querying user or a weighted round robin based routing policy. A ResourceRecordSet should only have either rrdata (static) or routing_policy(dynamic). An error is returned otherwise.",
				Type:        proto.ColumnType_JSON,
				Transform: transform.FromField("RecordSet.RoutingPolicy"),
			},
			{
				Name:        "rrdatas",
				Description: "As defined in RFC 1035 (section 5) and RFC 1034 (section 3.6.1)",
				Type:        proto.ColumnType_JSON,
				Transform: transform.FromField("RecordSet.Rrdatas"),
			},
			{
				Name:        "signature_rrdatas",
				Description: "As defined in RFC 4034 (section 3.2).",
				Type:        proto.ColumnType_JSON,
				Transform: transform.FromField("RecordSet.SignatureRrdatas"),
			},
			{
				Name:        "ttl",
				Description: "Number of seconds that this ResourceRecordSet can be cached by resolvers.",
				Type:        proto.ColumnType_INT,
				Transform: transform.FromField("RecordSet.Ttl"),
			},
			{
				Name:        "force_send_fields",
				Description: "ForceSendFields is a list of field names (e.g. 'Kind') to unconditionally include in API requests. By default, fields with empty values are omitted from API requests. However, any non-pointer, non-interface field appearing in ForceSendFields will be sent to the server regardless of whether the field is empty or not. This may be used to include empty fields in Patch requests.",
				Type:        proto.ColumnType_JSON,
				Transform: transform.FromField("RecordSet.ForceSendFields"),
			},
			{
				Name:        "null_fields",
				Description: "NullFields is a list of field names (e.g. 'Kind') to include in API requests with the JSON null value. By default, fields with empty values are omitted from API requests. However, any field with an empty value appearing in NullFields will be sent to the server as null. It is an error if a field in this list has a non-empty value. This may be used to include null fields in Patch requests.",
				Type:        proto.ColumnType_JSON,
				Transform: transform.FromField("RecordSet.NullFields"),
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

//// LIST FUNCTIONS

func listDnsRecordSet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listDnsRecordSet")

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

func getDnsRecordSet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getDnsRecordSet")

	// Get the details of Cloud DNS Managed Zone
	managedZone := h.Item.(*dns.ManagedZone)

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
	managedZoneName := managedZone.Name

	resp, err := service.ResourceRecordSets.List(project, managedZoneName).Name(name).Type(rrset_type).Do()
	if err != nil {
		return nil, err
	}

	if len(resp.Rrsets) < 1 {
		return nil, nil
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
