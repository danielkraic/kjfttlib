package components

import (
	"testing"

	"github.com/danielkraic/kjfttlib/pkg/book"
)

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

func Test_getInstanceTooltipText(t *testing.T) {
	testBook := &book.Model{
		ID:     "test-book",
		Title:  "Test Book",
		Author: "Test Author",
		Instances: []*book.Instance{
			{Location: "Library A", Status: "Voľný"},
			{Location: "Library B", Status: "Voľný"},
			{Location: "Library C", Status: "Požičaný do 23.09.2024"},
		},
	}

	tests := []struct {
		name         string
		book         *book.Model
		targetStatus string
		want         string
	}{
		{
			name:         "Available instances",
			book:         testBook,
			targetStatus: "Voľný",
			want:         "📚 Voľný (2):\n• Library A: Voľný\n• Library B: Voľný",
		},
		{
			name:         "Borrowed instances",
			book:         testBook,
			targetStatus: "Požičaný",
			want:         "📚 Požičaný (1):\n• Library C: Požičaný do 23.09.2024",
		},
		{
			name:         "No instances",
			book:         testBook,
			targetStatus: "Non-existent",
			want:         "No instances found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getInstanceTooltipText(tt.book, tt.targetStatus); got != tt.want {
				t.Errorf("getInstanceTooltipText() = %v, want %v", got, tt.want)
			}
		})
	}
}
