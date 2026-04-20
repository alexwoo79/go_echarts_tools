// Package model defines shared data types used across the application.
package model

// Task represents a single Gantt chart task row.
type Task struct {
	TaskName      string `json:"taskName"`
	Project       string `json:"project"`
	ColorGroup    string `json:"colorGroup"`
	StartISO      string `json:"startISO"`
	EndISO        string `json:"endISO"`
	PlanStartISO  string `json:"planStartISO"`
	PlanEndISO    string `json:"planEndISO"`
	DurationDays  int    `json:"durationDays"`
	Description   string `json:"description"`
	MilestoneName string `json:"milestoneName"`
	MilestoneISO  string `json:"milestoneISO"`
	Owner         string `json:"owner"`
}

// Stats holds summary statistics for a set of tasks.
type Stats struct {
	TaskCount            int     `json:"taskCount"`
	AvgDurationDays      float64 `json:"avgDurationDays"`
	TotalDurationDay     int     `json:"totalDurationDay"`
	MaxDurationDay       int     `json:"maxDurationDay"`
	PlanTotalDurationDay int     `json:"planTotalDurationDay"`
	HasPlanTotalDuration bool    `json:"hasPlanTotalDuration"`
}

// ChartOptions holds display/rendering options for a chart.
type ChartOptions struct {
	HierarchicalView bool   `json:"hierarchicalView"`
	ShowTaskDetails  bool   `json:"showTaskDetails"`
	ShowDuration     bool   `json:"showDuration"`
	DarkTheme        bool   `json:"darkTheme"`
	ChartTheme       string `json:"chartTheme"`
	TimeGranularity  string `json:"timeGranularity"`
}

// Dataset holds the raw tabular data parsed from an uploaded file.
type Dataset struct {
	ID      string
	Name    string
	Headers []string
	Rows    [][]string
}

// MappingDefaults holds auto-inferred column mapping suggestions.
type MappingDefaults struct {
	TaskCol          string
	StartCol         string
	EndCol           string
	ProjectCol       string
	ColorCol         string
	DescCol          string
	MilestoneCol     string
	MilestoneDateCol string
	PlanStartCol     string
	PlanEndCol       string
	OwnerCol         string
}

// MappingConfig holds the user-confirmed column mapping for chart building.
type MappingConfig struct {
	TaskCol          string
	StartCol         string
	EndCol           string
	ProjectCol       string
	ColorCol         string
	DescCol          string
	MilestoneCol     string
	MilestoneDateCol string
	PlanStartCol     string
	PlanEndCol       string
	OwnerCol         string
	SortByStart      bool
	ShowTaskNumber   bool
}
