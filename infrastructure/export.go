package infrastructure

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"pos/application"
	"strings"
)

type GradeFormatter interface {
	Format([]application.GradeRow) ([]byte, error)
}

type FileExporter struct{ gradeFormatter GradeFormatter }

func NewFileExporter(formatter GradeFormatter) FileExporter {
	return FileExporter{gradeFormatter: formatter}
}
func (fe FileExporter) Export(rows []application.GradeRow, path string) error {
	data, err := fe.gradeFormatter.Format(rows)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

type CSVFormatter struct{}

func (CSVFormatter) Format(rows []application.GradeRow) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Comma = ';'
	if err := writer.Write([]string{"студент", "итог"}); err != nil {
		return nil, err
	}
	for _, row := range rows {
		if err := writer.Write([]string{row.StudentName, fmt.Sprintf("%.1f", row.Grade)}); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	return buf.Bytes(), writer.Error()
}

type JSONFormatter struct{}

func (JSONFormatter) Format(rows []application.GradeRow) ([]byte, error) {
	type entry struct {
		Student string  `json:"студент"`
		Grade   float64 `json:"итог"`
	}
	entries := make([]entry, 0, len(rows))
	for _, row := range rows {
		entries = append(entries, entry{row.StudentName, row.Grade})
	}
	return json.MarshalIndent(entries, "", "  ")
}

type ExtensionExporter struct{}

func (ExtensionExporter) Export(rows []application.GradeRow, path string) error {
	var formatter GradeFormatter
	switch strings.ToLower(filepath.Ext(path)) {
	case ".csv":
		formatter = CSVFormatter{}
	case ".json":
		formatter = JSONFormatter{}
	default:
		return fmt.Errorf("укажите расширение .csv или .json")
	}
	return NewFileExporter(formatter).Export(rows, path)
}

var _ application.GradeExporter = FileExporter{}
var _ application.GradeExporter = ExtensionExporter{}
