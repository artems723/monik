package domain

import "testing"

func TestMetric_Validate(t *testing.T) {
	type fields struct {
		ID    string
		MType MetricType
		Delta *int64
		Value *float64
		Hash  string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "no value gauge",
			fields: fields{
				ID:    "test",
				MType: MetricTypeGauge,
				Delta: nil,
				Value: nil,
				Hash:  "",
			},
			args: args{
				key: "test",
			},
			wantErr: true,
		},
		{
			name: "no value counter",
			fields: fields{
				ID:    "test",
				MType: MetricTypeCounter,
				Delta: nil,
				Value: nil,
				Hash:  "",
			},
			args: args{
				key: "test",
			},
			wantErr: true,
		},
		{
			name: "test hash",
			fields: fields{
				ID:    "test",
				MType: MetricTypeCounter,
				Delta: int64Pointer(2),
				Value: nil,
				Hash:  "989f2d948235ebb42de49d06abefd188f94b8989a447cb3b1dabe354b46d040f",
			},
			args: args{
				key: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metric{
				ID:    tt.fields.ID,
				MType: tt.fields.MType,
				Delta: tt.fields.Delta,
				Value: tt.fields.Value,
				Hash:  tt.fields.Hash,
			}
			if err := m.Validate(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func int64Pointer(i int64) *int64 {
	return &i
}
