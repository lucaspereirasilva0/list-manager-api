package domain_test

import (
	"testing"

	"github.com/lucaspereirasilva0/list-manager-api/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestItem_IsActive(t *testing.T) {
	tests := []struct {
		name        string
		givenName   string
		givenStatus bool
		wantStatus  bool
	}{
		{
			name:        "Give_NameAndStatus_When_IsActivate_Then_ExpectedActiveItem",
			givenName:   "any item",
			givenStatus: true,
			wantStatus:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			givenItem := domain.NewItem(tt.givenName, tt.givenStatus)
			gotStatus := givenItem.IsActive()
			require.Equal(t, tt.wantStatus, gotStatus)
		})
	}
}
