package models

import "testing"

func TestStringToUpdatePriority(t *testing.T) {
	type args struct {
		prio string
	}
	tests := []struct {
		name    string
		args    args
		want    UpdatePriority
		wantErr bool
	}{
		{
			name:    "Valid update priority",
			args:    args{prio: "URGENT"},
			want:    URGENT,
			wantErr: false,
		},
		{
			name:    "Invalid update priority",
			args:    args{prio: "invalid"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StringToUpdatePriority(tt.args.prio)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringToUpdatePriority() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StringToUpdatePriority() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringToFeedType(t *testing.T) {
	type args struct {
		feedType string
	}
	tests := []struct {
		name    string
		args    args
		want    FeedType
		wantErr bool
	}{
		{
			name:    "Valid feed type",
			args:    args{feedType: "rss"},
			want:    RSS,
			wantErr: false,
		},
		{
			name:    "Invalid feed type",
			args:    args{feedType: "invalid"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StringToFeedType(tt.args.feedType)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringToFeedType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StringToFeedType() = %v, want %v", got, tt.want)
			}
		})
	}
}
