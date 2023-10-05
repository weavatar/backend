// This file is auto-generated, don't edit it. Thanks.
package summary

import (
	"github.com/alibabacloud-go/tea/tea"
)

type Paths struct {
}

func (s Paths) String() string {
	return tea.Prettify(s)
}

func (s Paths) GoString() string {
	return s.String()
}

type Parameters struct {
	// {"en" : "RFC 3339 date indicating the beginning of the time period. The time must be specified using the UTC timezone; it cannot be an offset. Example: startdate=2019-10-30T00:00:00Z ", "zh_CN": "查询范围的开始时间，以RFC 3339日期格式表示。必须使用UTC时区指定时间。示例：startdate=2019-10-30T00:00:00Z。"}
	Startdate *string `json:"startdate,omitempty" xml:"startdate,omitempty" require:"true"`
	// {"en" : "RFC 3339 date indicating the end of the time period. The time must be specified using the UTC timezone; it cannot be an offset. Example: enddate=2019-11-14T00:00:00Z Your enddate may be rounded up to the nearest minute, hour, or day depending on the type parameter. For example, if you enter enddate=2019-09-05T03:14:01Z&type=hourly, the response includes data ending 2019-09-05T04:00:00Z. Due to latency associated with new traffic data, enddate should be no later than five minutes before the current time. This ensures you get the most accurate results.", "zh_CN": "查询范围的结束时间，以RFC 3339日期格式表示。必须使用UTC时区指定时间。示例：enddate=2019-11-14T00:00:00Z。由于数据处理存在延迟，所指定的结束时间必须至少比当前时间早5分钟，否则返回的数据可能不准确。"}
	Enddate *string `json:"enddate,omitempty" xml:"enddate,omitempty" require:"true"`
	// {"en" : "Enum: http https all
	// Default: all
	// Limit the results to the specified scheme. By default, data from HTTPS and HTTP requests is returned.", "zh_CN": "[ 0 .. 5 ] 字符
	// 取值范围: http, https, all
	// 默认值: all
	// 指定查询HTTP与/或HTTPS协议的数据。默认查询全部2种协议的数据。"}
	Scheme *string `json:"scheme,omitempty" xml:"scheme,omitempty"`
}

func (s Parameters) String() string {
	return tea.Prettify(s)
}

func (s Parameters) GoString() string {
	return s.String()
}

func (s *Parameters) SetStartdate(v string) *Parameters {
	s.Startdate = &v
	return s
}

func (s *Parameters) SetEnddate(v string) *Parameters {
	s.Enddate = &v
	return s
}

func (s *Parameters) SetScheme(v string) *Parameters {
	s.Scheme = &v
	return s
}

type RequestHeader struct {
}

func (s RequestHeader) String() string {
	return tea.Prettify(s)
}

func (s RequestHeader) GoString() string {
	return s.String()
}

type GetASummaryOfRequestsRequest struct {
	// {"en" : "Specify conditions to filter report data.", "zh_CN": "指定查询条件过滤报表数据。"}
	Filters *GetASummaryOfRequestsRequestFilters `json:"filters,omitempty" xml:"filters,omitempty" type:"Struct"`
	// {"en" : "<= 2 items
	// items Enum: hostnames, serverGroups
	// You can group results using a combination of up to two of the following: 'hostnames', and 'serverGroups'.", "zh_CN": "<= 2 条目
	// 取值范围: hostnames, serverGroups
	// 指定分组依据对数据进行分组汇总。支持按'hostnames'，'serverGroups'单独进行分组汇总，也支持同时指定这2个参数进行分组汇总。"}
	GroupBy []*string `json:"groupBy,omitempty" xml:"groupBy,omitempty" type:"Repeated"`
}

func (s GetASummaryOfRequestsRequest) String() string {
	return tea.Prettify(s)
}

func (s GetASummaryOfRequestsRequest) GoString() string {
	return s.String()
}

func (s *GetASummaryOfRequestsRequest) SetFilters(v *GetASummaryOfRequestsRequestFilters) *GetASummaryOfRequestsRequest {
	s.Filters = v
	return s
}

func (s *GetASummaryOfRequestsRequest) SetGroupBy(v []*string) *GetASummaryOfRequestsRequest {
	s.GroupBy = v
	return s
}

type GetASummaryOfRequestsRequestFilters struct {
	// {"en" : "List of hostnames for which to return data. Wildcard hostnames such as *.domain.com are also permitted. If unspecified, data from all hostnames will be returned.", "zh_CN": "指定加速域名进行查询。可使用泛域名，如*.domain.com。如果未指定，将返回所有加速域名的数据。"}
	Hostnames []*string `json:"hostnames,omitempty" xml:"hostnames,omitempty" type:"Repeated"`
	// {"en" : "Indicates one or more server groups.", "zh_CN": "指定serverGroups（节点组）进行查询。"}
	ServerGroups []*string `json:"serverGroups,omitempty" xml:"serverGroups,omitempty" type:"Repeated"`
}

func (s GetASummaryOfRequestsRequestFilters) String() string {
	return tea.Prettify(s)
}

func (s GetASummaryOfRequestsRequestFilters) GoString() string {
	return s.String()
}

func (s *GetASummaryOfRequestsRequestFilters) SetHostnames(v []*string) *GetASummaryOfRequestsRequestFilters {
	s.Hostnames = v
	return s
}

func (s *GetASummaryOfRequestsRequestFilters) SetServerGroups(v []*string) *GetASummaryOfRequestsRequestFilters {
	s.ServerGroups = v
	return s
}

type ResponseHeader struct {
}

func (s ResponseHeader) String() string {
	return tea.Prettify(s)
}

func (s ResponseHeader) GoString() string {
	return s.String()
}

type GetASummaryOfRequestsResponse struct {
	// {"en" : "This object contains fields describing the data returned in the groups object.", "zh_CN": "此对象包含的字段是对groups对象中返回数据的描述。"}
	MetaData *GetASummaryOfRequestsResponseMetaData `json:"metaData,omitempty" xml:"metaData,omitempty" require:"true" type:"Struct"`
	// {"en" : "This object contains the breakdown of requests by group.", "zh_CN": "按分组返回数据。"}
	Groups []*GetASummaryOfRequestsResponseGroups `json:"groups,omitempty" xml:"groups,omitempty" require:"true" type:"Repeated"`
}

func (s GetASummaryOfRequestsResponse) String() string {
	return tea.Prettify(s)
}

func (s GetASummaryOfRequestsResponse) GoString() string {
	return s.String()
}

func (s *GetASummaryOfRequestsResponse) SetMetaData(v *GetASummaryOfRequestsResponseMetaData) *GetASummaryOfRequestsResponse {
	s.MetaData = v
	return s
}

func (s *GetASummaryOfRequestsResponse) SetGroups(v []*GetASummaryOfRequestsResponseGroups) *GetASummaryOfRequestsResponse {
	s.Groups = v
	return s
}

type GetASummaryOfRequestsResponseMetaData struct {
	// {"en" : "RFC 3339 date indicating the beginning of the period.", "zh_CN": "RFC 3339格式的日期，表示查询的起始时间。"}
	StartTime *string `json:"startTime,omitempty" xml:"startTime,omitempty"`
	// {"en" : "RFC 3339 date indicating the end of the period.", "zh_CN": "RFC 3339格式的日期，表示查询的结束时间。"}
	EndTime *string `json:"endTime,omitempty" xml:"endTime,omitempty"`
	// {"en" : "The response can contain up to 10000 groups. If there are more groups, isComplete will be false.", "zh_CN": "该接口最多返回10000个分组的数据。如果实际分组数量大于10000，则isComplete将为false。"}
	IsComplete *bool `json:"isComplete,omitempty" xml:"isComplete,omitempty" require:"true"`
	// {"en" : "Indicates the type of data returned. 'edge response' represents edge traffic. 'fast route response' refers to traffic from your origin accelerated through the HDT product. The order of the entries in dataNames array corresponds to the order of values returned in the data data array in the response.", "zh_CN": "表示返回的数据类型。'edge response'表示边缘请求数，'fast route response'表示快速回源请求数。dataNames数组中条目的顺序与groups[].data中返回值的顺序一一对应。"}
	DataNames []*string `json:"dataNames,omitempty" xml:"dataNames,omitempty" type:"Repeated"`
	// {"en" : "Indicates the unit of measurement of the returned values.", "zh_CN": "返回值的计量单位。"}
	DataUnit *string `json:"dataUnit,omitempty" xml:"dataUnit,omitempty"`
}

func (s GetASummaryOfRequestsResponseMetaData) String() string {
	return tea.Prettify(s)
}

func (s GetASummaryOfRequestsResponseMetaData) GoString() string {
	return s.String()
}

func (s *GetASummaryOfRequestsResponseMetaData) SetStartTime(v string) *GetASummaryOfRequestsResponseMetaData {
	s.StartTime = &v
	return s
}

func (s *GetASummaryOfRequestsResponseMetaData) SetEndTime(v string) *GetASummaryOfRequestsResponseMetaData {
	s.EndTime = &v
	return s
}

func (s *GetASummaryOfRequestsResponseMetaData) SetIsComplete(v bool) *GetASummaryOfRequestsResponseMetaData {
	s.IsComplete = &v
	return s
}

func (s *GetASummaryOfRequestsResponseMetaData) SetDataNames(v []*string) *GetASummaryOfRequestsResponseMetaData {
	s.DataNames = v
	return s
}

func (s *GetASummaryOfRequestsResponseMetaData) SetDataUnit(v string) *GetASummaryOfRequestsResponseMetaData {
	s.DataUnit = &v
	return s
}

type GetASummaryOfRequestsResponseGroups struct {
	// {"en" : "Name of a group.  '__all__' is a special group encompassing all groups.", "zh_CN": "分组名称。'__all__' 是一个特殊分组，包含其它所有分组的数据。"}
	Group *string `json:"group,omitempty" xml:"group,omitempty"`
	// {"en" : "Data values. The units of measurement are determined by the dataUnit field.", "zh_CN": "请求数。"}
	Data []*float64 `json:"data,omitempty" xml:"data,omitempty" type:"Repeated"`
}

func (s GetASummaryOfRequestsResponseGroups) String() string {
	return tea.Prettify(s)
}

func (s GetASummaryOfRequestsResponseGroups) GoString() string {
	return s.String()
}

func (s *GetASummaryOfRequestsResponseGroups) SetGroup(v string) *GetASummaryOfRequestsResponseGroups {
	s.Group = &v
	return s
}

func (s *GetASummaryOfRequestsResponseGroups) SetData(v []*float64) *GetASummaryOfRequestsResponseGroups {
	s.Data = v
	return s
}
