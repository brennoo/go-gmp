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
	XMLName    xml.Name   `xml:"get_aggregates_response"`
	Status     string     `xml:"status,attr"`
	StatusText string     `xml:"status_text,attr"`
	Aggregate  *Aggregate `xml:"aggregate,omitempty"`
}

// GetStatus returns the status and status text from the response.
func (r *GetAggregatesResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}

// Aggregate represents an aggregate block containing grouped data.
type Aggregate struct {
	DataType       string               `xml:"data_type,omitempty"`
	DataColumn     string               `xml:"data_column,omitempty"`
	GroupColumn    string               `xml:"group_column,omitempty"`
	SubgroupColumn string               `xml:"subgroup_column,omitempty"`
	TextColumn     []string             `xml:"text_column,omitempty"`
	Groups         []AggregateGroup     `xml:"group"`
	Overall        *AggregateOverall    `xml:"overall,omitempty"`
	Subgroups      *AggregateSubgroups  `xml:"subgroups,omitempty"`
	ColumnInfo     *AggregateColumnInfo `xml:"column_info,omitempty"`
}

// AggregateGroup represents a group within an aggregate with statistical data.
type AggregateGroup struct {
	Value    string              `xml:"value,omitempty"`
	Subgroup []AggregateSubgroup `xml:"subgroup,omitempty"`
	Count    int                 `xml:"count,omitempty"`
	CCount   int                 `xml:"c_count,omitempty"`
	Stats    []AggregateStats    `xml:"stats,omitempty"`
	Text     []AggregateText     `xml:"text,omitempty"`
	Min      string              `xml:"min,omitempty"`
	Max      string              `xml:"max,omitempty"`
	Mean     string              `xml:"mean,omitempty"`
	Sum      string              `xml:"sum,omitempty"`
	CSum     string              `xml:"c_sum,omitempty"`
}

// AggregateColumnInfo represents information about aggregate columns.
type AggregateColumnInfo struct {
	Columns []AggregateColumn `xml:"aggregate_column"`
}

// AggregateColumn represents information about a single aggregate column.
type AggregateColumn struct {
	Name     string `xml:"name,omitempty"`
	Stat     string `xml:"stat,omitempty"`
	Type     string `xml:"type,omitempty"`
	Column   string `xml:"column,omitempty"`
	DataType string `xml:"data_type,omitempty"`
}

// AggregateSubgroup represents a subgroup within an aggregate group.
type AggregateSubgroup struct {
	Value  string           `xml:"value,omitempty"`
	Count  int              `xml:"count,omitempty"`
	CCount int              `xml:"c_count,omitempty"`
	Stats  []AggregateStats `xml:"stats,omitempty"`
}

// AggregateStats represents statistical data for a column.
type AggregateStats struct {
	Column string `xml:"column,attr,omitempty"`
	Min    string `xml:"min,omitempty"`
	Max    string `xml:"max,omitempty"`
	Mean   string `xml:"mean,omitempty"`
	Sum    string `xml:"sum,omitempty"`
	CSum   string `xml:"c_sum,omitempty"`
}

// AggregateText represents a text column value.
type AggregateText struct {
	XMLName xml.Name `xml:"text"`
	Name    string   `xml:"name,attr,omitempty"`
	Value   string   `xml:",chardata"`
}

// AggregateOverall represents overall statistics for the entire aggregate.
type AggregateOverall struct {
	Count  int              `xml:"count,omitempty"`
	CCount int              `xml:"c_count,omitempty"`
	Stats  []AggregateStats `xml:"stats,omitempty"`
}

// AggregateSubgroups represents an overview of all subgroup values.
type AggregateSubgroups struct {
	Value []string `xml:"value,omitempty"`
}
