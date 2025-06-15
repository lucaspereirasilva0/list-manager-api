package handlers

import "net/http"

type ItemHandler interface {
	CreateItem(w http.ResponseWriter, r *http.Request) error
	GetItem(w http.ResponseWriter, r *http.Request) error
	UpdateItem(w http.ResponseWriter, r *http.Request) error
	DeleteItem(w http.ResponseWriter, r *http.Request) error
	ListItems(w http.ResponseWriter, r *http.Request) error
}
