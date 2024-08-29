package components

import "testing"

func Test_trimDateFromStatus(t *testing.T) {
	type args struct {
		status string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test 1",
			args: args{
				status: "Voľný",
			},
			want: "Voľný",
		},
		{
			name: "Test 2",
			args: args{
				status: "Požičaný do 23.09.2024",
			},
			want: "Požičaný",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trimDateFromStatus(tt.args.status); got != tt.want {
				t.Errorf("trimDateFromStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
