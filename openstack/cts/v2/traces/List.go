package traces

import (
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/internal/extract"
)

type ListTracesOpts struct {
	// Status of a trace. The value can be normal, warning, or incident.
	TraceStatus string `q:"trace_status,omitempty"`
	// Number of traces returned to the trace list. The default value is 50 and the maximum value is 200.
	Limit string `q:"limit,omitempty"`
	// UTC timestamp of the start time of the query time range. The value is in milliseconds and contains 13 digits.
	// Traces generated on the specified timestamp are not returned. The parameters from and to should be used together.
	From string `q:"from,omitempty"`
	// This parameter is used to query traces generated earlier than its specified value. The value can be that of marker in Table 2-34.
	// next can be used with from and to.
	// Traces generated in the overlap of the two time ranges specified respectively by next and by from and to will be returned.
	Next string `q:"next,omitempty"`
	// UTC timestamp of the end time of the query time range. The value is in milliseconds and contains 13 digits.
	// Traces generated on the specified timestamp are not returned. The parameters to and from should be used together.
	To string `q:"to,omitempty"`
	// Type of service whose traces are to be queried.
	// The value must be the acronym of a cloud service that has been connected with CTS.
	// It is a word composed of uppercase letters.
	// For cloud services that can be connected with CTS,
	// see section "Supported Services and Operations" in the Cloud Trace Service User Guide.
	ServiceType string `q:"service_type,omitempty"`
	// Name of the user whose traces are to be queried.
	// NOTE The username is case-sensitive.
	User string `q:"user,omitempty"`
	// ID of a cloud resource whose traces are to be queried.
	ResourceId string `q:"resource_id,omitempty"`
	// Name of a resource whose traces are to be queried.
	// NOTE The resource name is case-sensitive.
	ResourceName string `q:"resource_name,omitempty"`
	// Type of resource whose traces are to be queried. The value can contain 1 to 64 characters, including letters,
	// digits, hyphens (-), underscores (_), and periods (.). It must start with a letter.
	// For cloud services that can be connected with CTS, see section "Supported Services and Operations" in the Cloud Trace Service User Guide.
	ResourceType string `q:"resource_type,omitempty"`
	// Trace ID.
	// If this parameter is specified, other query criteria will not take effect.
	TraceId string `q:"trace_id,omitempty"`
	// Trace name. It indicates the operation recorded by this trace.
	// NOTE The trace name is case-sensitive.
	TraceName string `q:"trace_name,omitempty"`
}

func List(client *golangsdk.ServiceClient, trackerName string, opts ListTracesOpts) (*ListTracesResponse, error) {
	url, err := golangsdk.NewURLBuilder().WithEndpoints(trackerName, "trace").WithQueryParams(&opts).Build()
	if err != nil {
		return nil, err
	}

	// GET /v2.0/{project_id}/{tracker_name}/trace
	raw, err := client.Get(client.ServiceURL(url.String()), nil, nil)
	if err != nil {
		return nil, err
	}

	var res ListTracesResponse
	err = extract.Into(raw.Body, &res)
	return &res, err
}

type ListTracesResponse struct {
	Traces   []Traces `json:"traces,omitempty"`
	MetaData MetaData `json:"meta_data,omitempty"`
}

type Traces struct {
	// ID of a cloud resource on which the recorded operation was performed.
	ResourceId string `json:"resource_id,omitempty"`
	// Name of a trace. The value can contain 1 to 64 characters,
	// including letters, digits, hyphens (-), underscores (_), and periods (.). It must start with a letter.
	TraceName string `json:"trace_name,omitempty"`
	// Trace status. The value can be normal, warning, or incident.
	TraceStatus string `json:"trace_status,omitempty"`
	// Trace source. The value can be ApiCall, ConsoleAction, or SystemAction.
	TraceType string `json:"trace_type,omitempty"`
	// Request of an operation on resources.
	Request string `json:"request,omitempty"`
	// Response to a user request, that is, the returned information for an operation on resources.
	Response string `json:"response,omitempty"`
	// HTTP status code returned by the associated API.
	Code string `json:"code,omitempty"`
	// Version of the API.
	ApiVersion string `json:"api_version,omitempty"`
	// Remarks added by other cloud services to a trace.
	Message string `json:"message,omitempty"`
	// Timestamp when an operation was recorded by CTS.
	RecordTime int64 `json:"record_time,omitempty"`
	// Trace ID. The value is the UUID generated by the system.
	TraceId string `json:"trace_id,omitempty"`
	// Timestamp when a trace was generated.
	Time int64 `json:"time,omitempty"`
	// Information of the user who performed the operation that triggered the trace.
	User UserInfo `json:"user,omitempty"`
	// Type of service whose traces are to be queried.
	// The value must be the acronym of a cloud service that has been connected with CTS.
	// It is a word composed of uppercase letters.
	ServiceType string `json:"service_type,omitempty"`
	// Type of resource whose traces are to be queried. The value can contain 1 to 64 characters,
	// including letters, digits, hyphens (-), underscores (_), and periods (.). It must start with a letter.
	ResourceType string `json:"resource_type,omitempty"`
	// IP address of the tenant who performed the operation that triggered the trace.
	SourceIp string `json:"source_ip,omitempty"`
	// Name of a resource on which the recorded operation was performed.
	ResourceName string `json:"resource_name,omitempty"`
	// Request ID.
	RequestId string `json:"request_id,omitempty"`
	// Additional information required for fault locating after a request error.
	LocationInfo string `json:"location_info,omitempty"`
	// Endpoint in the details page URL of the cloud resource on which the recorded operation was performed.
	Endpoint string `json:"endpoint,omitempty"`
	// Details page URL (excluding the endpoint) of the cloud resource on which the recorded operation was performed.
	ResourceUrl string `json:"resource_url,omitempty"`
}

type UserInfo struct {
	Id     string   `json:"id,omitempty"`
	Name   string   `json:"name,omitempty"`
	Domain BaseUser `json:"domain,omitempty"`
}

type BaseUser struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type MetaData struct {
	// Number of returned traces.
	Count int `json:"count,omitempty"`
	// ID of the last trace in the returned trace list. The value of this parameter can be used as the next value.
	// If the value of marker is null, all traces have been returned.
	Marker string `json:"marker,omitempty"`
}
