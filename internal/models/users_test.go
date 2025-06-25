package models

import (
	"testing"

	"github.com/mohafarman/snippetbox/internal/assert"
)

func TestExists(t *testing.T) {
	tests := []struct {
		name   string
		userId int
		want   bool
	}{
		{
			name:   "User exists",
			userId: 1,
			want:   true,
		},
		{
			name:   "Zero ID",
			userId: 0,
			want:   false,
		},
		{
			name:   "User does not exist",
			userId: 2,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)
			m := UserModel{db}

			exists, err := m.Exists(tt.userId)

			assert.NilError(t, err)
			assert.Equal(t, exists, tt.want)
		})
	}

}
