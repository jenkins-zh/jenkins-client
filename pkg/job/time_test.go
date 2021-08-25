package job

import (
	"reflect"
	"testing"
	"time"
)

func TestTime_UnmarshalJSON(t1 *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{{
		name: "normal Jenkins time",
		args: args{
			data: []byte(`"2021-08-23T17:03:04.144+0000"`),
		},
		want: time.Date(2021, 8, 23, 17, 3, 4, 144000000, time.UTC).Local(),
	}, {
		name: "null time",
		args: args{
			data: []byte(`null`),
		},
		want: time.Time{},
	}, {
		name: "invalid Jenkins time",
		args: args{
			data: []byte(`"2021-08-23"`),
		},
		wantErr: true,
	}, {
		name: "invalid JSON",
		args: args{
			data: []byte(`invalid`),
		},
		wantErr: true,
	},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t *testing.T) {
			got := &Time{}
			if err := got.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got.Time, tt.want) {
				t.Errorf("UmarshalJSON() got = %v, want = %v", got.Time, tt.want)
			}
		})
	}
}

func TestTime_IsZero(t1 *testing.T) {
	tests := []struct {
		name string
		time *Time
		want bool
	}{{
		name: "nil time",
		time: nil,
		want: true,
	}, {
		name: "zero time",
		time: &Time{Time: time.Time{}},
		want: true,
	}, {
		name: "normal time",
		time: &Time{Time: time.Now()},
		want: false,
	},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			if got := tt.time.IsZero(); got != tt.want {
				t1.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}
