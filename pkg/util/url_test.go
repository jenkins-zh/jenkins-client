package util

import "testing"

func TestURLJoinAsString(t *testing.T) {
	type args struct {
		host string
		api  string
	}
	tests := []struct {
		name             string
		args             args
		wantTargetURLStr string
		wantErr          bool
	}{{
		name: "Path ending with slash",
		args: args{
			host: "https://devops/",
			api:  "path/",
		},
		wantTargetURLStr: "https://devops/path/",
	}, {
		name: "Path that dose not end with slash",
		args: args{
			host: "https://devops/",
			api:  "path",
		},
		wantTargetURLStr: "https://devops/path",
	}, {
		name: "Host that dose not end with slash",
		args: args{
			host: "https://devops",
			api:  "path",
		},
		wantTargetURLStr: "https://devops/path",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTargetURLStr, err := URLJoinAsString(tt.args.host, tt.args.api)
			if (err != nil) != tt.wantErr {
				t.Errorf("URLJoinAsString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTargetURLStr != tt.wantTargetURLStr {
				t.Errorf("URLJoinAsString() gotTargetURLStr = %v, want %v", gotTargetURLStr, tt.wantTargetURLStr)
			}
		})
	}
}
