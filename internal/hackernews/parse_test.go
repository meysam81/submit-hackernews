package hackernews

import (
	"errors"
	"strings"
	"testing"
)

func TestParseHiddenInputs(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		names   []string
		want    map[string]string
		wantErr error
	}{
		{
			name: "extracts requested fields",
			html: `<form>
				<input type="hidden" name="fnid" value="abc123">
				<input type="hidden" name="fnop" value="submit-page">
				<input type="text" name="title" value="">
			</form>`,
			names: []string{"fnid", "fnop"},
			want:  map[string]string{"fnid": "abc123", "fnop": "submit-page"},
		},
		{
			name:    "missing field errors",
			html:    `<form><input type="hidden" name="fnid" value="abc123"></form>`,
			names:   []string{"fnid", "fnop"},
			wantErr: ErrFormFieldMissing,
		},
		{
			name:  "empty value is still found",
			html:  `<input name="fnid" value=""><input name="fnop" value="x">`,
			names: []string{"fnid", "fnop"},
			want:  map[string]string{"fnid": "", "fnop": "x"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseHiddenInputs(strings.NewReader(tt.html), tt.names...)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("err = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %d fields, want %d", len(got), len(tt.want))
			}
			for k, v := range tt.want {
				if got[k] != v {
					t.Errorf("field %q = %q, want %q", k, got[k], v)
				}
			}
		})
	}
}
