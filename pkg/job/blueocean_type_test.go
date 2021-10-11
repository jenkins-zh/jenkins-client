package job

import (
	"testing"
)

func TestCause_getShortDescription(t *testing.T) {
	tests := []struct {
		name  string
		cause Cause
		want  string
	}{{
		name:  "Nil cause",
		cause: nil,
		want:  "",
	}, {
		name: "Nil value for short description",
		cause: Cause{
			"shortDescription": nil,
		},
		want: "",
	}, {
		name: "Non-nil value for short description",
		cause: Cause{
			"shortDescription": "testDesc",
		},
		want: "testDesc",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cause.GetShortDescription(); got != tt.want {
				t.Errorf("Cause.getShortDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}
