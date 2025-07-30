package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRootCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		wantOut []string
	}{
		{
			name:    "help flag",
			args:    []string{"--help"},
			wantErr: false,
			wantOut: []string{"gh-inspector", "Usage:", "Available Commands:"},
		},
		{
			name:    "no args shows help",
			args:    []string{},
			wantErr: false,
			wantOut: []string{"gh-inspector", "Usage:"},
		},
		{
			name:    "invalid command",
			args:    []string{"invalid"},
			wantErr: true,
			wantOut: []string{"unknown command"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)
			rootCmd.SetArgs(tt.args)

			err := rootCmd.Execute()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			output := buf.String()
			for _, want := range tt.wantOut {
				require.Contains(t, output, want, "Output missing expected content")
			}
		})
	}
}
