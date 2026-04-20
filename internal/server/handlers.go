// Package server contains HTTP request handlers.
package server

import (
	"io/fs"
	"strings"

	"github.com/gin-gonic/gin"

	"gantt/internal/charts"
	chartgantt "gantt/internal/charts/gantt"
	"gantt/internal/data"
	"gantt/internal/model"
)

const demoCSV = `Project,Task,StartDate,EndDate,Description,Milestone,MilestoneDate,PlanStartDate,PlanEndDate,Owner
Website,Requirements,2026-04-01,2026-04-05,Define scope,Scope Locked,2026-04-05
Website,Design,2026-04-06,2026-04-12,UI and UX design,Design Review,2026-04-11,2026-04-06,2026-04-10,Alice
Website,Frontend,2026-04-13,2026-04-24,Build pages,,,2026-04-12,2026-04-22,Bob
Website,Backend,2026-04-14,2026-04-26,Build APIs,,,2026-04-13,2026-04-23,Chris
Website,Integration,2026-04-27,2026-05-02,Integrate FE and BE,Integration Complete,2026-05-02,2026-04-25,2026-05-01,Diana
Website,Launch,2026-05-03,2026-05-05,Release to production,Go Live,2026-05-05,2026-05-02,2026-05-04,Eric
`

type handlers struct {
	assets fs.FS
}

func (h *handlers) home(c *gin.Context) {
	c.HTML(200, "index.tmpl", gin.H{"Title": "Gantt - Go + Gin + ECharts"})
}

func (h *handlers) demo(c *gin.Context) {
	dataset, err := data.ParseCSV("demo.csv", strings.NewReader(demoCSV))
	if err != nil {
		c.HTML(200, "index.tmpl", gin.H{"Title": "Gantt - Go + Gin + ECharts", "Error": err.Error()})
		return
	}
	data.Store(dataset)
	renderMapper(c, dataset, "")
}

func (h *handlers) clear(c *gin.Context) {
	c.HTML(200, "index.tmpl", gin.H{"Title": "Gantt - Go + Gin + ECharts"})
}

func (h *handlers) upload(c *gin.Context) {
	fileHeader, err := c.FormFile("data_file")
	if err != nil {
		c.HTML(200, "index.tmpl", gin.H{"Title": "Gantt - Go + Gin + ECharts", "Error": "请选择要上传的 CSV 或 XLSX 文件。"})
		return
	}
	dataset, err := data.ParseUploadedFile(fileHeader)
	if err != nil {
		c.HTML(200, "index.tmpl", gin.H{"Title": "Gantt - Go + Gin + ECharts", "Error": err.Error()})
		return
	}
	data.Store(dataset)
	renderMapper(c, dataset, "")
}

func (h *handlers) chart(c *gin.Context) {
	datasetID := c.PostForm("dataset_id")
	dataset, ok := data.Load(datasetID)
	if !ok {
		c.HTML(200, "index.tmpl", gin.H{"Title": "Gantt - Go + Gin + ECharts", "Error": "数据已过期，请重新上传文件。"})
		return
	}

	cfg := model.MappingConfig{
		TaskCol:          c.PostForm("task_col"),
		StartCol:         c.PostForm("start_col"),
		EndCol:           c.PostForm("end_col"),
		ProjectCol:       c.PostForm("project_col"),
		ColorCol:         c.PostForm("color_col"),
		DescCol:          c.PostForm("desc_col"),
		MilestoneCol:     c.PostForm("milestone_col"),
		MilestoneDateCol: c.PostForm("milestone_date_col"),
		PlanStartCol:     c.PostForm("plan_start_col"),
		PlanEndCol:       c.PostForm("plan_end_col"),
		OwnerCol:         c.PostForm("owner_col"),
		SortByStart:      c.PostForm("sort_by_start") == "on",
		ShowTaskNumber:   c.PostForm("show_task_number") == "on",
	}

	opts := model.ChartOptions{
		HierarchicalView: true,
		ShowTaskDetails:  c.PostForm("show_task_details") == "on",
		ShowDuration:     c.PostForm("show_duration") == "on",
		DarkTheme:        c.PostForm("dark_theme") == "on",
		ChartTheme:       strings.TrimSpace(c.PostForm("chart_theme")),
		TimeGranularity:  strings.TrimSpace(c.PostForm("time_granularity")),
	}
	if opts.ChartTheme == "" {
		opts.ChartTheme = "default"
	}
	if opts.TimeGranularity == "" {
		opts.TimeGranularity = "month"
	}

	// chart_type defaults to "gantt"; extensible via the registry.
	chartType := c.PostForm("chart_type")
	if chartType == "" {
		chartType = "gantt"
	}
	builder, ok := charts.Get(chartType)
	if !ok {
		renderMapper(c, dataset, "未知图表类型: "+chartType)
		return
	}

	result, err := builder.Build(dataset, cfg, opts)
	if err != nil {
		renderMapper(c, dataset, err.Error())
		return
	}

	ganttResult, ok := result.(chartgantt.Result)
	if !ok {
		renderMapper(c, dataset, "图表数据类型错误")
		return
	}

	renderWorkspace(c, dataset, cfg, opts, ganttResult.Tasks, ganttResult.Stats, "")
}
