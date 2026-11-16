package infrastructure

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"pos/application"
	"testing"
)

func TestExportRoundTrip(t *testing.T) {
	rows := []application.GradeRow{{StudentName: "Иван; \"Иванов\"\nСтудент", Grade: 8}}
	dir := t.TempDir()
	for _, ext := range []string{"csv", "json"} {
		path := filepath.Join(dir, "grades."+ext)
		if err := (ExtensionExporter{}).Export(rows, path); err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if ext == "csv" {
			reader := csv.NewReader(bytes.NewReader(data))
			reader.Comma = ';'
			records, err := reader.ReadAll()
			if err != nil || len(records) != 2 || records[1][0] != rows[0].StudentName || records[1][1] != "8.0" {
				t.Fatalf("csv: %v, %v", records, err)
			}
		} else {
			var records []struct {
				Student string  `json:"студент"`
				Grade   float64 `json:"итог"`
			}
			if err := json.Unmarshal(data, &records); err != nil {
				t.Fatal(err)
			}
			if len(records) != 1 || records[0].Student != rows[0].StudentName || records[0].Grade != 8 {
				t.Fatalf("json: %+v", records)
			}
		}
	}
	if err := (ExtensionExporter{}).Export(rows, filepath.Join(dir, "bad.txt")); err == nil {
		t.Fatal("unknown format accepted")
	}
	if err := (NewFileExporter(JSONFormatter{})).Export([]application.GradeRow{{Grade: math.NaN()}}, filepath.Join(dir, "bad.json")); err == nil {
		t.Fatal("invalid json number accepted")
	}
	data, err := (JSONFormatter{}).Format(nil)
	if err != nil || string(data) != "[]" {
		t.Fatalf("empty json: %s, %v", data, err)
	}
}
