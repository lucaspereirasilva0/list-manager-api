package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/lucaspereirasilva0/list-manager-api/internal/service"
)

// handler encapsulates the service and logger needed for the handlers
type handler struct {
	service service.ItemService
	parser  parser
}

// NewHandler creates a new instance of handler
func NewHandler(service service.ItemService) ItemHandler {
	return &handler{
		service: service,
		parser:  parser{},
	}
}

// CreateItem handles the creation of a new item
func (h *handler) CreateItem(w http.ResponseWriter, r *http.Request) error {
	var item Item

	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		return NewDecodeRequestError(err)
	}

	createdItem, err := h.service.CreateItem(ctx, h.parser.toDomainModel(item))
	if err != nil {
		return err
	}

	itemAPI := h.parser.toApiModel(createdItem)

	return writeJSONResponse(w, http.StatusCreated, itemAPI)
}

// GetItem handles the retrieval of an item by ID
func (h *handler) GetItem(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	id := r.URL.Query().Get("id")
	if id == "" {
		return NewDecodeRequestError(ErrIDRequired)
	}

	item, err := h.service.GetItem(ctx, id)
	if err != nil {
		return err
	}

	itemAPI := h.parser.toApiModel(item)

	return writeJSONResponse(w, http.StatusOK, itemAPI)
}

// UpdateItem handles the update of an item
func (h *handler) UpdateItem(w http.ResponseWriter, r *http.Request) error {
	var item Item

	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		return NewDecodeRequestError(err)
	}

	updatedItem, err := h.service.UpdateItem(ctx, h.parser.toDomainModel(item))
	if err != nil {
		return err
	}

	itemAPI := h.parser.toApiModel(updatedItem)

	return writeJSONResponse(w, http.StatusOK, itemAPI)
}

// DeleteItem handles the removal of an item
func (h *handler) DeleteItem(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	id := r.URL.Query().Get("id")
	if id == "" {
		return NewDecodeRequestError(ErrIDRequired)
	}

	err := h.service.DeleteItem(ctx, id)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

// ListItems handles the listing of all items
func (h *handler) ListItems(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	items, err := h.service.ListItems(ctx)
	if err != nil {
		return err
	}

	apiItems := make([]Item, len(items))
	for i, item := range items {
		apiItems[i] = h.parser.toApiModel(item)
	}

	return writeJSONResponse(w, http.StatusOK, apiItems)
}

// writeJSONResponse escreve uma resposta HTTP com o status code e corpo JSON.
func writeJSONResponse(w http.ResponseWriter, statusCode int, data any) error {
	rawData, err := json.Marshal(data)
	if err != nil {
		return NewInternalServerError(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_, err = w.Write(rawData)
	if err != nil {
		return NewInternalServerError(err)
	}

	return nil
}
