package web

import "testing"

func Test_getBookIDFromURL(t *testing.T) {
	type args struct {
		bookID string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Plain bookID without url",
			args: args{
				bookID: "136267",
			},
			want:    "136267",
			wantErr: false,
		},
		{
			name: "UID inside URL",
			args: args{
				bookID: "https://ttkjf.dawinci.sk/?fn=*recview&uid=136267&pageId=resultform&full=0&focusName=bsktchRZ1",
			},
			want:    "136267",
			wantErr: false,
		},
		{
			name: "UID at the end of URL",
			args: args{
				bookID: "https://ttkjf.dawinci.sk/?fn=*recview&pageId=resultform&full=0&focusName=bsktchRZ1&uid=136267",
			},
			want:    "136267",
			wantErr: false,
		},
		{
			name: "Missing UID in URL",
			args: args{
				bookID: "https://ttkjf.dawinci.sk/?fn=*recview&pageId=resultform&full=0&focusName=bsktchRZ1",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getBookIDFromURL(tt.args.bookID)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBookIDFromURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getBookIDFromURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
