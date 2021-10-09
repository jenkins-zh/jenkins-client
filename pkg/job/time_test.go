package job

import (
	"encoding/json"
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
		name: "RFC3339 time",
		args: args{
			data: []byte(`"2021-10-08T09:08:21.523Z"`),
		},
		want: time.Date(2021, 10, 8, 9, 8, 21, 523000000, time.UTC).Local(),
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

func TestTime_Marshaler(t *testing.T) {
	tests := []struct {
		name    string
		time    string
		wantErr bool
	}{{
		name: "Jenkins time layout",
		time: `"2021-08-23T17:03:04.144+0000"`,
	}, {
		name: "RFC3339 layout",
		time: `"2021-10-08T09:08:21.523Z"`,
	}, {
		name:    "Invalid time",
		time:    `"2021-10-09"`,
		wantErr: true,
	}, {
		name: "Null time",
		time: `null`,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time := &Time{}
			var err error
			if err = json.Unmarshal([]byte(tt.time), time); (err != nil) != tt.wantErr {
				t.Errorf("Failed to unmarshal the time: %v, and error = %v", tt.time, err)
			}
			if err != nil {
				// skip if error was occurred before
				return
			}
			var timeBytes []byte
			timeBytes, err = json.Marshal(time)
			if (err != nil) != tt.wantErr {
				t.Errorf("Failed to marshal the time: %v, and error = %v", time, err)
			}
			if err != nil {
				// skip if error was occurred before
				return
			}
			if err := json.Unmarshal(timeBytes, time); (err != nil) != tt.wantErr {
				t.Errorf("Failed to unmarshal the time: %v, and error = %v", string(timeBytes), err)
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

func TestTime_MarshalJSON(t *testing.T) {
	gmtZone, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Errorf("unable to prepare GMT zone, error = %v", err)
	}

	type fields struct {
		Time time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{{
		name: "UTC time",
		fields: fields{
			Time: time.Date(2021, 10, 9, 15, 23, 59, 123456789, time.UTC),
		},
		want: []byte(`"2021-10-09T15:23:59Z"`),
	}, {
		name: "GTM time",
		fields: fields{
			Time: time.Date(2021, 10, 9, 15, 23, 59, 123456789, gmtZone),
		},
		want: []byte(`"2021-10-09T07:23:59Z"`),
	}, {
		name:   "Nil time",
		fields: fields{},
		want:   []byte("null"),
	}, {
		name: "Zero time",
		fields: fields{
			Time: time.Time{},
		},
		want: []byte("null"),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Time{
				Time: tt.fields.Time,
			}
			got, err := tr.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Time.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Time.MarshalJSON() = %s, want %s", got, tt.want)
			}
		})
	}
}
