package tracer

import (
	"context"
	"testing"
)

func TestExtractTraceInfo(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
		want2 bool
	}{
		{ "ExtractTraceInfo",
			args{
			nil,
			},
		"",
		"",
		false},
	}
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := ExtractTraceInfo(tt.args.ctx)
			if got != tt.want {
				t.Errorf("ExtractTraceInfo() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ExtractTraceInfo() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("ExtractTraceInfo() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
