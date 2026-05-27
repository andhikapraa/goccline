// Package input parses the JSON payload that Claude Code pipes to the
// statusLine command on every refresh.
package input

import (
	"encoding/json"
	"io"
)

// Payload mirrors the JSON Claude Code sends. Fields are best-effort;
// missing pieces are tolerated.
type Payload struct {
	SessionID      string    `json:"session_id"`
	TranscriptPath string    `json:"transcript_path"`
	CWD            string    `json:"cwd"`
	Model          Model     `json:"model"`
	Workspace      Workspace `json:"workspace"`
	Cost           Cost      `json:"cost"`
	Version        string    `json:"version"`
	OutputStyle    Style     `json:"output_style"`
}

type Model struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

type Workspace struct {
	CurrentDir string `json:"current_dir"`
	ProjectDir string `json:"project_dir"`
}

type Cost struct {
	TotalCostUSD     float64 `json:"total_cost_usd"`
	TotalDurationMS  int64   `json:"total_duration_ms"`
	TotalAPIDuration int64   `json:"total_api_duration_ms"`
	TotalLinesAdded  int     `json:"total_lines_added"`
	TotalLinesRemvd  int     `json:"total_lines_removed"`
}

type Style struct {
	Name string `json:"name"`
}

// Parse reads a JSON payload from r. Empty payloads return an empty Payload
// rather than an error so the statusline can still render fallback content.
func Parse(r io.Reader) (Payload, error) {
	var p Payload
	dec := json.NewDecoder(r)
	if err := dec.Decode(&p); err != nil {
		if err == io.EOF {
			return p, nil
		}
		return p, err
	}
	// Backfill CWD from workspace if Claude Code sends both.
	if p.CWD == "" && p.Workspace.CurrentDir != "" {
		p.CWD = p.Workspace.CurrentDir
	}
	return p, nil
}
