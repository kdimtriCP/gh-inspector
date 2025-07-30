package formatter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	
	"github.com/kdimtriCP/gh-inspector/internal/metrics"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{name: "table format", format: FormatTable, wantErr: false},
		{name: "json format", format: FormatJSON, wantErr: false},
		{name: "json-compact format", format: FormatJSONCompact, wantErr: false},
		{name: "csv format", format: FormatCSV, wantErr: false},
		{name: "invalid format", format: "invalid", wantErr: true},
		{name: "empty format", format: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter, err := New(tt.format)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, formatter)
			} else {
				require.NoError(t, err)
				require.NotNil(t, formatter)
			}
		})
	}
}

func TestTableFormatter(t *testing.T) {
	repos := []*metrics.Repository{
		{
			Owner:           "test",
			Name:            "repo1",
			Stars:           100,
			Forks:           10,
			OpenIssues:      5,
			OpenPRs:         2,
			LastCommitDate:  time.Now().AddDate(0, 0, -5),
			Score:           85.5,
			PrimaryLanguage: "Go",
			HasLicense:      true,
			HasCICD:         true,
			Description:     "Test repository 1",
			IsArchived:      false,
		},
		{
			Owner:          "test",
			Name:           "repo2",
			Stars:          50,
			Forks:          5,
			OpenIssues:     10,
			OpenPRs:        5,
			LastCommitDate: time.Now().AddDate(0, -1, 0),
			Score:          65.0,
			Description:    "Test repository 2",
			IsArchived:     true,
		},
	}

	formatter := &TableFormatter{}
	buf := &bytes.Buffer{}
	err := formatter.Format(buf, repos)
	require.NoError(t, err)

	output := buf.String()
	
	expectedContents := []string{
		"REPOSITORY",
		"SCORE",
		"STARS",
		"FORKS",
		"test/repo1",
		"test/repo2",
		"85.5",
		"65.0",
		"100",
		"50",
		"Yes",
		"No",
	}

	for _, expected := range expectedContents {
		require.Contains(t, output, expected, "Table output missing expected content")
	}
}

func TestJSONFormatter(t *testing.T) {
	repos := []*metrics.Repository{
		{
			Owner:      "test",
			Name:       "repo",
			Stars:      100,
			Score:      75.0,
			HasLicense: true,
		},
	}

	t.Run("pretty JSON", func(t *testing.T) {
		formatter := &JSONFormatter{indent: true}
		buf := &bytes.Buffer{}
		err := formatter.Format(buf, repos)
		require.NoError(t, err)

		var result []map[string]interface{}
		err = json.Unmarshal(buf.Bytes(), &result)
		require.NoError(t, err)

		require.Len(t, result, 1, "Expected 1 repository")

		output := buf.String()
		require.Contains(t, output, "\n", "JSON output should have newlines")
		require.Contains(t, output, "  ", "JSON output should be indented")
	})

	t.Run("compact JSON", func(t *testing.T) {
		formatter := &JSONFormatter{indent: false}
		buf := &bytes.Buffer{}
		err := formatter.Format(buf, repos)
		require.NoError(t, err)

		var result []map[string]interface{}
		err = json.Unmarshal(buf.Bytes(), &result)
		require.NoError(t, err)

		output := strings.TrimSpace(buf.String())
		lines := strings.Split(output, "\n")
		require.Len(t, lines, 1, "Compact JSON should be on a single line")
	})
}

func TestCSVFormatter(t *testing.T) {
	repos := []*metrics.Repository{
		{
			Owner:           "test",
			Name:            "repo1",
			Stars:           100,
			Forks:           10,
			OpenIssues:      5,
			OpenPRs:         2,
			Score:           85.5,
			PrimaryLanguage: "Go",
			HasLicense:      true,
			HasCICD:         true,
			HasContributing: false,
			Description:     "Test repository with, comma",
			IsArchived:      false,
		},
	}

	formatter := &CSVFormatter{}
	buf := &bytes.Buffer{}
	err := formatter.Format(buf, repos)
	require.NoError(t, err)

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	
	require.GreaterOrEqual(t, len(lines), 2, "CSV output should have at least header and one data row")

	header := lines[0]
	expectedHeaders := []string{
		"repository",
		"score",
		"stars",
		"forks",
		"open_issues",
		"open_prs",
		"last_commit",
		"language",
		"ci_cd",
		"license",
		"description",
		"archived",
	}

	for _, h := range expectedHeaders {
		require.Contains(t, header, h, "CSV header missing expected column")
	}

	dataRow := lines[1]
	expectedData := []string{
		"test/repo1",
		"85.5",
		"100",
		"10",
		"5",
		"2",
		"Go",
		"Yes",
		"Yes",
		"No",
		`"Test repository with, comma"`,
	}

	for _, d := range expectedData {
		require.Contains(t, dataRow, d, "CSV data row missing expected value")
	}
}