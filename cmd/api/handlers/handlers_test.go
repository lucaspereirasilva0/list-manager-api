package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lucaspereirasilva0/list-manager-api/cmd/api/handlers"
	"github.com/lucaspereirasilva0/list-manager-api/cmd/api/handlers/middleware"
	"github.com/lucaspereirasilva0/list-manager-api/internal/domain"
	"github.com/lucaspereirasilva0/list-manager-api/internal/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	errDummy = errors.New("dummy error")
)

func TestCreateItem(t *testing.T) {
	tests := []struct {
		name                   string
		givenRequestBody       any
		givenServiceErr        error
		givenMockedServiceItem domain.Item
		wantAPIItem            handlers.Item
		wantHTTPStatus         int
		wantErr                error
	}{
		{
			name:                   "Given_Item_When_CreateItem_Then_ExpectedHTTPStatusCreated",
			givenRequestBody:       mockItem(),
			givenMockedServiceItem: mockServiceItem(),
			wantAPIItem:            mockAPIItem(),
			wantHTTPStatus:         http.StatusCreated,
		},
		{
			name:             "Given_InvalidJson_When_CreateItem_Then_ExpectedHTTPStatusBadRequest",
			givenRequestBody: mockInvalidJson(),
			wantHTTPStatus:   http.StatusBadRequest,
			wantErr:          handlers.NewDecodeRequestError(errors.New("invalid json")),
		},
		{
			name:             "Given_ItemWithMockedServiceError_When_CreateItem_Then_ExpectedHTTPStatusInternalServerError",
			givenRequestBody: mockItem(),
			givenServiceErr:  errDummy,
			wantHTTPStatus:   http.StatusInternalServerError,
			wantErr:          handlers.NewInternalServerError(errDummy),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			serviceMock := new(service.ItemServiceMock)
			serviceMock.On("CreateItem", mock.Anything, mock.Anything, mock.Anything).Return(tt.givenMockedServiceItem, tt.givenServiceErr)

			// Create handler with mock service
			h := handlers.NewHandler(serviceMock)

			// Create handler with middleware to simulate an HTTP request
			handlerWithMiddleware := middleware.ErrorHandlingMiddleware(h.CreateItem)

			// Create request
			body, err := json.Marshal(tt.givenRequestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			// Execute request using the handler with middleware
			handlerWithMiddleware.ServeHTTP(rec, req)

			// Assertions
			require.Equal(t, tt.wantHTTPStatus, rec.Code)

			if tt.wantErr != nil {
				// When an error is expected, the middleware writes the error message directly to the response body.
				require.Equal(t, parserAPIErr(t, tt.wantErr)+"\n", rec.Body.String())
			} else {
				// No error is expected, so the body should be empty for a 201 Created.
				require.Equal(t, tt.wantAPIItem, parserAPIItem(t, rec.Body.Bytes()))
			}
		})
	}
}

func TestGetItem(t *testing.T) {
	tests := []struct {
		name                   string
		givenItemID            string
		givenServiceErr        error
		givenMockedServiceItem domain.Item
		wantAPIItem            handlers.Item
		wantHTTPStatus         int
		wantErr                error
	}{
		{
			name:                   "Given_ItemID_When_GetItem_Then_ExpectedHTTPStatusOK",
			givenItemID:            "any-id",
			givenMockedServiceItem: mockServiceItem(),
			wantAPIItem:            mockAPIItem(),
			wantHTTPStatus:         http.StatusOK,
		},
		{
			name:           "Given_EmptyItemID_When_GetItem_Then_ExpectedHTTPStatusBadRequest",
			givenItemID:    "",
			wantHTTPStatus: http.StatusBadRequest,
			wantErr:        handlers.NewDecodeRequestError(nil),
		},
		{
			name:            "Given_ItemIDWithMockedServiceError_When_GetItem_Then_ExpectedHTTPStatusInternalServerError",
			givenItemID:     "any-id",
			givenServiceErr: errDummy,
			wantHTTPStatus:  http.StatusInternalServerError,
			wantErr:         handlers.NewInternalServerError(errDummy),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			serviceMock := new(service.ItemServiceMock)
			serviceMock.On("GetItem", mock.Anything, tt.givenItemID).Return(tt.givenMockedServiceItem, tt.givenServiceErr)

			// Create handler with mock service
			h := handlers.NewHandler(serviceMock)

			// Create handler with middleware to simulate an HTTP request
			handlerWithMiddleware := middleware.ErrorHandlingMiddleware(h.GetItem)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/item?id="+tt.givenItemID, nil)
			rec := httptest.NewRecorder()

			// Execute request using the handler with middleware
			handlerWithMiddleware.ServeHTTP(rec, req)

			// Assertions
			require.Equal(t, tt.wantHTTPStatus, rec.Code)

			if tt.wantErr != nil {
				// When an error is expected, the middleware writes the error message directly to the response body.
				require.Equal(t, parserAPIErr(t, tt.wantErr)+"\n", rec.Body.String())
			} else {
				// No error is expected, so the body should contain the item.
				require.Equal(t, tt.wantAPIItem, parserAPIItem(t, rec.Body.Bytes()))
			}
		})
	}
}

func TestUpdateItem(t *testing.T) {
	tests := []struct {
		name                   string
		givenItemID            string
		givenRequestBody       any
		givenServiceErr        error
		givenMockedServiceItem domain.Item
		wantAPIItem            handlers.Item
		wantHTTPStatus         int
		wantErr                error
	}{
		{
			name:                   "Given_ValidItem_When_UpdateItem_Then_ExpectedHTTPStatusOK",
			givenRequestBody:       mockItem(),
			givenMockedServiceItem: mockServiceItem(),
			wantAPIItem:            mockAPIItem(),
			wantHTTPStatus:         http.StatusOK,
		},
		{
			name:             "Given_InvalidJson_When_UpdateItem_Then_ExpectedHTTPStatusBadRequest",
			givenRequestBody: mockInvalidJson(),
			wantHTTPStatus:   http.StatusBadRequest,
			wantErr:          handlers.NewDecodeRequestError(errors.New("invalid json")),
		},
		{
			name:             "Given_ServiceError_When_UpdateItem_Then_ExpectedHTTPStatusInternalServerError",
			givenRequestBody: mockItem(),
			givenServiceErr:  errDummy,
			wantHTTPStatus:   http.StatusInternalServerError,
			wantErr:          handlers.NewInternalServerError(errDummy),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			serviceMock := new(service.ItemServiceMock)
			serviceMock.On("UpdateItem", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tt.givenMockedServiceItem, tt.givenServiceErr)

			// Create handler with mock service
			h := handlers.NewHandler(serviceMock)

			// Create handler with middleware to simulate an HTTP request
			handlerWithMiddleware := middleware.ErrorHandlingMiddleware(h.UpdateItem)

			// Create request
			body, err := json.Marshal(tt.givenRequestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/item", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			// Execute request using the handler with middleware
			handlerWithMiddleware.ServeHTTP(rec, req)

			// Assertions
			require.Equal(t, tt.wantHTTPStatus, rec.Code)

			if tt.wantErr != nil {
				require.Equal(t, parserAPIErr(t, tt.wantErr)+"\n", rec.Body.String())
			} else {
				require.Equal(t, tt.wantAPIItem, parserAPIItem(t, rec.Body.Bytes()))
			}
		})
	}
}

func TestDeleteItem(t *testing.T) {
	tests := []struct {
		name            string
		givenItemID     string
		givenServiceErr error
		wantHTTPStatus  int
		wantErr         error
	}{
		{
			name:           "Given_ItemID_When_DeleteItem_Then_ExpectedHTTPStatusNoContent",
			givenItemID:    "any-id",
			wantHTTPStatus: http.StatusNoContent,
		},
		{
			name:           "Given_EmptyItemID_When_DeleteItem_Then_ExpectedHTTPStatusBadRequest",
			givenItemID:    "",
			wantHTTPStatus: http.StatusBadRequest,
			wantErr:        handlers.NewDecodeRequestError(nil),
		},
		{
			name:            "Given_ServiceError_When_DeleteItem_Then_ExpectedHTTPStatusInternalServerError",
			givenItemID:     "any-id",
			givenServiceErr: errDummy,
			wantHTTPStatus:  http.StatusInternalServerError,
			wantErr:         handlers.NewInternalServerError(errDummy),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			serviceMock := new(service.ItemServiceMock)
			serviceMock.On("DeleteItem", mock.Anything, tt.givenItemID).Return(tt.givenServiceErr)

			// Create handler with mock service
			h := handlers.NewHandler(serviceMock)

			// Create handler with middleware to simulate an HTTP request
			handlerWithMiddleware := middleware.ErrorHandlingMiddleware(h.DeleteItem)

			// Create request
			req := httptest.NewRequest(http.MethodDelete, "/items?id="+tt.givenItemID, nil)
			rec := httptest.NewRecorder()

			// Execute request using the handler with middleware
			handlerWithMiddleware.ServeHTTP(rec, req)

			// Assertions
			require.Equal(t, tt.wantHTTPStatus, rec.Code)

			if tt.wantErr != nil {
				require.Equal(t, parserAPIErr(t, tt.wantErr)+"\n", rec.Body.String())
			} else {
				require.Empty(t, rec.Body.String())
			}
		})
	}
}

func TestListItems(t *testing.T) {
	tests := []struct {
		name                    string
		givenServiceErr         error
		givenMockedServiceItems []domain.Item
		wantAPIItems            []handlers.Item
		wantHTTPStatus          int
		wantErr                 error
	}{
		{
			name:                    "Given_Items_When_ListItems_Then_ExpectedHTTPStatusOK",
			givenMockedServiceItems: []domain.Item{mockServiceItem()},
			wantAPIItems:            []handlers.Item{mockAPIItem()},
			wantHTTPStatus:          http.StatusOK,
		},
		{
			name:                    "Given_NoItems_When_ListItems_Then_ExpectedHTTPStatusOK",
			givenMockedServiceItems: []domain.Item{},
			wantAPIItems:            []handlers.Item{},
			wantHTTPStatus:          http.StatusOK,
		},
		{
			name:            "Given_ServiceError_When_ListItems_Then_ExpectedHTTPStatusInternalServerError",
			givenServiceErr: errDummy,
			wantHTTPStatus:  http.StatusInternalServerError,
			wantErr:         handlers.NewInternalServerError(errDummy),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			serviceMock := new(service.ItemServiceMock)
			serviceMock.On("ListItems", mock.Anything).Return(tt.givenMockedServiceItems, tt.givenServiceErr)

			// Create handler with mock service
			h := handlers.NewHandler(serviceMock)

			// Create handler with middleware to simulate an HTTP request
			handlerWithMiddleware := middleware.ErrorHandlingMiddleware(h.ListItems)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/items", nil)
			rec := httptest.NewRecorder()

			// Execute request using the handler with middleware
			handlerWithMiddleware.ServeHTTP(rec, req)

			// Assertions
			require.Equal(t, tt.wantHTTPStatus, rec.Code)

			if tt.wantErr != nil {
				require.Equal(t, parserAPIErr(t, tt.wantErr)+"\n", rec.Body.String())
			} else {
				require.Equal(t, tt.wantAPIItems, parserAPIItems(t, rec.Body.Bytes()))
			}
		})
	}
}

func mockItem() handlers.Item {
	return handlers.Item{
		ID:     "any-id",
		Name:   "any name",
		Active: true,
	}
}

func mockServiceItem() domain.Item {
	return domain.Item{
		ID:     "any-id",
		Name:   "any name",
		Active: true,
	}
}

func mockAPIItem() handlers.Item {
	return handlers.Item{
		ID:     "any-id",
		Name:   "any name",
		Active: true,
	}
}

func mockInvalidJson() any {
	return "invalid json"
}

func parserAPIItem(t *testing.T, body []byte) handlers.Item {
	var item handlers.Item
	err := json.Unmarshal(body, &item)
	require.NoError(t, err)
	return item
}

func parserAPIErr(t *testing.T, err error) string {
	return err.(handlers.ErrorAPI).Message
}

func parserAPIItems(t *testing.T, body []byte) []handlers.Item {
	var items []handlers.Item
	err := json.Unmarshal(body, &items)
	require.NoError(t, err)
	return items
}
