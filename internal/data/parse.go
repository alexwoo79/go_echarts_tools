// Package data handles dataset storage and file parsing.
package data

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"

	"gantt/internal/model"
)

// ParseCSV parses a CSV reader into a Dataset.
func ParseCSV(name string, reader io.Reader) (model.Dataset, error) {
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1

	records, err := csvReader.ReadAll()
	if err != nil {
		return model.Dataset{}, fmt.Errorf("CSV 读取失败: %w", err)
	}
	if len(records) < 2 {
		return model.Dataset{}, fmt.Errorf("CSV 至少需要一行表头和一行数据")
	}

	headers := make([]string, 0, len(records[0]))
	for i, col := range records[0] {
		v := strings.TrimSpace(col)
		if v == "" {
			v = fmt.Sprintf("Column_%d", i+1)
		}
		headers = append(headers, v)
	}

	rows := make([][]string, 0, len(records)-1)
	for _, row := range records[1:] {
		normalized := make([]string, len(headers))
		nonEmpty := false
		for i := range headers {
			if i < len(row) {
				normalized[i] = strings.TrimSpace(row[i])
				if normalized[i] != "" {
					nonEmpty = true
				}
			}
		}
		if nonEmpty {
			rows = append(rows, normalized)
		}
	}

	if len(rows) == 0 {
		return model.Dataset{}, fmt.Errorf("上传文件没有可用数据行")
	}

	return model.Dataset{ID: NewID(), Name: name, Headers: headers, Rows: rows}, nil
}

// ParseXLSX parses an Excel reader into a Dataset.
func ParseXLSX(name string, reader io.Reader) (model.Dataset, error) {
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return model.Dataset{}, fmt.Errorf("Excel 解析失败: %w", err)
	}
	defer func() { _ = f.Close() }()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return model.Dataset{}, fmt.Errorf("Excel 中没有工作表")
	}
	rows, err := f.GetRows(sheets[0])
	if err != nil || len(rows) < 2 {
		return model.Dataset{}, fmt.Errorf("Excel 数据不足，至少需要一行表头和一行数据")
	}

	headers := make([]string, len(rows[0]))
	for i, col := range rows[0] {
		v := strings.TrimSpace(col)
		if v == "" {
			v = fmt.Sprintf("Column_%d", i+1)
		}
		headers[i] = v
	}

	dataRows := make([][]string, 0, len(rows)-1)
	for _, row := range rows[1:] {
		normalized := make([]string, len(headers))
		nonEmpty := false
		for i := range headers {
			if i < len(row) {
				normalized[i] = strings.TrimSpace(row[i])
				if normalized[i] != "" {
					nonEmpty = true
				}
			}
		}
		if nonEmpty {
			dataRows = append(dataRows, normalized)
		}
	}

	if len(dataRows) == 0 {
		return model.Dataset{}, fmt.Errorf("Excel 没有可用数据行")
	}

	return model.Dataset{ID: NewID(), Name: name, Headers: headers, Rows: dataRows}, nil
}

// ParseUploadedFile auto-detects file type and parses to a Dataset.
func ParseUploadedFile(fileHeader *multipart.FileHeader) (model.Dataset, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return model.Dataset{}, fmt.Errorf("无法打开上传文件")
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	switch ext {
	case ".csv":
		return ParseCSV(fileHeader.Filename, file)
	case ".xlsx", ".xlsm", ".xltx", ".xltm":
		raw, err := io.ReadAll(file)
		if err != nil {
			return model.Dataset{}, fmt.Errorf("读取 Excel 文件失败")
		}
		return ParseXLSX(fileHeader.Filename, bytes.NewReader(raw))
	default:
		return model.Dataset{}, fmt.Errorf("仅支持 CSV 或 XLSX 文件")
	}
}

// ParseDate parses a date string in many common formats.
func ParseDate(value string) (time.Time, error) {
	v := strings.TrimSpace(value)
	if v == "" {
		return time.Time{}, fmt.Errorf("empty date")
	}

	layouts := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"01-02-06",
		"1-2-06",
		"01/02/06",
		"1/2/06",
		"01.02.06",
		"1.2.06",
		"01-02-2006",
		"1-2-2006",
		"01/02/2006",
		"1/2/2006",
		"01-02-06 15:04",
		"1-2-06 15:04",
		"01-02-06 15:04:05",
		"1-2-06 15:04:05",
		"2006-01-02",
		"2006-1-2",
		"2006/01/02",
		"2006/1/2",
		"2006.01.02",
		"2006.1.2",
		"2006年01月02日",
		"2006年1月2日",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-1-2 15:04:05",
		"2006-1-2 15:04",
		"2006/01/02 15:04:05",
		"2006/01/02 15:04",
		"2006/1/2 15:04:05",
		"2006/1/2 15:04",
		"1/2/2006",
		"1/2/2006 15:04",
		"1/2/2006 15:04:05",
	}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, v, time.Local); err == nil {
			return t, nil
		}
	}

	// Excel serial date (e.g. 45291 or 45291.5).
	if serial, err := strconv.ParseFloat(v, 64); err == nil {
		if t, err := excelize.ExcelDateToTime(serial, false); err == nil {
			return t, nil
		}
		if t, err := excelize.ExcelDateToTime(serial, true); err == nil {
			return t, nil
		}
	}

	if unixMS, err := strconv.ParseInt(v, 10, 64); err == nil && len(v) >= 12 {
		return time.UnixMilli(unixMS), nil
	}
	return time.Time{}, fmt.Errorf("unsupported date format: %s", value)
}

// Cell safely retrieves a cell value from a row by index.
func Cell(row []string, idx int) string {
	if idx < 0 || idx >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[idx])
}
