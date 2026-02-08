package domain_test

import (
	"testing"

	"github.com/lucaspereirasilva0/list-manager-api/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestItem_IsActive(t *testing.T) {
	tests := []struct {
		name       string
		givenItem  domain.Item
		wantStatus bool
	}{
		{
			name:       "Given_ActiveItem_When_IsActive_Then_ReturnsTrue",
			givenItem:  mockActiveItem(),
			wantStatus: true,
		},
		{
			name:       "Given_InactiveItem_When_IsActive_Then_ReturnsFalse",
			givenItem:  mockInactiveItem(),
			wantStatus: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStatus := tt.givenItem.IsActive()
			require.Equal(t, tt.wantStatus, gotStatus)
		})
	}
}

func TestItem_IsEmpty(t *testing.T) {
	tests := []struct {
		name      string
		givenItem domain.Item
		wantEmpty bool
	}{
		{
			name:      "Given_ItemWithID_When_IsEmpty_Then_ReturnsFalse",
			givenItem: mockItemWithID(),
			wantEmpty: false,
		},
		{
			name:      "Given_ItemWithoutID_When_IsEmpty_Then_ReturnsTrue",
			givenItem: mockEmptyItem(),
			wantEmpty: true,
		},
		{
			name:      "Given_ItemWithWhitespaceID_When_IsEmpty_Then_ReturnsFalse",
			givenItem: mockItemWithWhitespaceID(),
			wantEmpty: false, // Espaços em branco não são considerados vazios
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEmpty := tt.givenItem.IsEmpty()
			require.Equal(t, tt.wantEmpty, gotEmpty)
		})
	}
}

func TestNewItem(t *testing.T) {
	tests := []struct {
		name        string
		givenName   string
		givenActive bool
		wantName    string
		wantActive  bool
	}{
		{
			name:        "Given_ValidNameAndActive_When_NewItem_Then_CreatesActiveItem",
			givenName:   "test item",
			givenActive: true,
			wantName:    "test item",
			wantActive:  true,
		},
		{
			name:        "Given_ValidNameAndInactive_When_NewItem_Then_CreatesInactiveItem",
			givenName:   "test item",
			givenActive: false,
			wantName:    "test item",
			wantActive:  false,
		},
		{
			name:        "Given_EmptyName_When_NewItem_Then_CreatesItemWithEmptyName",
			givenName:   "",
			givenActive: true,
			wantName:    "",
			wantActive:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := domain.NewItem(tt.givenName, tt.givenActive)

			// Verifica se o ID foi gerado (não vazio)
			require.NotEmpty(t, item.ID, "Item ID should not be empty")
			require.Len(t, item.ID, 24, "Item ID should be 24 characters (hex of 12 bytes)")

			// Verifica se os outros campos estão corretos
			require.Equal(t, tt.wantName, item.Name)
			require.Equal(t, tt.wantActive, item.Active)
		})
	}
}

func TestNewItem_GeneratesUniqueIDs(t *testing.T) {
	// Testa se IDs gerados são únicos
	item1 := domain.NewItem("item 1", true)
	item2 := domain.NewItem("item 2", true)
	item3 := domain.NewItem("item 3", false)

	require.NotEqual(t, item1.ID, item2.ID, "Generated IDs should be unique")
	require.NotEqual(t, item1.ID, item3.ID, "Generated IDs should be unique")
	require.NotEqual(t, item2.ID, item3.ID, "Generated IDs should be unique")

	// Verifica se todos os IDs têm o formato correto (24 caracteres hex)
	require.Len(t, item1.ID, 24, "Item ID should be 24 characters")
	require.Len(t, item2.ID, 24, "Item ID should be 24 characters")
	require.Len(t, item3.ID, 24, "Item ID should be 24 characters")
}

// Mock helper functions
func mockActiveItem() domain.Item {
	return domain.NewItem("any item", true)
}

func mockInactiveItem() domain.Item {
	return domain.NewItem("any item", false)
}

func mockItemWithID() domain.Item {
	return domain.Item{
		ID:     "some-id",
		Name:   "any name",
		Active: true,
	}
}

func mockEmptyItem() domain.Item {
	return domain.Item{
		ID:     "",
		Name:   "any name",
		Active: true,
	}
}

func mockItemWithWhitespaceID() domain.Item {
	return domain.Item{
		ID:     "   ",
		Name:   "any name",
		Active: true,
	}
}
