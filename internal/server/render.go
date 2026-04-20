// Package server contains rendering helpers.
package server

import (
	"encoding/json"
	"html/template"

	"github.com/gin-gonic/gin"

	"gantt/internal/charts"
	"gantt/internal/model"
)

// renderWorkspace renders the full workspace page with an optional chart.
func renderWorkspace(
	c *gin.Context,
	dataset model.Dataset,
	mapping model.MappingConfig,
	options model.ChartOptions,
	tasks []model.Task,
	stats model.Stats,
	errMsg string,
) {
	preview := dataset.Rows
	if len(preview) > 8 {
		preview = preview[:8]
	}

	view := gin.H{
		"Title":      "Gantt - 列映射与图表",
		"DatasetID":  dataset.ID,
		"FileName":   dataset.Name,
		"Headers":    dataset.Headers,
		"Preview":    preview,
		"Mapping":    mapping,
		"Options":    options,
		"Error":      errMsg,
		"HasChart":   len(tasks) > 0,
		"ChartTypes": charts.All(),
	}

	if len(tasks) > 0 {
		tasksJSON, _ := json.Marshal(tasks)
		statsJSON, _ := json.Marshal(stats)
		optionsJSON, _ := json.Marshal(options)
		view["TasksJSON"] = template.JS(string(tasksJSON))
		view["StatsJSON"] = template.JS(string(statsJSON))
		view["OptionsJSON"] = template.JS(string(optionsJSON))
	}

	c.HTML(200, "index.tmpl", view)
}

// renderMapper renders the workspace page in column-mapping mode
// (no chart data), with auto-inferred column defaults.
func renderMapper(c *gin.Context, dataset model.Dataset, errMsg string) {
	builder, ok := charts.Get("gantt")
	if !ok {
		c.HTML(200, "index.tmpl", gin.H{"Error": "gantt builder not registered"})
		return
	}
	cfg := builder.InferDefaults(dataset.Headers)
	opts := builder.DefaultOptions()
	renderWorkspace(c, dataset, cfg, opts, nil, model.Stats{}, errMsg)
}
