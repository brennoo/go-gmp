package commands

import "encoding/xml"

// GetAggregates represents a get_aggregates command request.
type GetAggregates struct {
	XMLName        xml.Name `xml:"get_aggregates"`
	Type           string   `xml:"type,attr"`
	Filter         string   `xml:"filter,attr,omitempty"`
	FiltID         string   `xml:"filt_id,attr,omitempty"`
	DataColumn     string   `xml:"data_column,attr,omitempty"`
	GroupColumn    string   `xml:"group_column,attr,omitempty"`
	SubgroupColumn string   `xml:"subgroup_column,attr,omitempty"`
	SortField      string   `xml:"sort_field,attr,omitempty"`
	SortOrder      string   `xml:"sort_order,attr,omitempty"`
	SortStat       string   `xml:"sort_stat,attr,omitempty"`
	FirstGroup     int      `xml:"first_group,attr,omitempty"`
	MaxGroups      int      `xml:"max_groups,attr,omitempty"`
	Mode           string   `xml:"mode,attr,omitempty"`
	UsageType      string   `xml:"usage_type,attr,omitempty"`
}

// GetAggregatesResponse represents a get_aggregates command response.
type GetAggregatesResponse struct {
	XMLName    xml.Name        `xml:"get_aggregates_response"`
	Status     string          `xml:"status,attr"`
	StatusText string          `xml:"status_text,attr"`
	Aggregate  *AggregateBlock `xml:"aggregate,omitempty"`
}

type AggregateBlock struct {
	DataType    string           `xml:"data_type,omitempty"`
	DataColumn  string           `xml:"data_column,omitempty"`
	GroupColumn string           `xml:"group_column,omitempty"`
	Groups      []AggregateGroup `xml:"group"`
	ColumnInfo  *ColumnInfo      `xml:"column_info,omitempty"`
}

type AggregateGroup struct {
	Value  string  `xml:"value,omitempty"`
	Count  int     `xml:"count,omitempty"`
	CCount int     `xml:"c_count,omitempty"`
	Min    float64 `xml:"min,omitempty"`
	Max    float64 `xml:"max,omitempty"`
	Mean   float64 `xml:"mean,omitempty"`
	Sum    float64 `xml:"sum,omitempty"`
	CSum   float64 `xml:"c_sum,omitempty"`
}

type ColumnInfo struct {
	Columns []AggregateColumn `xml:"aggregate_column"`
}

type AggregateColumn struct {
	Name     string `xml:"name,omitempty"`
	Stat     string `xml:"stat,omitempty"`
	Type     string `xml:"type,omitempty"`
	Column   string `xml:"column,omitempty"`
	DataType string `xml:"data_type,omitempty"`
}
