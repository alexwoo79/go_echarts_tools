// Package gantt implements the Gantt chart builder and auto-registers it.
package gantt

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"gantt/internal/charts"
	"gantt/internal/data"
	"gantt/internal/model"
)

func init() {
	charts.Register(New())
}

// Result is the JSON-serialisable output produced by the Gantt builder.
type Result struct {
	Tasks []model.Task `json:"tasks"`
	Stats model.Stats  `json:"stats"`
}

// Builder implements charts.ChartBuilder for Gantt charts.
type Builder struct{}

// New returns a new Gantt Builder.
func New() Builder { return Builder{} }

func (Builder) ID() string   { return "gantt" }
func (Builder) Name() string { return "甘特图" }

// DefaultOptions returns sensible defaults for Gantt rendering.
func (Builder) DefaultOptions() model.ChartOptions {
	return model.ChartOptions{
		HierarchicalView: true,
		ShowTaskDetails:  true,
		ShowDuration:     true,
		DarkTheme:        false,
		ChartTheme:       "default",
		TimeGranularity:  "month",
	}
}

// InferDefaults guesses column mapping from header names.
func (Builder) InferDefaults(headers []string) model.MappingConfig {
	def := inferDefaults(headers)
	return model.MappingConfig{
		TaskCol:          def.TaskCol,
		StartCol:         def.StartCol,
		EndCol:           def.EndCol,
		ProjectCol:       def.ProjectCol,
		ColorCol:         def.ColorCol,
		DescCol:          def.DescCol,
		MilestoneCol:     def.MilestoneCol,
		MilestoneDateCol: def.MilestoneDateCol,
		PlanStartCol:     def.PlanStartCol,
		PlanEndCol:       def.PlanEndCol,
		OwnerCol:         def.OwnerCol,
	}
}

// Build converts a dataset and mapping config into a Gantt Result.
func (b Builder) Build(dataset model.Dataset, cfg model.MappingConfig, _ model.ChartOptions) (interface{}, error) {
	tasks, err := buildTasks(dataset, cfg)
	if err != nil {
		return nil, err
	}
	return Result{Tasks: tasks, Stats: computeStats(tasks)}, nil
}

func buildTasks(dataset model.Dataset, cfg model.MappingConfig) ([]model.Task, error) {
	if cfg.TaskCol == "" || cfg.StartCol == "" || cfg.EndCol == "" {
		return nil, fmt.Errorf("请至少选择任务列、开始日期列、结束日期列")
	}

	headerIndex := make(map[string]int, len(dataset.Headers))
	for i, h := range dataset.Headers {
		headerIndex[h] = i
	}

	getIdx := func(col string) int {
		if col == "" {
			return -1
		}
		idx, ok := headerIndex[col]
		if !ok {
			return -1
		}
		return idx
	}

	taskIdx := getIdx(cfg.TaskCol)
	startIdx := getIdx(cfg.StartCol)
	endIdx := getIdx(cfg.EndCol)
	projectIdx := getIdx(cfg.ProjectCol)
	colorIdx := getIdx(cfg.ColorCol)
	descIdx := getIdx(cfg.DescCol)
	mileIdx := getIdx(cfg.MilestoneCol)
	mileDateIdx := getIdx(cfg.MilestoneDateCol)
	planStartIdx := getIdx(cfg.PlanStartCol)
	planEndIdx := getIdx(cfg.PlanEndCol)
	ownerIdx := getIdx(cfg.OwnerCol)

	if taskIdx < 0 || startIdx < 0 || endIdx < 0 {
		return nil, fmt.Errorf("列映射无效，请重新选择必填列")
	}

	tasks := make([]model.Task, 0, len(dataset.Rows))
	for _, row := range dataset.Rows {
		taskName := data.Cell(row, taskIdx)
		if taskName == "" {
			continue
		}

		startAt, err := data.ParseDate(data.Cell(row, startIdx))
		if err != nil {
			continue
		}
		endAt, err := data.ParseDate(data.Cell(row, endIdx))
		if err != nil {
			continue
		}
		if endAt.Before(startAt) {
			startAt, endAt = endAt, startAt
		}

		project := data.Cell(row, projectIdx)
		if project == "" {
			project = "未分组"
		}

		colorGroup := data.Cell(row, colorIdx)
		if colorGroup == "" {
			colorGroup = project
		}

		planStartISO := ""
		if t, err := data.ParseDate(data.Cell(row, planStartIdx)); err == nil {
			planStartISO = t.Format(time.RFC3339)
		}
		planEndISO := ""
		if t, err := data.ParseDate(data.Cell(row, planEndIdx)); err == nil {
			planEndISO = t.Format(time.RFC3339)
		}

		mileName := data.Cell(row, mileIdx)
		mileISO := ""
		if t, err := data.ParseDate(data.Cell(row, mileDateIdx)); err == nil {
			mileISO = t.Format(time.RFC3339)
		} else if mileName != "" {
			mileISO = startAt.Format(time.RFC3339)
		}

		days := int(endAt.Sub(startAt).Hours()/24) + 1
		if days < 1 {
			days = 1
		}

		tasks = append(tasks, model.Task{
			TaskName:      taskName,
			Project:       project,
			ColorGroup:    colorGroup,
			StartISO:      startAt.Format(time.RFC3339),
			EndISO:        endAt.Format(time.RFC3339),
			PlanStartISO:  planStartISO,
			PlanEndISO:    planEndISO,
			DurationDays:  days,
			Description:   data.Cell(row, descIdx),
			MilestoneName: mileName,
			MilestoneISO:  mileISO,
			Owner:         data.Cell(row, ownerIdx),
		})
	}

	if len(tasks) == 0 {
		return nil, fmt.Errorf("未解析出有效任务，请检查日期列格式")
	}

	if cfg.SortByStart {
		sort.SliceStable(tasks, func(i, j int) bool {
			if tasks[i].ColorGroup == tasks[j].ColorGroup {
				return tasks[i].StartISO < tasks[j].StartISO
			}
			return tasks[i].ColorGroup < tasks[j].ColorGroup
		})
	} else {
		sort.SliceStable(tasks, func(i, j int) bool {
			if tasks[i].ColorGroup == tasks[j].ColorGroup {
				return i < j
			}
			return tasks[i].ColorGroup < tasks[j].ColorGroup
		})
	}

	if cfg.ShowTaskNumber {
		for i := range tasks {
			tasks[i].TaskName = fmt.Sprintf("%02d  %s", i+1, tasks[i].TaskName)
		}
	}

	return tasks, nil
}

func computeStats(tasks []model.Task) model.Stats {
	if len(tasks) == 0 {
		return model.Stats{}
	}

	totalTaskDuration := 0
	maxDur := 0

	var actualMinStart, actualMaxEnd time.Time
	var planMinStart, planMaxEnd time.Time
	hasPlan := false

	for _, t := range tasks {
		totalTaskDuration += t.DurationDays
		if t.DurationDays > maxDur {
			maxDur = t.DurationDays
		}

		startAt, startErr := time.Parse(time.RFC3339, t.StartISO)
		endAt, endErr := time.Parse(time.RFC3339, t.EndISO)
		if startErr == nil && endErr == nil {
			if endAt.Before(startAt) {
				startAt, endAt = endAt, startAt
			}
			if actualMinStart.IsZero() || startAt.Before(actualMinStart) {
				actualMinStart = startAt
			}
			if actualMaxEnd.IsZero() || endAt.After(actualMaxEnd) {
				actualMaxEnd = endAt
			}
		}

		planStartAt, planStartErr := time.Parse(time.RFC3339, t.PlanStartISO)
		planEndAt, planEndErr := time.Parse(time.RFC3339, t.PlanEndISO)
		if planStartErr == nil && planEndErr == nil {
			if planEndAt.Before(planStartAt) {
				planStartAt, planEndAt = planEndAt, planStartAt
			}
			if !hasPlan {
				planMinStart = planStartAt
				planMaxEnd = planEndAt
				hasPlan = true
			} else {
				if planStartAt.Before(planMinStart) {
					planMinStart = planStartAt
				}
				if planEndAt.After(planMaxEnd) {
					planMaxEnd = planEndAt
				}
			}
		}
	}

	actualSpan := 0
	if !actualMinStart.IsZero() && !actualMaxEnd.IsZero() {
		actualSpan = int(actualMaxEnd.Sub(actualMinStart).Hours()/24) + 1
		if actualSpan < 1 {
			actualSpan = 1
		}
	}

	planSpan := 0
	if hasPlan {
		planSpan = int(planMaxEnd.Sub(planMinStart).Hours()/24) + 1
		if planSpan < 1 {
			planSpan = 1
		}
	}

	return model.Stats{
		TaskCount:            len(tasks),
		AvgDurationDays:      float64(totalTaskDuration) / float64(len(tasks)),
		TotalDurationDay:     actualSpan,
		MaxDurationDay:       maxDur,
		PlanTotalDurationDay: planSpan,
		HasPlanTotalDuration: hasPlan,
	}
}

func guessColumn(headers []string, keys ...string) string {
	lower := make([]string, len(headers))
	for i, h := range headers {
		lower[i] = strings.ToLower(h)
	}
	for _, key := range keys {
		for i := range headers {
			if strings.Contains(lower[i], strings.ToLower(key)) {
				return headers[i]
			}
		}
	}
	return ""
}

func inferDefaults(headers []string) model.MappingDefaults {
	return model.MappingDefaults{
		TaskCol:          guessColumn(headers, "task", "任务", "name"),
		StartCol:         guessColumn(headers, "start", "开始", "date"),
		EndCol:           guessColumn(headers, "end", "结束", "date"),
		ProjectCol:       guessColumn(headers, "project", "项目", "group"),
		ColorCol:         guessColumn(headers, "project", "分类", "group", "color"),
		DescCol:          guessColumn(headers, "desc", "description", "detail", "说明"),
		MilestoneCol:     guessColumn(headers, "milestone", "里程碑"),
		MilestoneDateCol: guessColumn(headers, "milestone date", "milestonedate", "里程碑日期"),
		PlanStartCol:     guessColumn(headers, "planstart", "计划开始", "baseline start"),
		PlanEndCol:       guessColumn(headers, "planend", "计划结束", "baseline end"),
		OwnerCol:         guessColumn(headers, "owner", "负责人", "assignee"),
	}
}
